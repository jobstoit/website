package serve

import (
	"fmt"

	"git.fuyu.moe/fuyu/router"
	"github.com/jobstoit/website/api"
	"github.com/jobstoit/website/model"
)

func Serve(cfg *model.Config) {
	rtr := router.New()

	api.Append(rtr, cfg.DBCS, cfg.OIDP)

	rtr.Start(fmt.Sprintf(":%d", cfg.Port))
}
