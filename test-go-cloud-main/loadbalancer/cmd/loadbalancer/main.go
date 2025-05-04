package main

import (
	"flag"
	"loabalancer/internal/balancer"
	"loabalancer/internal/config"
	"log"
	"path/filepath"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file (JSON or YAML)")

	flag.Parse()

	ext := filepath.Ext(*configPath)
	var conf *config.Config
	var err error

	switch ext {
	case ".json":
		conf, err = config.LoadConfigFromJSON(*configPath)
	case ".yaml", ".yml":
		conf, err = config.LoadConfigFromYAML(*configPath)
	default:
		log.Fatalf("Unsupported config file format : %s", ext)
	}

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	lb, err := balancer.NewLoadBalancer(conf)
	if err != nil {
		log.Fatal(err)
	}

	if err1 := lb.Start(); err1 != nil {
		log.Printf("Load balancer stopped: %v", err1)
	}
}
