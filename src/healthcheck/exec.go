package healthcheck

/**
 * exec.go - Exec healthcheck
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"time"

	"github.com/notional-labs/gobetween/src/config"
	"github.com/notional-labs/gobetween/src/core"
	"github.com/notional-labs/gobetween/src/logging"
	"github.com/notional-labs/gobetween/src/utils"
)

/**
 * Exec healthcheck
 */
func exec(t core.Target, cfg config.HealthcheckConfig, result chan<- CheckResult) {
	log := logging.For("healthcheck/exec")

	execTimeout, _ := time.ParseDuration(cfg.Timeout)

	checkResult := CheckResult{
		Target: t,
	}

	out, err := utils.ExecTimeout(execTimeout, cfg.ExecCommand, t.Host, t.Port)
	if err != nil {
		// TODO: Decide better what to do in this case
		checkResult.Status = Unhealthy
		log.Warn(err)
	} else {
		if out == cfg.ExecExpectedPositiveOutput {
			checkResult.Status = Healthy
		} else if out == cfg.ExecExpectedNegativeOutput {
			checkResult.Status = Unhealthy
		} else {
			log.Warn("Unexpected output: ", out)
		}
	}

	select {
	case result <- checkResult:
	default:
		log.Warn("Channel is full. Discarding value")
	}
}
