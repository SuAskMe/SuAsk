-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: suask
-- ------------------------------------------------------
-- 为当前使用的suask.sql生成的测试数据

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Records of users
-- ----------------------------

DELETE FROM users;

INSERT INTO `users` VALUES (1, 'root', '1@suask.me', 'iqjbenzvfz', '6917ef8afa3ffeb3cb02643b9feb2a46', 'student', 'root', '# 这里是root', NULL, '2024-12-19 03:25:28', '2024-12-19 03:25:42', NULL);
INSERT INTO `users` VALUES (2, 'teacher', '2@suask.me', 'dDQCQ5zaNL', '26e6d44c96358cdfc6ea25d15fa2442a', 'student', 'teacher', '# 我是老师', NULL, '2024-12-19 03:27:03', '2024-12-19 03:27:03', NULL);
INSERT INTO `users` VALUES (11, '苏玉鑫', 'suyx35@mail.sysu.edu.cn', 'dDQCQ5zaNL', '26e6d44c96358cdfc6ea25d15fa2442a', 'teacher', '苏玉鑫', '', NULL, '2024-12-19 03:27:03', '2024-12-19 03:27:03', NULL);
INSERT INTO `users` VALUES (1000, 'student', '3@suask.me', 'Ettykmu1Zc', '4fd5724ce6fdbe4ac6df655eb1dc31dc', 'student', 'student', '# 我是学生', NULL, '2024-12-19 03:26:25', '2024-12-19 03:26:25', NULL);

-- ----------------------------
-- Records of teachers
-- ----------------------------

DELETE FROM teachers;

INSERT INTO `teachers` VALUES (11, 0, '苏玉鑫', 'https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC02668_%285%29.jpg?itok=J5DfVj1v', '苏玉鑫 (B站ID：鸭大坑导)，软件工程学院副院长，副教授，博士生导师，国家重点研发计划青年科学家项目负责人，中国计算机学会（CCF）高级会员、服务计算专委会执行委员。2021年7月入选中山大学百人计划，加入软件工程学院。主要研究方向为系统软件可靠性分析与运行时性能优化，具体包括操作系统、分布式系统、云计算、云原生系统、日志分析、云系统可靠性与智能运维(AIOps)等。近年来在国际会议和期刊共发表近30篇论文，其中24篇发表于ICSE、ASE、ISSTA、SOSP、FAST、ICDE、CVPR、SIGIR、AAAI、IJCAI、CSUR、TKDE等软件工程、操作系统、分布式系统、人工智能等领域CCF A类顶级会议与期刊。', 'suyx35@mail.sysu.edu.cn', 'public');

-- ----------------------------
-- Records of questions
-- ----------------------------

DELETE FROM questions;

INSERT INTO `questions` VALUES (1,1,NULL,'你认为是什么导致人们对旅行的热爱？','不同的文化、美丽的风景，以及逃离日常生活的压力。',_binary '\0','2016-04-04 14:51:12',652, 0);
INSERT INTO `questions` VALUES (2,1,NULL,'如果可以选择任何一个历史人物共进晚餐，你会选择谁？','达·芬奇，想听听他的创意和思想。',_binary '\0','2009-02-07 03:34:10',195, 0);
INSERT INTO `questions` VALUES (3,1,NULL,'你最喜欢的书是什么，为什么？','《1984》，因为它对权力和自由的深刻思考令人警醒。',_binary '\0','2014-01-08 18:45:44',885, 0);
INSERT INTO `questions` VALUES (4,1,NULL,'有哪些你认为值得一试的爱好或活动？','尝试陶艺或摄影，可以提升创造力。',_binary '\0','2001-06-22 13:21:38',427, 0);
INSERT INTO `questions` VALUES (5,1,NULL,'你认为科技会如何改变未来的工作方式？','更多的远程工作、自动化和人工智能将改变传统的工作模式。',_binary '\0','2020-02-08 14:57:07',99, 0);
INSERT INTO `questions` VALUES (6,1,NULL,'你最难忘的假期是去哪里，为什么？','去了一次日本，丰富的文化和美味的食物给我留下深刻印象。',_binary '\0','2008-03-26 17:43:25',833, 0);
INSERT INTO `questions` VALUES (7,1,NULL,'如果能学会一种乐器，你希望学习什么？','钢琴，因为它的音色丰富多样。',_binary '\0','2012-01-02 13:50:47',255, 0);
INSERT INTO `questions` VALUES (8,1,NULL,'你对环保有什么建议或想法？','大力推广可再生能源和减少塑料使用是关键。',_binary '\0','2006-06-17 13:53:20',756, 0);
INSERT INTO `questions` VALUES (9,1,NULL,'在你的生活中，有没有什么事情让你感到特别感恩？','感恩家人和朋友的支持，他们在我困难时给予了我力量。',_binary '\0','2002-12-07 01:39:24',398, 0);
INSERT INTO `questions` VALUES (10,1,NULL,'你觉得友情和爱情之间最大的区别是什么？','友情更倾向于无条件的支持，而爱情通常伴随着深度的情感和浪漫。',_binary '\0','2011-02-13 00:29:18',273, 0);

-- ----------------------------
-- Records of favorites
-- ----------------------------

DELETE FROM favorites;

INSERT INTO `favorites` VALUES (1,1,1,'2024-12-18 13:25:35','默认收藏夹');
INSERT INTO `favorites` VALUES (2,1,2,'2024-12-15 09:36:32','默认收藏夹');
INSERT INTO `favorites` VALUES (3,1,3,'2024-12-15 08:10:24','默认收藏夹');
INSERT INTO `favorites` VALUES (4,1,4,'2024-12-16 02:16:41','默认收藏夹');
INSERT INTO `favorites` VALUES (5,1,5,'2024-12-16 02:17:04','默认收藏夹');

-- ----------------------------
-- Records of files
-- ----------------------------

DELETE FROM files;

INSERT INTO `files` VALUES (1, '1.png', 0x2CBDFB377C60CC77EB810E7134BCD01E1C3FC1ACD3EAA9BE815863B997428D36, 1, '2024-12-19 13:14:24');
INSERT INTO `files` VALUES (2, '2.png', 0xDAAA509C33DC6CA956ED3DD802AD0C6D851351D0D6C392BC81E1B4A56D1B2707, 1, '2024-12-19 13:14:34');
INSERT INTO `files` VALUES (3, '3.png', 0x2677AB5DEFC59056446EBDD43006F2DE27603660EFCE7BCAF01F2001C5CFBFD4, 1, '2024-12-19 13:15:55');
INSERT INTO `files` VALUES (4, '4.png', 0xAEB7BC6B3353AC89CE865DCAA4EFA8B2B5948DE57223537C8FBD1B8C295CC5A7, 1, '2024-12-19 13:16:12');
INSERT INTO `files` VALUES (5, '5.png', 0x5F36045C9D8C2FDF15B22DBB0F3ABAD11C1CCD3D8D1291587D54DC101BEE07B3, 1, '2024-12-19 13:16:17');
INSERT INTO `files` VALUES (6, '6.png', 0x7F43D95221A94AABE063D8F5435E069ACC23840F8DE6D8899BA46B56B8F6168D, 1, '2024-12-19 13:16:23');
INSERT INTO `files` VALUES (7, '7.png', 0x716A3C19E9BB451AC30AB010A509C8873D447E2AD5A940A53B42CBA1F61D0C1F, 1, '2024-12-19 13:16:31');
INSERT INTO `files` VALUES (8, '8.png', 0x47D6B5BA427F1DDB3EB885E406598E2BA568F2E1F53CF07DC86948DAD07C1C8C, 1, '2024-12-19 13:16:36');
INSERT INTO `files` VALUES (9, '9.png', 0xF770232CEF955B7BAAFB009B899D5775E991D7EB7EDA43D2D59B004CCCE473A3, 1, '2024-12-19 13:16:41');
INSERT INTO `files` VALUES (10, '10.png', 0x007045B2F2001DF842E0FD448587FA00CCCD04807BC26EEC7D01252C299EA137, 1, '2024-12-19 13:16:46');

-- ----------------------------
-- Records of config
-- ----------------------------

DELETE FROM config;

INSERT INTO `config`(`default_avatar_path`, `default_theme_id`) VALUES ('/assets/default_avatar.png', 1);

-- ----------------------------
-- Records of settings
-- ----------------------------

DELETE FROM settings;

INSERT INTO `settings` VALUES (1, 1, 'public');
INSERT INTO `settings` VALUES (2, NULL, 'protected');
INSERT INTO `settings` VALUES (11, NULL, 'public');

SET FOREIGN_KEY_CHECKS = 1;