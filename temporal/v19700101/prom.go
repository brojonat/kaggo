package temporal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"go.temporal.io/sdk/log"
)

func (a *ActivityRequester) setRedditPromMetrics(l log.Logger, labels prometheus.Labels, h http.Header) {

	// Set Prometheus metrics. The ones we're interested in for Reddit are
	// the X-Requestlimit-* header values. Range over metrics and set them.
	mnames := []string{
		PromMetricXRatelimitUsed,
		PromMetricXRatelimitRemaining,
		PromMetricXRatelimitReset,
	}
	for _, mk := range mnames {
		gv, ok := a.Metrics[mk].(*prometheus.GaugeVec)
		if !ok {
			l.Error(fmt.Sprintf("failed to locate prom metric %s, skipping", mk))
			continue
		}

		c, err := gv.GetMetricWith(labels)
		if err != nil {
			// GetMetricWith is a get-or-create operation, this should never happen
			l.Error(fmt.Sprintf("failed to get prom metric %s with labels: %s", mk, labels))
			continue
		}

		var val float64

		switch mk {
		case PromMetricXRatelimitUsed:
			val, err = strconv.ParseFloat(h.Get("X-Ratelimit-Used"), 64)
			if err != nil {
				// debug only, not all requests will include rate limit headers
				l.Debug(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("X-Ratelimit-Used")))
				continue
			}
		case PromMetricXRatelimitRemaining:
			val, err = strconv.ParseFloat(h.Get("X-Ratelimit-Remaining"), 64)
			if err != nil {
				// debug only, not all requests will include rate limit headers
				l.Debug(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("X-Ratelimit-Remaining")))
				continue
			}
		case PromMetricXRatelimitReset:
			val, err = strconv.ParseFloat(h.Get("X-Ratelimit-Reset"), 64)
			if err != nil {
				// debug only, not all requests will include rate limit headers
				l.Debug(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("X-Ratelimit-Reset")))
				continue
			}
		}

		c.Set(val)
	}
}

func (a *ActivityRequester) setTwitchPromMetrics(l log.Logger, labels prometheus.Labels, h http.Header) {

	// Set Prometheus metrics. The ones for Twitch are here:
	// https://dev.twitch.tv/docs/api/guide/#how-it-works
	mnames := []string{
		PromMetricXRatelimitLimit,
		PromMetricXRatelimitRemaining,
		PromMetricXRatelimitReset,
	}
	for _, mk := range mnames {
		gv, ok := a.Metrics[mk].(*prometheus.GaugeVec)
		if !ok {
			l.Error(fmt.Sprintf("failed to locate prom metric %s, skipping", mk))
			continue
		}

		c, err := gv.GetMetricWith(labels)
		if err != nil {
			// GetMetricWith is a get-or-create operation, this should never happen
			l.Error(fmt.Sprintf("failed to get prom metric %s with labels: %s", mk, labels))
			continue
		}

		var val float64

		// NOTE: twitch deviates a little from how these are conventionally supplied,
		// but we're fudging them a bit here to reduce the number of metrics and keep
		// things simple on our end. The main things are that Twitch ratelimit headers
		// don't have the `X-` prefix, and the first one is Limit rather than Used.
		// https://dev.twitch.tv/docs/api/guide/#how-it-works

		switch mk {
		case PromMetricXRatelimitLimit:
			val, err = strconv.ParseFloat(h.Get("Ratelimit-Limit"), 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("Ratelimit-Limit")))
				continue
			}
		case PromMetricXRatelimitRemaining:
			val, err = strconv.ParseFloat(h.Get("Ratelimit-Remaining"), 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("Ratelimit-Remaining")))
				continue
			}
		case PromMetricXRatelimitReset:
			val, err = strconv.ParseFloat(h.Get("Ratelimit-Reset"), 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("Ratelimit-Reset")))
				continue
			}
		}

		c.Set(val)
	}
}
