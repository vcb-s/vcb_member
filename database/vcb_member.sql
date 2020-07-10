/*
 Navicat Premium Data Transfer

 Source Server         : 本地
 Source Server Type    : MariaDB
 Source Server Version : 100504
 Source Host           : ::1:3306
 Source Schema         : vcb_member

 Target Server Type    : MariaDB
 Target Server Version : 100504
 File Encoding         : 65001

 Date: 10/07/2020 18:59:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for login_association
-- ----------------------------
DROP TABLE IF EXISTS `login_association`;
CREATE TABLE `login_association`  (
  `id` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'id',
  `uid` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户id',
  `type` tinyint(1) NOT NULL COMMENT '绑定类型\r\n1. 主站',
  `auth` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '第三方认证code',
  `create_at` timestamp(0) NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `update_at` timestamp(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'uid',
  `avast` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `pass` char(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '登录密码',
  `admin` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '组长身份，管理什么组，逗号分隔；普通组员为空字符串',
  `ban` int(1) NOT NULL DEFAULT 2 COMMENT '是否禁止登录，1是2否',
  `group` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户所属分组，逗号分隔',
  `last_token_key_id` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录使用的token keyID',
  `create_at` timestamp(0) NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `update_at` timestamp(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `main_key`(`id`) USING BTREE COMMENT '主键'
) ENGINE = MyISAM CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_card_group
-- ----------------------------
DROP TABLE IF EXISTS `user_card_group`;
CREATE TABLE `user_card_group`  (
  `id` int(255) NOT NULL COMMENT '主键',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组名',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_crad
-- ----------------------------
DROP TABLE IF EXISTS `user_crad`;
CREATE TABLE `user_crad`  (
  `id` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '卡片id',
  `avast` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '昵称',
  `bio` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '卡片简介',
  `job` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '职位',
  `order` int(11) NOT NULL DEFAULT 0 COMMENT '排序权重',
  `group` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '5' COMMENT '所在组 具体参见group表',
  `retired` int(1) NOT NULL DEFAULT 1 COMMENT '是否退休,\r\n是1\r\n否2',
  `create_at` timestamp(0) NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `update_at` timestamp(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '最近修改时间',
  `deleted_at` timestamp(0) NULL DEFAULT NULL COMMENT '删除时间',
  `uid` char(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '关联uid',
  `hide` int(1) NOT NULL DEFAULT 2 COMMENT '是否隐藏，1是2否',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `main_key`(`id`) USING BTREE COMMENT '主键'
) ENGINE = MyISAM CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
