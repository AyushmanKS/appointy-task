package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RedirectsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "morphlink_redirects_total",
	Help: "The total number of redirects processed",
})
