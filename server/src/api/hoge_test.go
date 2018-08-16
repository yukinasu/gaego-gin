package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gaego-gin/server/src/api"
	"gaego-gin/server/src/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestHogeAPI_Get(t *testing.T) {
	inst, err := aetest.NewInstance(&aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err.Error())
	}

	adminHelper := NewAdminTestHelper(inst)
	helper := newHogeTestHelper(inst, adminHelper.ctx)

	t.Run("Hogeが取得できること", func(t *testing.T) {
		defer adminHelper.ClearEntity(t, model.Hoge{})

		v := adminHelper.createHoge(t, &model.Hoge{
			ID:    "hoge",
			Value: "hogehoge",
		})

		code, resp, body := helper.requestGet(t, v.ID)

		AssertHTTPStatusCodeEquals(t, code, http.StatusOK, body)
		AssertEquals(t, "ID", resp.ID, "hoge")
		AssertEquals(t, "Value", resp.Value, "hogehoge")
	})

	t.Run("entityが存在しない場合、404エラーとなること", func(t *testing.T) {
		code, _, body := helper.requestGet(t, "hoge")

		AssertHTTPStatusCodeEquals(t, code, http.StatusNotFound, body)
	})
}

func TestHogeAPI_List(t *testing.T) {
	inst, err := aetest.NewInstance(&aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err.Error())
	}

	adminHelper := NewAdminTestHelper(inst)
	helper := newHogeTestHelper(inst, adminHelper.ctx)

	adminHelper.createHoge(t, &model.Hoge{ID: "hoge0", Value: "hogehoge0"})
	adminHelper.createHoge(t, &model.Hoge{ID: "hoge1", Value: "hogehoge1"})
	adminHelper.createHoge(t, &model.Hoge{ID: "hoge2", Value: "hogehoge2"})
	adminHelper.createHoge(t, &model.Hoge{ID: "hoge3", Value: "hogehoge3"})
	adminHelper.createHoge(t, &model.Hoge{ID: "hoge4", Value: "hogehoge4"})

	t.Run("全件取得できること", func(t *testing.T) {
		code, resp, body := helper.requestList(t, "", 0)

		AssertHTTPStatusCodeEquals(t, code, http.StatusOK, body)
		AssertEquals(t, "len(resp.List)", len(resp.List), 5)
		AssertEquals(t, "resp.Cursor", resp.Cursor, "")

		for idx, v := range resp.List {
			AssertEquals(t, fmt.Sprintf("resp.List[%d].ID", idx), v.ID, fmt.Sprintf("hoge%d", idx))
			AssertEquals(t, fmt.Sprintf("resp.List[%d].Value", idx), v.Value, fmt.Sprintf("hogehoge%d", idx))
		}
	})

	t.Run("limitよりも多く存在する場合、Cursorが返ること", func(t *testing.T) {
		code, resp, body := helper.requestList(t, "", 3)

		AssertHTTPStatusCodeEquals(t, code, http.StatusOK, body)
		AssertEquals(t, "len(resp.List)", len(resp.List), 3)
		AssertEquals(t, "resp.Cursor", resp.Cursor != "", true)

		code, resp, body = helper.requestList(t, resp.Cursor, 3)
		AssertEquals(t, "len(resp.List)", len(resp.List), 2)
		AssertEquals(t, "resp.Cursor", resp.Cursor, "")
	})
}

func TestHogeAPI_Insert(t *testing.T) {
	inst, err := aetest.NewInstance(&aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err.Error())
	}

	adminHelper := NewAdminTestHelper(inst)
	helper := newHogeTestHelper(inst, adminHelper.ctx)

	t.Run("Hogeが新規作成されること", func(t *testing.T) {
		defer adminHelper.ClearEntity(t, model.Hoge{})

		v := &model.Hoge{
			ID:    "hoge",
			Value: "hogehoge",
		}

		code, resp, body := helper.requestInsert(t, v)

		AssertHTTPStatusCodeEquals(t, code, http.StatusOK, body)
		AssertEquals(t, "resp.ID", resp.ID, "hoge")
		AssertEquals(t, "resp.Value", resp.Value, "hogehoge")

		g := goon.FromContext(adminHelper.ctx)

		store := &model.HogeStore{}
		hoge, err := store.Get(g, resp.ID)
		if err != nil {
			t.Fatal(err.Error())
		}

		AssertEquals(t, "hoge.ID", hoge.ID, "hoge")
		AssertEquals(t, "hoge.Value", hoge.Value, "hogehoge")
	})

	t.Run("Hoge.IDが未設定の場合、400エラーとなること", func(t *testing.T) {
		v := &model.Hoge{
			Value: "hogehoge",
		}

		code, _, body := helper.requestInsert(t, v)

		AssertHTTPStatusCodeEquals(t, code, http.StatusBadRequest, body)
	})
}

func TestHogeAPI_Update(t *testing.T) {
	inst, err := aetest.NewInstance(&aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err.Error())
	}

	adminHelper := NewAdminTestHelper(inst)
	helper := newHogeTestHelper(inst, adminHelper.ctx)

	t.Run("Hogeが更新されること", func(t *testing.T) {
		defer adminHelper.ClearEntity(t, model.Hoge{})

		v := adminHelper.createHoge(t, &model.Hoge{
			ID:    "hoge",
			Value: "hogehoge",
		})

		v.Value = "updated"

		code, resp, body := helper.requestUpdate(t, v)

		AssertHTTPStatusCodeEquals(t, code, http.StatusOK, body)
		AssertEquals(t, "resp.ID", resp.ID, "hoge")
		AssertEquals(t, "resp.Value", resp.Value, "updated")

		g := goon.FromContext(adminHelper.ctx)

		store := &model.HogeStore{}
		hoge, err := store.Get(g, resp.ID)
		if err != nil {
			t.Fatal(err.Error())
		}

		AssertEquals(t, "hoge.ID", hoge.ID, "hoge")
		AssertEquals(t, "hoge.Value", hoge.Value, "updated")
	})

	t.Run("対象IDのentityが存在しない場合、404エラーとなること", func(t *testing.T) {
		defer adminHelper.ClearEntity(t, model.Hoge{})

		v := &model.Hoge{
			ID:    "hoge",
			Value: "hogehoge",
		}

		code, _, body := helper.requestUpdate(t, v)

		AssertHTTPStatusCodeEquals(t, code, http.StatusNotFound, body)
	})
}

func TestHogeAPI_Delete(t *testing.T) {
	inst, err := aetest.NewInstance(&aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err.Error())
	}

	adminHelper := NewAdminTestHelper(inst)
	helper := newHogeTestHelper(inst, adminHelper.ctx)

	t.Run("Hogeが削除されること", func(t *testing.T) {
		defer adminHelper.ClearEntity(t, model.Hoge{})

		v := adminHelper.createHoge(t, &model.Hoge{
			ID:    "hoge",
			Value: "hogehoge",
		})

		code, body := helper.requestDelete(t, v.ID)

		AssertHTTPStatusCodeEquals(t, code, http.StatusOK, body)

		g := goon.FromContext(adminHelper.ctx)

		store := &model.HogeStore{}
		if _, err := store.Get(g, v.ID); err != nil {
			if err == datastore.ErrNoSuchEntity {
				// OK!
			} else {
				t.Fatal(err.Error())
			}
		} else {
			t.Fatal("unexpected")
		}
	})
}

/* Helper */

type hogeTestHelper struct {
	inst aetest.Instance
	ctx  context.Context
}

func newHogeTestHelper(inst aetest.Instance, ctx context.Context) *hogeTestHelper {
	return &hogeTestHelper{
		inst: inst,
		ctx:  ctx,
	}
}

func (h *hogeTestHelper) initializeHandler() http.Handler {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api.SetupHoge(r.Group("/api"))

	return r
}

func (h *hogeTestHelper) requestGet(t *testing.T, id string) (code int, v *model.Hoge, body []byte) {
	path := fmt.Sprintf("/api/hoge/%s", id)

	r, err := h.inst.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	handler := h.initializeHandler()

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if w.Code != http.StatusOK {
		return w.Code, nil, body
	}

	v = &model.Hoge{}
	err = json.Unmarshal(body, v)
	if err != nil {
		t.Fatal(err.Error())
	}

	return w.Code, v, body
}

func (h *hogeTestHelper) requestList(t *testing.T, cursor string, limit int) (code int, v *model.HogeListResp, body []byte) {
	vs := url.Values{}
	vs.Add("cursor", cursor)
	vs.Add("limit", fmt.Sprintf("%d", limit))
	path := fmt.Sprintf("/api/hoge?%s", vs.Encode())

	r, err := h.inst.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	handler := h.initializeHandler()

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if w.Code != http.StatusOK {
		return w.Code, nil, body
	}

	v = &model.HogeListResp{}
	err = json.Unmarshal(body, v)
	if err != nil {
		t.Fatal(err.Error())
	}

	return w.Code, v, body
}

func (h *hogeTestHelper) requestInsert(t *testing.T, dst *model.Hoge) (code int, v *model.Hoge, body []byte) {
	body, err := json.Marshal(dst)
	if err != nil {
		t.Fatal(err.Error())
	}

	r, err := h.inst.NewRequest("POST", "/api/hoge", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}

	handler := h.initializeHandler()

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if w.Code != http.StatusOK {
		return w.Code, nil, body
	}

	v = &model.Hoge{}
	err = json.Unmarshal(body, v)
	if err != nil {
		t.Fatal(err.Error())
	}

	return w.Code, v, body
}

func (h *hogeTestHelper) requestUpdate(t *testing.T, dst *model.Hoge) (code int, v *model.Hoge, body []byte) {
	path := fmt.Sprintf("/api/hoge/%s", dst.ID)

	body, err := json.Marshal(dst)
	if err != nil {
		t.Fatal(err.Error())
	}

	r, err := h.inst.NewRequest("PUT", path, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}

	handler := h.initializeHandler()

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if w.Code != http.StatusOK {
		return w.Code, nil, body
	}

	v = &model.Hoge{}
	err = json.Unmarshal(body, v)
	if err != nil {
		t.Fatal(err.Error())
	}

	return w.Code, v, body
}

func (h *hogeTestHelper) requestDelete(t *testing.T, id string) (code int, body []byte) {
	path := fmt.Sprintf("/api/hoge/%s", id)

	r, err := h.inst.NewRequest("DELETE", path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	handler := h.initializeHandler()

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	if w.Code != http.StatusOK {
		return w.Code, body
	}

	return w.Code, body
}
