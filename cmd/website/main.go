// Application website serves an easy configuratble website builder/host.
//
package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/jobstoit/tags/defaults"
	"github.com/jobstoit/tags/env"
	"github.com/jobstoit/tags/flags"
	"github.com/jobstoit/website/cmd/serve"
	"github.com/jobstoit/website/model"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg := readConfig()

	serve.Serve(cfg)
}

func readConfig() *model.Config {
	x := new(model.Config)

	defaultPath := `/etc/website/config.yaml`
	f := flag.String(`c`, defaultPath, `configuration file`)
	if !flag.Parsed() {
		flag.Parse()
	}

	if *f != `` {
		fb, err := ioutil.ReadFile(*f)
		if err != nil && *f != defaultPath {
			log.Fatalf("couldn't find file: %s", *f)
		}

		if len(fb) > 0 {
			if err := yaml.Unmarshal(fb, x); err != nil {
				log.Fatalf("error in given configuraion: %v", err)
			}
		}
	}

	if err := flags.Parse(x); err != nil {
		log.Fatalf("unrecognized formated flag: %v", err)
	}

	if err := env.Parse(x); err != nil {
		log.Fatalf("error getting environment: %v", err)
	}

	if err := defaults.Parse(x); err != nil {
		log.Fatalf("error parsing default: %v", err)
	}

	return x
}
