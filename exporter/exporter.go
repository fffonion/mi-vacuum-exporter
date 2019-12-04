package exporter

import (
	"time"

	"github.com/fffonion/mi-vacuum-exporter/miio"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	target string
	client *miio.MiioClient

	metricsUp,
	metricsBattery,
	metricsFanPower,
	metricsInCleaning,
	metricsState,
	metricsCleanTime,
	metricsConsumables *prometheus.Desc
}

type ExporterTarget struct {
	Host  string
	Token string
}

func NewExporter(t *ExporterTarget) (*Exporter, error) {
	c, err := miio.New(&miio.MiioClientConfig{
		Host:  t.Host,
		Token: t.Token,
	})
	if err != nil {
		return nil, err
	}
	err = c.Init()
	if err != nil {
		return nil, err
	}

	var (
		constLabels = prometheus.Labels{
			"device_id": c.ID(),
		}
	)

	e := &Exporter{
		target: t.Host,
		client: c,
		metricsUp: prometheus.NewDesc("vacuum_online",
			"Device online.",
			nil, constLabels,
		),
		metricsBattery: prometheus.NewDesc("vacuum_battery",
			"Device battery level.",
			nil, constLabels,
		),
		metricsFanPower: prometheus.NewDesc("vacuum_fan_power",
			"Device fan power.",
			nil, constLabels,
		),
		metricsInCleaning: prometheus.NewDesc("vacuum_cleaning",
			"Device is doing cleaning.",
			nil, constLabels,
		),
		metricsState: prometheus.NewDesc("vacuum_state",
			"Device state.",
			[]string{"state"}, constLabels,
		),

		metricsConsumables: prometheus.NewDesc("vacuum_consumables",
			"Device consumables use time.",
			[]string{"type"}, constLabels,
		),
	}
	return e, nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metricsUp

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	status := &miio.GetStatusResponse{}
	err := e.client.RPC("get_status", nil, status)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.metricsUp, prometheus.GaugeValue,
			0)
		return
	}
	ch <- prometheus.MustNewConstMetric(e.metricsBattery, prometheus.GaugeValue,
		float64(status.Battery))
	ch <- prometheus.MustNewConstMetric(e.metricsFanPower, prometheus.GaugeValue,
		float64(status.FanPower))

	ch <- prometheus.MustNewConstMetric(e.metricsInCleaning, prometheus.GaugeValue,
		float64(status.InCleaning))

	ch <- prometheus.MustNewConstMetric(e.metricsState, prometheus.GaugeValue,
		1, status.State.String())

	// slow down, vacuum will cry if request comes in too fast
	time.Sleep(time.Second)

	consumables := &miio.GetConsumableResponse{}
	err = e.client.RPC("get_consumable", nil, consumables)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.metricsUp, prometheus.GaugeValue,
			0)
		return
	}
	ch <- prometheus.MustNewConstMetric(e.metricsConsumables, prometheus.GaugeValue,
		float64(consumables.MainBrushWorkTime), "main_brush")
	ch <- prometheus.MustNewConstMetric(e.metricsConsumables, prometheus.GaugeValue,
		float64(consumables.SideBrushWorkTime), "side_brush")
	ch <- prometheus.MustNewConstMetric(e.metricsConsumables, prometheus.GaugeValue,
		float64(consumables.SensorDirtyTime), "sensor")
	ch <- prometheus.MustNewConstMetric(e.metricsConsumables, prometheus.GaugeValue,
		float64(consumables.FilterWorkTime), "filter")

	ch <- prometheus.MustNewConstMetric(e.metricsUp, prometheus.GaugeValue,
		1)
}
