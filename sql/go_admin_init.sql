-- Go Admin 初始化 SQL
-- 适用: MySQL 8.x
-- 导入后默认账号: admin / 123456

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS `go_admin`
DEFAULT CHARACTER SET utf8mb4
DEFAULT COLLATE utf8mb4_unicode_ci;

USE `go_admin`;

DROP TABLE IF EXISTS `operation_logs`;
DROP TABLE IF EXISTS `role_menus`;
DROP TABLE IF EXISTS `role_permissions`;
DROP TABLE IF EXISTS `user_roles`;
DROP TABLE IF EXISTS `menus`;
DROP TABLE IF EXISTS `permissions`;
DROP TABLE IF EXISTS `roles`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `username` varchar(64) NOT NULL,
  `nickname` varchar(64) NOT NULL,
  `password` varchar(255) NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `name` varchar(64) NOT NULL,
  `code` varchar(64) NOT NULL,
  `description` varchar(255) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_roles_name` (`name`),
  UNIQUE KEY `uk_roles_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `name` varchar(64) NOT NULL,
  `code` varchar(64) NOT NULL,
  `module` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_permissions_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL,
  `title` varchar(64) NOT NULL,
  `path` varchar(128) NOT NULL,
  `component` varchar(128) NOT NULL,
  `icon` varchar(64) NOT NULL DEFAULT '',
  `menu_type` varchar(16) NOT NULL DEFAULT 'menu',
  `permission` varchar(64) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `hidden` tinyint(1) NOT NULL DEFAULT 0,
  `keep_alive` tinyint(1) NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_menus_parent_id` (`parent_id`),
  UNIQUE KEY `uk_menus_path` (`path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_user_roles_role_id` (`role_id`),
  CONSTRAINT `fk_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `role_permissions` (
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`),
  KEY `idx_role_permissions_permission_id` (`permission_id`),
  CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `role_menus` (
  `role_id` bigint unsigned NOT NULL,
  `menu_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`, `menu_id`),
  KEY `idx_role_menus_menu_id` (`menu_id`),
  CONSTRAINT `fk_role_menus_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_role_menus_menu_id` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `user_id` bigint unsigned NOT NULL DEFAULT 0,
  `username` varchar(64) NOT NULL DEFAULT '',
  `method` varchar(16) NOT NULL DEFAULT '',
  `path` varchar(255) NOT NULL DEFAULT '',
  `status_code` int NOT NULL DEFAULT 0,
  `ip` varchar(64) NOT NULL DEFAULT '',
  `user_agent` varchar(255) NOT NULL DEFAULT '',
  `action` varchar(255) NOT NULL DEFAULT '',
  `latency_ms` bigint NOT NULL DEFAULT 0,
  `request_id` varchar(64) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `permissions` (`id`, `name`, `code`, `module`) VALUES
  (1, '用户列表', 'user:list', 'user'),
  (2, '用户保存', 'user:save', 'user'),
  (7, '用户删除', 'user:delete', 'user'),
  (8, '用户状态', 'user:status', 'user'),
  (3, '角色列表', 'role:list', 'role'),
  (4, '角色保存', 'role:save', 'role'),
  (9, '角色删除', 'role:delete', 'role'),
  (10, '角色状态', 'role:status', 'role'),
  (5, '菜单保存', 'menu:save', 'menu'),
  (6, '日志查看', 'log:list', 'log');

INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`) VALUES
  (1, '超级管理员', 'super-admin', '拥有全部后台权限', 1),
  (2, '普通管理员', 'admin', '可按需分配后台权限', 1);

INSERT INTO `menus` (`id`, `parent_id`, `name`, `title`, `path`, `component`, `icon`, `menu_type`, `permission`, `sort`, `hidden`, `keep_alive`, `status`) VALUES
  (1, 0, 'Dashboard', '工作台', '/dashboard', 'dashboard/index', 'House', 'menu', '', 1, 0, 0, 1),
  (2, 0, 'System', '系统管理', '/system', 'layout/router-view', 'Setting', 'menu', '', 2, 0, 0, 1),
  (3, 2, 'User', '用户管理', '/system/user', 'system/user/index', 'User', 'menu', 'user:list', 1, 0, 0, 1),
  (4, 2, 'Role', '角色管理', '/system/role', 'system/role/index', 'Avatar', 'menu', 'role:list', 2, 0, 0, 1),
  (5, 2, 'Menu', '菜单管理', '/system/menu', 'system/menu/index', 'Menu', 'menu', 'menu:save', 3, 0, 0, 1),
  (6, 2, 'Log', '操作日志', '/system/log', 'system/log/index', 'Document', 'menu', 'log:list', 4, 0, 0, 1);

INSERT INTO `users` (`id`, `username`, `nickname`, `password`, `status`) VALUES
  (1, 'admin', '系统管理员', '$2a$10$PcO31nlWwhvuWYuinCJfyurJ1/z0Ri7FMpiyH/5ezIaYq9jC9irM.', 1);

INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
  (1, 1);

INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
  (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
  (2, 1), (2, 3), (2, 6);

INSERT INTO `role_menus` (`role_id`, `menu_id`) VALUES
  (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6),
  (2, 1), (2, 2), (2, 3), (2, 4), (2, 6);

SET FOREIGN_KEY_CHECKS = 1;
