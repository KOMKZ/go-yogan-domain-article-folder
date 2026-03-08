-- 文章文件夹表（使用 folder 领域的通用 schema）
CREATE TABLE IF NOT EXISTS `article_folders` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL COMMENT '文件夹名称',
    `parent_id` bigint unsigned COMMENT '父文件夹ID',
    `sort_order` int DEFAULT 0 COMMENT '排序号',
    `depth` int DEFAULT 0 COMMENT '层级深度',
    `path` varchar(1000) COMMENT '物化路径，如 /1/3/5/',
    `item_count` int unsigned NOT NULL DEFAULT 0 COMMENT '直接子项数量',
    `total_item_count` int unsigned NOT NULL DEFAULT 0 COMMENT '所有子孙的总数量',
    `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` timestamp NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    INDEX `idx_parent_id` (`parent_id`),
    INDEX `idx_path` (`path`(255)),
    INDEX `idx_deleted_at` (`deleted_at`),
    INDEX `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章文件夹表';
