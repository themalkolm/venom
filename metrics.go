package venom

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type metricsConfig struct {
	Export        bool   `mapstructure:"metrics-export"`
	ListenAddress string `mapstructure:"metrics-listen-address"`
	HttpPath      string `mapstructure:"metrics-path"`
}

func initMetricsFlags(flags *pflag.FlagSet) error {
	flags.Bool("metrics-export", false, "Turn on metrics exporter.")
	flags.String("metrics-listen-address", ":9236", "Address on which to expose metrics and web interface.")
	flags.String("metrics-path", "/metrics", "Path under which to expose metrics.")
	return nil
}

func ListenAndServeMetrics(name string, v *viper.Viper) error {
	var cfg metricsConfig
	err := Unmarshal(&cfg, v)
	if err != nil {
		return err
	}

	if !cfg.Export {
		return nil
	}

	handler := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{ErrorLog: logrus.StandardLogger()},
	)

	mux := http.NewServeMux()
	mux.Handle(cfg.HttpPath, handler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			fmt.Sprintf(`<html>
				<head><title>%s exporter</title></head>
				<body>
				<h1>%s exporter</h1>
				<p><a href="`+cfg.HttpPath+`">metrics</a></p>
				</body>
				</html>`, name, name),
		))
	})

	err = http.ListenAndServe(cfg.ListenAddress, mux)
	if err != nil {
		logrus.WithError(err).Error("Failed to start metrics http server.")
	}
	return err
}
