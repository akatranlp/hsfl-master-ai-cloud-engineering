package strategy

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
	"github.com/stretchr/testify/assert"
)

func TestIPHashStrategy(t *testing.T) {
	strategyImpl := NewIPHashStrategy([]*target.Target{})

	t.Run("SetTargets", func(t *testing.T) {
		t.Run("Should set the targets", func(t *testing.T) {
			// given
			URL, _ := url.Parse("http://localhost:8080")
			target1 := &target.Target{Url: URL, Health: 0, CurrentRequests: 0}
			targets := []*target.Target{target1}

			// when
			strategyImpl.SetTargets(targets)

			// then
			assert.Equal(t, targets, strategyImpl.targets)
		})
	})

	t.Run("NextServer", func(t *testing.T) {
		t.Run("Should return the server", func(t *testing.T) {
			// given
			URL, _ := url.Parse("http://localhost:8080")
			target1 := &target.Target{Url: URL, Health: 0}
			targets := []*target.Target{target1}

			r := &http.Request{RemoteAddr: "192.168.7.10:8080"}

			strategyImpl.SetTargets(targets)

			// when
			target := strategyImpl.NextTarget(r)

			// then
			assert.Equal(t, target1, target)
		})

		t.Run("Should return the first server", func(t *testing.T) {
			// given
			url1, _ := url.Parse("http://first-server:8080")
			target1 := &target.Target{Url: url1, Health: 0}

			url2, _ := url.Parse("http://second-server:8080")
			target2 := &target.Target{Url: url2, Health: 0}

			targets := []*target.Target{target1, target2}

			r := &http.Request{RemoteAddr: "192.168.7.10:8080"}

			strategyImpl.SetTargets(targets)

			// when
			target := strategyImpl.NextTarget(r)

			// then
			assert.Equal(t, target1, target)
		})

		t.Run("Should return the second server", func(t *testing.T) {
			// given
			url1, _ := url.Parse("http://first-server:8080")
			target1 := &target.Target{Url: url1, Health: 0}

			url2, _ := url.Parse("http://second-server:8080")
			target2 := &target.Target{Url: url2, Health: 0}

			targets := []*target.Target{target1, target2}

			strategyImpl.SetTargets(targets)

			r := &http.Request{RemoteAddr: "192.168.7.11:8080"}

			// when
			target := strategyImpl.NextTarget(r)

			// then
			assert.Equal(t, target2, target)
		})

		t.Run("Should return the first server", func(t *testing.T) {
			// given
			url1, _ := url.Parse("http://first-server:8080")
			target1 := &target.Target{Url: url1, Health: 0}

			url2, _ := url.Parse("http://second-server:8080")
			target2 := &target.Target{Url: url2, Health: 1}

			targets := []*target.Target{target1, target2}

			strategyImpl.SetTargets(targets)

			r := &http.Request{RemoteAddr: "192.168.7.11:8080"}

			// when
			target := strategyImpl.NextTarget(r)

			// then
			assert.Equal(t, target1, target)
		})
	})
}
