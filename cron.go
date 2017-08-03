package venom

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/robfig/cron.v2"
)

type cronConfig struct {
	Schedule     string `mapstructure:"schedule"`
	ScheduleHttp string `mapstructure:"schedule-http"`
}

func initCronFlags(flags *pflag.FlagSet) error {
	flags.String("schedule", "", "Schedule spec to schedule command for (e.g. every 1s = */1 * * * * *)")
	flags.String("schedule-http", "", "Address to run scheduler controls")
	return nil
}

type Func func(cmd *cobra.Command, args []string) error

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

		spec := cfg.Schedule
		exitOnError := true
		if strings.HasPrefix(spec, "-") {
			exitOnError = false
			spec = spec[1:]
		}

		jobs := sync.WaitGroup{}
		schedule := cron.New()

		_, err = schedule.AddFunc(spec, func() {
			jobs.Add(1)
			defer jobs.Done()

			err := runE(cmd, args)
			if err != nil {
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

		if cfg.ScheduleHttp != "" {
			go ListenAndServe(cfg.ScheduleHttp, schedule)
		}

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

func ListenAndServe(addr string, c *cron.Cron) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})
	mux.HandleFunc("/quitquitquit", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "OK\n")
		c.Stop()
		for _, e := range c.Entries() {
			c.Remove(e.ID)
		}
	})
	mux.HandleFunc("/abortabortabort", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})
	return http.ListenAndServe(addr, mux)
}
