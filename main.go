package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"git.jmbaur.com/gobuild/builder"
	"git.jmbaur.com/gobuild/config"
)

func main() {
	port := flag.Int("port", 8080, "Port to bind server to")
	configFile := flag.String("config", "", "Path to configuration file")
	flag.Parse()
	cfg, err := config.GetConfig(filepath.Join(*configFile))
	if err != nil {
		log.Fatal(err)
	}
	for _, cfgBuilder := range cfg.Builders {
		b := builder.NewBuilder(cfgBuilder.Name, cfgBuilder.Url, cfgBuilder.Commands)
		http.Handle(b.Url, b)
		go b.Run()
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
