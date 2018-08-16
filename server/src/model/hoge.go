package model

import (
	"errors"
	"time"

	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

// HogeStore はHogeを操作するメソッドをまとめる
type HogeStore struct{}

// Hoge はサンプル用の構造体
type Hoge struct {
	ID        string    `json:"id" datastore:"-" goon:"id"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Load はPropertyLoadSaverのインターフェースを実装する
func (src *Hoge) Load(p []datastore.Property) error {
	if err := datastore.LoadStruct(src, p); err != nil {
		return err
	}

	return nil
}

// Save はPropertyLoadSaverのインターフェースを実装する
func (src *Hoge) Save() ([]datastore.Property, error) {
	src.UpdatedAt = time.Now()

	p, err := datastore.SaveStruct(src)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Insert はHogeを新規登録する
func (src *Hoge) Insert(g *goon.Goon) error {
	if src.ID == "" {
		return errors.New("id is required")
	}

	old := &Hoge{
		ID: src.ID,
	}
	if err := g.Get(old); err != nil {
		if err == datastore.ErrNoSuchEntity {
			// OK!
		} else {
			return err
		}
	} else {
		return errors.New("already exist")
	}

	return src.put(g, nil)
}

// Update はHogeを更新する
func (src *Hoge) Update(g *goon.Goon) error {
	if src.ID == "" {
		return errors.New("id is required")
	}

	old := &Hoge{
		ID: src.ID,
	}
	if err := g.Get(old); err != nil {
		return err
	}

	return src.put(g, old)
}

func (src *Hoge) put(g *goon.Goon, old *Hoge) error {
	src.CreatedAt = time.Now()

	if old != nil {
		src.CreatedAt = old.CreatedAt
	}

	if _, err := g.Put(src); err != nil {
		return err
	}

	return nil
}

// Get はHogeを1件取得する
func (store *HogeStore) Get(g *goon.Goon, id string) (*Hoge, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	hoge := &Hoge{
		ID: id,
	}
	if err := g.Get(hoge); err != nil {
		return nil, err
	}

	return hoge, nil
}

// HogeListResp はHoge一覧取得のレスポンス
type HogeListResp struct {
	List   []*Hoge `json:"list"`
	Cursor string  `json:"cursor"`
}

// List はHogeの一覧を取得する
func (store *HogeStore) List(g *goon.Goon, cursor string, limit int) (*HogeListResp, error) {
	q := datastore.NewQuery(g.Kind(Hoge{})).KeysOnly()

	if limit == 0 {
		limit = 10
	}
	if limit != -1 {
		// 次の1件が存在するかを確認するため、1件多く取得する
		q = q.Limit(limit + 1)
	}

	if cursor != "" {
		start, err := datastore.DecodeCursor(cursor)
		if err != nil {
			return nil, err
		}

		q = q.Start(start)
	}

	it := g.Run(q)

	count := 0
	hasNext := false
	var cur datastore.Cursor
	list := make([]*Hoge, 0, limit)

	for {
		key, err := it.Next(nil)
		if err != nil {
			if err == datastore.Done {
				break
			}

			return nil, err
		}

		count++
		if limit != -1 && limit < count {
			hasNext = true
			break
		}

		list = append(list, &Hoge{ID: key.StringID()})

		// limitで指定した件数に到達したところでCursorを保存
		if limit == count {
			cur, err = it.Cursor()
			if err != nil {
				return nil, err
			}
		}
	}

	if err := g.GetMulti(list); err != nil {
		return nil, err
	}

	resp := &HogeListResp{
		List: list,
	}

	if hasNext {
		resp.Cursor = cur.String()
	}

	return resp, nil
}

// Delete はHogeを削除する
func (store *HogeStore) Delete(g *goon.Goon, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	hoge := &Hoge{
		ID: id,
	}
	if err := g.Delete(g.Key(hoge)); err != nil {
		return err
	}

	return nil
}
