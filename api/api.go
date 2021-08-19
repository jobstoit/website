// Package api contains all the API requests
//
package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.fuyu.moe/fuyu/router"
	"github.com/jobstoit/website/model"
	"github.com/jobstoit/website/repo"
)

const (
	paramSiteID = "sid"
	paramPageID = "pid"
	paramRowID  = "rid"
)

// Append adds the api routes to the router
func Append(rtr *router.Router, cfg *model.Config) {
	x := new(a)
	x.repo = repo.New(cfg.DBCS)

	oauth := newOauth(cfg.Port, cfg.OID.StateString, cfg.OID.ClientID, cfg.OID.ClientSecret, cfg.OID.URL)
	x.oauth = oauth

	rtr.GET("/login", oauth.login)
	rtr.GET("/callback", oauth.oauthCallback)

	api := rtr.Group("/api")
	site := api.Group("/site")

	site.GET("/active", x.getActiveSite)

	admin := api.Group("/admin", x.isAdminMiddleware)
	adminSite := admin.Group("/site")

	adminSite.GET("/:"+paramSiteID, x.getSiteByID)
	adminSite.GET("/list", x.adminListSites)
	adminSite.POST("/create", x.adminAddSite)
	adminSite.POST("/:"+paramPageID+"/navigation", x.adminAddNavigation)
	adminSite.PUT("/:"+paramPageID+"/navigation", x.adminUpdateNavigationSequence)
	adminSite.POST("/:"+paramSiteID+"/page", x.adminAddPage)
	adminSite.POST("/page/:"+paramPageID+"/row", x.adminAddRow)
	adminSite.PUT("/page/:"+paramPageID+"/row", x.adminUpdateRowSequence)
	adminSite.PATCH("/page/:"+paramPageID+"/row/:"+paramRowID, x.adminUpdateRow)
	adminSite.DELETE("/page/:"+paramPageID+"/row/:"+paramRowID, x.adminDeleteRow)
}

// An api struct to hold the repo for all the router functions
type a struct {
	repo  *repo.Repo
	oauth *oa
}

func (x a) isAdminMiddleware(f router.Handle) router.Handle {
	return func(ctx *router.Context) error {
		user, err := x.oauth.GetUserInfo(ctx)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}

		if !strings.Contains(strings.Join(user.Roles, " "), "admin") {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		return f(ctx)
	}
}

func (x a) getActiveSite(ctx *router.Context) error {
	s := x.repo.GetActiveSite(ctx.Request.Context())
	return ctx.JSON(http.StatusOK, s)
}

func (x a) getSiteByID(ctx *router.Context) error {
	id, err := strconv.Atoi(ctx.Param(paramSiteID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	s := x.repo.GetSiteByID(ctx.Request.Context(), id)

	return ctx.JSON(http.StatusOK, s)
}

type adminAddNavigationRequest struct {
	URI      string `json:"uri"`
	Label    string `json:"label"`
	Position string `json:"position"`
}

func (x a) adminAddNavigation(ctx *router.Context, reqBody adminAddNavigationRequest) error {
	siteID, err := strconv.Atoi(ctx.Param(paramSiteID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	var res idResp
	res.ID = x.repo.AddNavigationLink(ctx.Request.Context(), siteID, reqBody.URI, reqBody.Label, reqBody.Position)

	return ctx.JSON(http.StatusCreated, res)
}

type adminUpdateNavigationSequenceRequest struct {
	Position string `json:"position"`
	IDs      []int  `json:"ids"`
}

func (x a) adminUpdateNavigationSequence(ctx *router.Context, reqBody adminUpdateNavigationSequenceRequest) error {
	siteID, err := strconv.Atoi(ctx.Param(paramSiteID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	x.repo.ChangeNavigationSequence(ctx.Request.Context(), siteID, reqBody.Position, reqBody.IDs)

	return ctx.NoContent(http.StatusOK)
}

func (x a) adminListSites(ctx *router.Context) error {
	sts := x.repo.ListSites(ctx.Request.Context())
	return ctx.JSON(http.StatusOK, sts)
}

type idResp struct {
	ID int
}

type adminAddSiteReq struct {
	Name string `json:"name"`
}

func (x a) adminAddSite(ctx *router.Context, reqBody adminAddSiteReq) error {
	user, err := x.oauth.GetUserInfo(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	x.repo.CreateSite(ctx.Request.Context(), reqBody.Name, user.Username)

	return ctx.NoContent(http.StatusCreated)
}

type adminAddPageReq struct {
	URI   string `json:"uri"`
	Label string `json:"label"`
}

func (x a) adminAddPage(ctx *router.Context, reqBody adminAddPageReq) error {
	ssiteID := ctx.Param(paramSiteID)
	siteID, err := strconv.Atoi(ssiteID)
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	var res idResp
	res.ID = x.repo.CreatePage(ctx.Request.Context(), siteID, reqBody.URI, reqBody.Label)

	return ctx.JSON(http.StatusCreated, res)
}

type adminAddRowReq struct {
	Titles  []string       `json:"titles"`
	Texts   []string       `json:"texts"`
	Media   []model.Medium `json:"media"`
	Buttons []model.Button `json:"buttons"`
}

func (x a) adminAddRow(ctx *router.Context, reqBody adminAddRowReq) error {
	pageID, err := strconv.Atoi(ctx.Param(paramPageID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	var res idResp
	res.ID = x.repo.CreateRow(ctx.Request.Context(), pageID, reqBody.Titles, reqBody.Texts, reqBody.Media, reqBody.Buttons)

	return ctx.JSON(http.StatusCreated, res)
}

type adminUpdateRowSequenceReq struct {
	RowIDs []int `json:"row_ids"`
}

func (x a) adminUpdateRowSequence(ctx *router.Context, reqBody adminUpdateRowSequenceReq) error {
	pageID, err := strconv.Atoi(ctx.Param(paramPageID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	x.repo.ChangeRowSequence(ctx.Request.Context(), pageID, reqBody.RowIDs)

	return ctx.NoContent(http.StatusOK)
}

type adminUpdateRowReq struct {
	Titles  []string       `json:"titles"`
	Texts   []string       `json:"texts"`
	Media   []model.Medium `json:"media"`
	Buttons []model.Button `json:"buttons"`
}

func (x a) adminUpdateRow(ctx *router.Context, reqBody adminUpdateRowReq) error {
	rowID, err := strconv.Atoi(ctx.Param(paramRowID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	x.repo.UpdateRow(ctx.Request.Context(), rowID, reqBody.Titles, reqBody.Texts, reqBody.Media, reqBody.Buttons)

	return ctx.NoContent(http.StatusOK)
}

func (x a) adminDeleteRow(ctx *router.Context) error {
	rowID, err := strconv.Atoi(ctx.Param(paramRowID))
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	x.repo.DeleteRow(ctx.Request.Context(), rowID)

	return ctx.NoContent(http.StatusOK)
}
