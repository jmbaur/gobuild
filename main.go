package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"git.jmbaur.com/gobuild/builder"
	"git.jmbaur.com/gobuild/config"
)

var (
	//go:embed favicon.ico
	favicon string
	//go:embed index.html
	index string
	templ = template.Must(template.New("index").Parse(index))
)

func main() {
	listen := flag.String("listen", "localhost:8080", "Port and address to bind server")
	configFile := flag.String("config", "", "Path to configuration file")
	flag.Parse()
	cfg, err := config.GetConfig(filepath.Join(*configFile))
	if err != nil {
		log.Fatal(err)
	}
	builders := []map[string]interface{}{}
	for _, cfgBuilder := range cfg.Builders {
		b := builder.NewBuilder(cfgBuilder.Name, cfgBuilder.Url, cfgBuilder.Commands)
		builders = append(builders, map[string]interface{}{
			"name": b.Name,
			"url":  b.Url,
		})
		http.Handle(b.Url, b)
		go b.Run()
	}
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, favicon); err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templ.Execute(w,
			map[string]interface{}{
				"builders": builders,
			}); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	log.Printf("Started server, listening on '%s'", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
