package service

import (
	"context"
	"testing"

	articleModel "github.com/KOMKZ/go-yogan-domain-article/model"
	articleService "github.com/KOMKZ/go-yogan-domain-article/service"
	folderErrors "github.com/KOMKZ/go-yogan-domain-folder/errors"
	folderModel "github.com/KOMKZ/go-yogan-domain-folder/model"
	folderService "github.com/KOMKZ/go-yogan-domain-folder/service"
	"github.com/KOMKZ/go-yogan-framework/logger"
)

// ==================== Mock Repositories ====================

// MockArticleRepository mock 文章仓储
type MockArticleRepository struct {
	articles map[uint]*articleModel.Article
	nextID   uint
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		articles: make(map[uint]*articleModel.Article),
		nextID:   1,
	}
}

func (m *MockArticleRepository) Create(ctx context.Context, article *articleModel.Article) error {
	article.ID = m.nextID
	m.nextID++
	m.articles[article.ID] = article
	return nil
}

func (m *MockArticleRepository) Update(ctx context.Context, article *articleModel.Article) error {
	m.articles[article.ID] = article
	return nil
}

func (m *MockArticleRepository) FindByID(ctx context.Context, id uint) (*articleModel.Article, error) {
	if a, ok := m.articles[id]; ok {
		return a, nil
	}
	return nil, nil
}

func (m *MockArticleRepository) Delete(ctx context.Context, id uint) error {
	if a, ok := m.articles[id]; ok {
		a.Status = articleModel.StatusDeleted
	}
	return nil
}

func (m *MockArticleRepository) Paginate(ctx context.Context, page, pageSize int, ownerId *uint, ownerType, articleType, title string, folderID *uint) ([]articleModel.Article, int64, error) {
	var result []articleModel.Article
	for _, a := range m.articles {
		if a.Status != articleModel.StatusDeleted {
			if folderID == nil || (a.FolderID != nil && *a.FolderID == *folderID) {
				result = append(result, *a)
			}
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockArticleRepository) PaginateByFolderIDs(ctx context.Context, page, pageSize int, ownerId *uint, ownerType, articleType, title string, folderIDs []uint) ([]articleModel.Article, int64, error) {
	var result []articleModel.Article
	folderSet := make(map[uint]bool)
	for _, id := range folderIDs {
		folderSet[id] = true
	}
	for _, a := range m.articles {
		if a.Status != articleModel.StatusDeleted {
			if a.FolderID != nil && folderSet[*a.FolderID] {
				result = append(result, *a)
			}
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockArticleRepository) CountByFolderID(ctx context.Context, folderID uint) (int64, error) {
	var count int64
	for _, a := range m.articles {
		if a.FolderID != nil && *a.FolderID == folderID && a.Status != articleModel.StatusDeleted {
			count++
		}
	}
	return count, nil
}

func (m *MockArticleRepository) FindByFolderID(ctx context.Context, folderID uint) ([]articleModel.Article, error) {
	var result []articleModel.Article
	for _, a := range m.articles {
		if a.FolderID != nil && *a.FolderID == folderID && a.Status != articleModel.StatusDeleted {
			result = append(result, *a)
		}
	}
	return result, nil
}

// MockMarkdownRepository mock Markdown仓储
type MockMarkdownRepository struct{}

func NewMockMarkdownRepository() *MockMarkdownRepository { return &MockMarkdownRepository{} }
func (m *MockMarkdownRepository) Create(ctx context.Context, article *articleModel.MarkdownArticle) error {
	return nil
}
func (m *MockMarkdownRepository) Update(ctx context.Context, article *articleModel.MarkdownArticle) error {
	return nil
}
func (m *MockMarkdownRepository) FindByArticleID(ctx context.Context, articleID uint) (*articleModel.MarkdownArticle, error) {
	return nil, nil
}
func (m *MockMarkdownRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return nil
}

// MockRichTextRepository mock 富文本仓储
type MockRichTextRepository struct{}

func NewMockRichTextRepository() *MockRichTextRepository { return &MockRichTextRepository{} }
func (m *MockRichTextRepository) Create(ctx context.Context, article *articleModel.RichTextArticle) error {
	return nil
}
func (m *MockRichTextRepository) Update(ctx context.Context, article *articleModel.RichTextArticle) error {
	return nil
}
func (m *MockRichTextRepository) FindByArticleID(ctx context.Context, articleID uint) (*articleModel.RichTextArticle, error) {
	return nil, nil
}
func (m *MockRichTextRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return nil
}

// MockTableRepository mock 表格仓储
type MockTableRepository struct{}

func NewMockTableRepository() *MockTableRepository { return &MockTableRepository{} }
func (m *MockTableRepository) Create(ctx context.Context, article *articleModel.TableArticle) error {
	return nil
}
func (m *MockTableRepository) Update(ctx context.Context, article *articleModel.TableArticle) error {
	return nil
}
func (m *MockTableRepository) FindByArticleID(ctx context.Context, articleID uint) (*articleModel.TableArticle, error) {
	return nil, nil
}
func (m *MockTableRepository) FindByTableID(ctx context.Context, tableID string) (*articleModel.TableArticle, error) {
	return nil, nil
}
func (m *MockTableRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return nil
}

// MockTableRowRepository mock 表格行仓储
type MockTableRowRepository struct{}

func NewMockTableRowRepository() *MockTableRowRepository { return &MockTableRowRepository{} }
func (m *MockTableRowRepository) Create(ctx context.Context, row *articleModel.TableArticleRow) error {
	return nil
}
func (m *MockTableRowRepository) BatchCreate(ctx context.Context, rows []articleModel.TableArticleRow) error {
	return nil
}
func (m *MockTableRowRepository) FindByArticleID(ctx context.Context, articleID uint) ([]articleModel.TableArticleRow, error) {
	return nil, nil
}
func (m *MockTableRowRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return nil
}
func (m *MockTableRowRepository) ReplaceAll(ctx context.Context, articleID uint, rows []articleModel.TableArticleRow) error {
	return nil
}

// MockFolderRepository mock 文件夹仓储
type MockFolderRepository struct {
	folders      map[uint]*folderModel.Folder
	nextID       uint
	maxSortOrder int
}

func NewMockFolderRepository() *MockFolderRepository {
	return &MockFolderRepository{
		folders: make(map[uint]*folderModel.Folder),
		nextID:  1,
	}
}

func (m *MockFolderRepository) Create(ctx context.Context, folder *folderModel.Folder) error {
	folder.ID = m.nextID
	m.nextID++
	m.folders[folder.ID] = folder
	return nil
}

func (m *MockFolderRepository) Update(ctx context.Context, folder *folderModel.Folder) error {
	m.folders[folder.ID] = folder
	return nil
}

func (m *MockFolderRepository) Delete(ctx context.Context, id uint) error {
	delete(m.folders, id)
	return nil
}

func (m *MockFolderRepository) FindByID(ctx context.Context, id uint) (*folderModel.Folder, error) {
	if f, ok := m.folders[id]; ok {
		return f, nil
	}
	return nil, nil
}

func (m *MockFolderRepository) FindByParentID(ctx context.Context, parentID *uint) ([]*folderModel.Folder, error) {
	var result []*folderModel.Folder
	for _, f := range m.folders {
		if (parentID == nil && f.ParentID == nil) || (parentID != nil && f.ParentID != nil && *f.ParentID == *parentID) {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *MockFolderRepository) FindChildren(ctx context.Context, parentID uint) ([]*folderModel.Folder, error) {
	return m.FindByParentID(ctx, &parentID)
}

func (m *MockFolderRepository) FindRoots(ctx context.Context) ([]*folderModel.Folder, error) {
	return m.FindByParentID(ctx, nil)
}

func (m *MockFolderRepository) FindByPath(ctx context.Context, pathPrefix string) ([]*folderModel.Folder, error) {
	var result []*folderModel.Folder
	for _, f := range m.folders {
		if len(f.Path) >= len(pathPrefix) && f.Path[:len(pathPrefix)] == pathPrefix {
			result = append(result, f)
		}
	}
	return result, nil
}

func (m *MockFolderRepository) FindAncestors(ctx context.Context, path string) ([]*folderModel.Folder, error) {
	return []*folderModel.Folder{}, nil
}

func (m *MockFolderRepository) FindAll(ctx context.Context) ([]*folderModel.Folder, error) {
	var result []*folderModel.Folder
	for _, f := range m.folders {
		result = append(result, f)
	}
	return result, nil
}

func (m *MockFolderRepository) UpdateSortOrder(ctx context.Context, id uint, sortOrder int) error {
	return nil
}

func (m *MockFolderRepository) FindMaxSortOrder(ctx context.Context, parentID *uint) (int, error) {
	m.maxSortOrder++
	return m.maxSortOrder, nil
}

func (m *MockFolderRepository) UpdatePathAndDepth(ctx context.Context, id uint, path string, depth int) error {
	return nil
}

func (m *MockFolderRepository) UpdateChildrenPathAndDepth(ctx context.Context, oldPathPrefix, newPathPrefix string, depthDiff int) error {
	return nil
}

func (m *MockFolderRepository) ExistsByNameAndParent(ctx context.Context, name string, parentID *uint, excludeID *uint) (bool, error) {
	for _, f := range m.folders {
		if f.Name == name {
			if (parentID == nil && f.ParentID == nil) || (parentID != nil && f.ParentID != nil && *f.ParentID == *parentID) {
				if excludeID == nil || f.ID != *excludeID {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (m *MockFolderRepository) HasChildren(ctx context.Context, id uint) (bool, error) {
	for _, f := range m.folders {
		if f.ParentID != nil && *f.ParentID == id {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockFolderRepository) IncrementItemCount(ctx context.Context, id uint, delta int) error {
	if f, ok := m.folders[id]; ok {
		f.ItemCount += delta
	}
	return nil
}

func (m *MockFolderRepository) IncrementTotalItemCount(ctx context.Context, path string, delta int) error {
	return nil
}

// ==================== Test Helper ====================

func createTestServices() (*ArticleFolderService, *articleService.ArticleService, *folderService.FolderService, *MockArticleRepository, *MockFolderRepository) {
	articleRepo := NewMockArticleRepository()
	folderRepo := NewMockFolderRepository()

	artSvc := articleService.NewArticleService(
		articleRepo,
		NewMockMarkdownRepository(),
		NewMockRichTextRepository(),
		NewMockTableRepository(),
		NewMockTableRowRepository(),
		logger.GetLogger("article_test"),
	)
	fldSvc := folderService.NewFolderService(folderRepo, logger.GetLogger("folder_test"))
	aggSvc := NewArticleFolderService(artSvc, fldSvc)

	return aggSvc, artSvc, fldSvc, articleRepo, folderRepo
}

// ==================== Tests ====================

func TestArticleFolderService_GetArticleWithFolder(t *testing.T) {
	aggSvc, artSvc, fldSvc, _, _ := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "测试文件夹",
		ParentID: nil,
	})

	article, _ := artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "测试文章",
		ArticleType: "markdown",
		FolderID:    &folder.ID,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	result, err := aggSvc.GetArticleWithFolder(ctx, article.ID)
	if err != nil {
		t.Fatalf("GetArticleWithFolder failed: %v", err)
	}

	if result.Article.ID != article.ID {
		t.Errorf("Expected article ID %d, got %d", article.ID, result.Article.ID)
	}

	if result.Folder == nil {
		t.Error("Expected folder info, got nil")
	} else if result.Folder.ID != folder.ID {
		t.Errorf("Expected folder ID %d, got %d", folder.ID, result.Folder.ID)
	}
}

func TestArticleFolderService_GetArticleWithFolder_NoFolder(t *testing.T) {
	aggSvc, artSvc, _, _, _ := createTestServices()
	ctx := context.Background()

	article, _ := artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "无文件夹文章",
		ArticleType: "markdown",
		FolderID:    nil,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	result, err := aggSvc.GetArticleWithFolder(ctx, article.ID)
	if err != nil {
		t.Fatalf("GetArticleWithFolder failed: %v", err)
	}

	if result.Folder != nil {
		t.Error("Expected nil folder, got folder info")
	}
}

func TestArticleFolderService_MoveArticleToFolder(t *testing.T) {
	aggSvc, artSvc, fldSvc, _, _ := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "目标文件夹",
		ParentID: nil,
	})

	article, _ := artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "待移动文章",
		ArticleType: "markdown",
		FolderID:    nil,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	err := aggSvc.MoveArticleToFolder(ctx, article.ID, &folder.ID)
	if err != nil {
		t.Fatalf("MoveArticleToFolder failed: %v", err)
	}

	updated, _ := artSvc.GetArticle(ctx, article.ID)
	if updated.FolderID == nil || *updated.FolderID != folder.ID {
		t.Errorf("Expected FolderID %d, got %v", folder.ID, updated.FolderID)
	}
}

func TestArticleFolderService_MoveArticleToFolder_FolderNotFound(t *testing.T) {
	aggSvc, artSvc, _, _, _ := createTestServices()
	ctx := context.Background()

	article, _ := artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "文章",
		ArticleType: "markdown",
		FolderID:    nil,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	nonExistentID := uint(9999)
	err := aggSvc.MoveArticleToFolder(ctx, article.ID, &nonExistentID)
	if err != folderErrors.ErrFolderNotFound {
		t.Errorf("Expected ErrFolderNotFound, got %v", err)
	}
}

func TestArticleFolderService_CanDeleteFolder(t *testing.T) {
	aggSvc, artSvc, fldSvc, _, _ := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "测试文件夹",
		ParentID: nil,
	})

	canDelete, count, err := aggSvc.CanDeleteFolder(ctx, folder.ID)
	if err != nil {
		t.Fatalf("CanDeleteFolder failed: %v", err)
	}
	if !canDelete {
		t.Error("Expected canDelete=true for empty folder")
	}
	if count != 0 {
		t.Errorf("Expected count=0, got %d", count)
	}

	_, _ = artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "文章",
		ArticleType: "markdown",
		FolderID:    &folder.ID,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	canDelete, count, err = aggSvc.CanDeleteFolder(ctx, folder.ID)
	if err != nil {
		t.Fatalf("CanDeleteFolder failed: %v", err)
	}
	if canDelete {
		t.Error("Expected canDelete=false for folder with articles")
	}
	if count != 1 {
		t.Errorf("Expected count=1, got %d", count)
	}
}

func TestArticleFolderService_OnArticleCreated(t *testing.T) {
	aggSvc, _, fldSvc, _, _ := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "计数测试",
		ParentID: nil,
	})

	err := aggSvc.OnArticleCreated(ctx, &folder.ID)
	if err != nil {
		t.Fatalf("OnArticleCreated failed: %v", err)
	}

	updated, _ := fldSvc.GetFolder(ctx, folder.ID)
	if updated.ItemCount != 1 {
		t.Errorf("Expected ItemCount=1, got %d", updated.ItemCount)
	}
}

func TestArticleFolderService_OnArticleCreated_NilFolder(t *testing.T) {
	aggSvc, _, _, _, _ := createTestServices()
	ctx := context.Background()

	err := aggSvc.OnArticleCreated(ctx, nil)
	if err != nil {
		t.Errorf("Expected no error for nil folder, got %v", err)
	}
}

func TestArticleFolderService_OnArticleDeleted(t *testing.T) {
	aggSvc, _, fldSvc, _, folderRepo := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "计数测试",
		ParentID: nil,
	})
	folderRepo.folders[folder.ID].ItemCount = 5

	err := aggSvc.OnArticleDeleted(ctx, &folder.ID)
	if err != nil {
		t.Fatalf("OnArticleDeleted failed: %v", err)
	}

	updated, _ := fldSvc.GetFolder(ctx, folder.ID)
	if updated.ItemCount != 4 {
		t.Errorf("Expected ItemCount=4, got %d", updated.ItemCount)
	}
}

func TestArticleFolderService_OnArticleMoved(t *testing.T) {
	aggSvc, _, fldSvc, _, folderRepo := createTestServices()
	ctx := context.Background()

	oldFolder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "旧文件夹",
		ParentID: nil,
	})
	folderRepo.folders[oldFolder.ID].ItemCount = 3

	newFolder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "新文件夹",
		ParentID: nil,
	})
	folderRepo.folders[newFolder.ID].ItemCount = 2

	err := aggSvc.OnArticleMoved(ctx, &oldFolder.ID, &newFolder.ID)
	if err != nil {
		t.Fatalf("OnArticleMoved failed: %v", err)
	}

	oldUpdated, _ := fldSvc.GetFolder(ctx, oldFolder.ID)
	if oldUpdated.ItemCount != 2 {
		t.Errorf("Expected old folder ItemCount=2, got %d", oldUpdated.ItemCount)
	}

	newUpdated, _ := fldSvc.GetFolder(ctx, newFolder.ID)
	if newUpdated.ItemCount != 3 {
		t.Errorf("Expected new folder ItemCount=3, got %d", newUpdated.ItemCount)
	}
}

func TestArticleFolderService_ValidateFolderExists(t *testing.T) {
	aggSvc, _, fldSvc, _, _ := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "存在的文件夹",
		ParentID: nil,
	})

	result, err := aggSvc.ValidateFolderExists(ctx, folder.ID)
	if err != nil {
		t.Fatalf("ValidateFolderExists failed: %v", err)
	}
	if result.ID != folder.ID {
		t.Errorf("Expected folder ID %d, got %d", folder.ID, result.ID)
	}

	_, err = aggSvc.ValidateFolderExists(ctx, 9999)
	if err != folderErrors.ErrFolderNotFound {
		t.Errorf("Expected ErrFolderNotFound, got %v", err)
	}
}

func TestArticleFolderService_ListArticlesWithFolder(t *testing.T) {
	aggSvc, artSvc, fldSvc, _, _ := createTestServices()
	ctx := context.Background()

	folder, _ := fldSvc.CreateFolder(ctx, &folderService.CreateFolderInput{
		Name:     "列表测试",
		ParentID: nil,
	})

	_, _ = artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "文章1",
		ArticleType: "markdown",
		FolderID:    &folder.ID,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	_, _ = artSvc.CreateArticle(ctx, &articleService.CreateArticleInput{
		Title:       "文章2",
		ArticleType: "markdown",
		FolderID:    &folder.ID,
		OwnerID:     1,
		OwnerType:   "admin",
	})

	result, err := aggSvc.ListArticlesWithFolder(ctx, 1, 10, nil, "", "", "", &folder.ID)
	if err != nil {
		t.Fatalf("ListArticlesWithFolder failed: %v", err)
	}

	if result.Total != 2 {
		t.Errorf("Expected total=2, got %d", result.Total)
	}

	for _, record := range result.Records {
		if record.Folder == nil {
			t.Error("Expected folder info in records")
		}
	}
}
