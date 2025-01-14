package discovery

/**
 * static.go - static list discovery implementation
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"github.com/notional-labs/gobetween/src/config"
	"github.com/notional-labs/gobetween/src/core"
	"github.com/notional-labs/gobetween/src/logging"
	"github.com/notional-labs/gobetween/src/utils/parsers"
)

/**
 * Creates new static discovery
 */
func NewStaticDiscovery(cfg config.DiscoveryConfig) interface{} {
	d := Discovery{
		opts:  DiscoveryOpts{0},
		cfg:   cfg,
		fetch: staticFetch,
	}

	return &d
}

/**
 * Start discovery
 */
func staticFetch(cfg config.DiscoveryConfig) (*[]core.Backend, error) {
	log := logging.For("discovery/static")

	var backends []core.Backend
	for _, s := range cfg.StaticList {
		backend, err := parsers.ParseBackendDefault(s)
		if err != nil {
			log.Warn(err)
			continue
		}
		backends = append(backends, *backend)
	}

	return &backends, nil
}
