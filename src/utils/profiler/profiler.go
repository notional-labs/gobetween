package profiler

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/notional-labs/gobetween/src/logging"
)

func Start(bind string) {
	log := logging.For("profiler")

	log.Infof("Starting profiler: %v", bind)

	go func() {
		log.Error(http.ListenAndServe(bind, nil))
	}()
}
