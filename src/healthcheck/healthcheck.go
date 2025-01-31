package healthcheck

/**
 * healthcheck.go - Healtheck
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"github.com/notional-labs/gobetween/src/config"
	"github.com/notional-labs/gobetween/src/core"
)

/**
 * Health Check function
 * Returns channel in which only one check result will be delivered
 */
type CheckFunc func(core.Target, config.HealthcheckConfig, chan<- CheckResult)

type HealthCheckStatus int32

const (
	Initial HealthCheckStatus = iota
	Unhealthy
	Healthy
)

/**
 * Check result
 * Handles target and it's live status
 */
type CheckResult struct {
	/* Check target */
	Target core.Target

	/* Check live status */
	Status HealthCheckStatus
}

/**
 * Healthcheck
 */
type Healthcheck struct {
	/* Healthcheck function */
	check CheckFunc

	/* Healthcheck configuration */
	cfg config.HealthcheckConfig

	/* Input channel to accept targets */
	In chan []core.Target

	/* Output channel to send check results for individual target */
	Out chan CheckResult

	/* Current check workers */
	workers []*Worker

	/* Channel to handle stop */
	stop chan bool
}

/**
 * Registry of factory methods
 */
var registry = make(map[string]CheckFunc)

/**
 * Initialize type registry
 */
func init() {
	registry["ping"] = ping
	registry["probe"] = probe
	registry["exec"] = exec
	registry["none"] = nil
}

/**
 * Create new Discovery based on strategy
 */
func New(strategy string, cfg config.HealthcheckConfig) *Healthcheck {
	check := registry[strategy]

	/* Create healthcheck */

	h := Healthcheck{
		check:   check,
		cfg:     cfg,
		In:      make(chan []core.Target),
		Out:     make(chan CheckResult),
		workers: []*Worker{},
		stop:    make(chan bool),
	}

	return &h
}

/**
 * Start healthcheck
 */
func (this *Healthcheck) Start() {
	go func() {
		for {
			select {

			/* got new targets */
			case targets := <-this.In:
				this.UpdateWorkers(targets)

			/* got stop requst */
			case <-this.stop:

				// Stop all workers
				for i := range this.workers {
					this.workers[i].Stop()
				}

				// And free it's memory
				this.workers = []*Worker{}

				return
			}
		}
	}()
}

/**
 * Sync current workers to represent healtcheck on targets
 * Will remove not needed workers, and add needed
 */
func (this *Healthcheck) UpdateWorkers(targets []core.Target) {
	result := []*Worker{}

	// Keep or add needed workers
	for _, t := range targets {
		var keep *Worker
		for i := range this.workers {
			c := this.workers[i]
			if t.EqualTo(c.target) {
				keep = c
				break
			}
		}

		if keep == nil {
			keep = &Worker{
				target: t,
				stop:   make(chan bool),
				out:    this.Out,
				cfg:    this.cfg,
				check:  this.check,
				LastResult: CheckResult{
					Status: Initial,
				},
			}
			keep.Start()
		}
		result = append(result, keep)
	}

	// Stop needed workers
	for i := range this.workers {
		c := this.workers[i]
		remove := true
		for _, t := range targets {
			if c.target.EqualTo(t) {
				remove = false
				break
			}
		}

		if remove {
			c.Stop()
		}
	}

	this.workers = result
}

func (this *Healthcheck) HasCheck() bool {
	return this.cfg.Kind != "none"
}

func (this *Healthcheck) InitialBackendHealthCheckStatus() HealthCheckStatus {
	if !this.HasCheck() {
		return Healthy
	}
	if this.cfg.InitialStatus != nil {
		switch *this.cfg.InitialStatus {
		case "unhealthy":
			return Unhealthy
		case "healthy":
			return Healthy
		default:
			panic("Healthcheck invalid initial status, this should have been validated in manager, but has invalid value " + *this.cfg.InitialStatus)
		}
	}
	return Healthy
}

/**
 * Stop healthcheck
 */
func (this *Healthcheck) Stop() {
	this.stop <- true
}
