# go-yogan-domain-article-folder

文章-文件夹聚合领域包，负责跨领域的编排逻辑，解决 `article` 和 `folder` 领域之间的依赖关系。

## 设计理念

遵循 **依赖倒置原则（DIP）**：
- `article` 领域不依赖 `folder` 领域
- `folder` 领域不依赖 `article` 领域
- 本包作为聚合层，编排两个领域的交互

## 特性

- **跨领域编排**：协调 article 和 folder 的交互
- **数据组装**：为文章附加文件夹信息（名称、路径、面包屑）
- **批量优化**：列表查询时批量获取 folder 信息，避免 N+1 问题
- **验证逻辑**：文件夹存在性验证、删除前检查等

## 安装

```bash
go get github.com/KOMKZ/go-yogan-domain-article-folder
```

## 依赖

- `github.com/KOMKZ/go-yogan-domain-article`
- `github.com/KOMKZ/go-yogan-domain-folder`

## 使用示例

```go
package main

import (
    articlefolder "github.com/KOMKZ/go-yogan-domain-article-folder"
    article "github.com/KOMKZ/go-yogan-domain-article"
    folder "github.com/KOMKZ/go-yogan-domain-folder"
)

func main() {
    // 初始化底层领域服务
    articleSvc := article.NewService(articleRepo)
    folderSvc := folder.NewService(folderRepo)

    // 创建聚合服务
    svc := articlefolder.NewService(articleSvc, folderSvc)

    // 获取带文件夹信息的文章详情
    articleWithFolder, err := svc.GetArticleWithFolder(ctx, articleID)
    // 返回结构包含 folder 信息：
    // {
    //   "id": 1,
    //   "title": "文章标题",
    //   "folderId": 5,
    //   "folder": {
    //     "id": 5,
    //     "name": "Go语言",
    //     "path": "/1/5/",
    //     "breadcrumb": [
    //       {"id": 1, "name": "技术文章"},
    //       {"id": 5, "name": "Go语言"}
    //     ]
    //   }
    // }

    // 获取带文件夹信息的文章列表
    result, err := svc.ListArticlesWithFolder(ctx, 1, 10, nil, "", "", "", nil)

    // 移动文章到指定文件夹（带验证）
    err = svc.MoveArticleToFolder(ctx, articleID, &folderID)

    // 检查文件夹是否可删除（是否有文章）
    canDelete, count, err := svc.CanDeleteFolder(ctx, folderID)
}
```

## API

### Service

| 方法 | 说明 |
|------|------|
| `GetArticleWithFolder` | 获取文章详情（带文件夹信息） |
| `ListArticlesWithFolder` | 获取文章列表（带文件夹信息，批量优化） |
| `MoveArticleToFolder` | 移动文章到指定文件夹（带验证） |
| `ValidateFolderExists` | 验证文件夹是否存在 |
| `CanDeleteFolder` | 检查文件夹是否可删除 |

### 数据结构

```go
// FolderInfo 文件夹信息
type FolderInfo struct {
    ID         uint             `json:"id"`
    Name       string           `json:"name"`
    Path       string           `json:"path"`
    Breadcrumb []BreadcrumbItem `json:"breadcrumb"`
}

// BreadcrumbItem 面包屑项
type BreadcrumbItem struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

// ArticleWithFolder 带文件夹信息的文章
type ArticleWithFolder struct {
    *articleModel.Article
    Folder *FolderInfo `json:"folder,omitempty"`
}
```

## License

MIT
