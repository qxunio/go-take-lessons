/*
 Source Server Type    : MySQL
 Source Server Version : 80021
 Source Schema         : take_lessons

 Target Server Type    : MySQL
 Target Server Version : 80021
 File Encoding         : 65001

 Date: 19/12/2020 00:21:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for configuration_subject
-- ----------------------------
DROP TABLE IF EXISTS `configuration_subject`;
CREATE TABLE `configuration_subject`  (
                                          `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                                          `event_id` bigint(0) UNSIGNED NOT NULL COMMENT '选课ID',
                                          `subject_id` bigint(0) UNSIGNED NOT NULL COMMENT '学科ID',
                                          `class_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '课堂名称（应对篮球一班，篮球二班的情况）',
                                          `subject_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '学科名称',
                                          `num` int(0) UNSIGNED NOT NULL COMMENT '该课程的限制人数',
                                          `selected_places` int(0) UNSIGNED NULL DEFAULT NULL COMMENT '已选人数',
                                          `teach_address` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '教学地址',
                                          `teach_time` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '教学时间',
                                          `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                                          `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                                          `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                                          `introduction` varchar(912) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '介绍',
                                          PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '选课学科配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for configuration_teacher
-- ----------------------------
DROP TABLE IF EXISTS `configuration_teacher`;
CREATE TABLE `configuration_teacher`  (
                                          `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                                          `event_id` bigint(0) UNSIGNED NOT NULL COMMENT '选课ID',
                                          `cs_id` bigint(0) UNSIGNED NOT NULL COMMENT '配置学科表ID',
                                          `teacher_id` bigint(0) UNSIGNED NOT NULL COMMENT '教师ID(user表id)',
                                          `teacher_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '教师姓名',
                                          `teacher_account` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '教师账号',
                                          `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                                          `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                                          `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                                          PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '选课教师配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for event
-- ----------------------------
DROP TABLE IF EXISTS `event`;
CREATE TABLE `event`  (
                          `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                          `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
                          `can_update` int(0) NULL DEFAULT NULL COMMENT '学生是否可以修改已经选择的课程,1可以，0不可以',
                          `num` int(0) UNSIGNED NOT NULL COMMENT '每名学生最多可选择的课程数',
                          `school_year` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '选课对应的学生入学年（年级）',
                          `tag_ids` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
                          `stime` datetime(0) NOT NULL COMMENT '开始时间',
                          `etime` datetime(0) NOT NULL COMMENT '结束时间',
                          `status` tinyint(0) UNSIGNED NOT NULL COMMENT '0默认状态（学生不可查询），1表示即将或进行中（学生可查询），2表示已经结束，成历史状态',
                          `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                          `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                          `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                          `creator` bigint(0) UNSIGNED NOT NULL COMMENT '创建者',
                          PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '选课表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`  (
                         `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                         `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
                         `router` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '路由',
                         `priority` int(0) UNSIGNED NOT NULL COMMENT '优先级',
                         `parent_id` bigint(0) UNSIGNED NULL DEFAULT NULL COMMENT '父ID',
                         `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                         `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '菜单表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (1235809292771135488, '用户', '/home/user', 1, NULL, 1, '2020-03-06 14:14:43');
INSERT INTO `menu` VALUES (1235811124583075840, '公告', '/home/notice', 5, NULL, 1, '2020-03-06 14:15:33');
INSERT INTO `menu` VALUES (1235811288966238208, '帮助', '/home/help', 6, NULL, 1, '2020-03-06 14:16:13');
INSERT INTO `menu` VALUES (1235811453122908160, '学科', '/home/subject', 3, NULL, 1, '2020-03-06 14:16:44');
INSERT INTO `menu` VALUES (1235811453122908180, '选课', '/home/event', 4, NULL, 1, '2020-03-24 16:36:21');
INSERT INTO `menu` VALUES (1235811653122908187, '我的关注', '/home/stu-focus', 1, NULL, 1, '2020-04-02 16:36:21');
INSERT INTO `menu` VALUES (1235811653144608155, '选课', '/home/stu-event', 2, NULL, 1, '2020-04-02 16:36:21');
INSERT INTO `menu` VALUES (1335811653122908134, '历史', '/home/stu-history', 3, NULL, 1, '2020-04-02 16:36:21');
INSERT INTO `menu` VALUES (1335811653122938154, '我的学科', '/home/teacher-subject', 1, NULL, 1, '2020-04-02 16:36:21');
INSERT INTO `menu` VALUES (1335811653122938156, '标签', '/home/tag', 2, NULL, 1, '2020-09-14 20:51:08');

-- ----------------------------
-- Table structure for notice
-- ----------------------------
DROP TABLE IF EXISTS `notice`;
CREATE TABLE `notice`  (
                           `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                           `content` varchar(330) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知内容',
                           `title` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知主题',
                           `type` tinyint(0) UNSIGNED NOT NULL COMMENT '通知类型，0表示站内，1表示站外',
                           `expire_date` datetime(0) NOT NULL COMMENT '到期时间',
                           `status` tinyint(1) NOT NULL COMMENT '状态，0失效，1有效',
                           `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                           `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                           `utime` datetime(0) NOT NULL COMMENT '更新时间',
                           `creator` bigint(0) UNSIGNED NOT NULL COMMENT '创建者',
                           PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '公告' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
                         `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                         `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
                         `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                         `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (1235808272481521664, '管理员', 1, '2020-03-06 14:04:01');
INSERT INTO `role` VALUES (1235808465281093632, '教师', 1, '2020-03-06 14:04:40');
INSERT INTO `role` VALUES (1235808540581433344, '学生', 1, '2020-03-06 14:04:53');

-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
DROP TABLE IF EXISTS `role_menu`;
CREATE TABLE `role_menu`  (
                              `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                              `menu_id` bigint(0) UNSIGNED NOT NULL COMMENT '菜单ID',
                              `role_id` bigint(0) UNSIGNED NOT NULL COMMENT '角色ID',
                              `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                              `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                              `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间(禁用时间)',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '菜单角色中间表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_menu
-- ----------------------------
INSERT INTO `role_menu` VALUES (1235811647927357440, 1235809292771135488, 1235808272481521664, 1, '2020-03-06 14:17:50', NULL);
INSERT INTO `role_menu` VALUES (1235811647927359900, 1335811653122938156, 1235808272481521664, 1, '2020-09-14 20:17:50', NULL);
INSERT INTO `role_menu` VALUES (1235811838990487552, 1235811124583075840, 1235808272481521664, 0, '2020-03-06 14:17:50', NULL);
INSERT INTO `role_menu` VALUES (1235811934901637120, 1235811288966238208, 1235808272481521664, 1, '2020-03-06 14:17:50', NULL);
INSERT INTO `role_menu` VALUES (1235811938571951162, 1235811653144608155, 1235808540581433344, 1, '2020-04-02 14:17:50', NULL);
INSERT INTO `role_menu` VALUES (1235811988571951104, 1235811453122908160, 1235808272481521664, 1, '2020-03-06 14:17:50', NULL);
INSERT INTO `role_menu` VALUES (1235811988578512312, 1235811453122908180, 1235808272481521664, 1, '2020-03-24 16:42:59', NULL);
INSERT INTO `role_menu` VALUES (1335811388578512317, 1335811653122908134, 1235808540581433344, 1, '2020-04-02 16:42:59', NULL);
INSERT INTO `role_menu` VALUES (1635817388576512317, 1335811653122938154, 1235808465281093632, 1, '2020-04-02 16:42:59', NULL);
INSERT INTO `role_menu` VALUES (4235811934901537122, 1235811653122908187, 1235808540581433344, 1, '2020-04-02 14:17:50', NULL);

-- ----------------------------
-- Table structure for school_year
-- ----------------------------
DROP TABLE IF EXISTS `school_year`;
CREATE TABLE `school_year`  (
                                `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                                `school_year` int(0) NOT NULL COMMENT '入学年',
                                `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                                `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                                `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                                PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '入学年表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for stu_focus
-- ----------------------------
DROP TABLE IF EXISTS `stu_focus`;
CREATE TABLE `stu_focus`  (
                              `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                              `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '学生ID',
                              `event_id` bigint(0) UNSIGNED NOT NULL COMMENT '选课ID',
                              `event_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '选课事件',
                              `cs_id` bigint(0) UNSIGNED NOT NULL COMMENT '课程ID',
                              `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                              `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                              `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                              `school_year` int(0) UNSIGNED NOT NULL COMMENT '学生的入学年',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '学生关注' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for stu_subject
-- ----------------------------
DROP TABLE IF EXISTS `stu_subject`;
CREATE TABLE `stu_subject`  (
                                `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                                `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '学生ID',
                                `event_id` bigint(0) UNSIGNED NOT NULL COMMENT '选课ID',
                                `event_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '选课事件',
                                `cs_id` bigint(0) UNSIGNED NOT NULL COMMENT '课程ID',
                                `class` int(0) NULL DEFAULT NULL COMMENT '班级',
                                `school_year` int(0) NOT NULL COMMENT '学生的入学年（级）',
                                `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                                `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                                `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                                PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '学生选课表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for subject
-- ----------------------------
DROP TABLE IF EXISTS `subject`;
CREATE TABLE `subject`  (
                            `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                            `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
                            `introduction` varchar(912) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '介绍',
                            `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                            `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                            `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新⁮时间',
                            `creator` bigint(0) UNSIGNED NOT NULL COMMENT '创建者',
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '学科表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`  (
                        `id` bigint(0) NOT NULL,
                        `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                        `total_num` int(0) NOT NULL,
                        `ctime` datetime(0) NOT NULL,
                        `creator` bigint(0) NOT NULL,
                        `enable` int(0) NOT NULL,
                        `utime` datetime(0) NULL DEFAULT NULL,
                        PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for tag_stu
-- ----------------------------
DROP TABLE IF EXISTS `tag_stu`;
CREATE TABLE `tag_stu`  (
                            `id` bigint(0) NOT NULL,
                            `uid` bigint(0) NOT NULL,
                            `tag_id` bigint(0) NOT NULL,
                            `ctime` datetime(0) NOT NULL,
                            `creator` bigint(0) NOT NULL,
                            `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
                            `class` int(0) NULL DEFAULT NULL,
                            `school_year` int(0) NULL DEFAULT NULL,
                            `account` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
                         `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                         `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '姓名',
                         `account` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账号',
                         `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
                         `salt` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '盐',
                         `school_year` int(0) UNSIGNED NOT NULL COMMENT '学生的入学年',
                         `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '账号状态',
                         `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                         `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
                         `type` tinyint(0) UNSIGNED NOT NULL COMMENT '类型',
                         `class` int(0) NULL DEFAULT NULL COMMENT '班级',
                         `creator` bigint(0) UNSIGNED NULL DEFAULT NULL COMMENT '创建者',
                         PRIMARY KEY (`id`) USING BTREE,
                         UNIQUE INDEX `account`(`account`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1235806822435131392, 'domain', 'domain', 'b41842fedc9304cf2871fc5d684e8b3a4a679e6b89ba14cf3d675d1aa44f2c3b', 'N4ZyrqP+Tj', 0, 1, '2020-03-06 13:58:35', '2020-03-06 13:58:37', 1, NULL, 0);

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role`  (
                              `id` bigint(0) UNSIGNED NOT NULL COMMENT 'PK',
                              `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '用户ID',
                              `role_id` bigint(0) UNSIGNED NOT NULL COMMENT '角色ID',
                              `enable` tinyint(0) UNSIGNED NOT NULL COMMENT '是否启用',
                              `ctime` datetime(0) NOT NULL COMMENT '创建时间',
                              `utime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间(禁用时间)',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户角色中间表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_role
-- ----------------------------
INSERT INTO `user_role` VALUES (1235806822435131392, 1235806822435131392, 1235808272481521664, 1, '2020-09-24 10:56:10', '2020-09-24 10:56:12');

SET FOREIGN_KEY_CHECKS = 1;
