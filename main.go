package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ghodss/yaml"
	vegeta "github.com/tsenart/vegeta/lib"
)

type Config struct {
	BaseURL string `json:"base_url"`
	Tests   []Test `json:"tests"`
}

type Test struct {
	Name     string `json:"name"`
	Rate     uint64 `json:"rps"`
	Duration uint64 `json:"duration"`
	SLA      SLA    `json:"sla"`
	Target   Target `json:"target"`
}

type SLA struct {
	Latency     int64   `json:"latency"`
	SuccessRate float64 `json:"success_rate"`
}

type Target struct {
	Method  string   `json:"method"`
	Path    string   `json:"path"`
	Headers []Header `json:"headers"`
	Body    []byte   `json:"body"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	var baseURL string
	var configFileName string

	flag.StringVar(&baseURL, "baseURL", "", "Override baseURL config property")
	flag.StringVar(&configFileName, "config", "stress.test.json", "Config source file")

	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := Config{}

	configFromFile(configFileName, &config)

	if baseURL != "" {
		config.BaseURL = baseURL
	}

	for _, test := range config.Tests {
		duration := time.Duration(test.Duration) * time.Second

		headers := &http.Header{}

		for _, header := range test.Target.Headers {
			headers.Set(header.Name, header.Value)
		}

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: test.Target.Method,
			URL:    config.BaseURL + test.Target.Path,
			Header: *headers,
			Body:   test.Target.Body,
		})
		attacker := vegeta.NewAttacker()

		metrics := vegeta.Metrics{}
		for res := range attacker.Attack(targeter, test.Rate, duration, test.Name) {
			metrics.Add(res)
		}
		metrics.Close()

		reporter := vegeta.NewTextReporter(&metrics)
		reporter.Report(os.Stdout)

		if metrics.Success*100 < test.SLA.SuccessRate {
			log.Fatalf("SLA success rate not met: expected: %%%f, got: %%%f", metrics.Success*100, test.SLA.SuccessRate)
		}

		if metrics.Latencies.P99.Nanoseconds() > test.SLA.Latency*time.Millisecond.Nanoseconds() {
			log.Fatalf("SLA latency not met: expected: %d, got: %d", metrics.Latencies.P99.Nanoseconds(), test.SLA.Latency*time.Millisecond.Nanoseconds())
		}
	}
}

// configFromFile reads in a file path and unmarshalls to Config
func configFromFile(cfg string, config *Config) error {
	file, err := os.Open(cfg)
	if err != nil {
		return err
	}

	cfgBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	ext := filepath.Ext(cfg)
	if ext == ".json" {
		cfgBytes, err = yaml.JSONToYAML(cfgBytes)
		if err != nil {
			return err
		}
	}

	err = yaml.Unmarshal(cfgBytes, config)
	if err != nil {
		return err
	}
	return nil
}
