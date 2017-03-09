/*
 Navicat Premium Data Transfer

 Source Server         : columbus.unovo.com.cn
 Source Server Type    : MySQL
 Source Server Version : 50717
 Source Host           : columbus.unovo.com.cn
 Source Database       : mgr

 Target Server Type    : MySQL
 Target Server Version : 50717
 File Encoding         : utf-8

 Date: 03/09/2017 15:00:18 PM
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `t_mgr_admin`
-- ----------------------------
DROP TABLE IF EXISTS `t_mgr_admin`;
CREATE TABLE `t_mgr_admin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `admin_name` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `user_id` bigint(20) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `invalid` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
--  Records of `t_mgr_admin`
-- ----------------------------
BEGIN;
INSERT INTO `t_mgr_admin` VALUES ('1', '系统管理员', '1', '2017-02-17 02:07:35', '2017-02-17 02:07:35', '0'), ('2', 'fd', '2', '2017-02-17 04:04:29', '2017-02-17 04:04:29', '0'), ('3', 'a', '3', '2017-02-17 04:07:38', '2017-02-17 04:07:38', '0'), ('4', 'b', '4', '2017-02-17 12:38:24', '2017-02-17 12:38:24', '0');
COMMIT;

-- ----------------------------
--  Table structure for `t_mgr_admin_role_ref`
-- ----------------------------
DROP TABLE IF EXISTS `t_mgr_admin_role_ref`;
CREATE TABLE `t_mgr_admin_role_ref` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `admin_id` bigint(20) NOT NULL DEFAULT '0',
  `role_id` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `t_mgr_admin_role_ref_admin_id` (`admin_id`),
  KEY `t_mgr_admin_role_ref_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
--  Table structure for `t_mgr_res`
-- ----------------------------
DROP TABLE IF EXISTS `t_mgr_res`;
CREATE TABLE `t_mgr_res` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `res_name` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `path` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `level` int(11) NOT NULL DEFAULT '0',
  `pid` bigint(20) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `invalid` tinyint(1) NOT NULL DEFAULT '0',
  `seq` int(11) NOT NULL DEFAULT '0',
  `res_type` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `res_name` (`res_name`),
  KEY `t_mgr_res_res_name` (`res_name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
--  Records of `t_mgr_res`
-- ----------------------------
BEGIN;
INSERT INTO `t_mgr_res` VALUES ('1', '资源管理', '/static/html/res/home.html', '0', '6', '2017-02-10 01:38:33', '2017-03-06 04:55:31', '0', '1', '0'), ('2', '角色管理', '/static/html/role/home.html', '0', '6', '2017-02-09 17:42:08', '2017-03-06 04:55:25', '0', '2', '0'), ('3', '用户管理', '/static/html/admin/home.html', '0', '6', '2017-02-09 09:43:05', '2017-03-06 04:55:06', '0', '3', '0'), ('4', '测试菜单二', '/static/html/role/home.html', '0', '-1', '2017-02-24 21:26:03', '2017-03-06 04:56:37', '0', '4', '0'), ('5', '测试菜单一', '/static/html/role/home.html', '0', '-1', '2017-02-24 13:26:45', '2017-03-06 04:56:28', '0', '5', '0'), ('6', '系统管理', '', '0', '-1', '2017-03-06 04:54:59', '2017-03-06 04:54:59', '0', '0', '0');
COMMIT;

-- ----------------------------
--  Table structure for `t_mgr_role`
-- ----------------------------
DROP TABLE IF EXISTS `t_mgr_role`;
CREATE TABLE `t_mgr_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `invalid` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_name` (`role_name`),
  KEY `t_mgr_role_role_name` (`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
--  Table structure for `t_mgr_student`
-- ----------------------------
DROP TABLE IF EXISTS `t_mgr_student`;
CREATE TABLE `t_mgr_student` (
  `invalid` tinyint(1) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `last_update_time` datetime NOT NULL,
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
--  Records of `t_mgr_student`
-- ----------------------------
BEGIN;
INSERT INTO `t_mgr_student` VALUES ('0', '2017-02-16 18:47:03', '2017-02-16 18:47:03', '1', 'aaaaa');
COMMIT;

-- ----------------------------
--  Table structure for `t_mgr_user`
-- ----------------------------
DROP TABLE IF EXISTS `t_mgr_user`;
CREATE TABLE `t_mgr_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `password` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `invalid` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `username_2` (`username`),
  KEY `t_mgr_user_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
--  Records of `t_mgr_user`
-- ----------------------------
BEGIN;
INSERT INTO `t_mgr_user` VALUES ('1', 'admin', 'admin22', '2017-02-17 02:07:35', '2017-02-17 02:07:35', '0'), ('2', 'fd', 'fd', '2017-02-17 04:04:29', '2017-02-17 04:04:29', '0'), ('3', 'a', 'a', '2017-02-17 04:07:38', '2017-02-17 04:07:38', '0'), ('4', 'b', 'b', '2017-02-17 12:38:24', '2017-02-17 12:38:24', '0');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
