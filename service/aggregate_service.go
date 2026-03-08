package service

import (
	"context"

	articleModel "github.com/KOMKZ/go-yogan-domain-article/model"
	articleService "github.com/KOMKZ/go-yogan-domain-article/service"
	folderModel "github.com/KOMKZ/go-yogan-domain-folder/model"
	folderService "github.com/KOMKZ/go-yogan-domain-folder/service"
)

// ArticleFolderService 文章-文件夹聚合服务
// 负责跨领域的编排逻辑，解决 article 和 folder 之间的依赖关系
// 严格遵循 DIP：article 不依赖 folder，folder 不依赖 article
type ArticleFolderService struct {
	articleSvc *articleService.ArticleService
	folderSvc  *folderService.FolderService
}

// NewArticleFolderService 创建聚合服务
func NewArticleFolderService(articleSvc *articleService.ArticleService, folderSvc *folderService.FolderService) *ArticleFolderService {
	return &ArticleFolderService{
		articleSvc: articleSvc,
		folderSvc:  folderSvc,
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
func (s *ArticleFolderService) GetArticleWithFolder(ctx context.Context, articleID uint) (*ArticleWithFolder, error) {
	art, err := s.articleSvc.GetArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	result := &ArticleWithFolder{Article: art}

	if art.FolderID != nil {
		folderInfo, err := s.getFolderInfo(ctx, *art.FolderID)
		if err == nil {
			result.Folder = folderInfo
		}
	}

	return result, nil
}

// PageResultWithFolder 带文件夹信息的分页结果
type PageResultWithFolder struct {
	Records     []*ArticleWithFolder `json:"records"`
	Total       int64                `json:"total"`
	Size        int                  `json:"size"`
	Current     int                  `json:"current"`
	Pages       int                  `json:"pages"`
	HasPrevious bool                 `json:"has_previous"`
	HasNext     bool                 `json:"has_next"`
	IsFirst     bool                 `json:"is_first"`
	IsLast      bool                 `json:"is_last"`
}

// ListArticlesWithFolder 获取文章列表（带文件夹信息）
// 当指定 folderID 时，会自动查询该分类及其所有子分类下的文章
func (s *ArticleFolderService) ListArticlesWithFolder(ctx context.Context, page, size int, ownerId *uint, ownerType, articleType, title string, folderID *uint) (*PageResultWithFolder, error) {
	var result *articleService.PageResult
	var err error

	if folderID != nil {
		descendantIDs, err := s.folderSvc.GetDescendantIDs(ctx, *folderID)
		if err != nil {
			result, err = s.articleSvc.ListArticles(ctx, page, size, ownerId, ownerType, articleType, title, folderID)
		} else {
			result, err = s.articleSvc.ListArticlesByFolderIDs(ctx, page, size, ownerId, ownerType, articleType, title, descendantIDs)
		}
	} else {
		result, err = s.articleSvc.ListArticles(ctx, page, size, ownerId, ownerType, articleType, title, folderID)
	}

	if err != nil {
		return nil, err
	}

	folderIDs := make(map[uint]struct{})
	for _, art := range result.Records {
		if art.FolderID != nil {
			folderIDs[*art.FolderID] = struct{}{}
		}
	}

	folderMap := s.batchGetFolderInfo(ctx, folderIDs)

	records := make([]*ArticleWithFolder, len(result.Records))
	for i, art := range result.Records {
		artCopy := art
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

// MoveArticleToFolder 移动文章到指定文件夹（带验证）
func (s *ArticleFolderService) MoveArticleToFolder(ctx context.Context, articleID uint, folderID *uint) error {
	if folderID != nil {
		_, err := s.folderSvc.GetFolder(ctx, *folderID)
		if err != nil {
			return err
		}
	}

	return s.articleSvc.MoveToFolder(ctx, articleID, folderID)
}

// ValidateFolderExists 验证文件夹是否存在
func (s *ArticleFolderService) ValidateFolderExists(ctx context.Context, folderID uint) (*folderModel.Folder, error) {
	return s.folderSvc.GetFolder(ctx, folderID)
}

// CanDeleteFolder 检查文件夹是否可以删除（是否有文章）
func (s *ArticleFolderService) CanDeleteFolder(ctx context.Context, folderID uint) (bool, int64, error) {
	count, err := s.articleSvc.CountByFolder(ctx, folderID)
	if err != nil {
		return false, 0, err
	}
	return count == 0, count, nil
}

// ==================== 计数更新（事件驱动调用） ====================

// OnArticleCreated 文章创建时调用，增加文件夹计数
func (s *ArticleFolderService) OnArticleCreated(ctx context.Context, folderID *uint) error {
	if folderID == nil {
		return nil
	}
	return s.folderSvc.IncrementItemCount(ctx, *folderID, 1)
}

// OnArticleDeleted 文章删除时调用，减少文件夹计数
func (s *ArticleFolderService) OnArticleDeleted(ctx context.Context, folderID *uint) error {
	if folderID == nil {
		return nil
	}
	return s.folderSvc.IncrementItemCount(ctx, *folderID, -1)
}

// OnArticleMoved 文章移动时调用，更新新旧文件夹计数
func (s *ArticleFolderService) OnArticleMoved(ctx context.Context, oldFolderID, newFolderID *uint) error {
	if oldFolderID != nil {
		if err := s.folderSvc.IncrementItemCount(ctx, *oldFolderID, -1); err != nil {
			return err
		}
	}
	if newFolderID != nil {
		if err := s.folderSvc.IncrementItemCount(ctx, *newFolderID, 1); err != nil {
			return err
		}
	}
	return nil
}

// getFolderInfo 获取单个文件夹信息
func (s *ArticleFolderService) getFolderInfo(ctx context.Context, folderID uint) (*FolderInfo, error) {
	f, err := s.folderSvc.GetFolder(ctx, folderID)
	if err != nil {
		return nil, err
	}

	ancestors, _ := s.folderSvc.GetAncestors(ctx, folderID)
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
func (s *ArticleFolderService) batchGetFolderInfo(ctx context.Context, folderIDs map[uint]struct{}) map[uint]*FolderInfo {
	result := make(map[uint]*FolderInfo)

	for id := range folderIDs {
		info, err := s.getFolderInfo(ctx, id)
		if err == nil {
			result[id] = info
		}
	}

	return result
}
