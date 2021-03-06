package venom

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/robfig/cron.v2"
)

type cronConfig struct {
	Schedule         string `mapstructure:"schedule"`
	ScheduleAfterRun bool   `mapstructure:"schedule-after-run"`
	ScheduleHttp     string `mapstructure:"schedule-http"`
}

type cronMetrics struct {
	Runs   prometheus.Counter
	Errors prometheus.Counter
	Skips  prometheus.Counter
}

func initCronMetrics() *cronMetrics {
	metrics := &cronMetrics{
		Errors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "venom",
			Subsystem: "cron",
			Name:      "errors_total",
			Help:      "Number of times execution has failed.",
		}),
		Skips: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "venom",
			Subsystem: "cron",
			Name:      "skips_total",
			Help:      "Number of times execution was skipped.",
		}),
		Runs: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "venom",
			Subsystem: "cron",
			Name:      "runs_total",
			Help:      "Number of times execution was attempted.",
		}),
	}

	prometheus.MustRegister(metrics.Errors)
	prometheus.MustRegister(metrics.Skips)
	prometheus.MustRegister(metrics.Runs)

	return metrics
}

func initCronFlags(flags *pflag.FlagSet) error {
	flags.String("schedule", "", "Schedule spec to schedule command for (e.g. every 1s = */1 * * * * *)")
	flags.Bool("schedule-after-run", false, "Schedule only after a first successful run (or fail)")
	flags.String("schedule-http", "", "Address to run scheduler controls")
	return nil
}

type Func func(cmd *cobra.Command, args []string) error

type GetCronFunc func() *cron.Cron

// TODO(akrasnukhin) This is racy
var getCron GetCronFunc

func CronRunE(runE Func, v *viper.Viper) Func {
	return func(cmd *cobra.Command, args []string) error {
		var cfg cronConfig
		err := Unmarshal(&cfg, v)
		if err != nil {
			return err
		}

		//
		// No cron schedule is provided, use function as is.
		//
		if cfg.Schedule == "" {
			return runE(cmd, args)
		}

		//
		// Register all required metrics, must be run once only.
		//
		metrics := initCronMetrics()

		//
		// Start serving health checks as soon as possible, this is to make sure
		// we are not killed during the initial run if --schedule-after-run is set.
		//
		if cfg.ScheduleHttp != "" {
			go ListenAndServe(cfg.ScheduleHttp)
		}

		if cfg.ScheduleAfterRun {
			metrics.Runs.Add(1)
			err := runE(cmd, args)
			if err != nil {
				metrics.Errors.Add(1)
				return err
			}
		}

		//
		// If cron spec starts with "-" we don't exit on errors. Very neat when you schedule
		// a command that periodically fails but you don't halt the whole process i.e. some kind
		// of poor man's command manager.
		//
		spec := cfg.Schedule
		exitOnError := true
		if strings.HasPrefix(spec, "-") {
			exitOnError = false
			spec = spec[1:]
		}

		//
		// WaitGroup allows us to wait for all jobs to complete on exit.
		//
		jobs := sync.WaitGroup{}
		schedule := cron.New()
		getCron = func() *cron.Cron {
			return schedule
		}

		var jobStartTime int64 = 0

		_, err = schedule.AddFunc(spec, func() {
			//
			// This prevents us from running this function twice.
			//
			now := time.Now()
			if !atomic.CompareAndSwapInt64(&jobStartTime, 0, now.UnixNano()) {
				startTime := time.Unix(0, atomic.LoadInt64(&jobStartTime))
				if now.After(startTime) { // required if go < 1.9 is used, see https://golang.org/doc/go1.9#monotonic-time
					logrus.WithField("duration", now.Sub(startTime)).Info("Skipping as it is still running.")
				} else {
					logrus.Info("Skipping as it is still running.")
				}
				metrics.Skips.Add(1)
				return
			}
			defer atomic.StoreInt64(&jobStartTime, 0)

			//
			// This allows us to wait for the function return on exit.
			//
			jobs.Add(1)
			defer jobs.Done()

			metrics.Runs.Add(1)
			err := runE(cmd, args)
			if err != nil {
				metrics.Errors.Add(1)
				logrus.WithError(err).Errorf("Error")
				if exitOnError {
					schedule.Stop()
					for _, e := range schedule.Entries() {
						schedule.Remove(e.ID)
					}
				}
			}
		})
		if err != nil {
			return err
		}
		schedule.Start()

		for {
			select {
			case <-time.After(time.Second):
				if len(schedule.Entries()) == 0 {
					jobs.Wait()
					return nil
				}
			}
		}
		return nil
	}
}

func QuitQuitQuit() {
	if getCron != nil {
		c := getCron()
		c.Stop()
		for _, e := range c.Entries() {
			c.Remove(e.ID)
		}
	}
}

func AbortAbortAbort() {
	QuitQuitQuit()
}

func ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})
	mux.HandleFunc("/quitquitquit", func(w http.ResponseWriter, req *http.Request) {
		QuitQuitQuit()
		fmt.Fprintf(w, "OK\n")
	})
	mux.HandleFunc("/abortabortabort", func(w http.ResponseWriter, req *http.Request) {
		AbortAbortAbort()
		fmt.Fprintf(w, "OK\n")
	})
	return http.ListenAndServe(addr, mux)
}
