/*
 Navicat Premium Data Transfer

 Source Server         : 主从dockermaster
 Source Server Type    : MySQL
 Source Server Version : 80023
 Source Host           : 127.0.0.1:3307
 Source Schema         : fileserver

 Target Server Type    : MySQL
 Target Server Version : 80023
 File Encoding         : 65001

 Date: 26/04/2021 23:44:53
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for tbl_file
-- ----------------------------
DROP TABLE IF EXISTS `tbl_file`;
CREATE TABLE `tbl_file` (
  `id` int NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
  `ext1` int DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of tbl_file
-- ----------------------------
BEGIN;
INSERT INTO `tbl_file` VALUES (4, 'ce1155eb4ffb11741a7ae0dba630fcc8981472a2', '0.jpeg', 156394, 'static/file/0.jpeg', '2021-03-12 21:23:35', '2021-03-12 21:23:35', 1, 0, NULL);
INSERT INTO `tbl_file` VALUES (5, 'ceffcc6a4eb0cd8a6bd529723c7a3ecfabdcdeac', 'e0ebcde492f291aa01fd256a053306bd.jpg', 190275, 'static/file/e0ebcde492f291aa01fd256a053306bd.jpg', '2021-03-17 23:05:36', '2021-03-17 23:05:36', 1, 0, NULL);
COMMIT;

-- ----------------------------
-- Table structure for tbl_user
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `email` varchar(64) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(128) DEFAULT '' COMMENT '手机号',
  `email_validated` tinyint(1) DEFAULT '0' COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) DEFAULT '0' COMMENT '手机号是否已验证',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `PROFILE` text COMMENT '用户属性',
  `status` int NOT NULL DEFAULT '0' COMMENT '账户状态（启用/禁用/锁定/标记删除等）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of tbl_user
-- ----------------------------
BEGIN;
INSERT INTO `tbl_user` VALUES (26, 'uploadprocess', '1a48f352ffe1e15d1c4611b22615df0092ec194c', '', '', 0, 0, '2021-03-17 23:10:01', '2021-03-17 23:10:01', NULL, 0);
INSERT INTO `tbl_user` VALUES (33, 'admin', '1a48f352ffe1e15d1c4611b22615df0092ec194c', '', '1', 0, 0, '2021-03-17 23:25:34', '2021-03-17 23:25:34', NULL, 0);
INSERT INTO `tbl_user` VALUES (34, 'admin', '1a48f352ffe1e15d1c4611b22615df0092ec194c', '', '3', 0, 0, '2021-03-17 23:25:54', '2021-03-17 23:25:54', NULL, 0);
COMMIT;

-- ----------------------------
-- Table structure for tbl_user_token
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_token`;
CREATE TABLE `tbl_user_token` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登录token',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of tbl_user_token
-- ----------------------------
BEGIN;
INSERT INTO `tbl_user_token` VALUES (28, 'admin2', '2315365e35f98f4f5ae38e12893ed9c9605b6b99');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
