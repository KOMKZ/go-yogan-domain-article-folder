package articlefolder

import (
	"context"

	article "github.com/KOMKZ/go-yogan-domain-article"
	articleModel "github.com/KOMKZ/go-yogan-domain-article/model"
	folder "github.com/KOMKZ/go-yogan-domain-folder"
	folderModel "github.com/KOMKZ/go-yogan-domain-folder/model"
)

// Service 文章-文件夹聚合服务
// 负责跨领域的编排逻辑，解决 article 和 folder 之间的依赖关系
type Service struct {
	articleService *article.Service
	folderService  *folder.Service
}

// NewService 创建聚合服务
func NewService(articleService *article.Service, folderService *folder.Service) *Service {
	return &Service{
		articleService: articleService,
		folderService:  folderService,
	}
}

// BreadcrumbItem 面包屑项（包含 ID 和名称）
type BreadcrumbItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// FolderInfo 文件夹信息（用于组装到文章响应中）
type FolderInfo struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
	Breadcrumb []BreadcrumbItem `json:"breadcrumb"`
}

// ArticleWithFolder 带文件夹信息的文章
type ArticleWithFolder struct {
	*articleModel.Article
	Folder *FolderInfo `json:"folder,omitempty"`
}

// GetArticleWithFolder 获取文章详情（带文件夹信息）
func (s *Service) GetArticleWithFolder(ctx context.Context, articleID uint) (*ArticleWithFolder, error) {
	art, err := s.articleService.GetArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	result := &ArticleWithFolder{Article: art}

	// 如果有 folder_id，获取文件夹信息
	if art.FolderID != nil {
		folderInfo, err := s.getFolderInfo(ctx, *art.FolderID)
		if err == nil {
			result.Folder = folderInfo
		}
		// 忽略文件夹查询错误，文章依然返回
	}

	return result, nil
}

// ListArticlesWithFolder 获取文章列表（带文件夹信息）
// 当指定 folderID 时，会自动查询该分类及其所有子分类下的文章
func (s *Service) ListArticlesWithFolder(ctx context.Context, page, size int, ownerId *uint, ownerType, articleType, title string, folderID *uint) (*PageResultWithFolder, error) {
	var result *article.PageResult
	var err error

	// 1. 如果指定了 folderID，获取该分类及所有子分类的 ID
	if folderID != nil {
		descendantIDs, err := s.folderService.GetDescendantIDs(ctx, *folderID)
		if err != nil {
			// 如果获取子孙失败，回退到单个 folderID 查询
			result, err = s.articleService.ListArticles(ctx, page, size, ownerId, ownerType, articleType, title, folderID)
		} else {
			// 使用多个 folderID 查询
			result, err = s.articleService.ListArticlesByFolderIDs(ctx, page, size, ownerId, ownerType, articleType, title, descendantIDs)
		}
	} else {
		result, err = s.articleService.ListArticles(ctx, page, size, ownerId, ownerType, articleType, title, folderID)
	}

	if err != nil {
		return nil, err
	}

	// 2. 收集所有 folderID（去重）
	folderIDs := make(map[uint]struct{})
	for _, art := range result.Records {
		if art.FolderID != nil {
			folderIDs[*art.FolderID] = struct{}{}
		}
	}

	// 3. 批量获取 folder 信息
	folderMap := s.batchGetFolderInfo(ctx, folderIDs)

	// 4. 组装返回
	records := make([]*ArticleWithFolder, len(result.Records))
	for i, art := range result.Records {
		artCopy := art // 避免循环引用问题
		records[i] = &ArticleWithFolder{Article: &artCopy}
		if art.FolderID != nil {
			if info, ok := folderMap[*art.FolderID]; ok {
				records[i].Folder = info
			}
		}
	}

	return &PageResultWithFolder{
		Records:     records,
		Total:       result.Total,
		Size:        result.Size,
		Current:     result.Current,
		Pages:       result.Pages,
		HasPrevious: result.HasPrevious,
		HasNext:     result.HasNext,
		IsFirst:     result.IsFirst,
		IsLast:      result.IsLast,
	}, nil
}

// PageResultWithFolder 带文件夹信息的分页结果
type PageResultWithFolder struct {
	Records     []*ArticleWithFolder `json:"records"`
	Total       int64                `json:"total"`
	Size        int                  `json:"size"`
	Current     int                  `json:"current"`
	Pages       int                  `json:"pages"`
	HasPrevious bool                 `json:"hasPrevious"`
	HasNext     bool                 `json:"hasNext"`
	IsFirst     bool                 `json:"isFirst"`
	IsLast      bool                 `json:"isLast"`
}

// MoveArticleToFolder 移动文章到指定文件夹（带验证）
func (s *Service) MoveArticleToFolder(ctx context.Context, articleID uint, folderID *uint) error {
	// 1. 验证文件夹存在（如果 folderID 不为空）
	if folderID != nil {
		_, err := s.folderService.GetFolder(ctx, *folderID)
		if err != nil {
			return err
		}
	}

	// 2. 移动文章
	return s.articleService.MoveToFolder(ctx, articleID, folderID)
}

// ValidateFolderExists 验证文件夹是否存在
func (s *Service) ValidateFolderExists(ctx context.Context, folderID uint) (*folderModel.Folder, error) {
	return s.folderService.GetFolder(ctx, folderID)
}

// CanDeleteFolder 检查文件夹是否可以删除（是否有文章）
func (s *Service) CanDeleteFolder(ctx context.Context, folderID uint) (bool, int64, error) {
	count, err := s.articleService.CountByFolder(ctx, folderID)
	if err != nil {
		return false, 0, err
	}
	return count == 0, count, nil
}

// getFolderInfo 获取单个文件夹信息
func (s *Service) getFolderInfo(ctx context.Context, folderID uint) (*FolderInfo, error) {
	f, err := s.folderService.GetFolder(ctx, folderID)
	if err != nil {
		return nil, err
	}

	// GetAncestors 返回的结果已包含当前节点，直接使用
	ancestors, _ := s.folderService.GetAncestors(ctx, folderID)
	breadcrumb := make([]BreadcrumbItem, 0, len(ancestors))
	for _, a := range ancestors {
		breadcrumb = append(breadcrumb, BreadcrumbItem{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return &FolderInfo{
		ID:         f.ID,
		Name:       f.Name,
		Path:       f.Path,
		Breadcrumb: breadcrumb,
	}, nil
}

// batchGetFolderInfo 批量获取文件夹信息
func (s *Service) batchGetFolderInfo(ctx context.Context, folderIDs map[uint]struct{}) map[uint]*FolderInfo {
	result := make(map[uint]*FolderInfo)

	for id := range folderIDs {
		info, err := s.getFolderInfo(ctx, id)
		if err == nil {
			result[id] = info
		}
	}

	return result
}
