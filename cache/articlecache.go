// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (c) 2017-present, b3log.org
//
// Pipe is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

// Package cache includes caches.
package cache

import (
	"os"
	"time"

	"github.com/88250/gulu"
	"github.com/leimeng-go/pipe/model"
	"github.com/bluele/gcache"
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// Article cache.
var Article = &articleCache{
	idHolder: gcache.New(1024 * 10).LRU().Expiration(30 * time.Minute).Build(),
}

type articleCache struct {
	idHolder gcache.Cache
}

func (cache *articleCache) Put(article *model.Article) {
	if err := cache.idHolder.Set(article.ID, article); nil != err {
		logger.Errorf("put article [id=%d] into cache failed: %s", article.ID, err)
	}
}

func (cache *articleCache) Get(id uint) *model.Article {
	ret, err := cache.idHolder.Get(id)
	if nil != err && gcache.KeyNotFoundError != err {
		logger.Errorf("get article [id=%d] from cache failed: %s", id, err)

		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.Article)
}
