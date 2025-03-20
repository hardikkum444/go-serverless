package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

type Metrics struct {
	devices prometheus.Gauge
}

func NewMatrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "connected_devices",
			Help:      "Number of currently connected devices",
		}),
	}
	reg.MustRegister(m.devices)
	return m
}

var Dvs []Device

func init() {
	Dvs = []Device{
		{1, "5F-33-CC-1F-43-82", "2.1.6"},
		{2, "EF-98-SC-4E-89-90", "2.1.6"},
	}
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(Dvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
