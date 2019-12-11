// Package plugin is a netdata external plugin that
// scrapes Go Micro services
package plugin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/micro/go-micro/client"
	stats "github.com/micro/micro/debug/stats/proto"
)

// Collector scrapes the go.micro.debug service for snapshots of stats
// and writes them to stdout in netdata format.
type Collector struct {
	client client.Client

	charts        []chart
	chartUpdate   bool
	knownServices map[string]bool

	sync.RWMutex
}

// New returns a configured collector
func New(c client.Client) *Collector {
	return &Collector{
		client:        c,
		knownServices: make(map[string]bool),
		charts: []chart{
			chart{
				id:          netdataType + ".go_micro_service_started",
				title:       "Start Time",
				units:       "timestamp",
				family:      "uptime",
				context:     "micro.service.started",
				priority:    "70000",
				updateEvery: "1",
				plugin:      netdataModule,
				module:      netdataModule,
			},
			chart{
				id:          netdataType + ".go_micro_service_uptime",
				title:       "Uptime",
				units:       "seconds",
				family:      "uptime",
				context:     "micro.service.uptime",
				priority:    "70001",
				updateEvery: "1",
				plugin:      netdataModule,
				module:      netdataModule,
			},
			chart{
				id:          netdataType + ".go_micro_service_memory",
				title:       "Heap Allocated",
				units:       "B",
				family:      "memory",
				context:     "micro.service.memory",
				priority:    "700002",
				updateEvery: "1",
				plugin:      netdataModule,
				module:      netdataModule,
			},
			chart{
				id:          netdataType + ".go_micro_service_threads",
				title:       "Goroutines",
				units:       "goroutines",
				family:      "threads",
				context:     "micro.service.threads",
				priority:    "70000",
				updateEvery: "1",
				plugin:      netdataModule,
				module:      netdataModule,
			},
			chart{
				id:          netdataType + ".go_micro_service_gcrate",
				title:       "GC Pause rate",
				units:       "nanoseconds/s",
				family:      "gc",
				context:     "micro.service.gc",
				priority:    "70000",
				updateEvery: "1",
				plugin:      netdataModule,
				module:      netdataModule,
			},
		},
	}
}

// Start starts collecting at the specified interval until the channel is closed
func (c *Collector) Start(interval time.Duration, done <-chan struct{}) {
	// Initialise the netdata charts
	for _, chart := range c.charts {
		fmt.Printf("CHART '%s' '' '%s' '%s' '%s' '%s' '' '%s' '%s' '' '%s' '%s'\n\n", chart.id, chart.title, chart.units, chart.family, chart.context, chart.priority, chart.updateEvery, chart.plugin, chart.module)
	}

	for {
		select {
		case <-done:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			c.scrape(ctx)
			<-ctx.Done()
			cancel()
		}
	}
}

// scrape collects a single snapshot of all stats
func (c *Collector) scrape(ctx context.Context) error {
	req := &stats.ReadRequest{}
	rsp := &stats.ReadResponse{}
	err := c.client.Call(ctx, client.NewRequest("go.micro.debug", "Stats.Read", req), rsp)
	if err != nil {
		return err
	}
	c.Lock()
	for _, s := range rsp.Stats {
		if _, found := c.knownServices[key(s)]; !found {
			c.knownServices[key(s)] = true
			c.chartUpdate = true
		}
	}

	return nil
}

func key(s *stats.Snapshot) string {
	return s.Service.Name + s.Service.Version + s.Service.Node.Id
}
