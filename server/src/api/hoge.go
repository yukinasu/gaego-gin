package api

import (
	"errors"
	"gaego-gin/server/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjibson/goon"
	"google.golang.org/appengine"
)

// HogeAPI manages api for Hoge
type HogeAPI struct{}

// SetupHoge is
func SetupHoge(rg *gin.RouterGroup) {
	api := &HogeAPI{}

	rg.GET("/hoge/:id", api.Get)
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
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, hoge)
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
	if err := c.Bind(hoge); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if hoge.ID == "" {
		c.JSON(http.StatusBadRequest, errors.New("id is required"))
		return
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	if err := hoge.Insert(g); err != nil {
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
// @Failure 500 {string} string
// @Router /hoge/{id} [put]
func (api *HogeAPI) Update(c *gin.Context) {
	hoge := &model.Hoge{}
	if err := c.Bind(hoge); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if hoge.ID == "" {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	ctx := appengine.NewContext(c.Request)
	g := goon.FromContext(ctx)

	if err := hoge.Update(g); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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
	if err := store.Delete(g, id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
