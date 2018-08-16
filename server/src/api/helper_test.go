package api_test

import (
	"context"
	"gaego-gin/server/src/model"
	"testing"

	"github.com/mjibson/goon"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

type AdminTestHelper struct {
	inst aetest.Instance
	ctx  context.Context
}

func NewAdminTestHelper(inst aetest.Instance) *AdminTestHelper {
	r, err := inst.NewRequest("GET", "", nil)
	if err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)

	return &AdminTestHelper{
		inst: inst,
		ctx:  ctx,
	}
}

func (h *AdminTestHelper) ClearEntity(t *testing.T, src interface{}) {
	g := goon.FromContext(h.ctx)

	q := datastore.NewQuery(g.Kind(src)).KeysOnly()

	keys, err := g.GetAll(q, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err := g.DeleteMulti(keys); err != nil {
		t.Fatal(err.Error())
	}
}

/* Kind別のentity作成 */

func (h *AdminTestHelper) createHoge(t *testing.T, v *model.Hoge) *model.Hoge {
	g := goon.FromContext(h.ctx)

	if err := v.Insert(g); err != nil {
		t.Fatal(err.Error())
	}

	return v
}

/* assert */

// AssertEquals は実値と期待値が同値か判定する
func AssertEquals(t *testing.T, title string, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%s: unexpected, actual: `%v`, expected: `%v`", title, actual, expected)
	}
}

// AssertHTTPStatusCodeEquals はHTTPステータスコードの実値と期待値が同値か判定する
func AssertHTTPStatusCodeEquals(t *testing.T, actual, expected int, body []byte) {
	if actual != expected {
		// ステータスコードが異なる場合は、以降のassertは必ずFAILEDとなり、場合によってはpanicとなる
		// よって`t.Errorf()`ではなく、`t.Fatalf()`を利用し、テストを実行を終了する
		t.Fatalf("unexpected status code: actual: `%d`, expected: `%d`, body: `%s`", actual, expected, string(body))
	}
}
