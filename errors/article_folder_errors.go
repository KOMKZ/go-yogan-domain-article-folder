package errors

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

const ModuleArticleFolder = 34

var (
	ErrFolderNotFound = errcode.Register(errcode.New(
		ModuleArticleFolder, 1001, "article-folder",
		"error.article_folder.not_found", "文件夹不存在",
		http.StatusNotFound,
	))
	ErrFolderHasArticles = errcode.Register(errcode.New(
		ModuleArticleFolder, 1002, "article-folder",
		"error.article_folder.has_articles", "文件夹下有文章，无法删除",
		http.StatusConflict,
	))
)
