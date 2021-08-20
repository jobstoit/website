package serve

import (
	"fmt"
	"log"

	"git.fuyu.moe/fuyu/router"
	"github.com/jobstoit/website/api"
	"github.com/jobstoit/website/model"
)

func Serve(cfg *model.Config) {
	rtr := router.New()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	api.Append(rtr, cfg)

	rtr.Start(fmt.Sprintf(":%d", cfg.Port))
}
