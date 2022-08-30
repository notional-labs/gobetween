/**
 * main.go - entry point
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */
package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/notional-labs/gobetween/src/api"
	"github.com/notional-labs/gobetween/src/cmd"
	"github.com/notional-labs/gobetween/src/config"
	"github.com/notional-labs/gobetween/src/info"
	"github.com/notional-labs/gobetween/src/logging"
	"github.com/notional-labs/gobetween/src/manager"
	"github.com/notional-labs/gobetween/src/metrics"
	"github.com/notional-labs/gobetween/src/utils/codec"
)

/**
 * version,revision,branch should be set while build using ldflags (see Makefile)
 */
var (
	version  string
	revision string
	branch   string
)

/**
 * Initialize package
 */
func init() {
	// Set GOMAXPROCS if not set
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// Init random seed
	rand.Seed(time.Now().UnixNano())

	// Save info
	info.Version = version
	info.Revision = revision
	info.Branch = branch
	info.StartTime = time.Now()
}

/**
 * Entry point
 */
func main() {
	log.Printf("gobetween v%s", version)

	env := os.Getenv("GOBETWEEN")
	if env != "" && len(os.Args) > 1 {
		log.Fatal("Passed GOBETWEEN env var and command-line arguments: only one allowed")
	}

	// Try parse env var to args
	if env != "" {
		a := []string{}
		if err := codec.Decode(env, &a, "json"); err != nil {
			log.Fatal("Error converting env var to parameters: ", err, " ", env)
		}
		os.Args = append([]string{""}, a...)
		log.Println("Using parameters from env var: ", os.Args)
	}

	// Process flags and start
	cmd.Execute(func(cfg *config.Config) {
		// Configure logging
		logging.Configure(cfg.Logging.Output, cfg.Logging.Level, cfg.Logging.Format)

		// Start manager
		manager.Initialize(*cfg)

		/* setup metrics */
		metrics.Start((*cfg).Metrics)

		// Start API
		api.Start((*cfg).Api)

		// block forever
		<-(chan string)(nil)
	})
}
