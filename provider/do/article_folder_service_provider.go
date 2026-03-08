package do

import (
	articleService "github.com/KOMKZ/go-yogan-domain-article/service"
	"github.com/KOMKZ/go-yogan-domain-article-folder/service"
	folderService "github.com/KOMKZ/go-yogan-domain-folder/service"
	"github.com/samber/do/v2"
)

// ProvideArticleFolderService 提供 ArticleFolderService
func ProvideArticleFolderService(i do.Injector) (*service.ArticleFolderService, error) {
	articleSvc, err := do.Invoke[*articleService.ArticleService](i)
	if err != nil {
		return nil, err
	}
	folderSvc, err := do.Invoke[*folderService.FolderService](i)
	if err != nil {
		return nil, err
	}
	return service.NewArticleFolderService(articleSvc, folderSvc), nil
}
