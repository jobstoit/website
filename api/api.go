// Copyright 2021 Job Stoit. All rights reserved.

package api

import (
	"fmt"
	"net/http"
	"strconv"

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
func Append(rtr *router.Router, dbcs string) {
	x := new(a)
	x.repo = repo.New(dbcs)

	api := rtr.Group("/api")
	site := api.Group("/site")

	site.GET("/active", x.getActiveSite)

	admin := api.Group("/admin")
	adminSite := admin.Group("/site")

	adminSite.GET("/:"+paramSiteID, x.getSiteByID)
	adminSite.GET("/list", x.adminListSites)
	adminSite.POST("/create", x.adminAddSite)
	adminSite.POST("/:"+paramSiteID+"/page", x.adminAddPage)
	adminSite.POST("/page/:"+paramPageID+"/row", x.adminAddRow)
	adminSite.PUT("/page/:"+paramPageID+"/row", x.adminUpdateRowSequence)
	adminSite.PATCH("/page/:"+paramPageID+"/row/:"+paramRowID, x.adminUpdateRow)
	adminSite.DELETE("/page/:"+paramPageID+"/row/:"+paramRowID, x.adminDeleteRow)
}

// An api struct to hold the repo for all the router functions
type a struct {
	repo *repo.Repo
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
	userID, err := x.getAdmin(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
	}

	x.repo.CreateSite(ctx.Request.Context(), reqBody.Name, userID)

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

	_, err = x.getAdmin(ctx)
	if err != nil {
		return ctx.NoContent(http.StatusUnauthorized)
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

// Helper functions
func (x a) getAdmin(ctx *router.Context) (id int, err error) {
	// TODO ctx.Request.Cookie("bearer-token")

	return 0, fmt.Errorf("not implemented")
}
