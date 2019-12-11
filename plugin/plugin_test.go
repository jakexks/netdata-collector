package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
)

func TestScrape(t *testing.T) {
	collector := New(client.DefaultClient)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := collector.scrape(ctx); err != nil {
		t.Error(err)
	}
}
