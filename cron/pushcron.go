// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (c) 2017-present, b3log.org
//
// Pipe is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package cron

import (
	"net/url"
	"time"

	"github.com/88250/gulu"
	"github.com/leimeng-go/pipe/model"
	"github.com/leimeng-go/pipe/service"
	"github.com/leimeng-go/pipe/util"
)

func pushArticlesPeriodically() {
	go pushArticles()

	go func() {
		for range time.Tick(time.Second * 30) {
			pushArticles()
		}
	}()
}

func pushArticles() {
	defer gulu.Panic.Recover(nil)

	server, _ := url.Parse(model.Conf.Server)
	if !util.IsDomain(server.Hostname()) {
		return
	}

	articles := service.Article.GetUnpushedArticles()
	for _, article := range articles {
		service.Article.ConsolePushArticle(article)
	}
}
