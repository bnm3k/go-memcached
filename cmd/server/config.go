package main

import "flag"

type config struct {
	addr string
}

func getConfig() *config {
	cfg := config{}
	//cfg := new(config)
	flag.StringVar(&cfg.addr, "addr", ":4000", "http network address")
	flag.Parse()
	return &cfg
}
