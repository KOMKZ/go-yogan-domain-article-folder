package articlefolder

import "errors"

var (
	// ErrFolderNotFound 文件夹不存在
	ErrFolderNotFound = errors.New("文件夹不存在")

	// ErrFolderHasArticles 文件夹下有文章
	ErrFolderHasArticles = errors.New("文件夹下有文章，无法删除")
)
