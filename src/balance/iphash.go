package balance

/**
 * iphash.go - iphash balance impl
 *
 * @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
 */

import (
	"errors"
	"hash/fnv"
	"sort"

	"github.com/notional-labs/gobetween/src/core"
)

/**
 * Iphash balancer
 */
type IphashBalancer struct{}

/**
 * Elect backend using iphash strategy
 * Using fnv1a for speed
 *
 * TODO: Improve as needed
 */
func (b *IphashBalancer) Elect(context core.Context, backends []*core.Backend) (*core.Backend, error) {
	if len(backends) == 0 {
		return nil, errors.New("Can't elect backend, Backends empty")
	}

	sort.SliceStable(backends, func(i, j int) bool {
		return backends[i].Target.String() < backends[j].Target.String()
	})

	hash := fnv.New32a()
	hash.Write(context.Ip())
	backend := backends[hash.Sum32()%uint32(len(backends))]

	return backend, nil
}
