package temporal

import (
	"fmt"
	"net/http"
	"strconv"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

func (a *ActivityRequester) setRedditPromMetrics(l log.Logger, mh client.MetricsHandler, h http.Header) {

	// Set Prometheus metrics. The ones we're interested in for Reddit are
	// the X-Requestlimit-* header values. Range over metrics and set them.
	mnames := []string{
		MetricXRatelimitUsed,
		MetricXRatelimitRemaining,
		MetricXRatelimitReset,
	}
	for _, mk := range mnames {
		g := mh.Gauge(mk)

		var val float64
		var err error

		switch mk {
		case MetricXRatelimitUsed:
			val, err = strconv.ParseFloat(h.Get("X-Ratelimit-Used"), 64)
			if err != nil {
				// debug only, not all requests will include rate limit headers
				l.Debug(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("X-Ratelimit-Used")))
				continue
			}
		case MetricXRatelimitRemaining:
			val, err = strconv.ParseFloat(h.Get("X-Ratelimit-Remaining"), 64)
			if err != nil {
				// debug only, not all requests will include rate limit headers
				l.Debug(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("X-Ratelimit-Remaining")))
				continue
			}
		case MetricXRatelimitReset:
			val, err = strconv.ParseFloat(h.Get("X-Ratelimit-Reset"), 64)
			if err != nil {
				// debug only, not all requests will include rate limit headers
				l.Debug(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("X-Ratelimit-Reset")))
				continue
			}
		}

		g.Update(val)
	}
}

func (a *ActivityRequester) setTwitchPromMetrics(l log.Logger, mh client.MetricsHandler, h http.Header) {

	// Set Prometheus metrics. The ones for Twitch are here:
	// https://dev.twitch.tv/docs/api/guide/#how-it-works
	mnames := []string{
		MetricXRatelimitLimit,
		MetricXRatelimitRemaining,
		MetricXRatelimitReset,
	}
	for _, mk := range mnames {
		g := mh.Gauge(mk)

		var val float64
		var err error

		// NOTE: twitch deviates a little from how these are conventionally supplied,
		// but we're fudging them a bit here to reduce the number of metrics and keep
		// things simple on our end. The main things are that Twitch ratelimit headers
		// don't have the `X-` prefix, and the first one is Limit rather than Used.
		// https://dev.twitch.tv/docs/api/guide/#how-it-works

		switch mk {
		case MetricXRatelimitLimit:
			val, err = strconv.ParseFloat(h.Get("Ratelimit-Limit"), 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("Ratelimit-Limit")))
				continue
			}
		case MetricXRatelimitRemaining:
			val, err = strconv.ParseFloat(h.Get("Ratelimit-Remaining"), 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("Ratelimit-Remaining")))
				continue
			}
		case MetricXRatelimitReset:
			val, err = strconv.ParseFloat(h.Get("Ratelimit-Reset"), 64)
			if err != nil {
				l.Error(fmt.Sprintf("failed to parse %s float from %s", mk, h.Get("Ratelimit-Reset")))
				continue
			}
		}

		g.Update(val)
	}
}
