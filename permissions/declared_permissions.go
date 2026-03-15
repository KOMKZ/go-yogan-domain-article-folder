package permissions

type DeclaredPermission struct {
	PermissionCode string
	PermissionName string
	PermissionType string
	ResourceCode   string
	GroupCode      string
	Description    string
}

func DeclaredPermissions() []DeclaredPermission {
	return []DeclaredPermission{
		{
			PermissionCode: "article_folder:read",
			PermissionName: "查看文章分组",
			PermissionType: "READ",
			ResourceCode:   "article_folder",
			GroupCode:      "SYSTEM",
			Description:    "文章分组查询",
		},
		{
			PermissionCode: "article_folder:write",
			PermissionName: "管理文章分组",
			PermissionType: "WRITE",
			ResourceCode:   "article_folder",
			GroupCode:      "SYSTEM",
			Description:    "文章分组创建、更新、删除",
		},
	}
}
