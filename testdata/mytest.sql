SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_cron_daily
-- ----------------------------
DROP TABLE IF EXISTS `t_cron_daily`;
CREATE TABLE `t_cron_daily`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '任务ID',
  `is_active` bit(1) NOT NULL DEFAULT b'0' COMMENT '有效',
  `workday` bit(1) NOT NULL DEFAULT b'0' COMMENT '工作日',
  `weekday` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '周X|周Y...',
  `run_clock` char(8) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '具体时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `task_id`(`task_id`) USING BTREE,
  INDEX `run_clock`(`run_clock`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '日常执行' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of t_cron_daily
-- ----------------------------

-- ----------------------------
-- Table structure for t_cron_notice
-- ----------------------------
DROP TABLE IF EXISTS `t_cron_notice`;
CREATE TABLE `t_cron_notice`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '用户ID',
  `task_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '任务ID',
  `is_active` bit(1) NOT NULL DEFAULT b'0' COMMENT '有效',
  `important` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '重要程度',
  `message` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '消息内容',
  `read_time` datetime(0) NULL DEFAULT NULL COMMENT '阅读时间',
  `delay_start_time` datetime(0) NULL DEFAULT NULL COMMENT '推迟开始时间',
  `start_time` datetime(0) NULL DEFAULT NULL COMMENT '开始时间',
  `stop_time` datetime(0) NULL DEFAULT NULL COMMENT '结束时间',
  `start_clock` char(8) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '开始时刻',
  `stop_clock` char(8) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '结束时刻',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `task_id`(`task_id`) USING BTREE,
  INDEX `read_time`(`read_time`) USING BTREE,
  INDEX `delay_start_time`(`delay_start_time`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '消息提醒' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of t_cron_notice
-- ----------------------------

-- ----------------------------
-- Table structure for t_cron_task
-- ----------------------------
DROP TABLE IF EXISTS `t_cron_task`;
CREATE TABLE `t_cron_task`  (
  `task_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(10) UNSIGNED NULL DEFAULT 0 COMMENT '用户ID',
  `refer_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联任务ID',
  `is_active` bit(1) NOT NULL DEFAULT b'0' COMMENT '有效',
  `behind` smallint(6) NOT NULL DEFAULT 0 COMMENT '相对推迟/提前多少分钟',
  `action_type` enum('command','message','http_get','http_post','function') CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'command' COMMENT '动作类型',
  `cmd_url` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '指令或网址',
  `args_data` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '参数或消息体',
  `last_time` datetime(0) NULL DEFAULT NULL COMMENT '最后执行时间',
  `last_result` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '执行结果',
  `last_error` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '出错信息',
  PRIMARY KEY (`task_id`) USING BTREE,
  INDEX `refer_id`(`refer_id`) USING BTREE,
  INDEX `last_time`(`last_time`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '定时任务' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of t_cron_task
-- ----------------------------

-- ----------------------------
-- Table structure for t_cron_timer
-- ----------------------------
DROP TABLE IF EXISTS `t_cron_timer`;
CREATE TABLE `t_cron_timer`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '任务ID',
  `is_active` bit(1) NOT NULL DEFAULT b'0' COMMENT '有效',
  `run_date` date NULL DEFAULT NULL COMMENT '指定日期',
  `run_clock` char(8) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '具体时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `task_id`(`task_id`) USING BTREE,
  INDEX `run_date`(`run_date`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '定时执行' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of t_cron_timer
-- ----------------------------

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user`  (
  `user_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `realname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '昵称/称呼',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '手机号码',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '电子邮箱',
  `avatar` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '头像',
  `introduction` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '介绍说明',
  `created_at` timestamp(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp(0) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`user_id`) USING BTREE,
  INDEX `idx_t_user_username`(`username`) USING BTREE,
  INDEX `idx_t_user_mobile`(`mobile`) USING BTREE,
  INDEX `idx_t_user_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO `t_user` VALUES (1, 'admin', '09e8ff53$x1KWXASXqGRzA7YwipQhibg/0LMtkoU39VfW8EYtxAI=', '管理员', NULL, NULL, '/avatars/avatar-admin.jpg', '不受限的超管账号。', '2019-12-01 03:12:00', '2019-12-01 03:12:00', NULL);
INSERT INTO `t_user` VALUES (2, 'demo', 'acfd1f8b$o6ySKi7yaMmZrKIaT4O/oGUoei6n/xKOXik4PtXuvwk=', '演示用户', NULL, NULL, '/avatars/avatar-demo.jpg', '演示和测试账号。', '2019-12-01 03:12:00', '2019-12-01 03:12:00', NULL);

SET FOREIGN_KEY_CHECKS = 1;
