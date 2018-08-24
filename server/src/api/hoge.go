package api

import (
	"gaego-gin/server/src/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// HogeAPI はHogeのAPIを管理する
type HogeAPI struct{}

// SetupHoge はHogeのAPIのハンドリングを行う
func SetupHoge(rg *gin.RouterGroup) {
	api := &HogeAPI{}

	rg.GET("/hoge/:id", api.Get)
	rg.GET("/hoge", api.List)
	rg.POST("/hoge", api.Insert)
	rg.PUT("/hoge/:id", api.Update)
	rg.DELETE("/hoge/:id", api.Delete)
}

// Get はHogeを1件取得する
// @Description Hogeを1件取得する
// @Tags Hoge
// @Summary Hoge 1件取得
// @Accept  json
// @Produce  json
// @Param  id path string true "Hoge.ID"
// @Success 200 {object} model.Hoge
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /hoge/{id} [get]
func (api *HogeAPI) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	store := &model.HogeStore{}
	hoge, err := store.Get(g, id)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, hoge)
}

// List はHogeの一覧を取得する
// @Description Hogeの一覧を取得する
// @Tags Hoge
// @Summary Hoge 一覧取得
// @Accept  json
// @Produce  json
// @Param  cursor query string false "start cursor"
// @Param  limit query string false "query limit"
// @Success 200 {object} model.HogeListResp
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /hoge [get]
func (api *HogeAPI) List(c *gin.Context) {
	cursor := c.Query("cursor")

	limit := 0
	if c.Query("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	store := &model.HogeStore{}
	resp, err := store.List(g, cursor, limit)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Insert はHogeを新規作成する
// @Description Hogeを新規作成する
// @Tags Hoge
// @Summary Hoge 新規作成
// @Accept  json
// @Produce  json
// @Param  hoge body model.Hoge true "新規作成するHoge"
// @Success 200 {object} model.Hoge
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /hoge [post]
func (api *HogeAPI) Insert(c *gin.Context) {
	hoge := &model.Hoge{}
	if err := c.BindJSON(hoge); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if hoge.ID == "" {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	if err := g.RunInTransaction(func(tg *goon.Goon) error {
		hoge := hoge

		return hoge.Insert(tg)

	}, &datastore.TransactionOptions{XG: true}); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, hoge)
}

// Update はHogeを更新する
// @Description Hogeを更新する
// @Tags Hoge
// @Summary Hoge 更新
// @Accept  json
// @Produce  json
// @Param  hoge body model.Hoge true "更新するHoge"
// @Success 200 {object} model.Hoge
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /hoge/{id} [put]
func (api *HogeAPI) Update(c *gin.Context) {
	hoge := &model.Hoge{}
	if err := c.BindJSON(hoge); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if hoge.ID == "" {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	statusCode := 0
	if err := g.RunInTransaction(func(tg *goon.Goon) error {
		hoge := hoge

		if err := hoge.Update(tg); err != nil {
			if err == datastore.ErrNoSuchEntity {
				statusCode = http.StatusNotFound
				return err
			}

			statusCode = http.StatusInternalServerError
			return err
		}

		return nil

	}, &datastore.TransactionOptions{XG: true}); err != nil {
		c.String(statusCode, err.Error())
		return
	}

	c.JSON(http.StatusOK, hoge)
}

// Delete はHogeを削除する
// @Description Hogeを削除する
// @Tags Hoge
// @Summary Hoge 削除
// @Accept  json
// @Produce  json
// @Param  id path string true "Hoge.ID"
// @Success 200 {null} null
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /hoge/{id} [delete]
func (api *HogeAPI) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	store := &model.HogeStore{}

	if err := g.RunInTransaction(func(tg *goon.Goon) error {
		return store.Delete(tg, id)

	}, &datastore.TransactionOptions{XG: true}); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
