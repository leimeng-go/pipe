// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (c) 2017-present, b3log.org
//
// Pipe is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package console

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/88250/gulu"
	"github.com/leimeng-go/pipe/model"
	"github.com/leimeng-go/pipe/service"
	"github.com/leimeng-go/pipe/util"
	"github.com/gin-gonic/gin"
)

// GetCommentsAction gets comments
func GetCommentsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	commentModels, pagination := service.Comment.ConsoleGetComments(c.Query("key"), util.GetPage(c), session.BID)
	blogURLSetting := service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, session.BID)

	var comments []*ConsoleComment
	for _, commentModel := range commentModels {
		article := service.Article.ConsoleGetArticle(commentModel.ArticleID)
		articleAuthor := service.User.GetUser(article.AuthorID)
		consoleArticleAuthor := &ConsoleAuthor{
			URL:       blogURLSetting.Value + util.PathAuthors + "/" + articleAuthor.Name,
			Name:      articleAuthor.Name,
			AvatarURL: articleAuthor.AvatarURL,
		}

		author := &ConsoleAuthor{}
		if model.SyncCommentAuthorID == commentModel.AuthorID {
			author.URL = commentModel.AuthorURL
			author.Name = commentModel.AuthorName
			author.AvatarURL = commentModel.AuthorAvatarURL
		} else {
			commentAuthor := service.User.GetUser(commentModel.AuthorID)
			commentAuthorBlog := service.User.GetOwnBlog(commentModel.AuthorID)
			author.URL = service.Setting.GetSetting(model.SettingCategoryBasic, model.SettingNameBasicBlogURL, commentAuthorBlog.ID).Value + util.PathAuthors + "/" + commentAuthor.Name
			author.Name = commentAuthor.Name
			author.AvatarURL = commentAuthor.AvatarURL
		}

		page := service.Comment.GetCommentPage(commentModel.ArticleID, commentModel.ID, commentModel.BlogID)
		mdResult := util.Markdown(commentModel.Content)
		comment := &ConsoleComment{
			ID:            commentModel.ID,
			Author:        author,
			ArticleAuthor: consoleArticleAuthor,
			CreatedAt:     commentModel.CreatedAt.Format("2006-01-02"),
			Title:         article.Title,
			Content:       template.HTML(mdResult.ContentHTML),
			URL:           blogURLSetting.Value + article.Path + "?p=" + strconv.Itoa(page) + "#pipeComment" + strconv.Itoa(int(commentModel.ID)),
		}

		comments = append(comments, comment)
	}

	data := map[string]interface{}{}
	data["comments"] = comments
	data["pagination"] = pagination
	result.Data = data
}

// RemoveCommentAction removes a comment.
func RemoveCommentAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.ParseUint(idArg, 10, 64)
	if nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()

		return
	}

	session := util.GetSession(c)
	blogID := session.BID

	if err := service.Comment.RemoveComment(id, blogID); nil != err {
		result.Code = util.CodeErr
		result.Msg = err.Error()
	}
}

// RemoveCommentsAction removes comments.
func RemoveCommentsAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = util.CodeErr
		result.Msg = "parses batch remove comments request failed"

		return
	}

	session := util.GetSession(c)
	blogID := session.BID
	ids := arg["ids"].([]interface{})
	for _, id := range ids {
		if err := service.Comment.RemoveComment(uint64(id.(float64)), blogID); nil != err {
			logger.Errorf("remove comment failed: " + err.Error())
		}
	}
}
