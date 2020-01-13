package main

import "flag"

type config struct {
	addr          string
	cacheType     string
	cacheCapacity int
}

func getConfig() *config {
	cfg := config{}
	//cfg := new(config)
	flag.StringVar(&cfg.addr, "addr", ":4000", "http network address")
	flag.StringVar(&cfg.cacheType, "cacheType", "lfu", "underlying cache type: [lfu, lru, lfu-lru-t]")
	flag.IntVar(&cfg.cacheCapacity, "cacheCapacity", 100, "cache capacity")
	flag.Parse()
	return &cfg
}
