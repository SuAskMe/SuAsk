-- MySQL dump 10.13  Distrib 8.0.18, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: suask
-- ------------------------------------------------------
-- Server version	8.0.18

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `answers`
--

DROP TABLE IF EXISTS `answers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `answers` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '回答ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `question_id` int(11) NOT NULL COMMENT '问题ID',
  `in_reply_to` int(11) DEFAULT NULL COMMENT '回复的回答ID，可为空',
  `contents` text CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '回答内容',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `upvotes` int(11) NOT NULL DEFAULT '0' COMMENT '点赞量',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_id` (`user_id`) USING BTREE,
  KEY `question_id` (`question_id`) USING BTREE,
  KEY `in_reply_to` (`in_reply_to`) USING BTREE,
  KEY `upvotes` (`upvotes` DESC) USING BTREE COMMENT '按点赞量降序索引',
  FULLTEXT KEY `contents` (`contents`) COMMENT '内容支持全文搜索，使用ngram parser以支持中文，默认token size为2',
  CONSTRAINT `answers_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `answers_ibfk_2` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `answers_ibfk_3` FOREIGN KEY (`in_reply_to`) REFERENCES `answers` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `answers`
--

LOCK TABLES `answers` WRITE;
/*!40000 ALTER TABLE `answers` DISABLE KEYS */;
/*!40000 ALTER TABLE `answers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `attachments`
--

DROP TABLE IF EXISTS `attachments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `attachments` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '附件ID',
  `question_id` int(11) DEFAULT NULL COMMENT '问题ID',
  `answer_id` int(11) DEFAULT NULL COMMENT '回答ID',
  `type` enum('picture') CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '附件类型（目前仅支持图片）',
  `file_id` int(11) NOT NULL COMMENT '文件ID',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `question_id` (`question_id`) USING BTREE,
  KEY `answer_id` (`answer_id`) USING BTREE,
  KEY `file_id` (`file_id`) USING BTREE,
  CONSTRAINT `attachments_ibfk_1` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `attachments_ibfk_2` FOREIGN KEY (`answer_id`) REFERENCES `answers` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `attachments_ibfk_3` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `attachments_chk_1` CHECK ((((`question_id` is not null) + (`answer_id` is not null)) = 1))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attachments`
--

LOCK TABLES `attachments` WRITE;
/*!40000 ALTER TABLE `attachments` DISABLE KEYS */;
/*!40000 ALTER TABLE `attachments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config`
--

DROP TABLE IF EXISTS `config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `config` (
  `id` bit(1) NOT NULL DEFAULT b'0' COMMENT '配置ID，限制为0',
  `default_avatar_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '默认头像文件路径',
  `default_theme_id` int(11) NOT NULL COMMENT '默认主题ID',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `default_theme_id` (`default_theme_id`) USING BTREE,
  CONSTRAINT `config_ibfk_1` FOREIGN KEY (`default_theme_id`) REFERENCES `themes` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `config_chk_1` CHECK ((`id` = 0))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config`
--

LOCK TABLES `config` WRITE;
/*!40000 ALTER TABLE `config` DISABLE KEYS */;
/*!40000 ALTER TABLE `config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `favorites`
--

DROP TABLE IF EXISTS `favorites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `favorites` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '收藏（置顶）ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `question_id` int(11) NOT NULL COMMENT '问题ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `package` varchar(50) COLLATE utf8mb4_zh_0900_as_cs NOT NULL DEFAULT '默认收藏夹' COMMENT '收藏夹',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `user_id` (`user_id`,`question_id`) USING BTREE COMMENT '每个用户收藏同个问题最多一次',
  KEY `question_id` (`question_id`) USING BTREE,
  CONSTRAINT `favorites_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `favorites_ibfk_2` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `favorites`
--

LOCK TABLES `favorites` WRITE;
/*!40000 ALTER TABLE `favorites` DISABLE KEYS */;
INSERT INTO `favorites` VALUES (1,1,12,'2024-12-18 13:25:35','默认收藏夹');
INSERT INTO `favorites` VALUES (2,1,13,'2024-12-15 09:36:32','默认收藏夹');
INSERT INTO `favorites` VALUES (3,1,14,'2024-12-15 08:10:24','默认收藏夹');
INSERT INTO `favorites` VALUES (4,1,15,'2024-12-16 02:16:41','默认收藏夹');
INSERT INTO `favorites` VALUES (5,1,16,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (6,1,17,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (7,1,18,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (8,1,19,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (9,1,20,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (10,1,21,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (11,1,22,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (12,1,23,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (13,1,24,'2024-12-16 02:17:04','默认收藏夹');
INSERT INTO `favorites` VALUES (14,1,25,'2024-12-16 07:15:58','默认收藏夹');
INSERT INTO `favorites` VALUES (15,1,26,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (16,1,27,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (17,1,28,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (18,1,29,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (19,1,30,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (20,1,31,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (21,1,32,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (22,1,33,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (23,1,34,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (24,1,35,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (25,1,36,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (26,1,37,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (27,1,38,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (28,1,39,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (29,1,40,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (30,1,41,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (31,1,42,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (32,1,43,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (33,1,44,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (34,1,45,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (35,1,46,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (36,1,47,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (37,1,48,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (38,1,49,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (39,1,50,'2024-12-18 10:55:59','默认收藏夹');
INSERT INTO `favorites` VALUES (40,1,51,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (41,1,52,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (42,1,53,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (43,1,54,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (44,1,55,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (45,1,56,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (46,1,57,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (47,1,58,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (48,1,59,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (49,1,60,'2024-12-18 12:19:24','默认收藏夹');
INSERT INTO `favorites` VALUES (50,1,61,'2024-12-18 13:23:19','默认收藏夹');
INSERT INTO `favorites` VALUES (51,1,62,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (52,1,63,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (53,1,64,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (54,1,65,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (55,1,66,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (56,1,67,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (57,1,68,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (58,1,69,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (59,1,70,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (60,1,71,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (61,1,72,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (62,1,73,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (63,1,74,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (64,1,75,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (65,1,76,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (66,1,77,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (67,1,78,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (68,1,79,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (69,1,80,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (70,1,81,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (71,1,82,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (72,1,83,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (73,1,84,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (74,1,85,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (75,1,86,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (76,1,87,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (77,1,88,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (78,1,89,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (79,1,90,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (80,1,91,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (81,1,92,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (82,1,93,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (83,1,94,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (84,1,95,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (85,1,96,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (86,1,97,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (87,1,98,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (88,1,99,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (89,1,100,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (90,1,101,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (91,1,102,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (92,1,103,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (93,1,104,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (94,1,105,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (95,1,106,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (96,1,107,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (97,1,108,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (98,1,109,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (99,1,110,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (100,1,111,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (101,1,112,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (102,1,113,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (103,1,114,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (104,1,115,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (105,1,116,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (106,1,117,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (107,1,118,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (108,1,119,'2024-12-18 13:25:08','默认收藏夹');
INSERT INTO `favorites` VALUES (109,1,120,'2024-12-18 13:25:08','默认收藏夹');
/*!40000 ALTER TABLE `favorites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `files` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '文件名，不得包含非法字符例如斜杠',
  `hash` binary(32) NOT NULL COMMENT '文件哈希，算法暂定为BLAKE2b',
  `uploader_id` int(11) DEFAULT NULL COMMENT '上传者用户ID',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '文件上传时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `uploader_id` (`uploader_id`) USING BTREE,
  CONSTRAINT `files_ibfk_1` FOREIGN KEY (`uploader_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `files`
--

LOCK TABLES `files` WRITE;
/*!40000 ALTER TABLE `files` DISABLE KEYS */;
/*!40000 ALTER TABLE `files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '提醒ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `question_id` int(11) NOT NULL COMMENT '问题ID',
  `type` enum('new_question','new_reply') CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '提醒类型（新提问或新回复）',
  `created_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `user_id_2` (`user_id`,`question_id`) USING BTREE COMMENT '每个用户只能收到关于同一个问题的一条提醒',
  KEY `user_id` (`user_id`) USING BTREE,
  KEY `question_id` (`question_id`) USING BTREE,
  CONSTRAINT `notifications_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `notifications_ibfk_2` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questions`
--

DROP TABLE IF EXISTS `questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questions` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '问题ID',
  `src_user_id` int(11) NOT NULL COMMENT '发起提问的用户ID',
  `dst_user_id` int(11) DEFAULT NULL COMMENT '被提问的用户ID，为空时问大家，不为空时问教师',
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '问题标题',
  `contents` text CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '问题内容',
  `is_private` bit(1) NOT NULL COMMENT '是否私密提问，仅在问教师时可为是',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `views` int(11) NOT NULL DEFAULT '0' COMMENT '浏览量',
  `upvotes` int(11) NOT NULL DEFAULT '0' COMMENT '点赞量',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `src_user_id` (`src_user_id`) USING BTREE,
  KEY `dst_user_id` (`dst_user_id`) USING BTREE,
  KEY `views` (`views` DESC) USING BTREE COMMENT '按浏览量降序索引',
  KEY `upvotes` (`upvotes` DESC) USING BTREE COMMENT '按点赞量降序索引',
  FULLTEXT KEY `contents` (`contents`) COMMENT '内容支持全文搜索，使用ngram parser以支持中文，默认token size为2' /*!50100 WITH PARSER `ngram` */ ,
  CONSTRAINT `questions_ibfk_1` FOREIGN KEY (`src_user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `questions_ibfk_2` FOREIGN KEY (`dst_user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `questions_chk_1` CHECK (((`dst_user_id` is not null) or (`is_private` = 0)))
) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questions`
--

LOCK TABLES `questions` WRITE;
/*!40000 ALTER TABLE `questions` DISABLE KEYS */;
INSERT INTO `questions` VALUES (1,1,NULL,'你认为是什么导致人们对旅行的热爱？','不同的文化、美丽的风景，以及逃离日常生活的压力。',_binary '\0','2016-04-04 14:51:12',540,652);
INSERT INTO `questions` VALUES (2,1,NULL,'如果可以选择任何一个历史人物共进晚餐，你会选择谁？','达·芬奇，想听听他的创意和思想。',_binary '\0','2009-02-07 03:34:10',169,195);
INSERT INTO `questions` VALUES (3,1,NULL,'你最喜欢的书是什么，为什么？','《1984》，因为它对权力和自由的深刻思考令人警醒。',_binary '\0','2014-01-08 18:45:44',500,885);
INSERT INTO `questions` VALUES (4,1,NULL,'有哪些你认为值得一试的爱好或活动？','尝试陶艺或摄影，可以提升创造力。',_binary '\0','2001-06-22 13:21:38',852,427);
INSERT INTO `questions` VALUES (5,1,NULL,'你认为科技会如何改变未来的工作方式？','更多的远程工作、自动化和人工智能将改变传统的工作模式。',_binary '\0','2020-02-08 14:57:07',502,99);
INSERT INTO `questions` VALUES (6,1,NULL,'你最难忘的假期是去哪里，为什么？','去了一次日本，丰富的文化和美味的食物给我留下深刻印象。',_binary '\0','2008-03-26 17:43:25',878,833);
INSERT INTO `questions` VALUES (7,1,NULL,'如果能学会一种乐器，你希望学习什么？','钢琴，因为它的音色丰富多样。',_binary '\0','2012-01-02 13:50:47',216,255);
INSERT INTO `questions` VALUES (8,1,NULL,'你对环保有什么建议或想法？','大力推广可再生能源和减少塑料使用是关键。',_binary '\0','2006-06-17 13:53:20',756,811);
INSERT INTO `questions` VALUES (9,1,NULL,'在你的生活中，有没有什么事情让你感到特别感恩？','感恩家人和朋友的支持，他们在我困难时给予了我力量。',_binary '\0','2002-12-07 01:39:24',476,398);
INSERT INTO `questions` VALUES (10,1,NULL,'你觉得友情和爱情之间最大的区别是什么？','友情更倾向于无条件的支持，而爱情通常伴随着深度的情感和浪漫。',_binary '\0','2011-02-13 00:29:18',783,273);
INSERT INTO `questions` VALUES (11,1,NULL,'你认为幸福的定义是什么？','幸福是内心的满足和平静，能与人分享美好的时刻。',_binary '\0','2012-07-28 05:38:14',412,160);
INSERT INTO `questions` VALUES (12,1,NULL,'你最喜欢的电影是什么，有什么特别的理由吗？','《肖申克的救赎》，因为它传递了希望和坚持的力量。',_binary '\0','2019-11-16 10:18:27',108,20);
INSERT INTO `questions` VALUES (13,1,NULL,'如果有一天可以拥有超能力，你希望拥有哪种？','弗雷迪·默丘里，想听听他的音乐创作理念。',_binary '\0','2013-03-13 16:21:04',911,571);
INSERT INTO `questions` VALUES (14,1,NULL,'你如何看待素食主义的影响？','对环境和健康有积极影响，但选择要灵活，适合个人情况。',_binary '\0','2014-03-23 06:19:31',738,146);
INSERT INTO `questions` VALUES (15,1,NULL,'对于你的人生目标，有哪些计划或梦想？','希望能在职业和个人生活中找到平衡，同时实现自我成长。',_binary '\0','2024-08-18 20:44:58',558,687);
INSERT INTO `questions` VALUES (16,1,NULL,'你认为心理健康的重要性体现在哪些方面？','心理健康影响情绪、决策和人际关系，良好的心理状态是全面健康的重要部分。',_binary '\0','2011-11-27 03:08:41',402,712);
INSERT INTO `questions` VALUES (17,1,NULL,'如果复活任何一位已故的音乐家或艺术家，你会选择谁？','学习到新技能和帮助他人成功，让我感到自豪。',_binary '\0','2012-07-03 21:44:12',728,768);
INSERT INTO `questions` VALUES (18,1,NULL,'有哪些事情是你希望未来能改善的？','持续学习和自我反思是关键，可以帮助认清自己的目标和进步。',_binary '\0','2022-02-09 10:33:18',66,664);
INSERT INTO `questions` VALUES (19,1,NULL,'你觉得社交媒体对人际关系的影响是什么？','帮助保持联系，但也会导致表面化的关系和误解。',_binary '\0','2015-08-06 05:36:43',372,675);
INSERT INTO `questions` VALUES (20,1,NULL,'你对未来的科技有什么期待？','期待人工智能进一步改善人类生活，如医疗、交通等行业的发展。',_binary '\0','2013-01-13 17:43:06',312,243);
INSERT INTO `questions` VALUES (21,1,NULL,'你有哪个特别，不为人知的才能或爱好吗？','喜欢写作，偶尔会写短篇故事。',_binary '\0','2005-07-03 11:53:54',465,256);
INSERT INTO `questions` VALUES (22,1,NULL,'在你的生活中，有哪些人对你产生过重大影响？','我的父母和一些老师，他们的支持和教导对我的发展至关重要。',_binary '\0','2018-10-22 18:46:41',469,935);
INSERT INTO `questions` VALUES (23,1,NULL,'你最喜欢的季节是什么，又是什么原因？','秋天，天气宜人，树叶变色也很美。',_binary '\0','2005-02-18 13:35:29',867,154);
INSERT INTO `questions` VALUES (24,1,NULL,'有哪些书籍或电影改变了你的观点或价值观？','《小王子》让我思考人际关系的深度和纯真。',_binary '\0','2008-04-22 02:01:55',488,500);
INSERT INTO `questions` VALUES (25,1,NULL,'你怎么看待教育系统的变革？','需要更关注学生的个性化发展和批判性思维的培养。',_binary '\0','2000-12-05 06:09:54',35,156);
INSERT INTO `questions` VALUES (26,1,NULL,'如果未来可以住在任何地方，你希望选择哪里？','希望住在一个风景优美、生活节奏适中的城市，如阿姆斯特丹。',_binary '\0','2023-12-08 04:48:53',898,349);
INSERT INTO `questions` VALUES (27,496,NULL,'你认为自我提升的最有效方法是什么？','设定明确目标并坚持行动，持续学习和反思。',_binary '\0','2024-10-16 20:18:33',589,234);
INSERT INTO `questions` VALUES (28,89,NULL,'过去一年中，有哪些事情让你感到特别骄傲？','成功完成某个项目或帮助他人实现梦想，让我感到自豪。',_binary '\0','2002-02-13 14:53:04',36,218);
INSERT INTO `questions` VALUES (29,868,NULL,'你对未来的环保措施有什么看法？','需要更多的支持和行动，以减少污染并保护生态系统。',_binary '\0','2019-01-16 18:59:07',981,348);
INSERT INTO `questions` VALUES (30,518,NULL,'你认为人类未来面临的最大挑战是什么？','应对气候变化和资源短缺是人类未来最严峻的挑战。',_binary '\0','2023-03-06 15:32:49',749,825);
INSERT INTO `questions` VALUES (31,909,NULL,'如果你能选择一个超能力，你希望是什么，为什么？','能隐形，因为可以秘密探索。',_binary '\0','2015-10-28 15:01:41',934,870);
INSERT INTO `questions` VALUES (32,954,NULL,'你最喜欢的书是什么，它对你有什么影响？','《活着》，让我更珍惜生命。',_binary '\0','2004-06-02 22:06:32',523,236);
INSERT INTO `questions` VALUES (33,111,NULL,'如果你可以与历史上的任何人共进晚餐，你会选择谁？','爱因斯坦，想听他的见解。',_binary '\0','2011-06-20 07:52:55',197,400);
INSERT INTO `questions` VALUES (34,496,NULL,'描述一个你最难忘的旅行经历。','去日本赏樱花。',_binary '\0','2007-07-22 03:18:18',182,163);
INSERT INTO `questions` VALUES (35,89,NULL,'你认为人工智能在未来将如何改变我们的生活？','会使生活更便利。',_binary '\0','2023-12-11 01:15:56',809,577);
INSERT INTO `questions` VALUES (36,868,NULL,'你有没有一个人生座右铭？是什么？','\"活在当下\"。',_binary '\0','2021-11-18 01:49:00',978,169);
INSERT INTO `questions` VALUES (37,518,NULL,'在你心中，成功的定义是什么？','能实现自己的价值。',_binary '\0','2013-07-12 19:10:01',953,117);
INSERT INTO `questions` VALUES (38,909,NULL,'你觉得科技进步更多是帮助了人类，还是创造了新的问题？','带来了便利，但也带来了依赖。',_binary '\0','2011-12-15 06:48:33',879,374);
INSERT INTO `questions` VALUES (39,579,NULL,'如果你可以住在任何地方，你希望住在哪里，为什么？','纽约，充满活力和文化。',_binary '\0','2010-08-14 11:03:21',399,667);
INSERT INTO `questions` VALUES (40,954,NULL,'你最喜欢的季节是什么，是什么让你喜欢它？','秋天，喜欢它的色彩和气候。',_binary '\0','2007-01-03 09:18:51',939,809);
INSERT INTO `questions` VALUES (41,111,NULL,'有哪部电影或电视剧对你产生了深刻的影响？','《海上钢琴师》，让我思考人生选择。',_binary '\0','2021-06-01 18:09:04',494,248);
INSERT INTO `questions` VALUES (42,496,NULL,'如果时间旅行成为可能，你想去哪个时代，为什么？','维多利亚时代，想看历史的变化。',_binary '\0','2007-03-16 09:29:21',304,428);
INSERT INTO `questions` VALUES (43,307,NULL,'你在生活中遇到的最大挑战是什么，你是如何克服的？','克服焦虑，学会自我调节。',_binary '\0','2001-06-04 12:12:38',796,940);
INSERT INTO `questions` VALUES (44,62,NULL,'你觉得友谊最重要的品质是什么？','诚实和支持。',_binary '\0','2005-07-26 00:08:30',443,679);
INSERT INTO `questions` VALUES (45,223,NULL,'如果你可以学习任何技能而不需时间，你想学什么？','能说流利的多种语言。',_binary '\0','2020-10-19 16:28:02',494,237);
INSERT INTO `questions` VALUES (46,633,NULL,'你最喜欢的食物是什么，有什么特别的回忆吗？','披萨，小时候的美好回忆。',_binary '\0','2019-07-06 12:16:56',84,591);
INSERT INTO `questions` VALUES (47,799,NULL,'对你来说，家庭意味着什么？','给予支持和爱的地方。',_binary '\0','2002-04-11 05:47:27',439,495);
INSERT INTO `questions` VALUES (48,971,NULL,'你有没有梦寐以求的工作，是什么？','作家，能创造故事。',_binary '\0','2022-10-14 23:24:21',273,43);
INSERT INTO `questions` VALUES (49,44,NULL,'如果你能改变一件历史事件，你希望改变什么？','改变二战，避免无辜生命的牺牲。',_binary '\0','2000-07-25 14:08:27',763,758);
INSERT INTO `questions` VALUES (50,798,NULL,'你认为人们对幸福的定义是一个主观的概念吗？','是的，每个人的经历不同。',_binary '\0','2007-02-03 23:10:02',731,694);
INSERT INTO `questions` VALUES (51,846,NULL,'描述一下你的理想周末是怎样的。','放松、看电影、与朋友聚会。',_binary '\0','2016-09-15 10:01:20',84,195);
INSERT INTO `questions` VALUES (52,727,NULL,'你认为教育的未来会是什么样子？','更加个性化和在线化。',_binary '\0','2020-10-11 21:12:34',289,630);
INSERT INTO `questions` VALUES (53,611,NULL,'你有没有过任何让你转变观念的经历？','遇到挫折后重新审视目标。',_binary '\0','2019-03-09 04:11:41',161,669);
INSERT INTO `questions` VALUES (54,407,NULL,'你最喜欢的音乐类型是什么，为什么？','摇滚乐，充满激情。',_binary '\0','2000-10-16 08:38:52',551,618);
INSERT INTO `questions` VALUES (55,852,NULL,'有哪件事情是你一直想做但还没有实现的？','学习弹吉他。',_binary '\0','2009-09-23 17:48:14',413,174);
INSERT INTO `questions` VALUES (56,595,NULL,'你认为社交媒体对人际关系的影响是积极还是消极？','有积极和消极的影响。',_binary '\0','2001-01-27 03:20:36',719,694);
INSERT INTO `questions` VALUES (57,55,NULL,'你是否相信运气在生活中的作用？为什么？','相信，常常能碰到好运。',_binary '\0','2024-09-16 04:21:19',993,542);
INSERT INTO `questions` VALUES (58,852,NULL,'如果你能发明一个新的节日，你希望它是关于什么的？','感恩节，专门感谢身边的人。',_binary '\0','2009-07-04 19:54:52',462,685);
INSERT INTO `questions` VALUES (59,266,NULL,'哪些品质让一个人值得信任？','诚实和一贯性。',_binary '\0','2013-07-25 00:11:41',582,115);
INSERT INTO `questions` VALUES (60,225,NULL,'你觉得善良在现代社会中仍然重要吗？','是的，它让社会更温暖。',_binary '\0','2000-05-30 13:43:02',657,819);
INSERT INTO `questions` VALUES (61,751,NULL,'你最喜欢的假期是什么时候，为什么？','圣诞节，因为有家人团聚。',_binary '\0','2023-01-26 05:15:51',544,825);
INSERT INTO `questions` VALUES (62,780,NULL,'如果你可以访问任何一个国家，你会选择哪个？','日本，因为文化丰富。',_binary '\0','2015-12-08 09:25:45',805,498);
INSERT INTO `questions` VALUES (63,735,NULL,'对于你而言，什么是成功？','完成自己设定的目标。',_binary '\0','2012-08-30 06:54:01',731,338);
INSERT INTO `questions` VALUES (64,189,NULL,'你最喜欢的书是什么，为什么？','《活着》，因为反映了人性的坚韧。',_binary '\0','2013-11-17 12:22:28',615,675);
INSERT INTO `questions` VALUES (65,430,NULL,'有什么事情是你一直想尝试但还没做的？','学习一门外语。',_binary '\0','2005-09-05 18:28:05',179,181);
INSERT INTO `questions` VALUES (66,378,NULL,'描述一次改变你人生的经历。','大学毕业典礼。',_binary '\0','2005-01-20 15:24:27',225,144);
INSERT INTO `questions` VALUES (67,109,NULL,'如果你可以选择任何一种超能力，你会选择什么？','隐形。',_binary '\0','2016-10-23 22:17:03',971,755);
INSERT INTO `questions` VALUES (68,609,NULL,'对你来说，家庭的意义是什么？','互相支持和爱护。',_binary '\0','2007-01-26 05:52:59',860,325);
INSERT INTO `questions` VALUES (69,212,NULL,'你最喜欢的音乐类型是什么，为什么？','摇滚乐，因为能激发能量。',_binary '\0','2004-07-16 04:58:59',842,694);
INSERT INTO `questions` VALUES (70,451,NULL,'有没有一部电影让你感动得流泪？是哪一部？','《海上钢琴师》。',_binary '\0','2012-06-18 04:39:44',274,135);
INSERT INTO `questions` VALUES (71,759,NULL,'你如何面对压力？','深呼吸和运动。',_binary '\0','2023-10-08 00:06:41',124,572);
INSERT INTO `questions` VALUES (72,639,NULL,'谈谈一次让你感到特别自豪的成就。','完成马拉松比赛。',_binary '\0','2004-05-04 15:05:46',348,752);
INSERT INTO `questions` VALUES (73,162,NULL,'你觉得科技在生活中发挥了怎样的作用？','改善我们的生活质量。',_binary '\0','2003-11-17 16:48:12',764,590);
INSERT INTO `questions` VALUES (74,411,NULL,'什么样的食物让你觉得幸福？','巧克力蛋糕。',_binary '\0','2015-11-21 17:26:43',535,693);
INSERT INTO `questions` VALUES (75,102,NULL,'如果可以回到过去，你最想改变什么？','做得更好一些。',_binary '\0','2015-01-02 02:15:21',302,894);
INSERT INTO `questions` VALUES (76,667,NULL,'你认为友情有什么重要性？','朋友能提供情感支持。',_binary '\0','2019-05-05 11:22:00',593,803);
INSERT INTO `questions` VALUES (77,584,NULL,'你有没有什么特别的爱好？是什么？','17.摄影。',_binary '\0','2019-05-10 01:26:28',702,74);
INSERT INTO `questions` VALUES (78,845,NULL,'你最喜欢的季节是什么，为什么？','春天，因为万物复苏。',_binary '\0','2005-04-25 21:41:25',791,262);
INSERT INTO `questions` VALUES (79,335,NULL,'谈谈你心目中的理想工作。','设计师/艺术家。',_binary '\0','2023-08-03 14:46:37',520,470);
INSERT INTO `questions` VALUES (80,364,NULL,'有没有一位名人对你产生过影响？是谁？','马丁·路德·金。',_binary '\0','2008-10-24 23:24:06',393,298);
INSERT INTO `questions` VALUES (81,386,NULL,'你对当前的社会发展有什么看法？','有正面影响，也有负面影响。',_binary '\0','2021-01-12 01:17:47',632,816);
INSERT INTO `questions` VALUES (82,664,NULL,'你如何定义幸福？','能够内心平和。',_binary '\0','2024-09-19 12:45:31',276,880);
INSERT INTO `questions` VALUES (83,87,NULL,'描述一次你最难忘的旅行经历。','去巴黎的旅行。',_binary '\0','2019-02-28 03:21:36',999,358);
INSERT INTO `questions` VALUES (84,322,NULL,'你喜欢通过什么方式放松自己？','听音乐和散步。',_binary '\0','2009-01-11 04:18:37',32,701);
INSERT INTO `questions` VALUES (85,548,NULL,'有什么事情是你希望了解的，但没机会学习？','心理学。',_binary '\0','2020-05-11 19:27:54',2,602);
INSERT INTO `questions` VALUES (86,528,NULL,'如果你能与任何人共进晚餐，你会选择谁？','爱因斯坦。',_binary '\0','2001-11-01 09:17:58',578,535);
INSERT INTO `questions` VALUES (87,358,NULL,'你觉得自己在生活中最大的挑战是什么？','工作与生活的平衡。',_binary '\0','2001-08-23 02:05:09',408,709);
INSERT INTO `questions` VALUES (88,356,NULL,'你会如何向别人描述你的性格？','开朗和乐于助人。',_binary '\0','2019-10-11 20:11:42',113,348);
INSERT INTO `questions` VALUES (89,245,NULL,'如果有人要写一部关于你的书，你希望它的主题是什么？','追求梦想与勇气。',_binary '\0','2001-08-22 19:11:01',439,513);
INSERT INTO `questions` VALUES (90,2,NULL,'描述一下你最喜欢的天气和活动。','晴天，适合户外活动。',_binary '\0','2012-02-24 04:10:49',936,515);
INSERT INTO `questions` VALUES (91,554,NULL,'你觉得教育在个人发展中扮演什么角色？','提升个体的思维和能力。',_binary '\0','2020-08-10 19:35:45',843,313);
INSERT INTO `questions` VALUES (92,322,NULL,'如果你可以改变世界上的一件事，那会是什么？','消除饥饿与贫困。',_binary '\0','2021-12-14 16:15:43',154,71);
INSERT INTO `questions` VALUES (93,625,NULL,'你有没有一个特别的梦想？是什么？','开一家自己的咖啡馆。',_binary '\0','2020-03-07 16:05:36',992,303);
INSERT INTO `questions` VALUES (94,215,NULL,'你认为成年人应该有怎样的责任？','关心社会与他人。',_binary '\0','2008-11-17 07:35:08',781,551);
INSERT INTO `questions` VALUES (95,928,NULL,'描述一下理想中的居住环境。','安静、和平、靠近自然。',_binary '\0','2023-06-25 22:04:48',467,397);
INSERT INTO `questions` VALUES (96,412,NULL,'有没有一本书或一部电影让你改变了看法？是什么？','《百年孤独》，让我对时间有新理解。',_binary '\0','2019-11-20 09:41:49',116,114);
INSERT INTO `questions` VALUES (97,108,NULL,'你认为与朋友的关系如何影响你的生活？','影响很大，能给生活带来乐趣。',_binary '\0','2004-10-15 23:38:22',43,875);
INSERT INTO `questions` VALUES (98,28,NULL,'什么样的事情能激励到你？','看到他人的努力和决心。',_binary '\0','2002-11-18 01:07:58',289,850);
INSERT INTO `questions` VALUES (99,715,NULL,'你在生活中最感激的事情是什么？','有爱我的家人和朋友。',_binary '\0','2003-08-27 01:31:37',439,770);
INSERT INTO `questions` VALUES (100,603,NULL,'你最喜欢的运动是什么，为什么？','足球，因为可以锻炼体力。',_binary '\0','2002-05-04 15:42:32',944,296);
INSERT INTO `questions` VALUES (101,487,NULL,'如果你可以学会任何一种乐器，你想学什么？','吉他。',_binary '\0','2018-10-20 13:45:49',304,480);
INSERT INTO `questions` VALUES (102,656,NULL,'你觉得言语的重要性在哪里？','交流和表达思想。',_binary '\0','2014-04-07 23:47:53',196,693);
INSERT INTO `questions` VALUES (103,23,NULL,'你最尊敬的人是谁，为什么？','我的父母，因为他们支持我。',_binary '\0','2021-10-23 02:17:11',614,901);
INSERT INTO `questions` VALUES (104,381,NULL,'对于未来的你，有什么期望？','希望能有更好的事业发展。',_binary '\0','2015-11-08 20:13:32',116,24);
INSERT INTO `questions` VALUES (105,872,NULL,'你有最喜欢的格言或座右铭吗？','“活在当下”。',_binary '\0','2022-08-30 09:01:15',781,937);
INSERT INTO `questions` VALUES (106,764,NULL,'如何让自己保持积极的心态？','每天记录三件好事。',_binary '\0','2021-06-14 02:27:41',974,148);
INSERT INTO `questions` VALUES (107,170,NULL,'你认为未来科技的发展会带来什么变化？','人工智能将更加普及和便捷。',_binary '\0','2010-05-10 11:50:21',437,902);
INSERT INTO `questions` VALUES (108,981,NULL,'你有什么保持健康的秘诀？','均衡饮食和规律锻炼。',_binary '\0','2019-04-30 15:18:45',321,277);
INSERT INTO `questions` VALUES (109,767,NULL,'你对社交媒体的看法是什么？','能连接人，但也需要适度使用。',_binary '\0','2014-03-07 05:55:11',613,87);
INSERT INTO `questions` VALUES (110,438,NULL,'描述一个你希望重温的美好时刻。','毕业旅行的欢乐时光。',_binary '\0','2002-02-22 09:13:37',382,364);
INSERT INTO `questions` VALUES (111,730,NULL,'你认为爱是什么样的感觉？','深切的关爱和理解。',_binary '\0','2015-05-28 00:32:02',461,756);
INSERT INTO `questions` VALUES (112,889,NULL,'有什么事情是你一直想学，但没有机会？','学习弹吉他。',_binary '\0','2007-08-10 08:55:10',222,584);
INSERT INTO `questions` VALUES (113,587,NULL,'你最崇拜的历史人物是谁？','玛丽·居里。',_binary '\0','2024-11-26 20:39:14',589,399);
INSERT INTO `questions` VALUES (114,445,NULL,'如果能选择任何一项技能，你希望学会什么？','外语。',_binary '\0','2013-11-27 17:30:25',89,52);
INSERT INTO `questions` VALUES (115,652,NULL,'哪个地方是你心中理想的度假胜地？','巴哈马。',_binary '\0','2013-03-31 16:43:43',941,787);
INSERT INTO `questions` VALUES (116,998,NULL,'你如何看待失败？','学习和成长的机会。',_binary '\0','2019-03-21 09:13:21',271,733);
INSERT INTO `questions` VALUES (117,11,NULL,'你相信命运吗？为什么？','相信，人生受多种因素影响。',_binary '\0','2004-08-20 06:48:12',95,376);
INSERT INTO `questions` VALUES (118,982,NULL,'有什么事情能让你感到非常愤怒？','不公平现象。',_binary '\0','2009-04-26 13:56:18',937,58);
INSERT INTO `questions` VALUES (119,324,NULL,'你记忆中最早的一次冒险是什么？','第一次露营。',_binary '\0','2000-08-24 03:05:15',956,630);
INSERT INTO `questions` VALUES (120,836,NULL,'你觉得现代社会最大的挑战是什么？','环境变化和社会极端化。',_binary '\0','2012-05-21 02:17:23',872,405);
INSERT INTO `questions` VALUES (121,386,NULL,'描述一下最有趣的工作经历。','在咖啡馆担任服务员时的交流。',_binary '\0','2019-07-20 06:50:07',992,156);
INSERT INTO `questions` VALUES (122,87,NULL,'你有什么信仰或原则指引你的生活？','坚持诚实与善良。',_binary '\0','2017-04-30 03:27:59',744,810);
INSERT INTO `questions` VALUES (123,632,NULL,'如果你能重温一天，你会选择哪一天？','高中毕业典礼那天。',_binary '\0','2000-02-23 06:32:47',13,69);
INSERT INTO `questions` VALUES (124,590,NULL,'你最喜欢的乐器是什么，为什么？','钢琴，能表达情感。',_binary '\0','2024-10-08 03:00:00',369,700);
INSERT INTO `questions` VALUES (125,84,NULL,'如果可以与过去的自己面对面，你想说什么？','别怕失败，珍惜现在。',_binary '\0','2010-11-15 07:17:29',51,961);
INSERT INTO `questions` VALUES (126,875,NULL,'你对未来的职业有什么计划或想法？','持续学习，寻求挑战。',_binary '\0','2001-04-19 10:10:40',835,900);
INSERT INTO `questions` VALUES (127,160,NULL,'有哪种动物是你特别喜欢的，为什么？','猫，性格独立可爱。',_binary '\0','2008-06-05 15:30:24',753,604);
INSERT INTO `questions` VALUES (128,856,NULL,'描述一下你心中的完美周末。','与朋友一起户外活动和放松。',_binary '\0','2019-03-23 15:09:13',330,930);
INSERT INTO `questions` VALUES (129,447,NULL,'你最喜欢的运动是什么，做得好吗？','游泳，享受水中运动。',_binary '\0','2024-01-27 03:58:15',791,237);
INSERT INTO `questions` VALUES (130,242,NULL,'如果你能改变自己的一个习惯，是什么？','想早起。',_binary '\0','2000-08-05 00:23:47',225,601);
INSERT INTO `questions` VALUES (131,205,NULL,'你觉得理想的工作环境是什么样的？','开放、支持性强的环境。',_binary '\0','2003-09-23 19:23:10',413,908);
INSERT INTO `questions` VALUES (132,370,NULL,'你有什么特别的才能或天赋？','艺术创作。',_binary '\0','2019-08-13 20:23:40',324,493);
INSERT INTO `questions` VALUES (133,623,NULL,'你认为家庭的重要性在哪里？','是生活的基础，提供安全感。',_binary '\0','2023-08-26 16:42:03',286,739);
INSERT INTO `questions` VALUES (134,505,NULL,'有没有一本书对你影响很大？是哪一本？','《平凡的世界》。',_binary '\0','2016-06-23 08:22:23',944,219);
INSERT INTO `questions` VALUES (135,893,NULL,'描述一下你最喜欢的传统节日。','中秋，团圆的象征。',_binary '\0','2016-06-18 10:57:07',690,359);
INSERT INTO `questions` VALUES (136,11,NULL,'你如何处理人际关系中的冲突？','坦诚沟通，寻找共识。',_binary '\0','2000-10-27 03:01:46',185,811);
INSERT INTO `questions` VALUES (137,717,NULL,'最让你感动的一件小事是什么？','别人送来的暖心问候。',_binary '\0','2000-10-06 12:07:13',18,175);
INSERT INTO `questions` VALUES (138,130,NULL,'你对成功的定义是什么？','实现自我价值。',_binary '\0','2021-03-28 11:39:19',785,394);
INSERT INTO `questions` VALUES (139,373,NULL,'你觉得自己最大的优点是什么？','关心他人。',_binary '\0','2023-02-06 00:19:55',16,490);
INSERT INTO `questions` VALUES (140,226,NULL,'如果有机会参与任何历史事件，你希望参与何种事件？','改变社会历史进程的时刻。',_binary '\0','2018-03-23 07:34:58',90,889);
INSERT INTO `questions` VALUES (141,652,NULL,'你认为快乐是由什么决定的？','自我认知与人际关系。',_binary '\0','2008-03-16 20:11:36',235,192);
INSERT INTO `questions` VALUES (142,752,NULL,'有什么电影你反复观看？','《肖申克的救赎》。',_binary '\0','2018-12-21 04:52:45',653,934);
INSERT INTO `questions` VALUES (143,118,NULL,'你有何种生活目标或梦想？','成为一名作家。',_binary '\0','2019-12-09 07:02:55',811,899);
INSERT INTO `questions` VALUES (144,999,NULL,'你最尊敬的职业是什么？','医生，治愈他人。',_binary '\0','2010-05-09 01:24:04',703,600);
INSERT INTO `questions` VALUES (145,313,NULL,'你会如何定义友谊？','彼此信任和理解。',_binary '\0','2008-02-29 19:59:38',348,816);
INSERT INTO `questions` VALUES (146,572,NULL,'你认为如何保持身心健康？','定期锻炼和冥想。',_binary '\0','2003-06-11 00:13:23',605,359);
INSERT INTO `questions` VALUES (147,5,NULL,'你觉得哪一项技术对生活影响最大？','互联网。',_binary '\0','2002-01-08 04:21:40',569,984);
INSERT INTO `questions` VALUES (148,327,NULL,'有没有特别的地方让你感到归属感？','家。',_binary '\0','2000-12-04 20:12:39',15,197);
INSERT INTO `questions` VALUES (149,833,NULL,'你在生活中最大的不安是什么？','感情的迷茫。',_binary '\0','2010-04-15 13:24:15',530,593);
INSERT INTO `questions` VALUES (150,531,NULL,'有什么一见钟情的事情或地方吗？','一见钟情的书籍。',_binary '\0','2023-12-12 05:35:06',737,144);
INSERT INTO `questions` VALUES (151,958,NULL,'你觉得怎样才能提高自己的创造力？','多读书、多尝试新事物。',_binary '\0','2012-04-15 00:13:16',165,747);
INSERT INTO `questions` VALUES (152,345,NULL,'你最喜欢的名言是什么，为什么？','“活在当下”，珍惜每一刻。',_binary '\0','2013-11-03 00:07:45',587,976);
INSERT INTO `questions` VALUES (153,379,NULL,'你如何看待自我反思的重要性？','重要，有助于成长和改变。',_binary '\0','2003-12-12 14:31:37',680,851);
INSERT INTO `questions` VALUES (154,953,NULL,'你有没有一个希望实现的愿望？','旅行环游世界。',_binary '\0','2005-12-14 03:26:49',746,273);
INSERT INTO `questions` VALUES (155,227,NULL,'对你来说，什么是最有价值的品质？','诚信。',_binary '\0','2012-01-16 08:09:21',16,808);
INSERT INTO `questions` VALUES (156,612,NULL,'有没有一首歌对你有特别的意义？','《你是我的荣光》，表达情感共鸣。',_binary '\0','2008-07-17 23:43:18',485,744);
INSERT INTO `questions` VALUES (157,386,NULL,'你希望在五年后自己变成什么样？','希望更自信和成熟。',_binary '\0','2009-08-28 21:59:42',694,80);
INSERT INTO `questions` VALUES (158,20,NULL,'什么样的艺术形式让你感到最有吸引力？','电影与音乐。',_binary '\0','2007-06-02 06:08:45',861,42);
INSERT INTO `questions` VALUES (159,548,NULL,'你认为幽默在生活中扮演什么角色？','帮助人们减轻压力。',_binary '\0','2009-10-25 19:46:15',905,213);
INSERT INTO `questions` VALUES (160,405,NULL,'如果能够选择一种生活方式，你想要怎样的生活？','自由与简单的生活。',_binary '\0','2006-05-05 11:51:17',126,966);
INSERT INTO `questions` VALUES (161,447,NULL,'描述一下你最喜欢的城市及其原因。','上海，因为融合了传统和现代。',_binary '\0','2014-07-01 19:23:56',384,117);
INSERT INTO `questions` VALUES (162,840,NULL,'你认为自己最大的成就是什么？','完成我的学位。',_binary '\0','2019-12-07 01:08:04',528,321);
INSERT INTO `questions` VALUES (163,417,NULL,'有没有一个瞬间改变了你的生活？是什么？','外出交换学习的经历。',_binary '\0','2019-03-11 23:07:41',595,931);
INSERT INTO `questions` VALUES (164,824,NULL,'你最喜欢的食物是什么？','披萨。',_binary '\0','2019-02-21 13:46:37',479,467);
INSERT INTO `questions` VALUES (165,684,NULL,'你如何在逆境中保持积极心态？','寻找正面影响和保持乐观。',_binary '\0','2008-02-04 17:01:58',640,405);
INSERT INTO `questions` VALUES (166,177,NULL,'如果可以选择一种职业，你希望从事什么？','作家。',_binary '\0','2016-12-23 19:08:03',803,50);
INSERT INTO `questions` VALUES (167,447,NULL,'你喜欢的课外活动是什么？','阅读和绘画。',_binary '\0','2021-07-25 23:56:26',684,788);
INSERT INTO `questions` VALUES (168,466,NULL,'在你的生活中，谁是你最大的支持者？','我的父母。',_binary '\0','2006-02-15 14:50:07',159,882);
INSERT INTO `questions` VALUES (169,963,NULL,'你向往的未来生活是怎样的？','和平、幸福、充实。',_binary '\0','2024-01-28 18:34:03',255,162);
INSERT INTO `questions` VALUES (170,279,NULL,'描述一下你的理想家园。','一个靠近大海的小房子。',_binary '\0','2020-02-26 15:57:52',681,566);
INSERT INTO `questions` VALUES (171,321,NULL,'如果可以旅行到任何一个历史时期，你想去哪个？','文艺复兴时期。',_binary '\0','2017-01-29 08:59:56',468,808);
INSERT INTO `questions` VALUES (172,563,NULL,'你每天的晨间例行活动是什么？','早餐和锻炼。',_binary '\0','2012-06-23 16:04:31',228,706);
INSERT INTO `questions` VALUES (173,202,NULL,'你如何看待环境保护的重要性？','非常重要，需要全社会的参与。',_binary '\0','2004-09-14 08:10:27',474,858);
INSERT INTO `questions` VALUES (174,451,NULL,'有什么事情是你不愿意妥协的？','个人信念和价值观。',_binary '\0','2012-08-25 08:26:12',151,204);
INSERT INTO `questions` VALUES (175,929,NULL,'你觉得技术对人际关系的影响是什么？','促进了联系，但也增加了孤独感。',_binary '\0','2003-12-11 01:20:26',119,355);
INSERT INTO `questions` VALUES (176,297,NULL,'你对快乐有怎样的理解？','感觉内心平静、满足。',_binary '\0','2013-11-06 07:25:33',827,270);
INSERT INTO `questions` VALUES (177,406,NULL,'你有没有特别喜欢的历史事件或时期？','二战历史。',_binary '\0','2010-10-29 15:00:24',682,58);
INSERT INTO `questions` VALUES (178,230,NULL,'你有怎样的读书习惯？','每天晚上读一小时。',_binary '\0','2018-07-29 01:21:00',98,319);
INSERT INTO `questions` VALUES (179,11,NULL,'心理健康对你来说有多重要？','极其重要，影响生活质量。',_binary '\0','2011-01-04 22:31:55',743,967);
INSERT INTO `questions` VALUES (180,12,NULL,'你有没有一种特别想去探访的地方？是哪里？','日本，想体验文化和美食。',_binary '\0','2005-10-08 10:47:49',101,936);
INSERT INTO `questions` VALUES (181,566,NULL,'你认为怎样才能实现人生目标？','制定清晰可行的计划。',_binary '\0','2007-05-06 20:25:33',877,946);
INSERT INTO `questions` VALUES (182,941,NULL,'描述一下你与朋友的关系。','密切、互相支持的朋友关系。',_binary '\0','2013-06-28 08:56:33',378,54);
INSERT INTO `questions` VALUES (183,274,NULL,'对你来说，最珍惜的回忆是什么？','家庭聚会的温馨时光。',_binary '\0','2020-01-14 22:21:02',107,21);
INSERT INTO `questions` VALUES (184,739,NULL,'你认为自我教育对个人发展的意义何在？','是不断提高自己和适应变化的途径。',_binary '\0','2005-02-17 01:59:26',108,378);
INSERT INTO `questions` VALUES (185,384,NULL,'有没有一项运动对你有特别的吸引力？是什么？','游泳，放松又锻炼。',_binary '\0','2006-09-24 21:03:14',869,139);
INSERT INTO `questions` VALUES (186,27,NULL,'你最喜欢的电影类型是什么？','科幻电影，刺激想象力。',_binary '\0','2011-06-09 15:51:20',884,293);
INSERT INTO `questions` VALUES (187,606,NULL,'有哪些情感让你容易共鸣？','动画和感人的故事。',_binary '\0','2013-07-19 01:30:44',337,418);
INSERT INTO `questions` VALUES (188,570,NULL,'您如何定义幸福的生活？','能够与自己和谐相处。',_binary '\0','2010-01-17 09:52:44',383,670);
INSERT INTO `questions` VALUES (189,538,NULL,'如果有机会，你会选择改变哪个社会现象？','消除贫困与不平等。',_binary '\0','2003-03-12 21:41:59',222,890);
INSERT INTO `questions` VALUES (190,890,NULL,'你最喜欢的儿童读物是什么？','《小王子》。',_binary '\0','2014-05-30 18:49:31',210,970);
INSERT INTO `questions` VALUES (191,169,NULL,'有没有一个名人是你心目中的偶像？是谁？','朴树。',_binary '\0','2003-02-08 06:52:26',459,581);
INSERT INTO `questions` VALUES (192,509,NULL,'你如何看待不同文化之间的交流？','促进理解与包容。',_binary '\0','2006-02-17 13:52:48',803,756);
INSERT INTO `questions` VALUES (193,514,NULL,'你希望未来的科技给生活带来什么变化？','提高生活便利性和智能化。',_binary '\0','2009-06-23 21:29:18',848,911);
INSERT INTO `questions` VALUES (194,863,NULL,'描述一下你心中的和平世界。','相互理解与尊重。',_binary '\0','2019-12-22 04:13:24',573,270);
INSERT INTO `questions` VALUES (195,294,NULL,'在朋友中，你被认为是什么样的人？','乐观、幽默。',_binary '\0','2005-10-01 22:33:44',474,574);
INSERT INTO `questions` VALUES (196,773,NULL,'你对爱情的看法是怎样的？','是深厚的信任与理解。',_binary '\0','2020-09-28 13:20:48',459,805);
INSERT INTO `questions` VALUES (197,187,NULL,'有什么方法帮助你提高效率？','设定优先级和制定清单。',_binary '\0','2007-06-11 19:51:53',768,352);
INSERT INTO `questions` VALUES (198,527,NULL,'你最喜欢的方式是如何表达自己的情感？','写日记。',_binary '\0','2001-04-18 19:44:50',659,947);
INSERT INTO `questions` VALUES (199,234,NULL,'你觉得人与人之间的信任重要吗？为什么？','非常重要，建立稳定关系。',_binary '\0','2005-03-04 10:32:12',5,767);
INSERT INTO `questions` VALUES (200,439,NULL,'如果可以从事任何一种志愿活动，你想参与什么？','保护动物权益。',_binary '\0','2007-03-31 05:36:59',524,220);
INSERT INTO `questions` VALUES (201,87,NULL,'描述一下你喜欢的季节及原因。','秋天，温暖而丰收的季节。',_binary '\0','2014-10-24 12:38:25',800,11);
INSERT INTO `questions` VALUES (202,436,NULL,'有没有一首歌让你觉得特别有力量？','《勇气》。',_binary '\0','2019-03-16 03:53:54',1000,484);
INSERT INTO `questions` VALUES (203,653,NULL,'你认为教育在生活中扮演什么角色？','提供知识和技能基础。',_binary '\0','2017-05-07 09:07:39',690,290);
INSERT INTO `questions` VALUES (204,172,NULL,'什么事情让你感到灵感涌现？','艺术作品和自然美景。',_binary '\0','2002-03-09 16:10:03',287,402);
INSERT INTO `questions` VALUES (205,962,NULL,'你怎样看待失败的经历？','学习机会，应对挑战。',_binary '\0','2017-04-02 09:53:24',647,694);
INSERT INTO `questions` VALUES (206,230,NULL,'有没有一个故事始终萦绕在你的脑海中？是什么？','《海上钢琴师》的孤独。',_binary '\0','2001-01-09 12:55:58',164,394);
INSERT INTO `questions` VALUES (207,781,NULL,'你对未来五年有哪些计划？','深入发展自己的事业。',_binary '\0','2011-09-10 15:47:06',148,318);
INSERT INTO `questions` VALUES (208,945,NULL,'有什么是你希望再尝试的？','大胆尝试新的技能。',_binary '\0','2009-03-26 07:13:33',811,610);
INSERT INTO `questions` VALUES (209,451,NULL,'你有没有自己设计过的项目或计划？','写小说。',_binary '\0','2015-10-10 07:01:17',388,96);
INSERT INTO `questions` VALUES (210,762,NULL,'你认为幸福是否可持续？为什么？','是，取决于内心的认知与选择。',_binary '\0','2020-05-22 10:02:15',518,933);
INSERT INTO `questions` VALUES (211,64,NULL,'你最喜欢的休闲活动是什么？','阅读小说。',_binary '\0','2019-01-03 20:23:49',928,959);
INSERT INTO `questions` VALUES (212,950,NULL,'描述一下你理想的职业是什么样的。','写作和传播思想的职业。',_binary '\0','2018-06-09 04:12:29',282,156);
INSERT INTO `questions` VALUES (213,997,NULL,'如果你能邀请三个人参加晚宴，他们会是谁？','爱因斯坦、史蒂芬·霍金、我已故的祖父。',_binary '\0','2018-08-20 19:02:57',557,873);
INSERT INTO `questions` VALUES (214,78,NULL,'你有遗憾的事情吗？是什么？','没能与远方的朋友保持联系。',_binary '\0','2016-10-19 16:27:57',501,955);
INSERT INTO `questions` VALUES (215,356,NULL,'你最喜欢的季节是什么？为什么？','秋天，天气宜人，色彩斑斓。',_binary '\0','2020-03-04 15:06:11',19,102);
INSERT INTO `questions` VALUES (216,288,NULL,'有什么事情让你十分感动？','看到孩子们的善良行为。',_binary '\0','2009-08-16 15:30:07',422,984);
INSERT INTO `questions` VALUES (217,417,NULL,'你希望如何改变自己的生活方式？','增加锻炼和健康饮食。',_binary '\0','2015-07-25 20:18:06',832,869);
INSERT INTO `questions` VALUES (218,928,NULL,'你认为一个人最重要的品质是什么？','诚实。',_binary '\0','2019-11-21 21:23:48',646,154);
INSERT INTO `questions` VALUES (219,26,NULL,'有什么书籍你希望大家都能读？','《活着》。',_binary '\0','2006-05-09 11:54:31',9,907);
INSERT INTO `questions` VALUES (220,512,NULL,'你最喜欢的回忆是什么？','家庭聚会的欢乐时光。',_binary '\0','2008-08-18 21:04:45',3,161);
INSERT INTO `questions` VALUES (221,240,NULL,'如何定义成功？','实现个人价值和目标。',_binary '\0','2017-12-15 10:00:49',532,915);
INSERT INTO `questions` VALUES (222,475,NULL,'你对未来的科技发展有什么期待？','我期待人工智能和清洁能源的进步。',_binary '\0','2011-08-26 15:37:19',544,297);
INSERT INTO `questions` VALUES (223,616,NULL,'描述一下你最喜欢的食物和做法。','意大利面，简单又美味。',_binary '\0','2000-04-02 17:35:08',114,838);
INSERT INTO `questions` VALUES (224,403,NULL,'有没有经历过让你感到不公的事情？','受到不公正对待的经历。',_binary '\0','2016-10-31 17:16:34',350,740);
INSERT INTO `questions` VALUES (225,122,NULL,'你如何处理压力和焦虑？','通过冥想和运动来缓解。',_binary '\0','2024-07-07 14:54:31',831,459);
INSERT INTO `questions` VALUES (226,303,NULL,'描述一下你特别喜欢的音乐类型。','爵士乐，优雅又放松。',_binary '\0','2013-07-19 17:34:23',794,766);
INSERT INTO `questions` VALUES (227,366,NULL,'你认为友谊应该建立在什么基础上？','互信与支持。',_binary '\0','2008-01-25 07:32:56',377,193);
INSERT INTO `questions` VALUES (228,206,NULL,'你心目中理想的假期是哪些活动？','放松、旅行和探索美食。',_binary '\0','2021-01-02 21:56:59',454,680);
INSERT INTO `questions` VALUES (229,238,NULL,'有什么人是你希望能够再次见到的？','青少年时期的好朋友。',_binary '\0','2018-05-19 18:52:46',661,308);
INSERT INTO `questions` VALUES (230,157,NULL,'你最崇拜的艺术作品是什么？','《蒙娜丽莎》。',_binary '\0','2017-09-05 17:45:16',473,725);
INSERT INTO `questions` VALUES (231,95,NULL,'你在生活中遵循的信条或座右铭是什么？','“活在当下”。',_binary '\0','2016-08-18 13:28:16',497,56);
INSERT INTO `questions` VALUES (232,292,NULL,'你觉得人生的意义是什么？','发现自我和帮助他人。',_binary '\0','2007-06-02 06:43:07',41,46);
INSERT INTO `questions` VALUES (233,551,NULL,'有没有一个地方是你特别想去却未能到达的？','北极，想看极光。',_binary '\0','2014-03-09 20:53:54',764,713);
INSERT INTO `questions` VALUES (234,137,NULL,'有什么事情是你希望能够遗忘的？','遇到的痛苦经历。',_binary '\0','2005-07-10 01:39:24',839,99);
INSERT INTO `questions` VALUES (235,952,NULL,'描述一下你最喜欢的电影情节。','《泰坦尼克号》中最后的告别。',_binary '\0','2007-08-24 04:02:00',989,487);
INSERT INTO `questions` VALUES (236,118,NULL,'你认为家庭在生活中的作用是什么？','提供支持和归属感。',_binary '\0','2021-09-15 09:03:52',319,417);
INSERT INTO `questions` VALUES (237,839,NULL,'你对社交媒体的看法是什么？','有积极和消极两面。',_binary '\0','2015-08-06 05:15:26',629,684);
INSERT INTO `questions` VALUES (238,735,NULL,'有什么事情让你感到害怕？','失去亲人。',_binary '\0','2012-08-07 16:07:59',286,862);
INSERT INTO `questions` VALUES (239,514,NULL,'你如何理解“爱”？','深厚的相互关心与支持。',_binary '\0','2023-06-01 01:01:52',167,911);
INSERT INTO `questions` VALUES (240,103,NULL,'怎样才能保持一个健康的生活方式？','均衡饮食与规律运动。',_binary '\0','2006-05-09 06:23:08',697,34);
INSERT INTO `questions` VALUES (241,247,NULL,'你心目中的理想社区是什么样的？','安全、友好、充满活力。',_binary '\0','2012-07-22 05:14:15',748,297);
INSERT INTO `questions` VALUES (242,463,NULL,'你最大的梦想是什么？','成为作家，环游世界。',_binary '\0','2020-08-25 23:07:10',709,580);
INSERT INTO `questions` VALUES (243,96,NULL,'描述最近一次让你感到开心的经历。','和朋友一起徒步旅行。',_binary '\0','2014-09-18 14:51:25',974,158);
INSERT INTO `questions` VALUES (244,965,NULL,'你最喜欢的游戏是什么？','棋类游戏。',_binary '\0','2006-12-01 20:57:38',856,374);
INSERT INTO `questions` VALUES (245,713,NULL,'如何在工作和生活中找到平衡？','制定优先级和合理安排时间。',_binary '\0','2005-01-17 09:38:54',521,240);
INSERT INTO `questions` VALUES (246,311,NULL,'你认为是否可以做到完全的独立？','在某种程度上可以，但不完全。',_binary '\0','2015-09-01 08:50:05',542,680);
INSERT INTO `questions` VALUES (247,49,NULL,'有什么事情让你感到极大的放松？','阅读一本好书。',_binary '\0','2016-11-11 02:42:37',490,785);
INSERT INTO `questions` VALUES (248,324,NULL,'你理想的退休生活是什么样的？','旅行木屋，享受宁静。',_binary '\0','2013-08-12 03:02:27',600,560);
INSERT INTO `questions` VALUES (249,628,NULL,'有什么事情是你最近学到的有趣知识？','关于心理健康的知识。',_binary '\0','2001-12-22 07:42:20',291,724);
INSERT INTO `questions` VALUES (250,964,NULL,'你认为教育的重要性在哪里？','是基础，更是未来的关键。',_binary '\0','2008-01-18 17:53:30',140,827);
INSERT INTO `questions` VALUES (251,708,NULL,'有什么独特的传统你家庭中保留着？','每年的家族聚会。',_binary '\0','2019-11-14 02:40:17',70,379);
INSERT INTO `questions` VALUES (252,871,NULL,'描述一下你的一个梦想。','写一本畅销小说。',_binary '\0','2001-07-16 03:39:16',858,318);
INSERT INTO `questions` VALUES (253,893,NULL,'如果你能改变任何一项法律，你会选择什么？','限制垃圾食品的广告。',_binary '\0','2013-12-08 18:41:34',872,809);
INSERT INTO `questions` VALUES (254,594,NULL,'你最希望传递给他人的价值观是什么？','善良与包容。',_binary '\0','2022-01-05 02:05:57',14,301);
INSERT INTO `questions` VALUES (255,228,NULL,'你有追随的榜样吗？是谁？','李光耀，领导智慧。',_binary '\0','2007-06-16 08:06:30',201,383);
INSERT INTO `questions` VALUES (256,583,NULL,'你如何看待义务与责任？','两者同样重要，促进人类进步。',_binary '\0','2022-11-06 22:19:20',760,348);
INSERT INTO `questions` VALUES (257,629,NULL,'有什么事情让你感到骄傲？','完成学业时的成就感。',_binary '\0','2011-08-22 21:20:17',462,273);
INSERT INTO `questions` VALUES (258,950,NULL,'你认为感恩的重要性是什么？','促使我们珍惜生命中的美好。',_binary '\0','2004-03-13 14:18:52',977,369);
INSERT INTO `questions` VALUES (259,457,NULL,'你最喜欢的户外活动是什么？','登山，亲近自然。',_binary '\0','2009-07-05 13:42:28',522,32);
INSERT INTO `questions` VALUES (260,412,NULL,'描述一下你心目中完美的一天。','享受阳光，书读到天黑。',_binary '\0','2018-05-25 09:08:06',458,62);
INSERT INTO `questions` VALUES (261,554,NULL,'你常用的放松方法是什么？','听音乐和冥想。',_binary '\0','2020-10-13 15:58:29',499,73);
INSERT INTO `questions` VALUES (262,825,NULL,'在你的生活中，最重要的三件事是什么？','健康、家庭、事业。',_binary '\0','2015-10-15 05:15:43',544,322);
INSERT INTO `questions` VALUES (263,737,NULL,'第一次旅行的经历对你有什么影响？','让我学会珍惜每一刻。',_binary '\0','2016-11-10 15:56:28',193,636);
INSERT INTO `questions` VALUES (264,514,NULL,'有什么事情是你一直想尝试的？','学习滑雪。',_binary '\0','2009-10-16 21:08:11',401,847);
INSERT INTO `questions` VALUES (265,26,NULL,'你认为什么样的地方最能让人感到平静？','海边或山顶的自然环境。',_binary '\0','2016-04-13 02:37:52',25,14);
INSERT INTO `questions` VALUES (266,224,NULL,'你的梦想假期是怎样的？','去巴哈马的沙滩度假。',_binary '\0','2003-06-02 08:14:13',972,201);
INSERT INTO `questions` VALUES (267,445,NULL,'如果可以选择任何一种超能力，你希望拥有什么？','瞬间移动。',_binary '\0','2023-06-05 10:17:20',110,191);
INSERT INTO `questions` VALUES (268,145,NULL,'描述你最喜欢的一次家庭聚会。','全家一起聚餐，分享笑声。',_binary '\0','2022-01-01 11:37:24',601,886);
INSERT INTO `questions` VALUES (269,567,NULL,'你认为人与人之间最珍贵的纽带是什么？','互信与支持。',_binary '\0','2014-08-22 12:52:13',994,836);
INSERT INTO `questions` VALUES (270,283,NULL,'你常听的播客或节目是什么？','有关科技和心理学的播客。',_binary '\0','2021-12-04 01:17:52',919,401);
INSERT INTO `questions` VALUES (271,245,NULL,'你觉得好的朋友应该具备哪些特质？','诚实、幽默、善良。',_binary '\0','2022-11-17 19:26:53',886,454);
INSERT INTO `questions` VALUES (272,856,NULL,'描述一次让你感到骄傲的经历。','完成自己的第一个项目。',_binary '\0','2022-03-17 04:28:52',723,514);
INSERT INTO `questions` VALUES (273,867,NULL,'你认为生活中最值得追求的是什么？','真正的快乐与内心的安宁。',_binary '\0','2007-08-16 21:27:48',74,655);
INSERT INTO `questions` VALUES (274,207,NULL,'有没有一首歌让你印象深刻？是什么？','《岁月神偷》。',_binary '\0','2021-02-25 10:53:25',939,586);
INSERT INTO `questions` VALUES (275,981,NULL,'你会如何形容自己的个性？','外向、乐观、好奇。',_binary '\0','2020-06-06 08:06:44',735,147);
INSERT INTO `questions` VALUES (276,204,NULL,'你有什么希望改变的社会问题？','教育不平等和贫困问题。',_binary '\0','2012-11-01 16:35:37',399,466);
INSERT INTO `questions` VALUES (277,553,NULL,'你认为怎样才能实现内心的平静？','通过冥想和自然散步。',_binary '\0','2001-01-16 00:36:40',414,851);
INSERT INTO `questions` VALUES (278,261,NULL,'有什么事情让你感觉到幸福？','与家人朋友共度的时光。',_binary '\0','2013-03-12 08:21:26',343,946);
INSERT INTO `questions` VALUES (279,215,NULL,'你最感激的事情是什么？','家人的支持和爱。',_binary '\0','2023-02-03 14:31:23',172,299);
INSERT INTO `questions` VALUES (280,196,NULL,'有没有一本书改变了你的世界观？','《活着》，改变了我对生命的理解。',_binary '\0','2013-09-08 01:53:43',782,360);
INSERT INTO `questions` VALUES (281,513,NULL,'描述你理想的社交活动是什么。','一起烧烤、聊天的聚会。',_binary '\0','2009-06-25 20:58:18',239,556);
INSERT INTO `questions` VALUES (282,983,NULL,'你对孤独有什么看法？','孤独是一种必要的自我反思。',_binary '\0','2019-03-31 22:45:13',67,863);
INSERT INTO `questions` VALUES (283,28,NULL,'对你来说，教育最重要的方面是什么？','实际应用和批判性思维。',_binary '\0','2020-10-09 23:38:06',20,565);
INSERT INTO `questions` VALUES (284,813,NULL,'如果有一天你能环游世界，你最想去哪里？','日本，体验文化和美食。',_binary '\0','2016-07-10 05:17:20',464,110);
INSERT INTO `questions` VALUES (285,530,NULL,'描述一个你永远不会忘记的瞬间。','第一次独自旅行的激动时刻。',_binary '\0','2004-04-13 15:58:11',139,235);
INSERT INTO `questions` VALUES (286,871,NULL,'你有如何应对失败的策略？','从失败中学习，重新开始。',_binary '\0','2016-04-05 14:50:07',245,312);
INSERT INTO `questions` VALUES (287,483,NULL,'有哪些事情让你感到充实？','学习新事物和帮助他人。',_binary '\0','2009-10-29 18:22:22',458,813);
INSERT INTO `questions` VALUES (288,406,NULL,'朋友对你来说意味着什么？','朋友是支持我成长的人。',_binary '\0','2016-09-14 08:14:00',29,595);
INSERT INTO `questions` VALUES (289,776,NULL,'你认为团队合作的重要性在哪里？','能够集思广益和分担责任。',_binary '\0','2019-04-04 18:16:07',380,998);
INSERT INTO `questions` VALUES (290,274,NULL,'如果你能对自己的过去说一句话，你会说什么？','一切都会好起来。',_binary '\0','2005-02-12 18:55:28',369,732);
INSERT INTO `questions` VALUES (291,554,NULL,'你最想学习的东西是什么？','写作和心理学。',_binary '\0','2021-11-03 17:08:48',51,248);
INSERT INTO `questions` VALUES (292,719,NULL,'有什么事情让你觉得生活的希望？','看到他人的幸福与进步。',_binary '\0','2022-08-06 12:01:36',81,828);
INSERT INTO `questions` VALUES (293,416,NULL,'你最欣赏的人擅长什么？','对艺术的坚定与创新。',_binary '\0','2023-06-29 13:18:54',895,555);
INSERT INTO `questions` VALUES (294,213,NULL,'描述一下你的早晨例行程序。','起床、喝水、锻炼和吃早餐。',_binary '\0','2011-03-11 17:27:13',7,263);
INSERT INTO `questions` VALUES (295,811,NULL,'你认为多元文化的价值是什么？','推动理解和社会融合。',_binary '\0','2018-01-09 09:01:28',577,82);
INSERT INTO `questions` VALUES (296,108,NULL,'你心中的完美社会是什么样的？','安全、包容且公正的社会。',_binary '\0','2018-06-14 18:46:37',78,936);
INSERT INTO `questions` VALUES (297,625,NULL,'描述一次难忘的购物经历。','找到独特的手工艺品。',_binary '\0','2007-10-03 04:37:24',502,739);
INSERT INTO `questions` VALUES (298,351,NULL,'你认为人类最大的成就是什么？','医疗和科技的进步。',_binary '\0','2024-01-01 15:33:38',964,918);
INSERT INTO `questions` VALUES (299,85,NULL,'你对未来的职业有怎样的想法？','追求我热爱的事业。',_binary '\0','2001-09-05 22:36:08',980,682);
INSERT INTO `questions` VALUES (300,806,NULL,'有什么事情能让你很快开心起来？','看一部喜剧电影。',_binary '\0','2007-10-26 13:27:22',186,295);
INSERT INTO `questions` VALUES (301,654,NULL,'描述一个你希望能够实现的愿望。','可以退休后旅行世界。',_binary '\0','2001-11-20 17:58:03',335,396);
INSERT INTO `questions` VALUES (302,836,NULL,'你如何看待社交媒体对人际关系的影响？','促进交流，但也可能引发误解。',_binary '\0','2007-04-27 12:39:19',238,132);
INSERT INTO `questions` VALUES (303,296,NULL,'你最崇敬的艺术家是谁，为什么？','艾伦·德杰尼勒斯，因为她的幽默感和善良。',_binary '\0','2014-11-23 04:01:34',93,376);
INSERT INTO `questions` VALUES (304,121,NULL,'你是否相信运气在生活中会影响到人的命运？','相信，但努力更重要。',_binary '\0','2018-05-07 12:31:33',100,84);
INSERT INTO `questions` VALUES (305,185,NULL,'描述一下你和自己的关系。','有时需要反思，有时是支持者。',_binary '\0','2001-09-07 21:28:47',669,124);
INSERT INTO `questions` VALUES (306,993,NULL,'你觉得理想的政府应该具备哪些特质？','公正、透明和服务意识。',_binary '\0','2021-01-18 08:41:28',305,529);
INSERT INTO `questions` VALUES (307,800,NULL,'对你来说，成功的人生意味着什么？','能自由选择自己想要的生活。',_binary '\0','2002-08-14 12:01:44',937,149);
INSERT INTO `questions` VALUES (308,969,NULL,'有哪些小事情能让你感到幸福和满足？','一杯热茶和一本好书。',_binary '\0','2008-03-12 03:56:06',113,744);
INSERT INTO `questions` VALUES (309,496,NULL,'描述一下你最喜欢的运动。','跑步，能释放压力。',_binary '\0','2004-10-29 08:56:26',492,592);
INSERT INTO `questions` VALUES (310,280,NULL,'你曾经做过的最有挑战性的事情是什么？','攀登高峰的挑战和喜悦。',_binary '\0','2003-01-24 03:57:54',660,350);
INSERT INTO `questions` VALUES (311,342,NULL,'你心里最美丽的地方是哪里？','我的家乡的山水风光。',_binary '\0','2005-12-18 10:02:36',310,604);
INSERT INTO `questions` VALUES (312,465,NULL,'描述一下你最喜欢的书或作者。','《百年孤独》，加西亚·马尔克斯的作品。',_binary '\0','2011-03-27 07:03:52',881,319);
INSERT INTO `questions` VALUES (313,225,NULL,'你认为科技对人际关系的影响是什么？','科技可以拉近人际关系，但也可能造成疏离。',_binary '\0','2013-09-27 17:37:34',423,246);
INSERT INTO `questions` VALUES (314,288,NULL,'有没有什么特别的习惯或仪式，你每天都做？','每天早晨冥想和锻炼。',_binary '\0','2001-05-04 15:36:55',290,33);
INSERT INTO `questions` VALUES (315,476,NULL,'你最喜欢的运动队或运动员是谁？','我最喜欢的运动员是迈克尔·乔丹。',_binary '\0','2020-12-27 11:53:14',665,441);
INSERT INTO `questions` VALUES (316,862,NULL,'你如何看待环境保护的重要性？','非常重要，关系到地球的未来。',_binary '\0','2017-09-26 04:47:06',176,931);
INSERT INTO `questions` VALUES (317,507,NULL,'如果时间旅行是真的，你想去哪个年代？','希腊古代，体验哲学和文化的起源。',_binary '\0','2017-06-24 15:18:57',140,12);
INSERT INTO `questions` VALUES (318,352,NULL,'你认为社会对年轻人的期望是什么？','希望他们能够勇敢追求自己的梦想。',_binary '\0','2021-11-04 17:43:47',546,160);
INSERT INTO `questions` VALUES (319,546,NULL,'你认为生活中最重要的技能是什么？','辩证思维和解决问题的能力。',_binary '\0','2018-02-15 21:41:21',255,386);
INSERT INTO `questions` VALUES (320,23,NULL,'有什么事情让你感到特别的成就感？','完成一个长期项目或目标。',_binary '\0','2014-09-27 17:48:03',563,875);
INSERT INTO `questions` VALUES (321,402,NULL,'拥有宠物对你生活的影响是什么？','宠物让我感到陪伴和快乐。',_binary '\0','2003-03-04 03:56:45',314,459);
INSERT INTO `questions` VALUES (322,931,NULL,'有没有一件事让你最近感到内疚或后悔？','可能没有珍惜与某些朋友的时光。',_binary '\0','2008-03-13 03:08:32',553,514);
INSERT INTO `questions` VALUES (323,809,NULL,'描述一下你最喜欢的风格的音乐或歌手。','我喜欢听周杰伦的音乐。',_binary '\0','2009-04-16 15:59:37',832,980);
INSERT INTO `questions` VALUES (324,84,NULL,'你对未来的职业有什么规划或梦想？','成为一名作家或心理咨询师。',_binary '\0','2015-11-20 17:50:10',372,469);
INSERT INTO `questions` VALUES (325,792,NULL,'如何在日常生活中保持积极心态？','保持积极的自言自语和乐观的社交圈。',_binary '\0','2007-12-04 21:26:36',526,574);
INSERT INTO `questions` VALUES (326,847,NULL,'你认为家庭角色在现代社会中变化吗？','是的，更加灵活和多元。',_binary '\0','2015-09-15 08:48:16',329,832);
INSERT INTO `questions` VALUES (327,414,NULL,'如果能与任何历史人物对话，你会选择谁？','阿尔伯特·爱因斯坦，讨论科学与生活。',_binary '\0','2004-02-03 07:27:54',96,369);
INSERT INTO `questions` VALUES (328,395,NULL,'你觉得当今社会最需要改变的是什么？','减少环境污染和促进平等。',_binary '\0','2003-03-29 17:29:48',192,163);
INSERT INTO `questions` VALUES (329,435,NULL,'有什么事情是你希望能够教给年轻一代的？','学会关心他人和珍惜时间。',_binary '\0','2016-03-01 15:12:48',43,955);
INSERT INTO `questions` VALUES (330,283,NULL,'你做过的最有趣的事情是什么？','在活动中做志愿者，帮助别人。',_binary '\0','2002-01-13 02:07:04',168,190);
INSERT INTO `questions` VALUES (331,378,NULL,'你最喜欢的节日是什么？为什么？','春节，因为它代表团聚和希望。',_binary '\0','2004-12-17 14:29:47',7,179);
INSERT INTO `questions` VALUES (332,872,NULL,'有什么你从不敢尝试但想要尝试的活动？','跳伞，一直很想尝试但太害怕。',_binary '\0','2024-01-06 04:19:33',545,424);
INSERT INTO `questions` VALUES (333,555,NULL,'你认为理解和沟通在友谊中有多重要？','非常重要，能加深理解与信任。',_binary '\0','2008-08-20 10:21:35',73,663);
INSERT INTO `questions` VALUES (334,474,NULL,'如果你可以居住在任何地方，你想去哪里？','一个安静的海边村庄。',_binary '\0','2014-04-28 19:19:23',61,371);
INSERT INTO `questions` VALUES (335,272,NULL,'你对财富的看法是什么？','财富是一种工具，重要的是如何使用它。',_binary '\0','2017-12-07 12:58:55',111,117);
INSERT INTO `questions` VALUES (336,519,NULL,'如何平衡工作和个人生活？','设定明确的界限和目标。',_binary '\0','2004-02-17 20:49:30',827,749);
INSERT INTO `questions` VALUES (337,24,NULL,'有什么习惯对你的生活产生了积极影响？','每天写日记，帮助我反思和成长。',_binary '\0','2011-09-21 01:45:46',202,933);
INSERT INTO `questions` VALUES (338,225,NULL,'你对善良和同情心的理解是什么？','是对他人痛苦的敏感与理解。',_binary '\0','2014-04-07 11:30:56',386,297);
INSERT INTO `questions` VALUES (339,461,NULL,'有哪些事情是你希望在生活中留给自己或他人的遗产？','积累爱与善良的回忆。',_binary '\0','2022-04-09 16:14:50',188,881);
INSERT INTO `questions` VALUES (340,180,NULL,'你认为健康饮食对生活质量的影响是什么？','影响很大，能提升身体和心理健康。',_binary '\0','2016-09-08 16:18:03',585,694);
INSERT INTO `questions` VALUES (341,312,NULL,'描述一下你认为的理想工作环境。','开放、包容的工作环境，强调团队合作。',_binary '\0','2019-01-13 08:21:28',851,680);
INSERT INTO `questions` VALUES (342,924,NULL,'有什么事情让你感到害怕或不安？','对未来的未知和失败的恐惧。',_binary '\0','2020-02-28 03:51:28',608,109);
INSERT INTO `questions` VALUES (343,680,NULL,'你怎么看待当今社会的压力？','生活中压力普遍存在，但管理得当会有积极效果。',_binary '\0','2006-03-28 17:43:53',725,660);
INSERT INTO `questions` VALUES (344,54,NULL,'你的爱好是什么，它们带给你什么？','我喜欢阅读、写作和旅行。',_binary '\0','2006-03-13 17:08:48',223,437);
INSERT INTO `questions` VALUES (345,922,NULL,'你最喜欢的电影或电视剧是什么？','我最喜欢的电视剧是《老友记》。',_binary '\0','2008-02-28 21:44:19',635,597);
INSERT INTO `questions` VALUES (346,681,NULL,'你对团队合作的看法是什么？','团队合作可以带来更好的成果和创新。',_binary '\0','2019-10-20 12:24:29',402,120);
INSERT INTO `questions` VALUES (347,757,NULL,'有什么事情让你觉得生活充满希望？','听闻身边人的成长与幸福。',_binary '\0','2009-09-25 11:56:36',312,500);
INSERT INTO `questions` VALUES (348,116,NULL,'你认为未来的教育会是什么样子？','更加注重实践与个人发展。',_binary '\0','2002-10-24 02:33:36',246,504);
INSERT INTO `questions` VALUES (349,683,NULL,'如果可以给年轻人一个建议，你会说什么？','要勇敢并相信自己能够改变世界。',_binary '\0','2022-08-20 16:34:48',402,511);
INSERT INTO `questions` VALUES (350,326,NULL,'你有什么特别的才能或技能？','擅长写作和沟通。',_binary '\0','2011-04-29 01:55:41',496,257);
INSERT INTO `questions` VALUES (351,194,NULL,'有什么事情是你意料之外却让你开心的？','意外收到的朋友的支持信息。',_binary '\0','2005-09-14 20:02:26',330,569);
INSERT INTO `questions` VALUES (352,476,NULL,'描述一下你最喜欢的城市或乡村。','我最喜欢的城市是京都，自然与文化结合得很好。',_binary '\0','2012-06-18 01:48:45',566,975);
INSERT INTO `questions` VALUES (353,336,NULL,'你对压力管理有什么看法或策略？','学会放松技巧，如深呼吸和运动。',_binary '\0','2011-01-27 06:36:59',791,923);
INSERT INTO `questions` VALUES (354,250,NULL,'有什么事情在你的生活中产生了重大影响？','家庭的支持与教育经历。',_binary '\0','2015-11-21 09:14:24',228,335);
INSERT INTO `questions` VALUES (355,44,NULL,'你希望自己的生活中有什么样的变化？','希望生活更加平衡与充实。',_binary '\0','2018-10-23 19:19:50',508,70);
INSERT INTO `questions` VALUES (356,678,NULL,'如何提高自己的创造力和想象力？','多读书和尝试新的创意活动。',_binary '\0','2002-01-23 01:30:45',114,918);
INSERT INTO `questions` VALUES (357,72,NULL,'你如何看待传统与现代的关系？','传统为现代提供基础，但现代需要创新。',_binary '\0','2008-05-20 13:10:59',745,483);
INSERT INTO `questions` VALUES (358,372,NULL,'有没有一位老师或导师对你产生了重要影响？','是的，她教会了我如何思考。',_binary '\0','2014-07-29 02:46:22',782,121);
INSERT INTO `questions` VALUES (359,86,NULL,'描述一下你希望在人际关系中建立的连接。','建立真诚与支持的关系。',_binary '\0','2019-11-14 03:57:52',302,209);
INSERT INTO `questions` VALUES (360,476,NULL,'你觉得定义自己身份的因素有哪些？','文化背景、生活经历和个人价值观。',_binary '\0','2021-05-23 16:42:54',766,550);
INSERT INTO `questions` VALUES (361,409,NULL,'你觉得生活的意义是什么？','生活的意义在于体验、成长与贡献。',_binary '\0','2014-05-22 21:30:08',300,408);
INSERT INTO `questions` VALUES (362,97,NULL,'描述一下你理想的职业生涯。','理想的职业是在艺术与心理学领域结合创作。',_binary '\0','2002-03-19 03:04:11',360,940);
INSERT INTO `questions` VALUES (363,488,NULL,'有什么事情是你希望能改变的习惯？','希望能改善拖延的习惯。',_binary '\0','2023-01-09 19:40:36',682,22);
INSERT INTO `questions` VALUES (364,906,NULL,'你最喜欢的童年记忆是什么？','和家人在一起的快乐时光，特别是节日聚会。',_binary '\0','2022-03-20 09:26:33',664,819);
INSERT INTO `questions` VALUES (365,253,NULL,'对你影响最大的书籍是什么？','《活出生命的意义》，它改变了我对困境的看法。',_binary '\0','2017-01-10 08:43:37',793,93);
INSERT INTO `questions` VALUES (366,974,NULL,'描述一次你学到重要人生课的一次经历。','一次跌倒后重新站起来的经历，让我明白坚持的重要性。',_binary '\0','2012-03-08 21:35:48',138,741);
INSERT INTO `questions` VALUES (367,495,NULL,'你最喜欢的食物是什么？','我最喜欢的食物是意大利面。',_binary '\0','2023-09-21 03:15:02',900,125);
INSERT INTO `questions` VALUES (368,925,NULL,'如果可以选择任何一个地方生活，你会选哪个？','新西兰，享受自然的美丽与宁静。',_binary '\0','2008-10-18 12:45:36',583,489);
INSERT INTO `questions` VALUES (369,408,NULL,'有什么事情让你感到非常自豪？','完成学业并找到理想工作。',_binary '\0','2023-03-08 18:39:27',244,751);
INSERT INTO `questions` VALUES (370,603,NULL,'描述一下你的旅行愿望清单。','日本的樱花季、欧洲的文化之旅、探索南美的热带雨林。',_binary '\0','2020-11-10 09:20:07',39,570);
INSERT INTO `questions` VALUES (371,59,NULL,'你认为影响力最大的社会运动是什么？','人权运动，推动了全球的平等与公正。',_binary '\0','2000-07-13 21:34:12',836,133);
INSERT INTO `questions` VALUES (372,658,NULL,'你最喜欢的艺术形式是什么？','我最喜欢的艺术形式是绘画和音乐。',_binary '\0','2015-04-16 04:55:21',996,597);
INSERT INTO `questions` VALUES (373,97,NULL,'如果时间不限，你希望掌握哪种技能？','希望掌握流利的外语，例如西班牙语。',_binary '\0','2008-10-23 00:14:19',320,548);
INSERT INTO `questions` VALUES (374,869,NULL,'描述一下你心目中完美的一周。','和朋友旅行、阅读、学习新技能的完美一周。',_binary '\0','2008-09-26 02:28:39',171,779);
INSERT INTO `questions` VALUES (375,401,NULL,'你在生活中最大的挑战是什么？','学业与工作的平衡。',_binary '\0','2018-07-11 03:41:25',636,496);
INSERT INTO `questions` VALUES (376,359,NULL,'有什么事情是你想要实现的长期目标？','成为一名心理咨询师，帮助他人。',_binary '\0','2013-12-26 07:09:20',970,680);
INSERT INTO `questions` VALUES (377,44,NULL,'你觉得人们最常犯的错误是什么？','常常过于追求完美。',_binary '\0','2002-03-03 15:17:09',254,569);
INSERT INTO `questions` VALUES (378,152,NULL,'如果你能改变一件事，你会选择什么？','改变对待失败的态度，让它成为学习的机会。',_binary '\0','2019-09-17 00:55:18',734,277);
INSERT INTO `questions` VALUES (379,239,NULL,'你如何应对工作和生活的压力？','定期锻炼和冥想应对压力。',_binary '\0','2013-05-12 22:49:52',125,873);
INSERT INTO `questions` VALUES (380,616,NULL,'描述一下你的家乡。','我家乡的自然风光和传统文化。',_binary '\0','2020-07-17 01:43:30',371,396);
INSERT INTO `questions` VALUES (381,673,NULL,'你认为幸福的定义是什么？','幸福是一种内心的满足与宁静。',_binary '\0','2004-12-04 05:57:09',891,456);
INSERT INTO `questions` VALUES (382,793,NULL,'有什么事情是你希望大家都能学习的？','更加同理和理解他人。',_binary '\0','2020-06-11 16:57:34',411,757);
INSERT INTO `questions` VALUES (383,50,NULL,'你如何看待长途友谊？','是的，虽然很有挑战，但也很宝贵。',_binary '\0','2015-07-27 21:45:41',619,280);
INSERT INTO `questions` VALUES (384,713,NULL,'有什么事情让你感到被重视？','感觉到我为他人带来的帮助和影响。',_binary '\0','2018-03-03 14:55:17',464,839);
INSERT INTO `questions` VALUES (385,666,NULL,'你认为什么样的活动能促进团队合作？','团队建设活动和开放式沟通。',_binary '\0','2019-10-16 16:16:53',949,979);
INSERT INTO `questions` VALUES (386,342,NULL,'你希望为社会做出什么贡献？','希望通过心理健康倡导来促进社会福祉。',_binary '\0','2024-10-10 19:47:44',620,629);
INSERT INTO `questions` VALUES (387,420,NULL,'有什么事情是你认为人们应该更关注的？','心理健康和环境保护。',_binary '\0','2024-05-28 11:20:20',614,781);
INSERT INTO `questions` VALUES (388,849,NULL,'如何在生活中找到平衡？','学会设定优先级与边界。',_binary '\0','2015-07-13 10:09:39',976,750);
INSERT INTO `questions` VALUES (389,86,NULL,'描述一次让你改变看法的对话。','一次深夜的谈话让我对人生有了新的思考。',_binary '\0','2003-04-02 13:46:22',857,640);
INSERT INTO `questions` VALUES (390,280,NULL,'你对当前时事的看法是什么？','对社会发展的担忧与希望。',_binary '\0','2018-02-16 02:48:12',693,423);
INSERT INTO `questions` VALUES (391,54,NULL,'有什么事情是你想尝试而未曾尝试过的？','尝试极限运动，如攀岩。',_binary '\0','2024-02-14 05:00:05',161,658);
INSERT INTO `questions` VALUES (392,882,NULL,'你认为积极心态的重要性在哪里？','积极心态帮助我们应对挑战和困难。',_binary '\0','2006-07-25 21:10:31',395,514);
INSERT INTO `questions` VALUES (393,752,NULL,'描述一下你最喜欢的季节及原因。','我最喜欢秋天，气候宜人，万物丰收。',_binary '\0','2014-02-20 19:37:33',10,452);
INSERT INTO `questions` VALUES (394,141,NULL,'你最钦佩的人是谁，为什么？','我最钦佩的人成为社会变革者，如马丁·路德·金。',_binary '\0','2010-05-19 15:44:20',227,433);
INSERT INTO `questions` VALUES (395,720,NULL,'如何评价当前教育体系的优劣？','教育应该更注重实践和批判性思维。',_binary '\0','2011-04-22 08:17:09',84,122);
INSERT INTO `questions` VALUES (396,238,NULL,'你对待失败的态度是什么？','失败是成功之母，是学习的机会。',_binary '\0','2013-12-02 04:28:10',391,554);
INSERT INTO `questions` VALUES (397,681,NULL,'有什么值得期待的新趋势或技术？','人工智能和可持续发展的新趋势。',_binary '\0','2021-08-22 21:50:38',9,133);
INSERT INTO `questions` VALUES (398,870,NULL,'你认为人们应如何珍惜时光？','设定目标，合理安排时间。',_binary '\0','2001-07-25 06:33:59',92,611);
INSERT INTO `questions` VALUES (399,264,NULL,'描述一下你认为的理想家庭生活。','理想家庭生活是相互支持与理解。',_binary '\0','2018-12-08 17:47:43',864,521);
INSERT INTO `questions` VALUES (400,540,NULL,'你最喜欢的方式来庆祝成功是什么？','和家人朋友分享庆祝时刻，聚餐或旅行。',_binary '\0','2020-07-07 14:19:30',695,304);
INSERT INTO `questions` VALUES (401,726,NULL,'你喜欢的电影类型是什么？','我喜欢科幻和剧情片。',_binary '\0','2012-12-01 05:29:13',865,144);
INSERT INTO `questions` VALUES (402,946,NULL,'有什么事情是你一直想学但没机会的？','学习乐器，例如吉他或钢琴。',_binary '\0','2017-07-06 20:10:07',319,518);
INSERT INTO `questions` VALUES (403,425,NULL,'你认为内心的平静来源于哪里？','内心的平静来自自我反思和冥想。',_binary '\0','2016-02-06 19:47:47',213,85);
INSERT INTO `questions` VALUES (404,53,NULL,'描述你的生活哲学或信念。','生活中的每一次经历都是学习的机会。',_binary '\0','2019-06-05 09:51:38',647,212);
INSERT INTO `questions` VALUES (405,242,NULL,'你认为社交技巧对职业发展的重要性如何？','社交技巧是建立人际关系和职业发展的基础。',_binary '\0','2008-06-03 17:32:37',211,221);
INSERT INTO `questions` VALUES (406,179,NULL,'描述一次让你感到快乐的简单瞬间。','一杯热茶时，窗外阳光洒下的瞬间。',_binary '\0','2023-09-10 03:05:56',257,599);
INSERT INTO `questions` VALUES (407,906,NULL,'你认为什么是个人成长的关键？','明确的目标与持续的学习。',_binary '\0','2013-06-06 14:41:00',73,447);
INSERT INTO `questions` VALUES (408,602,NULL,'如何在生活中激励自己？','定期回顾目标与成就，保持动力。',_binary '\0','2019-12-23 07:35:05',722,990);
INSERT INTO `questions` VALUES (409,7,NULL,'你希望未来的自己是什么样子？','希望成为一个更成熟、智慧的人。',_binary '\0','2004-06-26 14:55:54',159,114);
INSERT INTO `questions` VALUES (410,999,NULL,'有什么事情是你一直渴望去做却没有去做的？','一直想去环球旅行，还没实现。',_binary '\0','2018-11-12 19:57:55',10,220);
INSERT INTO `questions` VALUES (411,984,NULL,'如果你可以拥有任何一种超能力，你希望选择什么？','能飞翔的超能力。',_binary '\0','2000-12-17 08:36:47',884,962);
INSERT INTO `questions` VALUES (412,410,NULL,'描述一次你感到无比快乐的经历。','和朋友一起在海边的日子。',_binary '\0','2009-10-12 22:30:41',627,759);
INSERT INTO `questions` VALUES (413,514,NULL,'你认为自己的个性中最特别的是什么？','我的好奇心与包容心。',_binary '\0','2021-05-13 12:52:11',461,839);
INSERT INTO `questions` VALUES (414,780,NULL,'有什么让你感到放松的方法？','听音乐和冥想。',_binary '\0','2008-01-21 01:15:42',462,732);
INSERT INTO `questions` VALUES (415,594,NULL,'你最敬佩的历史人物是谁？','马丁·路德·金，因其对平等的追求。',_binary '\0','2009-01-03 16:44:17',411,382);
INSERT INTO `questions` VALUES (416,164,NULL,'你认为交流中的非语言信息有多重要？','非语言信息如肢体语言和面部表情非常重要。',_binary '\0','2012-10-20 14:21:38',519,492);
INSERT INTO `questions` VALUES (417,986,NULL,'描述一下你理想的假期。','在一座海边别墅放松几周，远离喧嚣。',_binary '\0','2019-11-12 13:57:33',503,705);
INSERT INTO `questions` VALUES (418,360,NULL,'有什么事情让你感到恐惧？','失去亲人的恐惧。',_binary '\0','2023-04-22 00:40:15',420,170);
INSERT INTO `questions` VALUES (419,138,NULL,'你如何看待社交媒体对人际关系的影响？','社交媒体既拉近了联系，也增添了隔阂。',_binary '\0','2016-05-04 22:35:17',799,455);
INSERT INTO `questions` VALUES (420,450,NULL,'你最大的梦想是什么？','实现自身的理想与价值。',_binary '\0','2021-10-07 07:08:03',896,786);
INSERT INTO `questions` VALUES (421,179,NULL,'定义一下成功在你心中的意义。','成功在于实现自己的目标。',_binary '\0','2001-10-10 04:20:19',725,708);
INSERT INTO `questions` VALUES (422,384,NULL,'如果给你一个月的时间去做任何事，你会选择什么？','参加一次世界旅行，探索不同文化。',_binary '\0','2005-09-03 13:41:41',385,622);
INSERT INTO `questions` VALUES (423,747,NULL,'你觉得人们在年轻时最大的误区是什么？','过早追求成就，忽视了内心感受。',_binary '\0','2011-06-14 19:51:03',911,382);
INSERT INTO `questions` VALUES (424,253,NULL,'你对慈善事业有什么看法？','慈善事业是对社会的重要贡献。',_binary '\0','2024-03-17 18:56:14',965,164);
INSERT INTO `questions` VALUES (425,941,NULL,'你喜欢独处还是与朋友在一起？','我喜欢与朋友在一起，却也珍惜独处的时光。',_binary '\0','2002-11-08 02:26:41',495,952);
INSERT INTO `questions` VALUES (426,738,NULL,'如果可以重返过去的一天，你会选择哪一天？','我会选择某次假期中的愉快时光。',_binary '\0','2020-09-09 01:30:58',189,271);
INSERT INTO `questions` VALUES (427,892,NULL,'你认为构建自信的关键是什么？','自我反思和不断尝试新事物。',_binary '\0','2023-06-18 11:51:47',998,378);
INSERT INTO `questions` VALUES (428,239,NULL,'你觉得最重要的事是什么？','追求真诚和快乐。',_binary '\0','2016-12-11 22:37:15',14,305);
INSERT INTO `questions` VALUES (429,982,NULL,'描述一次改变你生活的决定。','放弃嫉妒，聚焦成长。',_binary '\0','2003-06-08 14:59:06',848,626);
INSERT INTO `questions` VALUES (430,473,NULL,'你最难忘的一次旅行经历是什么？','和家人一起去旅行的经历。',_binary '\0','2004-07-31 12:36:53',776,422);
INSERT INTO `questions` VALUES (431,783,NULL,'在生活中，你是如何处理冲突的？','通过沟通与理解解决冲突。',_binary '\0','2009-07-05 20:05:37',796,686);
INSERT INTO `questions` VALUES (432,525,NULL,'有什么技能是你希望现在就能掌握的？','编程或一门外语。',_binary '\0','2009-03-06 17:29:17',807,409);
INSERT INTO `questions` VALUES (433,291,NULL,'你最大的灵感来源是什么？','自然和艺术。',_binary '\0','2011-05-03 18:30:30',429,334);
INSERT INTO `questions` VALUES (434,62,NULL,'描述一下你理想的生活方式。','着重于家庭与事业的平衡。',_binary '\0','2015-02-18 21:10:59',64,814);
INSERT INTO `questions` VALUES (435,136,NULL,'你对未来十年的期望是什么？','接受变化，努力追求目标。',_binary '\0','2015-07-25 03:35:27',351,653);
INSERT INTO `questions` VALUES (436,517,NULL,'有什么事情你认为应该立刻改变？','提高心理健康意识。',_binary '\0','2024-04-17 18:23:37',991,75);
INSERT INTO `questions` VALUES (437,581,NULL,'你最爱的一句名言或经典语录是什么？','“The only way to do great work is to love what you do.”',_binary '\0','2013-08-12 21:35:07',419,500);
INSERT INTO `questions` VALUES (438,592,NULL,'你如何看待失败与成功的关系？','失败是成功的重要组成部分。',_binary '\0','2022-03-28 10:47:12',543,7);
INSERT INTO `questions` VALUES (439,370,NULL,'你认为有多重要去追求自己的激情？','追寻激情能带来生活的动力与快乐。',_binary '\0','2021-08-10 22:22:38',947,880);
INSERT INTO `questions` VALUES (440,214,NULL,'描述一下你最喜欢的家庭传统。','春节团圆饭是特别的传统。',_binary '\0','2015-08-12 14:29:10',25,110);
INSERT INTO `questions` VALUES (441,731,NULL,'有什么事情是你希望从未发生的？','与某些人分开的事。',_binary '\0','2016-12-08 13:15:19',394,113);
INSERT INTO `questions` VALUES (442,820,NULL,'你对教育的看法是什么？','教育需要更多实践和批判性思维。',_binary '\0','2008-02-05 19:25:37',947,900);
INSERT INTO `questions` VALUES (443,962,NULL,'描述你最喜欢的运动或运动队。','我喜欢篮球和湖人队。',_binary '\0','2008-04-28 00:06:05',221,200);
INSERT INTO `questions` VALUES (444,589,NULL,'在工作中，你觉得最重要的品质是什么？','诚信和团队合作精神。',_binary '\0','2014-12-09 04:11:36',399,309);
INSERT INTO `questions` VALUES (445,239,NULL,'你认为终身学习的价值在哪里？','终身学习是个人与社会进步的基础。',_binary '\0','2010-07-12 17:56:14',944,694);
INSERT INTO `questions` VALUES (446,240,NULL,'有什么事情是你希望能多花时间做的？','写作和学习新语言。',_binary '\0','2004-08-16 21:20:45',841,859);
INSERT INTO `questions` VALUES (447,868,NULL,'你最喜欢的季节是什么，为什么？','我喜欢秋天，色彩斑斓，气候宜人。',_binary '\0','2021-04-26 00:18:49',240,376);
INSERT INTO `questions` VALUES (448,848,NULL,'你认为现代社会最需要的品质是什么？','诚实和包容。',_binary '\0','2013-03-24 17:57:57',627,170);
INSERT INTO `questions` VALUES (449,809,NULL,'描述一次让你感到满足的经历。','志愿活动带来的成就感。',_binary '\0','2004-11-29 22:20:28',124,824);
INSERT INTO `questions` VALUES (450,631,NULL,'如果时间和金钱不是问题，你会为自己选择什么样的生活方式？','在海滩度假，享受美好的生活。',_binary '\0','2006-06-19 18:48:00',754,30);
INSERT INTO `questions` VALUES (451,337,NULL,'你认为个人财务管理的重要性如何？','理财与制定预算重要性不言而喻。',_binary '\0','2013-03-30 08:41:36',8,29);
INSERT INTO `questions` VALUES (452,206,NULL,'描述你心目中理想的伴侣是什么样的。','理想中的伴侣应该有智慧与幽默感。',_binary '\0','2020-05-05 23:53:39',65,443);
INSERT INTO `questions` VALUES (453,112,NULL,'有什么事情让你感到挫败？','处理复杂人际关系的挫败感。',_binary '\0','2019-11-05 07:25:11',804,527);
INSERT INTO `questions` VALUES (454,667,NULL,'你如何定义忠诚？','忠诚在于对他人的一贯支持和信任。',_binary '\0','2002-06-03 18:55:41',928,771);
INSERT INTO `questions` VALUES (455,765,NULL,'你对生活中的变化持有什么态度？','生活的变化不可避免，学会接纳。',_binary '\0','2016-01-12 12:27:38',223,10);
INSERT INTO `questions` VALUES (456,541,NULL,'描述一下你最喜欢的电影角色或书中人物。','斯皮尔伯格的角色，执着与创造力结合。',_binary '\0','2000-11-18 01:50:34',912,69);
INSERT INTO `questions` VALUES (457,513,NULL,'有什么事情是你认为值得坚持的？','坚持自我信念与价值。',_binary '\0','2021-01-26 07:59:59',223,324);
INSERT INTO `questions` VALUES (458,233,NULL,'你最喜欢的电视节目是什么，为什么？','我最喜欢的剧集是《黑镜》，因其引人深思。',_binary '\0','2014-09-28 17:32:26',140,89);
INSERT INTO `questions` VALUES (459,647,NULL,'关于保持健康，你有哪些个人经验或建议？','坚持锻炼和保持均衡饮食。',_binary '\0','2010-12-11 23:42:20',293,144);
INSERT INTO `questions` VALUES (460,459,NULL,'描述你最喜欢的儿童故事或童话。','《小王子》，它让我懂得了爱与责任。',_binary '\0','2020-10-02 01:22:01',56,969);
INSERT INTO `questions` VALUES (461,612,NULL,'有什么事情让你有成就感的？','完成学业的那一刻。',_binary '\0','2018-12-13 20:47:12',768,498);
INSERT INTO `questions` VALUES (462,611,NULL,'如何培养自己的创造力？','经常练习和开放思维。',_binary '\0','2004-12-25 22:21:45',407,130);
INSERT INTO `questions` VALUES (463,713,NULL,'你认为友谊在生活中的作用是什么？','友谊是支持和理解的纽带。',_binary '\0','2015-10-08 13:28:31',563,461);
INSERT INTO `questions` VALUES (464,11,NULL,'如果可以选择成为任何动物，你会选择什么？','我想成为只在自然中生活的鸟。',_binary '\0','2003-09-16 20:52:14',939,241);
INSERT INTO `questions` VALUES (465,350,NULL,'你对最喜欢的假期有什么特别的回忆？','与家人一起度过的圣诞节。',_binary '\0','2003-05-13 12:58:59',653,265);
INSERT INTO `questions` VALUES (466,451,NULL,'描述你心目中完美的一餐。','和好友们共同分享的晚餐。',_binary '\0','2022-12-07 09:58:43',207,774);
INSERT INTO `questions` VALUES (467,537,NULL,'有什么是你在生活中不愿妥协的？','对于健康的坚持与自律。',_binary '\0','2011-04-28 14:01:10',44,271);
INSERT INTO `questions` VALUES (468,564,NULL,'你想如何影响他人的生活？','影响他人的生活方式是传播知识和爱。',_binary '\0','2022-12-31 15:08:48',835,48);
INSERT INTO `questions` VALUES (469,272,NULL,'描述一下你最喜欢的玩具或游戏。','我最喜欢的游戏是象棋。',_binary '\0','2017-03-27 17:14:31',764,14);
INSERT INTO `questions` VALUES (470,517,NULL,'你认为应该如何处理压力？','通过运动和放松技巧应对压力。',_binary '\0','2002-05-18 14:37:42',937,875);
INSERT INTO `questions` VALUES (471,44,NULL,'有什么事是你希望大家都能理解的？','对心理健康的普遍重视。',_binary '\0','2008-01-06 14:58:46',406,690);
INSERT INTO `questions` VALUES (472,703,NULL,'描述一次你感到无奈的情况。','在失业时感到无奈。',_binary '\0','2011-02-05 06:09:47',389,55);
INSERT INTO `questions` VALUES (473,296,NULL,'你对爱与被爱的理解是什么？','爱是理解与支持的结合。',_binary '\0','2017-11-19 21:58:08',973,889);
INSERT INTO `questions` VALUES (474,108,NULL,'你如何评估自己的价值观？','反思自己的经历和学习。',_binary '\0','2002-06-06 12:49:49',611,158);
INSERT INTO `questions` VALUES (475,464,NULL,'有什么事情是你想尝试的冒险活动？','学习滑雪或潜水等冒险活动。',_binary '\0','2006-02-07 09:16:27',52,734);
INSERT INTO `questions` VALUES (476,725,NULL,'Describe how you celebrate your birthday.','我会和朋友们聚会庆祝。',_binary '\0','2012-05-05 06:13:42',38,229);
INSERT INTO `questions` VALUES (477,800,NULL,'你觉得家庭对个人成长的重要性如何？','家庭支持对我的成长至关重要。',_binary '\0','2010-02-07 19:38:39',451,728);
INSERT INTO `questions` VALUES (478,367,NULL,'你如何看待社会责任感？','社会责任感能推动积极的变化。',_binary '\0','2022-08-16 07:42:15',540,370);
INSERT INTO `questions` VALUES (479,264,NULL,'描述一下你心中的英雄。','我心中的英雄是无畏无惧的志愿者。',_binary '\0','2019-03-29 15:45:41',82,671);
INSERT INTO `questions` VALUES (480,572,NULL,'有什么事情会让你感到嫉妒？','看到身边他人的成功。',_binary '\0','2020-10-01 06:15:37',739,322);
INSERT INTO `questions` VALUES (481,44,NULL,'你认为自我反省的重要性在哪里？','自我反省帮助我成长与改善。',_binary '\0','2016-01-16 21:32:25',211,654);
INSERT INTO `questions` VALUES (482,531,NULL,'描述一次你制作或参与的特别活动。','举办家庭聚会时的欢乐。',_binary '\0','2023-02-25 05:03:33',113,204);
INSERT INTO `questions` VALUES (483,360,NULL,'你对生活的最大希望是什么？','希望能看到一个更和平的世界。',_binary '\0','2018-09-30 09:14:09',521,498);
INSERT INTO `questions` VALUES (484,639,NULL,'有什么事情是你希望能得到建议的？','希望能在心理领域获得更多指导。',_binary '\0','2021-05-19 05:27:53',768,390);
INSERT INTO `questions` VALUES (485,285,NULL,'如何应对生活中的不确定性？','保持积极的心态与灵活应对。',_binary '\0','2022-07-25 17:34:07',644,531);
INSERT INTO `questions` VALUES (486,513,NULL,'你对友情的期望是什么？','期待真诚和支持的友谊。',_binary '\0','2000-11-12 04:53:06',701,497);
INSERT INTO `questions` VALUES (487,121,NULL,'描述你的理想退休生活。','与旅行和写作相关的生活。',_binary '\0','2003-06-11 04:50:56',781,714);
INSERT INTO `questions` VALUES (488,778,NULL,'如果你可以拥有任何一项才能，你希望是什么？','希望能学会演奏乐器。',_binary '\0','2000-11-01 13:06:45',46,38);
INSERT INTO `questions` VALUES (489,765,NULL,'你认为自律与自由之间的关系是什么？','自律使我更加自由地追求目标。',_binary '\0','2004-08-29 13:27:26',942,350);
INSERT INTO `questions` VALUES (490,563,NULL,'有什么事情是你希望能消失的社会问题？','类似贫困与饥饿等问题。',_binary '\0','2003-09-04 11:30:31',386,852);
INSERT INTO `questions` VALUES (491,432,NULL,'你有过的最有趣的工作是什么？','做过的最有趣的工作是学校活动组织者。',_binary '\0','2006-08-02 07:14:46',63,873);
INSERT INTO `questions` VALUES (492,905,NULL,'你怎么看待尊重与信任的关系？','尊重能建立信任，反之亦然。',_binary '\0','2013-12-26 02:55:24',188,552);
INSERT INTO `questions` VALUES (493,826,NULL,'描述你最喜欢的学习方式。','我喜欢动手实践与互动学习。',_binary '\0','2000-04-11 19:41:02',484,887);
INSERT INTO `questions` VALUES (494,427,NULL,'有什么事情是你认为能影响你人生道路的重要事件？','每个教育经历都是深远的影响。',_binary '\0','2007-10-04 01:32:56',638,130);
INSERT INTO `questions` VALUES (495,609,NULL,'你勇于尝试新事物的原因是什么？','新体验带来成长与乐趣。',_binary '\0','2012-02-28 23:11:26',720,657);
INSERT INTO `questions` VALUES (496,857,NULL,'描述你心目中的完美周末。','与朋友聚会、外出探险。',_binary '\0','2010-02-01 16:26:12',263,258);
INSERT INTO `questions` VALUES (497,497,NULL,'有什么事情让你感到特别珍贵？','家人和深厚的人际关系。',_binary '\0','2022-05-25 02:45:20',965,732);
INSERT INTO `questions` VALUES (498,206,NULL,'你认为团体活动在社会中的作用是什么？','团体活动能增强凝聚力和友谊。',_binary '\0','2002-06-16 07:37:10',94,43);
INSERT INTO `questions` VALUES (499,400,NULL,'对你来说，真正的友谊是什么样的？','真正的友谊是相互理解与支持。',_binary '\0','2023-03-13 22:11:12',370,352);
INSERT INTO `questions` VALUES (500,730,NULL,'你希望给后代留下怎样的遗产？','给后代留下爱与知识的传承。',_binary '\0','2007-03-12 17:56:48',46,153);
/*!40000 ALTER TABLE `questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teachers`
--

DROP TABLE IF EXISTS `teachers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teachers` (
  `id` int(11) NOT NULL,
  `responses` int(11) DEFAULT NULL COMMENT '回复数',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '老师名字',
  `avatar_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '老师头像链接',
  `introduction` text COMMENT '老师简介',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '老师邮箱',
  PRIMARY KEY (`id`) USING BTREE,
  CONSTRAINT `teachers_ibfk_1` FOREIGN KEY (`id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teachers`
--

LOCK TABLES `teachers` WRITE;
/*!40000 ALTER TABLE `teachers` DISABLE KEYS */;
INSERT INTO `teachers` VALUES (1,NULL,'郑子彬','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E9%83%91%E5%AD%90%E5%BD%AC_0.jpg?itok=IFxeXGlr','郑子彬，教授、博导、IEEE Fellow、IET Fellow、ACM杰出科学家、全球高被引科学家、国家优秀青年科学基金获得者、国家数字家庭工程技术研究中心副主任、广东省区块链工程技术研究中心主任。发表论文多篇，论文谷歌学术引用超过37,000次，H指数为89。主持国家重点研发计划项目、自然科学基金重点项目等多个项目；获得教育部自然科学二等奖、 吴文俊人工智能自然科学二等奖、 ACM中国新星提名奖、IEEE TCSVC Rising Star Award、CCF服务计算专委会杰出青年奖、ACM SIGSOFT Distinguished Paper Award、ICWS最佳学生论文奖等奖项；担任TOSEM、TSC、TVT、OJCS等期刊的副编辑；担任ICSS2022、IEEE SMDS2021、BlockSys2019、CollaborateCom 2016等会议的General Co-Chair；担任ICSOC2023、SC22019、教师简介: \n\n郑子彬，教授、博导、IEEE Fellow、IET Fellow、ACM杰出科学家、全球高被引科学家、国家优秀青年科学基金获得者、国家数字家庭工程技术研究中心副主任、广东省区块链工程技术研究中心主任。发表论文多篇，论文谷歌学术引用超过37,000次，H指数为89。主持国家重点研发计划项目、自然科学基金重点项目等多个项目；获得教育部自然科学二等奖、 吴文俊人工智能自然科学二等奖、 ACM中国新星提名奖、IEEE TCSVC Rising Star Award、CCF服务计算专委会杰出青年奖、ACM SIGSOFT Distinguished Paper Award、ICWS最佳学生论文奖等奖项；担任TOSEM、TSC、TVT、OJCS等期刊的副编辑；担任ICSS2022、IEEE SMDS2021、BlockSys2019、CollaborateCom 2016等会议的General Co-Chair；担任ICSOC2023、SC22019、ICIOT2018 及IoV2014等会议的PC Co-Chair。\n\n \n\n研究领域: \n\n软件可靠性，程序分析，区块链，智能合约，可信软件。\n\n本科、硕士和博士招生：常年欢迎报名。 及IoV2014等会议的PC Co-Chair。','zhzibin@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (2,NULL,'余阳','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E4%BD%99%E9%98%B33.jpg?itok=pMaiveiN','余阳，博士，教授，博士生导师。国家数字家庭工程技术研究中心副主任，软件技术广东普通高校重点实验室主任。中国计算机学会（CCF）杰出会员，ACM会员，CCF协同计算专委会执行委员，广东省高性能计算学会理事，广东省软件工程教学指导委员会副主任。','yuy@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (3,NULL,'王国利','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E7%8E%8B%E5%9B%BD%E5%88%A9.jpg?itok=u4ijqHep','王国利，工学博士，教授，德国洪堡基金获得者。\n\n \n\n邮箱：isswgl@mail.sysu.edu.cn\n\n \n\n研究领域：普适计算、泛在感知与泛在智能\n\n \n\n教育背景：1982年至1992年期间在南开大学学习，于1986年、1989年和19992年分获南开大学理学学士、理学硕士和工学博士学位。\n\n \n\n工作经历：1992年至2003年期间在汕头大学工作，期间担任计算机系系主任，1999年晋升教授。2003年11月起受聘中山大学信息科学与技术学院教授和博士导师，曾分别担任自动化系系主任、信息科学与技术学院副院长。2016年至2023年受聘计算机学院，曾任智能科学与技术研究所所长、机器智能与先进计算教育部重点实验室主任。\n\n \n\n获奖及荣誉：分别获得汕头市先进劳动者称号、教育部科学奖励自然科学二等奖和广东省科学奖励自然科学二等奖、中国人工智能学会智能空天专委会 “杰出贡献奖”。\n\n \n\n学术兼职：目前担任中国系统仿真学会智能物联系统专业委员会副主任委员、广东省机械工程学会/自动化学会机器人与智能制造分会副理事长。\n\n \n\n科研项目：入职中山大学以来，持续不间断获得国家自然科学基金项目资助6项，目前主持在研项目1项。','isswgl@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (4,NULL,'王若梅','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E7%8E%8B%E8%8B%A5%E6%A2%85.jpg?itok=POTAeAHl','王若梅，教授，博士生导师，国家数字家庭工程技术研究中心副主任。主持国家重点研发计划、国家自然科学基金、国家科技2030重大专项（中大负责人）等一批重点科研项目和课题，理论方法创新方面的成果发表在顶级的学术期刊，获得国家教育部科技进步一等奖；广州市科技进步二等奖；广东省高教厅科技进步二等奖；广东省科技进步三等奖，国家科技进步二等奖，2022年教育部科技进步二等奖。','isswrm@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (5,NULL,'周育人','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E5%91%A8%E8%82%B2%E4%BA%BA.jpg?itok=Nj2dnpGl','周育人，博士，教授，博士生导师。\n\n \n\n邮箱：zhouyuren@mail.sysu.edu.cn\n\n \n\n主要研究方向为算法设计与分析、演化算法等。\n\n \n\n教育背景：1988年本科毕业于北京大学数学专业，2003年在武汉大学获计算机软件与理论专业博士学位。2007年—2008年美国加州大学尔湾分校访问学者。\n\n \n\n在国际国内计算机科学、软件工程重要刊物“Artificial Intelligence” 、“ACM Transactions on Software Engineering and Methodology ”、 “IEEE Transactions on Evolutionary Computation”、“IEEE Transactions on Cybernetics”等发表了多篇相关学术论文， 近年来承担国家、省部级项目多项。','zhouyuren@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (6,NULL,'陈武辉','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E9%99%88%E6%AD%A6%E8%BE%892_0.jpg?itok=WMVafkgR','陈武辉，软件工程学院教授，博士生导师，珠海市可信大模型重点实验室主任。入选广东省“珠江人才计划”青年拔尖人才项目、广东省“珠江人才计划”引进创新创业团队核心成员。承担国家重点研发计划课题、国家自然科学基金面上和青年项目等多个国家级省部级项目，以及来自华为、百川智能、华东院、广州金融科技、招联消费金融等多个校企合作项目，部分成果得到落地应用。成果发表在EuroSys、NDSS、SoCC、Infocom、VLDB、ICDE、IEEE TPDS、IEEE TC、IEEE TDSC等知名会议和期刊上，获吴文俊人工智能自然科学二等奖、CCF B类会议IEEE ICPP 2020最佳论文亚军奖。出版Springer著作《Blockchain Scalability》。目前担任国际期刊 International Journal of Systems and Service-Oriented Engineering主编。入选全球前 2% 顶尖科学家榜单。','chenwuh@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (7,NULL,'苏玉鑫','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC02668_%285%29.jpg?itok=J5DfVj1v','苏玉鑫 (B站ID：鸭大坑导)，软件工程学院副院长，副教授，博士生导师，国家重点研发计划青年科学家项目负责人，中国计算机学会（CCF）高级会员、服务计算专委会执行委员。2021年7月入选中山大学百人计划，加入软件工程学院。主要研究方向为系统软件可靠性分析与运行时性能优化，具体包括操作系统、分布式系统、云计算、云原生系统、日志分析、云系统可靠性与智能运维(AIOps)等。近年来在国际会议和期刊共发表近30篇论文，其中24篇发表于ICSE、ASE、ISSTA、SOSP、FAST、ICDE、CVPR、SIGIR、AAAI、IJCAI、CSUR、TKDE等软件工程、操作系统、分布式系统、人工智能等领域CCF A类顶级会议与期刊。','suyx35@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (8,NULL,'李丹','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E6%9D%8E%E4%B8%B92.jpg?itok=Q3rT9vZ7','李丹，副教授，硕士生导师，2021年2月入选中山大学百人计划青年学术骨干，加入软件工程学院。2018年至2021年于新加坡国立大学担任研究员，从事博士后研究工作。2013年至2017年就读于新加坡南洋理工大学，受新加坡与加州大学伯克利分校联合项目资助，获得博士学位。2008年至2012年就读于电子科技大学，获得学士学位。主要从事信息物理融合系统，工业互联网，预测性维护，时序数据分析，大模型垂域应用等方面的研究。目前于IEEE TII、IEEE TASE、Energy Build. ICDE等国际著名期刊和会议上发表20余篇论文。','lidan263@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (9,NULL,'毛明志','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT05440-1%E6%AF%9B%E8%80%81%E5%B8%88.jpg?itok=z3XzuTYb','毛明志，博士，中山大学副教授，硕士生导师。中国计算机学会（CCF）杰出会员，全国高校计算机研究会理事，CCF教育专委会执行委员,《软件导刊》和《软件工程与应用》编委和审稿人，广东省科协荣誉工作者。长期从事高等学校计算机与软件工程专业相关的教学和研究工作，获国家科学技术进步二等奖、中山大学卓越服务奖和中国高校计算机教育大会优秀论文一等奖（2021）。','mcsmmz@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (10,NULL,'南雨宏','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E5%8D%97%E9%9B%A8%E5%AE%8F.jpg?itok=kXjqet3q','南雨宏，副教授，硕士生导师。中山大学百人计划青年学术骨干，校青年拔尖人才。曾任美国普渡大学(Purdue University)计算机系博士后研究员，普渡CERIAS访问学者。博士毕业于复旦大学。博士期间曾获国家留学基金委资助，于美国印第安纳大学布卢明顿分校(Indiana University Bloomington)进行联合培养。\n\n主要研究方向为系统软件安全及隐私保护。包括面向大模型、移动互联网生态、智能合约等平台的研究。研究成果发表于USENIX Security、ACM CCS，NDSS, ICSE, FSE，ASE, ISSTA等系统安全及软件工程领域顶级会议。主持国家自然科学基金、广东省新一代电子信息（半导体）重点领域专项、广东省自然科学基金（面上项目）等省部级项目，曾作为科研骨干参与国家 973 计划、美国政府及企业资助的多项研究。目前担任广东省区块链工程技术研究中心智能合约安全研发负责人，CCF珠海委员。担任ACM CCS 2024, ACM/IEEE ASE 2024, ASIACCS 2021, 2022, ICICS 2021, 2022等国际会议程序委员会委员。担任IEEE TIFS, TDSC, TOPS，TMC, TSE，EMSE等期刊审稿人。研究发现的安全及隐私问题多次获得来自Google、Meta (Facebook)、X (Twitter)、Slack、国内三大电信运营商（移动、联通、电信）等厂商的官方确认及致谢。','nanyh@sysu.edu.cn');
INSERT INTO `teachers` VALUES (11,NULL,'黄袁','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC02566.jpg?itok=S-7eWZsa','黄袁，博士/博士后，副教授，博士生导师。主持国家重点研发计划子任务，国家自然科学基金青年项目，博士后面上项目（一等资助），广东省基金面上项目等。以项目骨干身份参与国家重点研发计划，国家基金重点项目，广东省重点研发计划等多个项目。2017年于中山大学获得博士学位，从事软件工程相关研究，重点关注软件缺陷与代码智能等研究方向。近年来在IEEE Transactions on Software Engineering, ACM Transactions on Software Engineering and Methodology, IEEE Transactions on Services Computing, FSE/ESEC, ASE, ICSE等软件工程领域CCF A/B期刊及会议上发表论文40余篇。同时担任多个国际期刊和会议的审稿人。获“2024优秀硕士学位论文指导教师”奖项。','huangyuan5@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (12,NULL,'陈建国','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC07884.jpg?itok=-qdAcDF6','陈建国，副教授，硕士生导师，中山大学百人计划青年学术骨干。湖南大学计算机科学与技术博士，美国伊利诺伊大学联合培养博士，曾任加拿大多伦多大学博士后，新加坡A*STAR科技研究局研究科学家。目前主要研究方向是分布式并行计算、分布式人工智能、联邦计算、计算机视觉、图计算及应用。\n\n目前在IEEE-TII、IEEE-TITS、IEEE-TPDS、IEEE-TKDE、IEEE/ACM-TCBB、ACM-TIST、ACM-TCPS等国际著名期刊和会议上发表学术论文50余篇。主持国家自然科学基金青年项目、博士后国际交流计划派出项目、广东省自然科学基金面上项目、湖南省自然科学基金青年项目等项目。作为科研骨干先后曾参与国家高技术研究发展计划（863计划）、国家重点基础研究发展计划（ 973 计划）等多项科研课题项目。担任国际学术期刊《International Journal of Embedded Systems》副主编、《Journal of Current Scientific Research》副主编、《Information Sciences》客座编辑、《Neural Computing and Applications》客座编辑，以及多个国际学术会议的技术委员会成员。','chenjg33@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (13,NULL,'吴嘉婧','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/2769e108-9ff3-431a-8e04-8a53b7a55abe_1.jpg?itok=VIZdKWGq','吴嘉婧，中山大学副教授，博士生导师。入选广东省重大人才工程、Elsevier全球前2%顶尖科学家榜单。现任ACM-中国珠海分会副主席，广东省区块链工程技术研究中心副主任，IEEE和中国计算机学会高级会员，中国计算机学会 (CCF) 区块链专委会执行委员、服务计算专委会执行委员。\n\n2014年于香港理工大学获博士学位，随后入职中山大学计算机学院，2023年8月调入中山大学软件工程学院。研究方向包括区块链数据监管与反欺诈、Web3、智能合约、复杂网络、图挖掘等。 主持国家重点研发计划课题、子课题，国家自然科学基金面上（2项）、青年项目，和企业横向（蚂蚁集团、招联金融等）等项目十余项。已在ISSTA、WWW、IEEE Trans. Software Engineering, IEEE Trans. Information Forensics and Security, IEEE Transactions on Systems, Man, and Cybernetics: Systems等A类或一区的国际学术期刊和会议上发表论文80余篇，其中IEEE Transactions论文50余篇，ESI高被引论文4篇。在Springer 出版区块链英文学术著作 3 部。组织团队同学负责区块链数据共享平台Xblock的开发与设计 www.xblock.pro。','wujiajing@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (14,NULL,'黄华威','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT07479.jpg?itok=hRTe5A7C','黄华威，中山大学副教授，软件工程学院 “区块链与可信软件研究中心” 副主任，广东省杰青，从2023年起连续入选年度全球前 2% 科学家榜单，IEEE Senior Member，中国计算机学会 (CCF) 高级会员。研究方向包括高性能区块链系统、Web3 基础设施与协议、DeFi 协议、分布式基础设施网络 (DePIN) 等。多篇研究成果发表在 CCF A 类期刊 ToN, JSAC, TPDS, TDSC, TMC, TC 等，以及 CCF A 类国际会议 INFOCOM。论文谷歌引用 6500+, H-index 34。曾担任十余个国内外学术会议论坛研讨会的主席。主持国家重点研发计划课题、国自然面上青年项目、广东省普通高校重点领域专项、CCF-华为胡杨林基金区块链专项、鹏城实验室区块链项目等十余项科研项目。出版区块链英文学术著作 2 部《From Blockchain to Web3 & Metaverse》与《Blockchain Scalability》，出版区块链科普书《从区块链到Web3：构建下一代互联网生态》；带领团队开源区块链实验平台 BlockEmulator，该平台为区块链方向的研究生提供了成熟优质的实验平台，可快速迭代出新成果。','huanghw28@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (15,NULL,'何笑雨','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT00703wz.jpg?itok=veP1vbCt','何笑雨，副教授，硕士生导师，中山大学百人计划青年学术骨干。博士毕业于中山大学。先后在中山大学和南洋理工大学从事博士后研究工作。智能优化方法及其在人工智能和控制论中的应用，包括\n\n演化学习方法：机器学习、联邦学习、数据挖掘的演化计算求解方法\n人工智能的底层优化问题：同步/异步并行优化、面向大数据的随机优化、流形优化、多目标/多层次优化等\n智能控制方法：演化多智能体系统、强化学习的黑盒策略搜索','hexy73@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (16,NULL,'孔树锋','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT00791-1wz.jpg?itok=EY7BSECK',NULL,NULL);
INSERT INTO `teachers` VALUES (17,NULL,'蒋子规','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E8%92%8B%E5%AD%90%E8%A7%84%202023.jpg?itok=Hxnxi2nX','蒋子规，副教授，硕士生导师，中国计算机学会（CCF）服务计算专委会执行委员。2015-2019就读于北京邮电大学网络与交换技术国家重点实验室，获得博士学位，并获北京邮电大学优秀博士学位论文、北京市优秀毕业生等荣誉。2019年入职中山大学计算机学院从事博士后研究工作，2020年调入中山大学软件工程学院。主要研究方向包括区块链数据分析、智能合约、代码推荐与补全、推荐算法、大数据分析、用户行为分析等。承担包括国家自然科学基金、博士后站前特别资助、广东省自然科学基金、山东省自然科学基金智慧计算联合基金以及企业横向等项目。近年来在国内外权威学术期刊与会议上发表论文二十余篇，并担任BlockSys2023/2020/2019、CollaborateCom 2022/2016、IEEE Services 2020、IEEE SCC 2016、ICCSA 2016等国际会议的组织者、审稿人和Session Chair，以及IEEE TSC、IEEE TCC、IEEE ACCESS、ACM TSAS等多个国际期刊的审稿人。','jiangzg3@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (18,NULL,'刘名威','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E5%88%98%E5%90%8D%E5%A8%81.jpg?itok=0VBU-Dmm','刘名威，副教授，博士生导师，中山大学“逸仙学者计划”新锐学者，中国计算机学会（CCF）软件工程专委会执行委员。2024年6月加入软件工程学院，此前在复旦大学获得本科学位（2017年）、博士学位（2022年）并完成博士后研究（2024年）。\n\n \n\n研究聚焦于软件工程（SE）与人工智能（AI）交叉领域，特别是在AI4SE和SE4AI方面。我的主要兴趣在于利用先进的AI技术，如大规模语言模型（LLM）和知识图谱（KG），解决软件工程中的挑战，并应对AI应用和场景中常见的软件工程和系统工程问题。具体研究方向包括但不限于基于大模型的智能化开发与维护、软件开发知识图谱、大模型可信评测、可信代码大模型等。近五年来，在软件工程领域顶级国际期刊和会议（如TSE、TOSEM、ICSE、FSE、ASE等）发表了20余篇论文，并荣获多项殊荣，包括IEEE TCSE杰出论文奖（ICSME 2018，CCF-B）和ACM SIGSOFT杰出论文奖（FSE 2023，CCF-A）。更多信息请见个人主页：https://mingwei-liu.github.io/。','liumw26@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (19,NULL,'娄坚','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E5%A8%84%E5%9D%9A.jpg?itok=VH2cW-fS','娄坚，副教授，博士生导师。曾于美国埃默里大学(Emory University)从事博士后研究工作。主要研究方向包括可信人工智能、可信大模型、人工智能隐私保护、数据隐私保护、数据质量评估等。近年在NeurIPS、ICCV、CVPR、SIGMOD、VLDB、WWW、ACM CCS、IEEE S&P、NDSS、TDSC等人工智能、数据库、安全与隐私保护领域的顶会顶刊上发表论文60余篇，并获得顶会ACM CCS 2024杰出论文奖(ACM SIGSAC Distinguished Paper Award)，IEEE/WIC/ACM WI-IAT 2020最佳理论论文奖(Best in Theoretical Paper Award)等，研究成果曾获国际知名科技媒体New Scientist采访报道。担任人工智能顶会ICML领域主席、AAAI资深程序委员，安全与隐私保护顶会ACM CCS程序委员，数据库顶会VLDB程序委员。','louj5@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (20,NULL,'廖国成','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC02617.jpg?itok=ITQjoCgm','廖国成，助理教授，硕士生导师。2021年4月入选中山大学百人计划，加入软件工程学院。2021年于香港中文大学信息工程系获得博士学位；2019 于加州理工学院进行交流访问；2016年于中山大学电子与信息工程学院获学士学位。主要研究方向为边缘计算、算力网络、联邦学习、隐私保护、群智感知、物联网、网络优化。主持国自然青基和广东省面上项目。近年来的工作多发表于国际顶级会议和期刊，包括 CCF A 类期刊ACM/IEEE Trans. on Networking','liaogch6@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (21,NULL,'吴炜滨','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/7-1.jpg?itok=4c2lNxCH','吴炜滨，中山大学“百人计划”助理教授，硕士生导师，中国计算机学会软件工程专业委员会执行委员。2021年于香港中文大学计算机科学与工程学系获得博士学位，师从ACM/IEEE/AAAS Fellow 吕荣聪教授与IEEE Fellow 金国庆教授。2017年于同济大学电子与信息工程学院获得学士学位。主要研究方向包括可信人工智能、深度学习、计算机视觉、自然语言处理、智能软件工程等，重点关注深度学习（智能软件）的可靠性、安全性、可解释性与隐私性。主持或参与国家自然科学基金、香港研资局、深圳科创委等多个基金项目。为多个国际顶级会议和期刊的审稿人，如AAAI、ICLR、ICCV、ECCV、ACL、EMNLP、TKDE等。近年来在NeurIPS、CVPR、FSE等人工智能、计算机视觉、软件工程领域的顶级会议上发表论文16篇。','wuwb36@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (22,NULL,'陈嘉弛','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC04954.jpg?itok=D7a2XJmo','陈嘉弛，中山大学“百人计划”助理教授，硕士生导师。2022年于澳大利亚蒙纳士大学获得博士学位，主要研究方向包括软件可靠性、智能合约安全、Web3安全、大模型技术及应用、代码分析、智能合约、经验软件工程等，重点关注大模型技术在代码分析上的应用、漏洞挖掘、区块链。近5年在软件工程领域四大顶会及两大顶刊发表论文30余篇，并获得3次ACM SIGSOFT Distinguished Best Paper, 3次Best Paper Award。','chenjch86@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (23,NULL,'陈文清','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT09654_%282%29.jpg?itok=f4jGQIHQ','陈文清，中山大学“百人计划”助理教授，硕士研究生导师，博士毕业于上海交通大学电子信息与电气工程学院、人工智能研究院，师从金耀辉教授。主要研究方向包括自然语言处理、大语言模型、代码智能、智慧医疗、智慧司法等，目前研究兴趣集中在大语言模型的复杂推理。主持纵向项目包括国自然青年基金、广东省面上项目、人工智能教育部重点实验室开放课题；主持两项大模型横向课题；参与国家重点研发计划、上海市级科技重大专项资助、上海市人工智能创新发展专项资助、上海市科技计划等多个基金项目。在人工智能领域、自然语言处理国际权威期刊及会议发表论文20+篇，获国家发明专利授权6项。担任多个国际顶级会议和期刊的审稿人，如ACL、EMNLP、IJCAI、AAAI、 Expert Syst App.等。','chenwq95@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (24,NULL,'王焱林','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT09685.jpg?itok=xt70AOOM','王焱林，助理教授，硕士生导师。2022年7月入选中山大学百人计划，加入软件工程学院。加入中山大学前，于微软亚洲研究院担任主管研究员。2014年至2019年就读于香港大学计算机科学系，师从Bruno Oliveira教授，获得博士学位。2010年至2014年就读于浙江大学，获学士学位。主要研究方向为大模型驱动的智能化软件工程、大模型技术。近5年来在国际会议和期刊共发表30余篇论文，发表于ICSE、ASE、AAAI、ACL、KDD、TKDE、EMSE、EMNLP、CIKM、ICSME等软件工程、人工智能、自然语言处理等领域CCF A/B类顶级会议与期刊。','wangylin36@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (25,NULL,'陈壮彬','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/TTT08695.jpg?itok=bEc7TfeI','陈壮彬，助理教授，硕士生导师，中国计算机学会（CCF）服务计算专委会执行委员。2023年1月入选中山大学百人计划，加入软件工程学院。2018年至2022年，于香港中文大学计算机科学与工程系攻读博士学位，师从ACM/IEEE/AAAS Fellow吕荣聪（Michael R. Lyu）教授。主攻研究领域为大型智能软件系统，包括云计算、机器学习、软件可靠性、智能运维、数据中心网络等方面。近年来，在ICSE、FSE、ASE、DSN、CIKM、ACM Computing Surveys、Information Sciences等软件工程、数据挖掘、计算机网络领域的CCF A类顶级会议和期刊发表论文20余篇。','chenzhb36@mail.sysu.edu.cn');
INSERT INTO `teachers` VALUES (26,NULL,'潘茂林','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/%E6%BD%98%E8%8C%82%E6%9E%97.jpg?itok=xCxLfA2r',NULL,NULL);
INSERT INTO `teachers` VALUES (27,NULL,'王乐球','https://sse.sysu.edu.cn/sites/default/files/styles/image_style_2/public/DSC02832%20%E7%8E%8B%E8%80%81%E5%B8%88%E5%B7%B2%E6%94%B9%E6%AF%94%E4%BE%8B%E5%86%85%E5%AD%98.jpg?itok=z6oERpe_','王乐球，硕士，中山大学软件工程学院专任教员。\n\n教育背景：\n1987-1991，浙江大学计算机软件学士\n1991-1994浙江大学计算机应用硕士\n\n工作情况：\n1994-1999珠海远方电脑公司从事软件开发\n1999-2001珠海电力工业局工作\n2001-2020中山大学资讯管理学院\n2020-至今，软件工程学院','isswlq@mail.sysu.edu.cn');
/*!40000 ALTER TABLE `teachers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `upvotes`
--

DROP TABLE IF EXISTS `upvotes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `upvotes` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `question_id` int(11) DEFAULT NULL COMMENT '问题ID',
  `answer_id` int(11) DEFAULT NULL COMMENT '回复ID',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `upvotes_ibfk_1` (`user_id`) USING BTREE,
  KEY `upvotes_ibfk_2` (`question_id`) USING BTREE,
  KEY `upvotes_ibfk_3` (`answer_id`) USING BTREE,
  CONSTRAINT `upvotes_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `upvotes_ibfk_2` FOREIGN KEY (`question_id`) REFERENCES `questions` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `upvotes_ibfk_3` FOREIGN KEY (`answer_id`) REFERENCES `answers` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `upvotes_chk_1` CHECK ((((`question_id` is not null) + (`answer_id` is not null)) = 1))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `upvotes`
--

LOCK TABLES `upvotes` WRITE;
/*!40000 ALTER TABLE `upvotes` DISABLE KEYS */;
/*!40000 ALTER TABLE `upvotes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '用户名',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '邮箱',
  `salt` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '加密盐',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '密码哈希，算法暂定为Argon2id',
  `role` enum('admin','teacher','student') CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '角色',
  `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '昵称',
  `introduction` text CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs NOT NULL COMMENT '简介',
  `avatar_file_id` int(11) DEFAULT NULL COMMENT '头像文件ID，为空时为配置的默认头像',
  `theme_id` int(11) DEFAULT NULL COMMENT '主题ID，为空时为配置的默认主题',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `name` (`name`) USING BTREE COMMENT '用户名唯一',
  UNIQUE KEY `email` (`email`) USING BTREE COMMENT '邮箱唯一',
  KEY `avatar_file_id` (`avatar_file_id`) USING BTREE,
  KEY `theme_id` (`theme_id`) USING BTREE,
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`avatar_file_id`) REFERENCES `files` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_zh_0900_as_cs ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'Fujiwara Ryota','ryotafujiwara4@gmail.com','ys3C2w7pfs','hqNfwrMKp8','admin','Fujiwara Ryota','追求自己，超越他人。',NULL,189,'2002-04-06 11:18:14',NULL,NULL);
INSERT INTO `users` VALUES (2,'Emma Payne','emma9@outlook.com','sdG5opQODx','EkaqXZq4EU','teacher','Emma Payne','随心而动，自由自在。',NULL,331,'2021-01-09 07:48:50',NULL,NULL);
INSERT INTO `users` VALUES (3,'Suzuki Yuito','suzukiy@icloud.com','jGUoTmU6Ni','95PjBnFV7l','teacher','Suzuki Yuito','每一步都走得坚定。',NULL,218,'2013-02-05 03:10:56',NULL,NULL);
INSERT INTO `users` VALUES (4,'Kong Ka Ling','kongkali@yahoo.com','GfOHioczUa','apTQL6L87t','student','Kong Ka Ling','做自己，世界会配合你。',NULL,499,'2013-08-05 17:56:30',NULL,NULL);
INSERT INTO `users` VALUES (5,'Maruyama Seiko','mseiko@icloud.com','HNAvFpaPmU','kTrSITjLtC','teacher','Maruyama Seiko','生活没有彩排，活出精彩。',NULL,637,'2019-08-20 14:37:01',NULL,NULL);
INSERT INTO `users` VALUES (6,'Alexander Clark','alexanderclark4@gmail.com','vsfzmvl4dy','UdzjfAJ8wY','teacher','Alexander Clark','梦想不设限，未来由我。',NULL,703,'2008-07-25 16:13:10',NULL,NULL);
INSERT INTO `users` VALUES (7,'Duan Zitao','dzi@outlook.com','tpApE24Udk','yx0iO63t3S','admin','Duan Zitao','不惧风雨，勇敢前行。',NULL,406,'2021-12-02 08:14:15',NULL,NULL);
INSERT INTO `users` VALUES (8,'Tsui Sum Wing','tsuisw@gmail.com','6xRer8EYlE','hHFgXWvQCB','teacher','Tsui Sum Wing','努力就会有回报。',NULL,593,'2003-12-03 13:24:40',NULL,NULL);
INSERT INTO `users` VALUES (9,'Deng Zitao','deng6@icloud.com','aMqthjTLWL','BKJCVF91vd','student','Deng Zitao','活得像自己，才不会后悔。',NULL,352,'2007-09-21 02:38:18',NULL,NULL);
INSERT INTO `users` VALUES (10,'Ando Riku','anriku@icloud.com','1VGWWjDWeb','pRrlxyar6Z','admin','Ando Riku','放下过去，迎接未来。',NULL,936,'2016-10-29 14:35:59',NULL,NULL);
INSERT INTO `users` VALUES (11,'Saito Rin','saitri3@yahoo.com','0O0xJVGotG','Vw4HGI0Uen','student','Saito Rin','心若向阳，无畏悲伤。',NULL,832,'2005-09-23 08:20:32',NULL,NULL);
INSERT INTO `users` VALUES (12,'Dong Ziyi','dong6@hotmail.com','Qn9wBvXFkX','ZuPcOM6H9e','admin','Dong Ziyi','每一个今天，都是未来的回忆。',NULL,524,'2020-04-09 20:39:37',NULL,NULL);
INSERT INTO `users` VALUES (13,'Yamaguchi Sakura','yasa@gmail.com','BFHdd1U4Nj','SpOhE8lqgg','student','Yamaguchi Sakura','只为你，我愿改变一切。',NULL,167,'2012-07-15 06:56:21',NULL,NULL);
INSERT INTO `users` VALUES (14,'Ono Seiko','seikoono10@mail.com','q8FOaTxjdN','5Y110z65nD','student','Ono Seiko','笑对生活，生活也会笑对你。',NULL,219,'2008-11-10 14:13:51',NULL,NULL);
INSERT INTO `users` VALUES (15,'Dennis Herrera','hede@mail.com','DEOZBzu5jA','RvTxBl6xLL','student','Dennis Herrera','做个有温度的人。',NULL,560,'2008-09-11 05:09:25',NULL,NULL);
INSERT INTO `users` VALUES (16,'Tin Sze Yu','szeyutin@gmail.com','de997HxwuZ','OJdJ9GpGh3','admin','Tin Sze Yu','逆风的方向，更适合飞翔。',NULL,408,'2024-02-20 13:39:06',NULL,NULL);
INSERT INTO `users` VALUES (17,'Lau Chung Yin','lau211@outlook.com','VVW9P9qPdV','s5UK3KJYps','admin','Lau Chung Yin','只有自己，才是最好的依靠。',NULL,788,'2007-09-17 02:48:53',NULL,NULL);
INSERT INTO `users` VALUES (18,'Lau Chiu Wai','lau2@yahoo.com','fWxyEbwxo2','B4yXe6UHvG','teacher','Lau Chiu Wai','过去已成烟，未来更精彩。',NULL,722,'2019-05-08 10:26:21',NULL,NULL);
INSERT INTO `users` VALUES (19,'Dorothy Stewart','dorothy56@gmail.com','Be6RlbhkCP','l0ZiWLgHiK','teacher','Dorothy Stewart','没有最好，只有更好。',NULL,426,'2016-02-28 08:00:41',NULL,NULL);
INSERT INTO `users` VALUES (20,'Fu Wai Man','fuwaiman2@gmail.com','hS9htVUxvn','V0PXj0PBBA','teacher','Fu Wai Man','随风而行，心无旁骛。',NULL,122,'2011-03-03 06:50:22',NULL,NULL);
INSERT INTO `users` VALUES (21,'Annie Jenkins','jenkiannie@mail.com','2nzgzuzmKO','yVc0DlLgdY','student','Annie Jenkins','未来可期，人生可贵。',NULL,325,'2000-06-24 07:51:49',NULL,NULL);
INSERT INTO `users` VALUES (22,'Yamazaki Akina','yakina7@icloud.com','xjB7ODDxvP','5l6FPdXyL9','teacher','Yamazaki Akina','我的世界，自己做主。',NULL,800,'2019-03-31 20:43:39',NULL,NULL);
INSERT INTO `users` VALUES (23,'Kato Shino','skato98@yahoo.com','QpKTxRvGkz','wMAUoZOcuc','admin','Kato Shino','只因梦想，才如此努力。',NULL,43,'2008-01-31 03:31:10',NULL,NULL);
INSERT INTO `users` VALUES (24,'Raymond Collins','raymocol@gmail.com','YEWTxTSQXW','pcUEd7j6wx','admin','Raymond Collins','平凡不等于平庸。',NULL,201,'2000-06-10 04:27:11',NULL,NULL);
INSERT INTO `users` VALUES (25,'Ito Sakura','itosakura@mail.com','1LiqRQWYE1','aPt3Ur148H','student','Ito Sakura','用微笑战胜一切困难。',NULL,883,'2017-02-05 01:19:27',NULL,NULL);
INSERT INTO `users` VALUES (26,'Tammy Griffin','griffin330@outlook.com','zVrSlooHKp','v9zm7TPHEL','teacher','Tammy Griffin','越努力，越幸运。',NULL,413,'2016-05-24 03:46:20',NULL,NULL);
INSERT INTO `users` VALUES (27,'Cho Wai Man','waimanc@yahoo.com','jpHbv9vKRv','nPYOKlplgd','student','Cho Wai Man','轻轻放下，快乐自由。',NULL,473,'2024-01-04 00:02:13',NULL,NULL);
INSERT INTO `users` VALUES (28,'Yuen Siu Wai','swyuen@gmail.com','Fe2Xmt82SH','93x6vMZpC0','admin','Yuen Siu Wai','每一个笑容背后都有故事。',NULL,698,'2008-06-11 13:26:14',NULL,NULL);
INSERT INTO `users` VALUES (29,'Koyama Ryota','koyama59@mail.com','xc7bcexNe5','ZwEd47jdZH','student','Koyama Ryota','成为你想要的那个人。',NULL,218,'2005-09-26 06:05:49',NULL,NULL);
INSERT INTO `users` VALUES (30,'Siu Chi Yuen','chiyuen20@icloud.com','IYr7NJxihk','nA0HI9TgNU','student','Siu Chi Yuen','平淡是福，简单是美。',NULL,696,'2020-06-17 14:13:59',NULL,NULL);
INSERT INTO `users` VALUES (31,'Fung Ling Ling','llf@icloud.com','nJwuFXQhgH','69KVx9IQ1a','teacher','Fung Ling Ling','没有梦就没有未来。',NULL,271,'2017-11-20 21:34:02',NULL,NULL);
INSERT INTO `users` VALUES (32,'Ho Tsz Hin','tszhin227@mail.com','PWjjMmKsfu','rOnk6fEKmj','teacher','Ho Tsz Hin','不畏将来，不念过往。',NULL,566,'2012-12-27 12:21:34',NULL,NULL);
INSERT INTO `users` VALUES (33,'Sit Hui Mei','sihuimei@outlook.com','tJ1XYbnW1H','9CbH6es4jZ','teacher','Sit Hui Mei','为梦而生，为爱而活。',NULL,897,'2007-03-01 17:59:14',NULL,NULL);
INSERT INTO `users` VALUES (34,'Hou Lu','luho@hotmail.com','0S4UQuKR6G','tq1EsLxOy2','teacher','Hou Lu','你若安好，便是晴天。',NULL,135,'2021-03-19 15:31:34',NULL,NULL);
INSERT INTO `users` VALUES (35,'Thelma Rodriguez','rthelma708@gmail.com','dVv7snyAOY','cO7TpYbVEA','teacher','Thelma Rodriguez','用心生活，快乐每一天。',NULL,285,'2024-04-07 11:37:16',NULL,NULL);
INSERT INTO `users` VALUES (36,'Tin Fat','fattin9@yahoo.com','aVfniL3qKf','kV6kPTONft','admin','Tin Fat','没有永远的朋友，只有永远的自己。',NULL,397,'2001-07-12 15:22:25',NULL,NULL);
INSERT INTO `users` VALUES (37,'Wan Wing Fat','wanwingfat8@outlook.com','CUl37zbOKU','AKCulF5iD7','student','Wan Wing Fat','做你自己，别人没资格批判。',NULL,564,'2011-07-31 10:31:56',NULL,NULL);
INSERT INTO `users` VALUES (38,'Carol Scott','scottcarol@outlook.com','43vnzFPoEH','kMSLkpDyXN','admin','Carol Scott','未来一片光明，照亮前行。',NULL,714,'2024-07-27 07:14:06',NULL,NULL);
INSERT INTO `users` VALUES (39,'Qian Zhennan','zhennanq@gmail.com','4MvYTxRz62','5gi1Oz3khU','teacher','Qian Zhennan','青春不再，奋斗未停。',NULL,373,'2020-12-15 09:01:30',NULL,NULL);
INSERT INTO `users` VALUES (40,'Mo Lu','molu@icloud.com','62VNIxP74G','o4f04uyUXq','admin','Mo Lu','向阳而生，不畏阴霾。',NULL,160,'2004-09-11 17:11:28',NULL,NULL);
INSERT INTO `users` VALUES (41,'Ando Kenta','andokent@outlook.com','kmRnq2EqsK','mJdHQ16DLc','student','Ando Kenta','不求与人相比，只求超越自己。',NULL,289,'2022-03-30 20:51:36',NULL,NULL);
INSERT INTO `users` VALUES (42,'Iwasaki Eita','iwaeita@icloud.com','8hoR2Q24su','uz66f17nc7','teacher','Iwasaki Eita','喜欢自己的样子，爱上孤独的美。',NULL,726,'2013-01-31 01:51:34',NULL,NULL);
INSERT INTO `users` VALUES (43,'Du Lan','ldu10@icloud.com','1hq5pWOCr9','4hIGOZFFDT','admin','Du Lan','做一个温暖的人。',NULL,799,'2006-05-03 08:58:43',NULL,NULL);
INSERT INTO `users` VALUES (44,'Shen Anqi','anqishen@gmail.com','nsRuURABaf','UOclCur9OA','teacher','Shen Anqi','我的原则：做自己，做真我。',NULL,323,'2019-08-04 15:20:23',NULL,NULL);
INSERT INTO `users` VALUES (45,'Lee Wai Han','waihanle1943@hotmail.com','GQHr4rNTyC','hwaY3FUqYc','admin','Lee Wai Han','人生没有如果，只有结果。',NULL,615,'2006-10-03 02:14:02',NULL,NULL);
INSERT INTO `users` VALUES (46,'Wu Kwok Wing','wu409@icloud.com','eggi8CBWnN','6UB2DggSvv','teacher','Wu Kwok Wing','不在乎结果，只要尽力。',NULL,371,'2019-11-19 10:23:17',NULL,NULL);
INSERT INTO `users` VALUES (47,'Jesse Sanders','jsan@outlook.com','00hYDJRL4Z','tPZ63PTgnS','teacher','Jesse Sanders','我有我的坚持，你有你的选择。',NULL,998,'2007-11-20 12:28:33',NULL,NULL);
INSERT INTO `users` VALUES (48,'Han Xiuying','xiuyingh716@icloud.com','oaQHH2bmKH','u7BKdyjOjX','student','Han Xiuying','世界那么大，我只想与你并肩。',NULL,297,'2004-08-29 22:55:23',NULL,NULL);
INSERT INTO `users` VALUES (49,'Francis Sanders','sandef@hotmail.com','NPr8zhgkTn','r9em7HJZak','teacher','Francis Sanders','生活不需要太多解释。',NULL,760,'2020-03-15 04:42:05',NULL,NULL);
INSERT INTO `users` VALUES (50,'Tsui Fat','tfat306@yahoo.com','jTKvLvumU6','tfDFQ3Dbas','student','Tsui Fat','永远年轻，永远热泪盈眶。',NULL,659,'2001-09-23 14:58:10',NULL,NULL);
INSERT INTO `users` VALUES (51,'Du Rui','rudu@gmail.com','Okw4rUrzTy','OxXR9ssIhn','admin','Du Rui','生活是最好的老师。',NULL,69,'2007-03-26 06:32:11',NULL,NULL);
INSERT INTO `users` VALUES (52,'Watanabe Akina','watanabe1980@icloud.com','CJxNPcvmuI','WSGI5jFJeT','admin','Watanabe Akina','世界这么大，不该只为一个人停留。',NULL,551,'2014-10-23 03:43:54',NULL,NULL);
INSERT INTO `users` VALUES (53,'Takada Eita','takada2@mail.com','uhbexIO5ge','H3OICR2ERj','admin','Takada Eita','走自己的路，让别人说去吧。',NULL,365,'2015-08-16 10:13:56',NULL,NULL);
INSERT INTO `users` VALUES (54,'Ono Mio','onom@yahoo.com','qrtJPiTYtM','CGTunYQfsM','teacher','Ono Mio','你的每一步，都是人生的精彩。',NULL,469,'2007-06-15 16:28:35',NULL,NULL);
INSERT INTO `users` VALUES (55,'Chung Wing Sze','chungws@icloud.com','dDNlnZdn8h','sU3FcUor1u','student','Chung Wing Sze','用心去感受，用爱去温暖。',NULL,403,'2018-05-30 16:23:23',NULL,NULL);
INSERT INTO `users` VALUES (56,'Shimada Yuna','yuns@hotmail.com','ym3ISsQiKL','1a6ogvmInM','student','Shimada Yuna','用奋斗书写未来。',NULL,409,'2005-02-15 18:43:50',NULL,NULL);
INSERT INTO `users` VALUES (57,'Chiba Yuna','yunachi@yahoo.com','ZkkAgj9VhS','JVBWzLxRZA','admin','Chiba Yuna','每一天都是新的开始。',NULL,190,'2007-08-22 23:36:48',NULL,NULL);
INSERT INTO `users` VALUES (58,'To Tsz Hin','ttszhin82@yahoo.com','LIUL1UOBvO','mHhtMXFDaM','teacher','To Tsz Hin','不言放弃，坚持到底。',NULL,832,'2002-07-04 19:34:54',NULL,NULL);
INSERT INTO `users` VALUES (59,'Duan Yuning','duan2@mail.com','G84Xg96sE8','Q3ZHitgtel','admin','Duan Yuning','远方的梦想，照亮我的脚步。',NULL,741,'2022-03-09 23:57:34',NULL,NULL);
INSERT INTO `users` VALUES (60,'Charlotte Fernandez','fernandez58@outlook.com','DS7TcndasE','GsR7uGmFH3','admin','Charlotte Fernandez','不畏困难，勇往直前。',NULL,852,'2010-10-19 08:09:11',NULL,NULL);
INSERT INTO `users` VALUES (61,'Feng Lan','fengla@mail.com','V25OOGH416','Lhrr2kNX3p','admin','Feng Lan','与其诋毁别人，不如完善自己。',NULL,637,'2005-10-30 01:30:14',NULL,NULL);
INSERT INTO `users` VALUES (62,'Kato Takuya','takuyakato@icloud.com','RPrTR7VnV9','2ZLKgnkn64','student','Kato Takuya','做一个温暖的人，温暖世界。',NULL,92,'2022-01-25 17:42:21',NULL,NULL);
INSERT INTO `users` VALUES (63,'Huang Jialun','jhu3@yahoo.com','t7gCTpOOYL','RuV3fhxPaB','student','Huang Jialun','向阳而生，随风而舞。',NULL,3,'2015-01-10 00:31:36',NULL,NULL);
INSERT INTO `users` VALUES (64,'Hui Wing Suen','wingsuenh@outlook.com','zRrsjrvl3X','Api3yQhmpH','admin','Hui Wing Suen','过去的已成历史，未来还需努力。',NULL,954,'2019-01-13 11:55:24',NULL,NULL);
INSERT INTO `users` VALUES (65,'Nakano Sakura','nsakura2017@gmail.com','yLhcgiesJx','FXgDUkBuSF','admin','Nakano Sakura','人生苦短，何不放肆一回。',NULL,536,'2020-12-03 18:26:12',NULL,NULL);
INSERT INTO `users` VALUES (66,'Jiang Yunxi','jiayunxi909@icloud.com','X6dtxkJbGC','qBeXhO8MvR','admin','Jiang Yunxi','生活的美好在于不完美。',NULL,504,'2022-05-01 18:47:21',NULL,NULL);
INSERT INTO `users` VALUES (67,'Tian Xiuying','tianxiuying@outlook.com','SwKZOpMRxl','7pUbpKA6Tt','admin','Tian Xiuying','未来总是掌握在自己手中。',NULL,756,'2023-03-24 04:25:34',NULL,NULL);
INSERT INTO `users` VALUES (68,'Gary Crawford','cragary3@outlook.com','N9NWCcQBz7','JNr7SbCRIy','student','Gary Crawford','生活本无意义，唯有你我共同赋予。',NULL,444,'2016-12-20 19:09:02',NULL,NULL);
INSERT INTO `users` VALUES (69,'Matsui Yuna','myuna@mail.com','D05KogsSlD','v9g59OZlb2','student','Matsui Yuna','不羡慕别人，只做更好的自己。',NULL,626,'2013-01-20 11:33:14',NULL,NULL);
INSERT INTO `users` VALUES (70,'Lisa Mcdonald','mclisa1962@outlook.com','9VYkV5xFOd','oH19wqPCVY','admin','Lisa Mcdonald','没有最好，只有更好。',NULL,257,'2000-04-21 09:18:43',NULL,NULL);
INSERT INTO `users` VALUES (71,'Yin Yuning','yyuning67@mail.com','UCW8I7sjel','TTsWlrbQj8','student','Yin Yuning','做自己的英雄，过自己想要的生活。',NULL,628,'2021-01-29 20:24:28',NULL,NULL);
INSERT INTO `users` VALUES (72,'Matsuda Ayano','ayano8@gmail.com','YyeqZWdb6Y','TigW5AMBgE','teacher','Matsuda Ayano','努力的意义，终会看到成果。',NULL,920,'2023-09-03 14:02:38',NULL,NULL);
INSERT INTO `users` VALUES (73,'Man Wing Sze','wingszem@gmail.com','T05v65jJ85','MNd77FW6r4','teacher','Man Wing Sze','青春，就是用来冒险的。',NULL,921,'2005-09-01 18:44:08',NULL,NULL);
INSERT INTO `users` VALUES (74,'Luo Jialun','jialluo@gmail.com','Y7oMBJVSsW','6u9GKAVYSq','admin','Luo Jialun','未来是自己的，不要留给别人。',NULL,194,'2010-02-03 15:20:52',NULL,NULL);
INSERT INTO `users` VALUES (75,'Tin Ling Ling','tin70@gmail.com','pIBqyRyZqM','1HHfokb171','teacher','Tin Ling Ling','世界上最远的距离，是你在我身边，我却不能靠近。',NULL,227,'2024-08-06 15:48:37',NULL,NULL);
INSERT INTO `users` VALUES (76,'Fu Ming Sze','fumin@mail.com','wjTF4l4Q0z','KDCkxQxZlB','admin','Fu Ming Sze','梦想，是人生最美的期待。',NULL,706,'2017-08-04 12:58:10',NULL,NULL);
INSERT INTO `users` VALUES (77,'Wan Kar Yan','wkaryan@hotmail.com','S3fCHbjN1b','PjHx941BR0','teacher','Wan Kar Yan','远方的诗和田野，都是我的归属。',NULL,414,'2012-08-19 09:33:34',NULL,NULL);
INSERT INTO `users` VALUES (78,'Fu Tak Wah','takwahfu3@outlook.com','eFdJjcZooP','G2kzI4SFzT','student','Fu Tak Wah','不怕路远，只怕心不坚。',NULL,628,'2002-09-14 00:13:10',NULL,NULL);
INSERT INTO `users` VALUES (79,'Koon Ka Fai','kafaikoon58@gmail.com','iWFN1KlJwo','wW0e0fIruz','teacher','Koon Ka Fai','有梦想的人生才有意义。',NULL,815,'2011-07-30 06:15:50',NULL,NULL);
INSERT INTO `users` VALUES (80,'Hou Xiaoming','houx93@hotmail.com','F0VI16TDzw','NexHcN82ZG','student','Hou Xiaoming','只为真心，不为虚伪。',NULL,886,'2015-09-20 02:39:38',NULL,NULL);
INSERT INTO `users` VALUES (81,'Dong Shihan','shihadong1984@icloud.com','eAXdfgiP4i','HQKV5PsLUk','teacher','Dong Shihan','活得自在，做个有趣的人。',NULL,26,'2003-06-14 13:24:12',NULL,NULL);
INSERT INTO `users` VALUES (82,'Sato Yuito','yuitos1969@hotmail.com','KL9Kc8vuyB','qk181cfhNV','teacher','Sato Yuito','一切皆有可能，心态决定未来。',NULL,84,'2010-05-16 22:49:30',NULL,NULL);
INSERT INTO `users` VALUES (83,'Matsuda Yota','yota8@gmail.com','oeR0qmtxpY','gdKkMxRUP0','teacher','Matsuda Yota','永远记住，不管走得多远，都不要忘记自己。',NULL,415,'2007-12-01 13:47:16',NULL,NULL);
INSERT INTO `users` VALUES (84,'Dai On Na','onnadai914@gmail.com','TynZnMKQ9g','PCvMxbzOvu','student','Dai On Na','人生如戏，我是主角。',NULL,611,'2023-10-24 05:35:55',NULL,NULL);
INSERT INTO `users` VALUES (85,'Nakagawa Mitsuki','nmitsuki@gmail.com','GfMSDPLKRJ','efSiP2QwEe','admin','Nakagawa Mitsuki','勇敢追梦，永不回头。',NULL,343,'2002-01-25 07:56:43',NULL,NULL);
INSERT INTO `users` VALUES (86,'Samuel Chavez','samuelc4@mail.com','TEBGGsDJcQ','CnSHHw5zWL','student','Samuel Chavez','微笑面对一切，勇敢做自己。',NULL,596,'2021-08-10 09:34:20',NULL,NULL);
INSERT INTO `users` VALUES (87,'Rose Stewart','stewartrose1945@icloud.com','KPBlGZTVqF','wvr2Bk4dfa','teacher','Rose Stewart','梦想就在前方，努力就是通行证。',NULL,323,'2002-12-22 02:29:42',NULL,NULL);
INSERT INTO `users` VALUES (88,'Keith Watson','kwa@outlook.com','rNeTIGoIy2','2YbzWEbVQX','admin','Keith Watson','你的坚持，终将成为你的人生底色。',NULL,815,'2005-07-11 01:09:57',NULL,NULL);
INSERT INTO `users` VALUES (89,'Wong Ka Ming','kmwo@icloud.com','4BhklPAXOM','Jd2s8c3soh','teacher','Wong Ka Ming','世界因你而美丽。',NULL,192,'2005-09-16 19:50:13',NULL,NULL);
INSERT INTO `users` VALUES (90,'Wanda Bennett','wandab1959@gmail.com','DD54KTnota','zBK7UZuUPD','teacher','Wanda Bennett','路是自己走出来的。',NULL,707,'2003-06-13 03:04:21',NULL,NULL);
INSERT INTO `users` VALUES (91,'Ono Aoshi','ono9@outlook.com','gttSGdBuxK','7bVEzWpnBj','admin','Ono Aoshi','即使生活不如意，依然要心怀感恩。',NULL,492,'2012-03-27 13:45:27',NULL,NULL);
INSERT INTO `users` VALUES (92,'Sato Sara','sarasato@gmail.com','9GiuIKmFkj','DAXjGOkRXP','admin','Sato Sara','做自己，不为谁而改变。',NULL,733,'2013-03-20 08:09:02',NULL,NULL);
INSERT INTO `users` VALUES (93,'Luis Robinson','luirobinson@outlook.com','kG17cl1czV','1Lioj1hBHP','student','Luis Robinson','生活本来就是一场冒险。',NULL,866,'2019-04-07 12:11:41',NULL,NULL);
INSERT INTO `users` VALUES (94,'Kathryn Hayes','kathrynha@hotmail.com','XNtXBkPOnM','8w6G40mtWW','student','Kathryn Hayes','阳光总在风雨后。',NULL,675,'2024-09-11 15:44:36',NULL,NULL);
INSERT INTO `users` VALUES (95,'Fong Ching Wan','chifong@gmail.com','5TidH8UrH8','JNjWPKYodu','admin','Fong Ching Wan','没有任何人能左右我的心情。',NULL,368,'2016-03-12 07:32:45',NULL,NULL);
INSERT INTO `users` VALUES (96,'David Roberts','roberts55@mail.com','sQA9dwicl0','pbKWkZcc0A','student','David Roberts','你有你的精彩，我有我的独特。',NULL,35,'2017-11-26 12:24:14',NULL,NULL);
INSERT INTO `users` VALUES (97,'Ueno Seiko','uenosei@gmail.com','hk5O6II4Gf','cXIZMVMLys','teacher','Ueno Seiko','每一份努力都会有回报。',NULL,624,'2002-05-24 14:49:36',NULL,NULL);
INSERT INTO `users` VALUES (98,'Fu Zhennan','fuzhennan1956@outlook.com','bzAhQgDVPx','HV2Acwm9yc','admin','Fu Zhennan','生命是自己拼搏出来的。',NULL,606,'2007-05-03 13:53:26',NULL,NULL);
INSERT INTO `users` VALUES (99,'Choi Chieh Lun','chchiehlun@mail.com','L3BjSwADsg','KLvPZEeDMP','student','Choi Chieh Lun','幸福就在身边，从未远离。',NULL,71,'2024-09-15 12:47:56',NULL,NULL);
INSERT INTO `users` VALUES (100,'Chow Chung Yin','chowcy5@gmail.com','1KQdtCKEY7','1YcAyu1oM7','student','Chow Chung Yin','每天都是新的起点。',NULL,175,'2010-07-03 16:31:16',NULL,NULL);
INSERT INTO `users` VALUES (101,'Amy Butler','amybutler1@hotmail.com','E0nrf2yBzy','VNjcPHvMPt','student','Amy Butler','勇敢去追，错过了不后悔。',NULL,84,'2018-06-22 20:15:15',NULL,NULL);
INSERT INTO `users` VALUES (102,'Zhao Zitao','zitaozhao@yahoo.com','RVzJaOHr3d','dEPZgobjyk','teacher','Zhao Zitao','人生苦短，不要等到失去才懂得珍惜。',NULL,437,'2015-09-29 07:35:39',NULL,NULL);
INSERT INTO `users` VALUES (103,'Teresa Watson','water@icloud.com','Sdk1JvdI76','596LaMOz6s','teacher','Teresa Watson','心若澄明，万物皆美。',NULL,222,'2015-02-18 17:03:15',NULL,NULL);
INSERT INTO `users` VALUES (104,'Dennis Moreno','dmor@icloud.com','e47dnnfsrk','StS8Dp0e0A','student','Dennis Moreno','勇敢一点，梦想更近一点。',NULL,666,'2017-11-19 21:26:47',NULL,NULL);
INSERT INTO `users` VALUES (105,'Feng Zitao','fengzita@yahoo.com','G2SC7VJcRj','zCDYVhvxkt','admin','Feng Zitao','爱自己，从此不再妥协。',NULL,362,'2004-07-19 04:36:09',NULL,NULL);
INSERT INTO `users` VALUES (106,'Xu Shihan','xus8@mail.com','gCUDvTkXKs','pfXhaVyD27','teacher','Xu Shihan','未来可期，万象更新。',NULL,803,'2013-05-15 23:56:08',NULL,NULL);
INSERT INTO `users` VALUES (107,'Matsui Hazuki','hazmatsui@gmail.com','1RgvitgsDl','n6ONvdj4LP','admin','Matsui Hazuki','为了梦想，不畏艰难险阻。',NULL,834,'2019-01-09 10:39:02',NULL,NULL);
INSERT INTO `users` VALUES (108,'Wong Chi Yuen','wong7@icloud.com','z9lq7yhFSg','BnHe4e59XK','teacher','Wong Chi Yuen','平凡中见伟大，坚韧成就辉煌。',NULL,894,'2017-05-19 02:13:23',NULL,NULL);
INSERT INTO `users` VALUES (109,'Sandra James','james2@icloud.com','DHKe34goO5','Jt4VSTeUgN','admin','Sandra James','永远做一个快乐的自己。',NULL,670,'2016-01-10 06:01:11',NULL,NULL);
INSERT INTO `users` VALUES (110,'Ma Ka Ling','klm217@yahoo.com','SGwJmeBuKe','gufqGNJIfx','teacher','Ma Ka Ling','给自己一个微笑，生活因此不同。',NULL,300,'2001-10-01 02:09:46',NULL,NULL);
INSERT INTO `users` VALUES (111,'Lucille Sullivan','sullivan8@mail.com','YKiRPL0y4H','wBFCuNoi8K','student','Lucille Sullivan','梦想的道路上，不怕风雨。',NULL,941,'2022-08-31 19:43:24',NULL,NULL);
INSERT INTO `users` VALUES (112,'Lam Chieh Lun','cllam7@mail.com','yKUsywelur','Cq1UEiHxHf','student','Lam Chieh Lun','让每一天都充满希望。',NULL,583,'2010-01-12 15:21:52',NULL,NULL);
INSERT INTO `users` VALUES (113,'Chris Salazar','chrsal41@hotmail.com','LEtx1aCFAJ','lKul9XxP6F','admin','Chris Salazar','每一个微小的努力，都值得骄傲。',NULL,775,'2017-04-25 17:10:58',NULL,NULL);
INSERT INTO `users` VALUES (114,'Lu Yuning','yuninglu@mail.com','4eMoGfHsvq','mGzZkvoKNW','teacher','Lu Yuning','不要等机会来找你，要去创造机会。',NULL,269,'2010-04-03 19:48:26',NULL,NULL);
INSERT INTO `users` VALUES (115,'Imai Ikki','ikkiimai@hotmail.com','5t6844LUrR','KcxXkOMyBK','student','Imai Ikki','心态决定高度。',NULL,612,'2003-05-05 21:38:20',NULL,NULL);
INSERT INTO `users` VALUES (116,'Xu Ziyi','xuziyi4@gmail.com','W468AiWRWB','XxDx7QWaF9','student','Xu Ziyi','让人生充满惊喜与挑战。',NULL,120,'2011-12-18 07:40:39',NULL,NULL);
INSERT INTO `users` VALUES (117,'Lu Anqi','anqilu903@gmail.com','JvuPzWSq1L','Z2OG1BDdyw','student','Lu Anqi','人生没有终点，只有不断的努力。',NULL,87,'2022-04-26 12:43:18',NULL,NULL);
INSERT INTO `users` VALUES (118,'Aoki Ayato','ayatoaok1@gmail.com','tpqDrP3pw3','gF3AJO1j4h','student','Aoki Ayato','向阳而生，永不低头。',NULL,320,'2007-01-08 08:22:55',NULL,NULL);
INSERT INTO `users` VALUES (119,'Shen Jialun','jialus10@gmail.com','dFYf2Ui29Q','9m51LCz8sm','teacher','Shen Jialun','每个小小的努力，都会成就更大的未来。',NULL,213,'2002-09-30 19:50:39',NULL,NULL);
INSERT INTO `users` VALUES (120,'Mui Ling Ling','mui1018@icloud.com','yguIJ0perw','oB0llUiN1t','student','Mui Ling Ling','做真实的自己，给世界最好的样子。',NULL,609,'2002-06-27 16:59:27',NULL,NULL);
INSERT INTO `users` VALUES (121,'Sun Anqi','sun1030@gmail.com','EKN3C5tpR4','oSYzrWuhpv','admin','Sun Anqi','风雨后，彩虹更美丽。',NULL,757,'2010-03-13 14:12:21',NULL,NULL);
INSERT INTO `users` VALUES (122,'Amanda Fisher','fishera@outlook.com','rIhbunNNje','q7huI0Bszj','student','Amanda Fisher','生活是自己给自己的最好的礼物。',NULL,697,'2005-11-06 15:20:12',NULL,NULL);
INSERT INTO `users` VALUES (123,'Nishimura Hikaru','nish81@outlook.com','bqUQq7PVMh','ZA8fYH5Ubg','teacher','Nishimura Hikaru','做一个不被定义的人。',NULL,133,'2007-11-22 04:15:55',NULL,NULL);
INSERT INTO `users` VALUES (124,'Xiong Jialun','jialunxiong1227@outlook.com','M1cVneQt4b','cFyKc95RGk','teacher','Xiong Jialun','跨越千山万水，只为追寻一个梦想。',NULL,134,'2009-11-22 07:09:43',NULL,NULL);
INSERT INTO `users` VALUES (125,'Fukuda Miu','miufukuda@icloud.com','A3kCLtAQA5','1xAjrz9wUh','admin','Fukuda Miu','不急不躁，步步为营。',NULL,812,'2004-09-19 20:00:24',NULL,NULL);
INSERT INTO `users` VALUES (126,'Sasaki Daisuke','sasadaisu@icloud.com','bskjv56Ccj','cEkoDogY0r','teacher','Sasaki Daisuke','人生的意义在于不断探索。',NULL,628,'2005-10-07 08:26:32',NULL,NULL);
INSERT INTO `users` VALUES (127,'Hirano Misaki','hirano923@icloud.com','ADqFsyrTO5','RmGPP50H8t','student','Hirano Misaki','一切皆有可能，心态决定一切。',NULL,663,'2017-11-06 19:18:28',NULL,NULL);
INSERT INTO `users` VALUES (128,'Kong Jialun','jialkong7@gmail.com','xngBNc1paM','N86ET6dsKi','student','Kong Jialun','我是我，独一无二。',NULL,127,'2003-01-18 02:32:33',NULL,NULL);
INSERT INTO `users` VALUES (129,'Yuen Suk Yee','yusuky6@gmail.com','Qa8HjhvQBp','ar2s4QbBDh','admin','Yuen Suk Yee','在平凡中活出不平凡。',NULL,694,'2023-01-02 02:16:18',NULL,NULL);
INSERT INTO `users` VALUES (130,'Mak Sze Yu','makszeyu@mail.com','SpaG3c6BBU','VpzPNxjhFw','admin','Mak Sze Yu','生活，不必完美，但要尽力。',NULL,524,'2019-08-07 18:56:17',NULL,NULL);
INSERT INTO `users` VALUES (131,'Kaneko Sara','saka@yahoo.com','W6hMtXuhea','n95DhYTU09','teacher','Kaneko Sara','从今天起，做自己的英雄。',NULL,478,'2001-09-14 02:24:02',NULL,NULL);
INSERT INTO `users` VALUES (132,'Fukuda Riku','rikuf@icloud.com','1dSX0VFtku','9P2yFirxIE','admin','Fukuda Riku','勇敢的心，无所畏惧。',NULL,993,'2007-03-28 21:20:58',NULL,NULL);
INSERT INTO `users` VALUES (133,'Fujiwara Hikari','fujiwarahikari86@icloud.com','cIdmXCimfb','532XXUpMW6','student','Fujiwara Hikari','做一颗温暖的心，照亮他人的人生。',NULL,169,'2018-08-28 13:05:05',NULL,NULL);
INSERT INTO `users` VALUES (134,'Qiu Jiehong','jieqiu2@hotmail.com','nfPre9hdgO','8PmP930f8i','teacher','Qiu Jiehong','生活就是不断前行，勇敢去追。',NULL,40,'2005-06-26 04:18:18',NULL,NULL);
INSERT INTO `users` VALUES (135,'Lam Wai Man','lam923@icloud.com','ezJNJOBuqd','1cNbzbnOc6','teacher','Lam Wai Man','用力去爱，用心去活。',NULL,210,'2023-02-12 03:29:54',NULL,NULL);
INSERT INTO `users` VALUES (136,'Gong Xiuying','xiuygong@icloud.com','KGWIqlbDAo','gyTPiXNUiP','student','Gong Xiuying','人生不怕挫折，怕的是放弃。',NULL,217,'2011-06-12 22:13:56',NULL,NULL);
INSERT INTO `users` VALUES (137,'Mui Hui Mei','muhm1978@outlook.com','VSIcgc8o42','2UGrIK4f06','admin','Mui Hui Mei','永远做一个微笑面对困难的人。',NULL,216,'2008-09-29 14:07:11',NULL,NULL);
INSERT INTO `users` VALUES (138,'Hui Chun Yu','huchunyu@yahoo.com','SO1oEPqCI6','NJLkorqu2s','admin','Hui Chun Yu','走过风雨，见彩虹。',NULL,927,'2002-09-20 09:32:24',NULL,NULL);
INSERT INTO `users` VALUES (139,'Ku Chieh Lun','clku45@gmail.com','Bmoo4tO5Y0','6JvMdalhaj','teacher','Ku Chieh Lun','做一个自由自在的人。',NULL,662,'2001-09-21 22:05:35',NULL,NULL);
INSERT INTO `users` VALUES (140,'Troy Long','long7@outlook.com','8bTZH5xprz','xXO8VKOYtB','admin','Troy Long','梦想是人生最美的力量。',NULL,869,'2014-05-09 06:25:43',NULL,NULL);
INSERT INTO `users` VALUES (141,'Ishikawa Ikki','ishikki1006@gmail.com','g46wlTOtxT','fd1IHU4B53','teacher','Ishikawa Ikki','生活需要点勇气，梦想需要点坚持。',NULL,69,'2004-07-31 18:15:08',NULL,NULL);
INSERT INTO `users` VALUES (142,'Huang Yunxi','yunxihuang@outlook.com','MWGMUWewIv','cQLEDDsjMg','teacher','Huang Yunxi','笑对人生，走好每一步。',NULL,556,'2008-12-27 10:07:34',NULL,NULL);
INSERT INTO `users` VALUES (143,'Helen Cole','helencole@icloud.com','WtlIzISqwG','BIn5hrgQlj','student','Helen Cole','坚持，是成功的最重要的秘诀。',NULL,280,'2011-11-01 20:49:06',NULL,NULL);
INSERT INTO `users` VALUES (144,'Gregory Woods','woodsg@icloud.com','pDzWTSKsXQ','KVzpmKS3OE','student','Gregory Woods','生活太短，何必为难自己。',NULL,46,'2010-09-15 18:38:52',NULL,NULL);
INSERT INTO `users` VALUES (145,'Ono Ryota','ryotaono@hotmail.com','5C6HC42CIi','iDGZtvK4ml','student','Ono Ryota','阳光总会在阴霾后洒满大地。',NULL,5,'2009-07-07 16:09:36',NULL,NULL);
INSERT INTO `users` VALUES (146,'Hashimoto Ryota','hasryota726@outlook.com','QSsNF6aXcG','QipFCgg7Xv','student','Hashimoto Ryota','路很远，风景很美，继续走。',NULL,592,'2017-09-01 20:54:16',NULL,NULL);
INSERT INTO `users` VALUES (147,'Li Xiaoming','lxiaoming1202@outlook.com','1QHbUPZ9w6','WISpNLp0st','admin','Li Xiaoming','相信自己的力量，成就不一样的人生。',NULL,147,'2006-08-19 07:15:56',NULL,NULL);
INSERT INTO `users` VALUES (148,'Yamaguchi Ryota','ryamaguchi@icloud.com','JNsUjn5nKC','Tlon3pGByE','student','Yamaguchi Ryota','幸福就是，做自己喜欢做的事。',NULL,59,'2018-10-08 04:45:42',NULL,NULL);
INSERT INTO `users` VALUES (149,'Masuda Sara','sm704@gmail.com','hGtizszzAa','0ueAltacyK','student','Masuda Sara','坚持自己所信，勇敢走自己的路。',NULL,398,'2009-01-03 19:28:41',NULL,NULL);
INSERT INTO `users` VALUES (150,'Fan Ziyi','ziyifa@mail.com','Ih8Xjy5Gn4','xxj2wmY1aY','teacher','Fan Ziyi','生活如诗，勇敢去写。',NULL,929,'2011-11-20 13:33:57',NULL,NULL);
INSERT INTO `users` VALUES (151,'Nakamura Hikaru','hikarnakamura@icloud.com','UHIo54qVZh','dZvDwOnw3K','teacher','Nakamura Hikaru','爱与努力，才是生命最美的语言。',NULL,355,'2004-09-25 08:34:24',NULL,NULL);
INSERT INTO `users` VALUES (152,'Yamamoto Riku','riku220@icloud.com','S9DbuTFBgM','BMRWB1vgMm','student','Yamamoto Riku','人生因坚持而美丽。',NULL,608,'2024-02-01 17:30:29',NULL,NULL);
INSERT INTO `users` VALUES (153,'Lam Chung Yin','lamchungyin7@gmail.com','5CkzMiH5Nx','aKVjF3LMKg','admin','Lam Chung Yin','过好今天，就是最好的明天。',NULL,857,'2017-08-27 15:35:21',NULL,NULL);
INSERT INTO `users` VALUES (154,'Li Lan','lila4@yahoo.com','vHZhIHGlcr','eCYZkP5TI2','teacher','Li Lan','生活是一场未知的冒险。',NULL,164,'2003-05-17 14:01:33',NULL,NULL);
INSERT INTO `users` VALUES (155,'Dale Coleman','coleman2019@icloud.com','2gxR7Ry3Ss','gnebm4PneA','admin','Dale Coleman','未来就是你现在努力的模样。',NULL,984,'2011-02-21 22:20:25',NULL,NULL);
INSERT INTO `users` VALUES (156,'Lui Chi Yuen','luchiyuen@gmail.com','O3SVb345nx','2ATWl94uHI','admin','Lui Chi Yuen','每一天，都要活得有意义。',NULL,203,'2019-12-25 09:20:28',NULL,NULL);
INSERT INTO `users` VALUES (157,'Loui Sau Man','losauman@outlook.com','XD9RYrblsu','V2cg92AYB3','teacher','Loui Sau Man','每个微小的选择，都决定着未来的样子。',NULL,634,'2005-04-08 19:57:44',NULL,NULL);
INSERT INTO `users` VALUES (158,'Cynthia Vasquez','vasc8@outlook.com','Bnf7nJQoAW','Z9sDC0abOn','admin','Cynthia Vasquez','不断前进，超越自我。',NULL,408,'2006-12-14 00:24:00',NULL,NULL);
INSERT INTO `users` VALUES (159,'Edna Evans','ee89@gmail.com','xjB6U24Yoo','qJmTPhW4w0','admin','Edna Evans','逆风飞扬，向着梦想出发。',NULL,312,'2022-10-30 00:49:04',NULL,NULL);
INSERT INTO `users` VALUES (160,'Harada Seiko','seikoharada@yahoo.com','8N9EfmGhju','kXCXufsU0C','student','Harada Seiko','用笑容面对困难，成功自然属于你。',NULL,622,'2021-09-26 23:48:57',NULL,NULL);
INSERT INTO `users` VALUES (161,'Ma On Na','onm1129@icloud.com','A9qBIvKZGU','xwcaZwaH5l','admin','Ma On Na','快乐源自内心的坚持。',NULL,519,'2006-12-05 19:48:45',NULL,NULL);
INSERT INTO `users` VALUES (162,'Jean Mendoza','jeanm@mail.com','6ZeDkFFaEa','oBRQeK56Rf','teacher','Jean Mendoza','做一个温暖的自己，温暖别人。',NULL,666,'2000-02-16 08:45:40',NULL,NULL);
INSERT INTO `users` VALUES (163,'Sakai Misaki','sakaimisaki@gmail.com','KILIqs7HgS','ocWTxO5G4Z','student','Sakai Misaki','梦想不远，未来就在前方。',NULL,873,'2022-07-10 15:21:39',NULL,NULL);
INSERT INTO `users` VALUES (164,'Okada Kazuma','okadakazuma@hotmail.com','MUTflBkscA','lBkmt1X3Hf','teacher','Okada Kazuma','世界因我而精彩，我因梦想而坚强。',NULL,520,'2007-12-22 01:23:45',NULL,NULL);
INSERT INTO `users` VALUES (165,'Ishii Kasumi','ki8@yahoo.com','vwnGrWj0VV','Tg79TfxCtH','admin','Ishii Kasumi','坚持下去，你的努力终会被看见。',NULL,695,'2023-12-28 13:57:59',NULL,NULL);
INSERT INTO `users` VALUES (166,'Mui Ka Fai','mukafai@hotmail.com','xKD1rJfZkB','Q3MzO1n518','teacher','Mui Ka Fai','不为未来担忧，只为今天努力。',NULL,240,'2013-05-14 09:16:01',NULL,NULL);
INSERT INTO `users` VALUES (167,'Zeng Anqi','zena@icloud.com','JxXMjeulJy','qejq1c5DGY','teacher','Zeng Anqi','抓住每一个机会，活出精彩人生。',NULL,740,'2018-04-28 21:54:33',NULL,NULL);
INSERT INTO `users` VALUES (168,'Maeda Yuna','ymaeda@gmail.com','BmgXfjwqHJ','kBHjq8Uz6G','teacher','Maeda Yuna','努力的过程比结果更重要。',NULL,328,'2016-05-02 19:32:22',NULL,NULL);
INSERT INTO `users` VALUES (169,'Wendy Aguilar','wendy76@mail.com','eiNGFAgCcQ','sB8ZctWvRo','teacher','Wendy Aguilar','阳光总会照进心里。',NULL,473,'2005-04-02 18:11:48',NULL,NULL);
INSERT INTO `users` VALUES (170,'Michael Kim','michaelki@icloud.com','gFBhUu2ZE2','mboP6erItg','student','Michael Kim','爱自己，就是最好的开始。',NULL,103,'2024-08-26 05:50:26',NULL,NULL);
INSERT INTO `users` VALUES (171,'Man Kwok Ming','kwokming1960@outlook.com','BHNYS8K3Zu','v2ERUSdipK','student','Man Kwok Ming','心若无畏，未来可期。',NULL,863,'2012-08-02 18:32:37',NULL,NULL);
INSERT INTO `users` VALUES (172,'Jia Ziyi','zjia@gmail.com','kRM26XY1nb','HhxU0gqxSX','admin','Jia Ziyi','活在当下，珍惜每一刻。',NULL,370,'2019-12-21 13:05:40',NULL,NULL);
INSERT INTO `users` VALUES (173,'Gong Anqi','ganqi621@mail.com','g7kCgbXWkr','RX17Cvi2tb','teacher','Gong Anqi','用心活出每一天的精彩。',NULL,167,'2021-08-30 13:58:35',NULL,NULL);
INSERT INTO `users` VALUES (174,'Sylvia Stevens','sylvistevens@yahoo.com','DA5i5ktq93','lbpvcSkFGs','teacher','Sylvia Stevens','每一天都是新的希望。',NULL,542,'2022-04-19 16:10:27',NULL,NULL);
INSERT INTO `users` VALUES (175,'Dorothy Hicks','hicks1992@yahoo.com','drfCXVNqot','GGEhpMMGV3','teacher','Dorothy Hicks','世界很大，我愿走下去。',NULL,991,'2023-09-24 23:57:36',NULL,NULL);
INSERT INTO `users` VALUES (176,'Sano Hikaru','hikas8@outlook.com','7A8PRZDc1w','uuweYnOMyA','student','Sano Hikaru','活成自己喜欢的样子。',NULL,561,'2011-12-15 15:59:51',NULL,NULL);
INSERT INTO `users` VALUES (177,'Joyce Gonzales','jgonza8@icloud.com','a9xgVpRsWO','Wt1xG5ipek','student','Joyce Gonzales','一切皆有可能，取决于你自己。',NULL,257,'2018-03-02 21:19:12',NULL,NULL);
INSERT INTO `users` VALUES (178,'Charles Stephens','chs@hotmail.com','CQYNuT6WHR','i1d9A7xzWm','admin','Charles Stephens','坚持，终有一天你会笑着看自己走过的路。',NULL,958,'2014-07-08 07:25:04',NULL,NULL);
INSERT INTO `users` VALUES (179,'Wong Yu Ling','ylwong@outlook.com','1xsGpTDupF','ZVNWOJRWFQ','admin','Wong Yu Ling','生活从不问你是否准备好，它只告诉你：走下去。',NULL,575,'2012-02-24 21:33:29',NULL,NULL);
INSERT INTO `users` VALUES (180,'Lai Sum Wing','swla@gmail.com','kS6TXEUt8F','hCDaMP7fcr','student','Lai Sum Wing','每天都是新的开始，给自己一个微笑。',NULL,287,'2022-01-07 12:14:26',NULL,NULL);
INSERT INTO `users` VALUES (181,'Ono Daisuke','daiso3@outlook.com','cIp0k6pMT2','hGlTtf2TaS','admin','Ono Daisuke','做最好的自己，不负韶华。',NULL,346,'2019-07-27 22:02:29',NULL,NULL);
INSERT INTO `users` VALUES (182,'Yin Lan','layin@outlook.com','ALZYvSHiOE','k63DBaZG92','student','Yin Lan','青春不常在，努力正当时。',NULL,648,'2013-01-10 00:04:19',NULL,NULL);
INSERT INTO `users` VALUES (183,'Chic Kwok Wing','kwokwingchic3@yahoo.com','6gg2LYYiOO','m5fFbmD0Wr','student','Chic Kwok Wing','心怀梦想，脚步从不停歇。',NULL,104,'2003-10-12 03:09:06',NULL,NULL);
INSERT INTO `users` VALUES (184,'Huang Shihan','shhuang1026@outlook.com','nMPaxQjvze','cExAWdqeLa','admin','Huang Shihan','生活不可能完美，但要有意义。',NULL,520,'2021-11-03 19:42:09',NULL,NULL);
INSERT INTO `users` VALUES (185,'Nomura Hikaru','hikarunomura@icloud.com','YJIqGa63h2','ECvvib7NoU','student','Nomura Hikaru','路虽远，行则必至。',NULL,513,'2000-05-09 13:08:05',NULL,NULL);
INSERT INTO `users` VALUES (186,'Yu Zhennan','yuz1102@outlook.com','5HoUrzOL77','BxZJaYmYn6','student','Yu Zhennan','活得精彩，为自己而活。',NULL,413,'2014-01-23 04:51:51',NULL,NULL);
INSERT INTO `users` VALUES (187,'Arthur Jenkins','artjenkins3@gmail.com','Lyea7MBKvx','8TcbHfQ2Fp','admin','Arthur Jenkins','未来属于勇敢的人。',NULL,901,'2015-08-31 16:13:43',NULL,NULL);
INSERT INTO `users` VALUES (188,'Maruyama Mai','maimaruyama8@hotmail.com','8tUuKcsUb9','JplT7mofjH','admin','Maruyama Mai','走自己的路，让别人去说吧。',NULL,917,'2014-12-07 00:52:28',NULL,NULL);
INSERT INTO `users` VALUES (189,'Philip Hughes','hughephil@outlook.com','Cfr4BSWVcz','IRRyaUJrpz','teacher','Philip Hughes','走得再远，也别忘了最初的自己。',NULL,974,'2003-11-03 16:14:31',NULL,NULL);
INSERT INTO `users` VALUES (190,'Beverly Evans','bevevans1225@gmail.com','MyV4IyicI1','pNrTSTsZuc','teacher','Beverly Evans','做自己喜欢做的事，不问结果。',NULL,524,'2012-03-18 06:41:19',NULL,NULL);
INSERT INTO `users` VALUES (191,'Kimberly Cox','kimberly902@icloud.com','Qo88Pqk7cX','Nxl9fR8WeZ','teacher','Kimberly Cox','世界那么大，我要去看看。',NULL,326,'2017-10-20 18:50:48',NULL,NULL);
INSERT INTO `users` VALUES (192,'Yu Yunxi','yunxiyu@gmail.com','jglxk7dAX1','hpQOkSCZNJ','teacher','Yu Yunxi','我的坚持，你看得见。',NULL,404,'2002-08-13 07:11:51',NULL,NULL);
INSERT INTO `users` VALUES (193,'Eleanor Bell','belle91@gmail.com','Uj1SXddPWr','TyUEMZZVKF','student','Eleanor Bell','用心生活，活得真诚。',NULL,909,'2018-11-03 04:57:58',NULL,NULL);
INSERT INTO `users` VALUES (194,'Tang Jialun','jialun716@outlook.com','QOKG69Sn7M','9hhDJnxjMw','admin','Tang Jialun','感恩每一段旅程，收获每一次成长。',NULL,646,'2009-02-09 08:06:28',NULL,NULL);
INSERT INTO `users` VALUES (195,'Hashimoto Misaki','hmisaki1941@gmail.com','QRzp13xXCL','6rh7f1wd7p','student','Hashimoto Misaki','不管风多大，我都要勇敢向前走。',NULL,182,'2013-11-23 22:50:49',NULL,NULL);
INSERT INTO `users` VALUES (196,'Jean Torres','jeantorres526@hotmail.com','WuQyVEHgpt','TAOJLCxadD','student','Jean Torres','时光不负努力的人。',NULL,331,'2018-05-10 12:10:34',NULL,NULL);
INSERT INTO `users` VALUES (197,'Bryan Coleman','bryacol802@icloud.com','nrof39RorL','gBzZa0OOjC','student','Bryan Coleman','你若盛开，清风自来。',NULL,532,'2012-03-23 17:16:32',NULL,NULL);
INSERT INTO `users` VALUES (198,'Jesus Campbell','campbelljesus@yahoo.com','pqldu9s3td','yRyTd8o8Ke','student','Jesus Campbell','生活不简单，活得要精彩。',NULL,523,'2012-07-01 18:18:31',NULL,NULL);
INSERT INTO `users` VALUES (199,'Du Anqi','da10@yahoo.com','pCGt5EdT8X','awQQPFebI3','admin','Du Anqi','未来不惧，过去不忆。',NULL,108,'2012-09-21 03:32:09',NULL,NULL);
INSERT INTO `users` VALUES (200,'Xu Jialun','xu92@hotmail.com','0JITNlK2qn','JrGNmyBW2r','teacher','Xu Jialun','做自己想做的事，过自己想过的生活。',NULL,971,'2006-02-04 19:12:48',NULL,NULL);
INSERT INTO `users` VALUES (201,'Fujii Mio','miof@yahoo.com','2edxQ5IwFG','t6b7MrB8za','admin','Fujii Mio','不忘初心，方得始终。',NULL,294,'2016-03-02 02:56:26',NULL,NULL);
INSERT INTO `users` VALUES (202,'Valerie Morales','valerie1@gmail.com','H9Ftkuya0o','TZCCsLSwWl','admin','Valerie Morales','做最好的自己，过最棒的生活。',NULL,857,'2006-07-19 14:50:53',NULL,NULL);
INSERT INTO `users` VALUES (203,'Wan Wing Kuen','wanwingkuen@gmail.com','TOs87FcpGb','PVmi6Ta5rm','teacher','Wan Wing Kuen','生活因奋斗而美丽。',NULL,731,'2020-08-29 14:42:48',NULL,NULL);
INSERT INTO `users` VALUES (204,'Gregory Walker','walkergregory@gmail.com','rivSy6dwW0','2ksyt8bhlR','admin','Gregory Walker','梦想从不嫌晚，努力就能改变。',NULL,512,'2021-09-08 16:16:09',NULL,NULL);
INSERT INTO `users` VALUES (205,'Gu Xiuying','xiuying225@icloud.com','tavPzmRzp1','cyUFhR0yU5','teacher','Gu Xiuying','用微笑面对生活的每一天。',NULL,640,'2016-01-22 11:28:16',NULL,NULL);
INSERT INTO `users` VALUES (206,'Yin Jialun','jiay@gmail.com','kMqAYapDr6','Tj5hKpUTAE','student','Yin Jialun','不管过去如何，未来依旧光明。',NULL,991,'2021-11-25 01:49:26',NULL,NULL);
INSERT INTO `users` VALUES (207,'Lui Cho Yee','choyeelui@outlook.com','bTZALcttUv','za0osv6Kcn','admin','Lui Cho Yee','每一步都算数，未来可期。',NULL,313,'2003-08-24 20:59:59',NULL,NULL);
INSERT INTO `users` VALUES (208,'Johnny Thomas','thomasj@icloud.com','bON5b18gdy','t2KPxx7zzV','teacher','Johnny Thomas','做一个有趣的人，活得充实。',NULL,645,'2003-05-18 06:49:31',NULL,NULL);
INSERT INTO `users` VALUES (209,'Tse Tsz Hin','tsetszhin516@icloud.com','1K2vnDpucO','1SaCna12Yj','admin','Tse Tsz Hin','生活，就是一场不断追寻的旅程。',NULL,367,'2014-05-06 04:28:19',NULL,NULL);
INSERT INTO `users` VALUES (210,'Watanabe Ryota','watanabe1951@icloud.com','clrnjkVg2W','IHHkvJQPgj','admin','Watanabe Ryota','因为有你，世界才更美丽。',NULL,824,'2006-01-08 16:05:42',NULL,NULL);
INSERT INTO `users` VALUES (211,'Ishikawa Yamato','ishikawaya908@outlook.com','OTkC4KmhmI','dF6UNWUz2G','student','Ishikawa Yamato','未来属于敢于追梦的人。',NULL,79,'2017-05-23 22:28:30',NULL,NULL);
INSERT INTO `users` VALUES (212,'Ueda Aoshi','ua1228@icloud.com','xdGQyH6qb8','r7O9N5NUaN','student','Ueda Aoshi','人生就是不断的超越自己。',NULL,461,'2004-01-21 13:48:19',NULL,NULL);
INSERT INTO `users` VALUES (213,'Ma Sum Wing','swma@icloud.com','CfPRom5VsI','6aDCssDPQo','student','Ma Sum Wing','你若勇敢，生活就会温柔以待。',NULL,440,'2016-08-09 20:02:38',NULL,NULL);
INSERT INTO `users` VALUES (214,'Ding Xiuying','dingxiuying6@icloud.com','nIGhqFJeqG','d7WNOApnCZ','teacher','Ding Xiuying','只为梦，不负自己。',NULL,903,'2013-02-21 00:04:11',NULL,NULL);
INSERT INTO `users` VALUES (215,'Nomura Kazuma','kazumano320@outlook.com','hzWK0PgWWp','wLzFLlFDCT','teacher','Nomura Kazuma','追寻梦想，无畏风雨。',NULL,215,'2017-11-25 03:47:39',NULL,NULL);
INSERT INTO `users` VALUES (216,'Jesus White','whjes@icloud.com','yxsz6bC6SS','i11lQqBWK3','teacher','Jesus White','用力去爱，用心去活。',NULL,770,'2010-10-31 18:21:36',NULL,NULL);
INSERT INTO `users` VALUES (217,'Kim Mitchell','kim5@mail.com','s8jovn4aHf','9D7TbaDwug','admin','Kim Mitchell','一切都会好起来的，未来可期。',NULL,7,'2019-08-24 06:46:11',NULL,NULL);
INSERT INTO `users` VALUES (218,'Kobayashi Rin','rinkobay@icloud.com','6xgc2VffUS','0BwZH1zMBe','teacher','Kobayashi Rin','不畏前路漫漫，勇敢追逐梦想。',NULL,144,'2004-03-12 13:16:23',NULL,NULL);
INSERT INTO `users` VALUES (219,'Tsui Wai San','waisantsui10@hotmail.com','hBdnOjt4bU','DLCmIJRXAt','teacher','Tsui Wai San','用行动书写属于自己的精彩。',NULL,272,'2002-01-24 12:38:15',NULL,NULL);
INSERT INTO `users` VALUES (220,'Uchida Seiko','seuchida@icloud.com','kQRiH8syUF','31B4TEl2lv','student','Uchida Seiko','只有前行，才是人生唯一的选择。',NULL,609,'2014-02-04 17:34:53',NULL,NULL);
INSERT INTO `users` VALUES (221,'Lok Ming Sze','lms829@icloud.com','4kgNO4KLpJ','d5oExhj9tF','admin','Lok Ming Sze','成功来自不懈努力。',NULL,762,'2009-08-29 01:29:16',NULL,NULL);
INSERT INTO `users` VALUES (222,'Han Fat','hfat@icloud.com','NKMzR5AERt','Q2FkBXW834','teacher','Han Fat','路的尽头，是未知的美好。',NULL,811,'2002-07-12 08:10:11',NULL,NULL);
INSERT INTO `users` VALUES (223,'Heung Hui Mei','huimeiheu8@mail.com','v3c72QADVE','MaHOefdTM8','admin','Heung Hui Mei','梦想不是空想，而是努力追逐的目标。',NULL,621,'2007-02-18 20:42:31',NULL,NULL);
INSERT INTO `users` VALUES (224,'Mary Butler','mbutler@outlook.com','8VbiDzqmpR','bxkJXkfNN6','admin','Mary Butler','只要心怀希望，风雨也能成彩虹。',NULL,811,'2023-04-29 23:53:56',NULL,NULL);
INSERT INTO `users` VALUES (225,'Nishimura Eita','nishimuraei@outlook.com','sF9tTgzFcg','QNBnJX0hjU','student','Nishimura Eita','感谢生活中的每一份惊喜。',NULL,253,'2021-02-04 16:36:54',NULL,NULL);
INSERT INTO `users` VALUES (226,'Yung Kar Yan','yung1105@icloud.com','4xx4pSi8ov','vQNVF2vc41','student','Yung Kar Yan','每一份努力都会有收获。',NULL,225,'2004-10-30 17:53:25',NULL,NULL);
INSERT INTO `users` VALUES (227,'Dorothy Dixon','dorothydix@yahoo.com','zNAKNm4eAI','V6ilTY357o','teacher','Dorothy Dixon','永远相信，自己是最棒的。',NULL,479,'2005-08-26 14:07:26',NULL,NULL);
INSERT INTO `users` VALUES (228,'Marie Hawkins','marie209@gmail.com','e52MaMGeBE','i3qze5rdqa','admin','Marie Hawkins','人生就是不断成长，不断超越。',NULL,833,'2023-06-05 07:55:24',NULL,NULL);
INSERT INTO `users` VALUES (229,'Cheng Jiehong','jiehongc@gmail.com','mJ71ZTViTz','mtf2qotHmC','admin','Cheng Jiehong','只有坚持，才会看到光明的未来。',NULL,849,'2009-03-02 13:33:30',NULL,NULL);
INSERT INTO `users` VALUES (230,'Amber Cole','cambe@outlook.com','PUmEEcpcMf','0xdiqNmJ5L','admin','Amber Cole','遇见你是我生命中的最美时光。',NULL,697,'2007-02-27 03:31:02',NULL,NULL);
INSERT INTO `users` VALUES (231,'Yuen Tsz Hin','yuen807@gmail.com','l3kmBQZTSi','2g68ucEdxe','admin','Yuen Tsz Hin','青春是一本书，写满奋斗与希望。',NULL,730,'2020-06-21 18:37:27',NULL,NULL);
INSERT INTO `users` VALUES (232,'Au Wai Lam','au2@hotmail.com','r6J1ucNvgx','UBTriizRsB','student','Au Wai Lam','生活不完美，但有你便无憾。',NULL,96,'2021-06-01 10:24:43',NULL,NULL);
INSERT INTO `users` VALUES (233,'Robert Vasquez','robeva57@hotmail.com','LtubqlP55A','Ygl0XR9PHF','admin','Robert Vasquez','生命的意义在于追求与实现。',NULL,323,'2022-01-28 15:27:10',NULL,NULL);
INSERT INTO `users` VALUES (234,'Pan Zitao','zip@outlook.com','0YAfNR2gnx','eZSwL7hhKt','teacher','Pan Zitao','永远不要放弃梦想，人生因此精彩。',NULL,746,'2003-07-03 17:09:58',NULL,NULL);
INSERT INTO `users` VALUES (235,'Emma Wallace','waem@icloud.com','j9stRtoIM8','p4uMlRTHHF','student','Emma Wallace','时光不老，我们不散。',NULL,121,'2005-08-14 10:28:25',NULL,NULL);
INSERT INTO `users` VALUES (236,'Edwin Gardner','gedwi2@mail.com','Laz2qnraWh','0iUfjlmr9W','teacher','Edwin Gardner','坚持自己的选择，走自己的路。',NULL,195,'2004-11-16 09:59:58',NULL,NULL);
INSERT INTO `users` VALUES (237,'Tong Ka Man','kamanto9@mail.com','X7mKsmacuB','bntokXuvfe','student','Tong Ka Man','每一天都做最好的自己。',NULL,367,'2019-06-15 01:24:08',NULL,NULL);
INSERT INTO `users` VALUES (238,'Chung Chung Yin','chungyinchung@icloud.com','GFZgePAS8D','tqrJGGfIMP','teacher','Chung Chung Yin','做自己最喜欢的事，过最幸福的生活。',NULL,281,'2003-06-16 01:31:42',NULL,NULL);
INSERT INTO `users` VALUES (239,'Tan Jialun','tanji@outlook.com','eJ7MSkpnZm','Tszjo1T7jS','teacher','Tan Jialun','未来无限可能，努力成就梦想。',NULL,160,'2005-06-15 19:24:07',NULL,NULL);
INSERT INTO `users` VALUES (240,'Liao Anqi','anqil@outlook.com','8WZ9M5adN2','x2zVkD1Eub','admin','Liao Anqi','生命短暂，做自己喜欢的事。',NULL,36,'2007-12-12 17:57:33',NULL,NULL);
INSERT INTO `users` VALUES (241,'Ruby Rice','ricer1101@hotmail.com','snvjxzqfNr','BsSNOMckTh','teacher','Ruby Rice','成长就是不断挑战自我。',NULL,56,'2023-07-16 17:16:50',NULL,NULL);
INSERT INTO `users` VALUES (242,'Gao Jialun','gao60@icloud.com','oSY6CEuppL','WFphBPy7gH','admin','Gao Jialun','我不是命运的奴隶，命运由我掌握。',NULL,387,'2014-07-21 07:51:43',NULL,NULL);
INSERT INTO `users` VALUES (243,'He Jialun','hejialun1959@outlook.com','VJTVN7YbsA','JNRABeN2p2','admin','He Jialun','勇敢追梦，人生才能精彩。',NULL,818,'2015-04-05 17:19:36',NULL,NULL);
INSERT INTO `users` VALUES (244,'Okamoto Nanami','no68@gmail.com','SmFzqjPkmx','3mjDDGgfCF','admin','Okamoto Nanami','明天会更好，今天只需要努力。',NULL,456,'2021-08-22 00:18:50',NULL,NULL);
INSERT INTO `users` VALUES (245,'Yung Wai San','yung625@yahoo.com','s2hLRyXbeM','asnshZiBkl','admin','Yung Wai San','生活不止眼前的苟且，还有诗和远方。',NULL,657,'2011-11-18 11:57:12',NULL,NULL);
INSERT INTO `users` VALUES (246,'Kyle Jackson','kyljacks@hotmail.com','SkO1RdVT3Z','oNdZxFZR8J','teacher','Kyle Jackson','世界那么大，我想去看看。',NULL,152,'2022-12-17 18:39:35',NULL,NULL);
INSERT INTO `users` VALUES (247,'Yin Zitao','ziyi@gmail.com','sCb6NMmt51','ELrty7izxd','admin','Yin Zitao','每天都是崭新的起点。',NULL,442,'2015-06-06 07:40:06',NULL,NULL);
INSERT INTO `users` VALUES (248,'Ho Tin Wing','tinwingho2017@gmail.com','nhqtj7YwzQ','fJnQ2Fb9Vw','teacher','Ho Tin Wing','幸福是奋斗出来的。',NULL,647,'2007-08-22 21:37:23',NULL,NULL);
INSERT INTO `users` VALUES (249,'Lo Cho Yee','choyeelo@gmail.com','eivTRrpEnP','rR6vPvTkse','student','Lo Cho Yee','每一步都走得踏实，每一滴汗水都有价值。',NULL,174,'2014-01-29 02:33:27',NULL,NULL);
INSERT INTO `users` VALUES (250,'Steve Martin','mars53@gmail.com','69hqJ90HLm','zxcE1aSih4','teacher','Steve Martin','梦想需要努力，生活需要勇气。',NULL,905,'2004-01-14 03:31:29',NULL,NULL);
INSERT INTO `users` VALUES (251,'Qian Xiuying','qx15@icloud.com','0UjVrt39k7','O2uyQSfZAa','admin','Qian Xiuying','每一次努力，都是未来的基石。',NULL,113,'2024-02-27 07:31:24',NULL,NULL);
INSERT INTO `users` VALUES (252,'Kono Yuna','yk11@icloud.com','OHTXyCVkdK','k84IQabipd','admin','Kono Yuna','不怕困难，勇敢前行。',NULL,821,'2002-06-14 22:41:26',NULL,NULL);
INSERT INTO `users` VALUES (253,'Ho Ka Man','hokaman812@outlook.com','j052NaeTIq','hA0B1P1gOi','teacher','Ho Ka Man','梦想不会辜负每一个努力的人。',NULL,837,'2012-06-19 07:55:41',NULL,NULL);
INSERT INTO `users` VALUES (254,'Saito Momoe','momoesai6@icloud.com','IMR76yoK4v','myFskNdNsL','admin','Saito Momoe','成为你想见到的那个人。',NULL,76,'2002-05-21 07:59:22',NULL,NULL);
INSERT INTO `users` VALUES (255,'Li Zhennan','li1@gmail.com','tA00gniezR','DYIR6Xw2IW','teacher','Li Zhennan','做自己最想做的事，成为最好的自己。',NULL,641,'2021-03-14 02:53:20',NULL,NULL);
INSERT INTO `users` VALUES (256,'Randy Dunn','randy5@icloud.com','eel7PEoWSh','lX0OC6kLQW','admin','Randy Dunn','人生最美好的事，就是为了梦想而努力。',NULL,115,'2019-06-16 20:05:29',NULL,NULL);
INSERT INTO `users` VALUES (257,'Hsuan Hui Mei','hmhsuan@hotmail.com','SAk1Qcf5w3','mpe4a3GrQB','admin','Hsuan Hui Mei','用心活出每一天的精彩。',NULL,6,'2003-10-20 17:44:44',NULL,NULL);
INSERT INTO `users` VALUES (258,'Pang Sum Wing','pang8@icloud.com','yRmlOVPVNo','lFnNQYXJ37','student','Pang Sum Wing','世界因你而美丽，我因梦想而坚强。',NULL,752,'2018-07-23 00:27:30',NULL,NULL);
INSERT INTO `users` VALUES (259,'Matsumoto Yuna','yunam@gmail.com','csgKNMkkHV','lxSzCLglKd','student','Matsumoto Yuna','青春不再，奋斗未停。',NULL,635,'2006-03-05 00:19:56',NULL,NULL);
INSERT INTO `users` VALUES (260,'Duan Shihan','duanshiha@gmail.com','XQUhdifmpM','N5Y1p4TBvX','student','Duan Shihan','一直走，直到走到光明。',NULL,566,'2014-01-06 00:35:06',NULL,NULL);
INSERT INTO `users` VALUES (261,'Chen Lu','luche7@icloud.com','G2Zx5EEF6k','TIUIczCXU7','student','Chen Lu','勇敢梦想，尽力追求。',NULL,181,'2002-04-13 13:27:52',NULL,NULL);
INSERT INTO `users` VALUES (262,'Shibata Ikki','ikkishibata88@icloud.com','seuVGyWk6N','dTLsbpGgvC','admin','Shibata Ikki','心怀梦想，脚踏实地。',NULL,424,'2007-05-13 19:33:32',NULL,NULL);
INSERT INTO `users` VALUES (263,'Song Jialun','songj406@outlook.com','gPq6pDLUzv','b70mhGKd1Z','teacher','Song Jialun','向阳而生，勇敢追梦。',NULL,101,'2016-05-19 14:53:43',NULL,NULL);
INSERT INTO `users` VALUES (264,'Kwong Tin Lok','kwontl@icloud.com','iEade5YTpR','wZDMFC07OZ','student','Kwong Tin Lok','努力，终会看见彩虹。',NULL,66,'2000-04-04 18:42:31',NULL,NULL);
INSERT INTO `users` VALUES (265,'Gladys Mendez','meglady@gmail.com','p65Dmk0zHg','9obIoDiRPu','student','Gladys Mendez','世界那么大，走出去看看。',NULL,114,'2018-06-10 19:43:59',NULL,NULL);
INSERT INTO `users` VALUES (266,'Lu Jiehong','jiehonglu402@yahoo.com','2D7kE6sIwq','8eHgdKl1F0','student','Lu Jiehong','每天都为梦想而努力。',NULL,813,'2010-02-11 23:12:43',NULL,NULL);
INSERT INTO `users` VALUES (267,'Tao Hiu Tung','taohiutung@mail.com','jzlbR5eef1','hZvMD14rCb','student','Tao Hiu Tung','做最好的自己，成就最好的未来。',NULL,620,'2013-09-25 04:35:47',NULL,NULL);
INSERT INTO `users` VALUES (268,'Yuen Kwok Wing','kwyuen322@outlook.com','Qjl6PxHxCf','8q2MIi189B','admin','Yuen Kwok Wing','每一天都是新的挑战。',NULL,534,'2021-02-12 14:49:19',NULL,NULL);
INSERT INTO `users` VALUES (269,'Fujii Momoka','fmomo1130@mail.com','tWXv3DAHxW','1QwEKtMFgW','teacher','Fujii Momoka','每一次拼搏，都是为了明天的辉煌。',NULL,84,'2010-02-09 14:00:23',NULL,NULL);
INSERT INTO `users` VALUES (270,'Matsui Seiko','seiko1@outlook.com','45wGtADEeF','0c49CdepF8','student','Matsui Seiko','每一份坚持，都会有回报。',NULL,166,'2010-06-04 13:29:50',NULL,NULL);
INSERT INTO `users` VALUES (271,'Nakayama Kenta','nakayamak1@outlook.com','Wy1ahY1Zvb','04x3knjMHU','admin','Nakayama Kenta','不断追求，不断超越。',NULL,140,'2013-08-11 07:35:26',NULL,NULL);
INSERT INTO `users` VALUES (272,'Wu Sze Yu','wusy@icloud.com','zjhjTOr0Fi','bm4QyWc1nu','teacher','Wu Sze Yu','做一个温暖的人，感动自己和他人。',NULL,7,'2008-10-11 15:13:07',NULL,NULL);
INSERT INTO `users` VALUES (273,'Chung Hiu Tung','hiutungch@yahoo.com','nIc4dTwbxi','16VAusCd4o','student','Chung Hiu Tung','走过风雨，迎接阳光。',NULL,973,'2017-05-27 21:30:18',NULL,NULL);
INSERT INTO `users` VALUES (274,'Hsuan Chun Yu','cyhsuan@mail.com','YqU2yJvXI3','yJxnFasCPj','student','Hsuan Chun Yu','勇敢面对挑战，成就不一样的自己。',NULL,276,'2020-11-04 17:38:09',NULL,NULL);
INSERT INTO `users` VALUES (275,'Jason Robinson','robijason@icloud.com','YPB8jV8Wzp','K7ihiQgyS1','teacher','Jason Robinson','梦想的背后，是努力的足迹。',NULL,92,'2002-10-27 18:13:14',NULL,NULL);
INSERT INTO `users` VALUES (276,'Ng Tin Wing','ngtw9@gmail.com','cbFGdRKm73','yFz3s8VVH6','teacher','Ng Tin Wing','人生没有终点，只有不断的奋斗。',NULL,884,'2020-12-31 02:56:23',NULL,NULL);
INSERT INTO `users` VALUES (277,'Nakano Hazuki','nakah@gmail.com','pQPe895z5L','qNwpnyGC6J','admin','Nakano Hazuki','每天努力一点，未来就更接近一点。',NULL,569,'2011-12-13 14:38:04',NULL,NULL);
INSERT INTO `users` VALUES (278,'Jessica Roberts','roberjess@icloud.com','DzziDuqqAI','8npRKbfCm7','student','Jessica Roberts','成功源自不懈的坚持与努力。',NULL,306,'2007-04-06 12:54:43',NULL,NULL);
INSERT INTO `users` VALUES (279,'Yang Zitao','zitaoyang@hotmail.com','8f0Q4Z5XCb','Ww2UPwzhok','student','Yang Zitao','做真实的自己，活出不一样的精彩。',NULL,999,'2006-06-14 22:47:33',NULL,NULL);
INSERT INTO `users` VALUES (280,'Liao Wai Man','liawaiman@yahoo.com','WkvtKka73O','L7JKOTlV3Y','admin','Liao Wai Man','用心生活，用爱去感受。',NULL,252,'2021-08-02 22:50:45',NULL,NULL);
INSERT INTO `users` VALUES (281,'Rebecca Ramos','rebecca7@outlook.com','k1geg1i14R','tY2UcQ4vpH','teacher','Rebecca Ramos','做自己，追求梦想。',NULL,111,'2017-05-15 20:38:56',NULL,NULL);
INSERT INTO `users` VALUES (282,'Sano Mio','sano8@outlook.com','Rbt6kKJmWu','w2OWVtu3gS','teacher','Sano Mio','不畏风雨，勇敢前行。',NULL,706,'2020-12-07 17:52:16',NULL,NULL);
INSERT INTO `users` VALUES (283,'Wendy Green','wgreen@icloud.com','FcQxwvg63n','pvQReU0KpW','student','Wendy Green','抓住每一刻，让生活更精彩。',NULL,390,'2006-08-20 21:22:20',NULL,NULL);
INSERT INTO `users` VALUES (284,'Luo Yuning','luyuning@outlook.com','nTyC2ph7uZ','3zhNvw0Ylf','admin','Luo Yuning','努力活出自我，做最好的自己。',NULL,1000,'2022-12-11 18:05:55',NULL,NULL);
INSERT INTO `users` VALUES (285,'Iwasaki Mio','mio1124@yahoo.com','MJu6pJb3Ee','7xtCF2yWUQ','student','Iwasaki Mio','今天的努力，成就明天的精彩。',NULL,153,'2024-04-27 09:09:55',NULL,NULL);
INSERT INTO `users` VALUES (286,'Abe Minato','abeminato3@outlook.com','kS6g6r05mV','13RPfof1Zg','student','Abe Minato','永远相信，自己能创造奇迹。',NULL,746,'2024-11-18 04:27:18',NULL,NULL);
INSERT INTO `users` VALUES (287,'Steven Vasquez','svasquez6@outlook.com','gHKxBUgpJU','qOcYOlTk05','teacher','Steven Vasquez','生活是一场马拉松，坚定前行。',NULL,737,'2009-10-22 03:49:56',NULL,NULL);
INSERT INTO `users` VALUES (288,'Du Xiuying','xiuying3@gmail.com','T708LWHUH1','sugPLbVTTl','admin','Du Xiuying','向阳而生，无畏艰难。',NULL,157,'2014-08-02 14:30:46',NULL,NULL);
INSERT INTO `users` VALUES (289,'Kimura Rin','rin126@yahoo.com','y1Z0l2xijK','XFAnmjwxcV','admin','Kimura Rin','不怕走得慢，只怕停下脚步。',NULL,836,'2004-07-29 03:14:14',NULL,NULL);
INSERT INTO `users` VALUES (290,'Zhao Xiuying','zhao1128@gmail.com','MTO9UknZQi','yicTYXWsnE','teacher','Zhao Xiuying','每一个微笑，都是对生活的热爱。',NULL,300,'2023-09-06 17:33:53',NULL,NULL);
INSERT INTO `users` VALUES (291,'Zhou Jialun','zhouj@outlook.com','wBZGexjwNP','tez7HfFVPN','student','Zhou Jialun','与其抱怨不如改变。',NULL,900,'2022-09-12 08:36:03',NULL,NULL);
INSERT INTO `users` VALUES (292,'Miura Ikki','ikki2@mail.com','X11dmGWXDP','gRE0VDUMzs','admin','Miura Ikki','让每一天都充满希望。',NULL,871,'2001-09-08 05:15:11',NULL,NULL);
INSERT INTO `users` VALUES (293,'Charlotte Grant','charlotteg10@outlook.com','oKKJwCt1ah','8CfdrhwUa0','teacher','Charlotte Grant','不放弃，直到成功。',NULL,577,'2019-09-29 04:58:54',NULL,NULL);
INSERT INTO `users` VALUES (294,'Paul Kennedy','kpaul@outlook.com','LnXOQNAOWl','MFVynpcdhi','student','Paul Kennedy','每天都是新的开始，做更好的自己。',NULL,150,'2008-09-22 17:07:13',NULL,NULL);
INSERT INTO `users` VALUES (295,'Hayashi Riku','hayriku2001@outlook.com','lkHOaI0QJE','GhUQJevdwr','student','Hayashi Riku','梦想在前方，未来在脚下。',NULL,823,'2000-06-18 16:35:44',NULL,NULL);
INSERT INTO `users` VALUES (296,'Cheng Zhennan','zhenncheng@gmail.com','omRT3Lecak','7Yw2BStcxe','admin','Cheng Zhennan','生活因奋斗而美丽。',NULL,509,'2021-07-26 18:59:35',NULL,NULL);
INSERT INTO `users` VALUES (297,'Ku Sai Wing','saiwingku@outlook.com','o0iHiCjWBY','52g6ChVUfz','student','Ku Sai Wing','每个小小的努力，都会改变人生的轨迹。',NULL,221,'2011-02-04 10:58:48',NULL,NULL);
INSERT INTO `users` VALUES (298,'Wei Lan','weil816@icloud.com','Ucv6R5lIsZ','Is92nhS84v','admin','Wei Lan','向前看，未来可期。',NULL,199,'2007-03-17 08:00:30',NULL,NULL);
INSERT INTO `users` VALUES (299,'Ryan Rose','rr507@mail.com','KQ6WMuS07m','lAU6vUiSjO','student','Ryan Rose','未来是属于勇敢的人的。',NULL,153,'2012-03-26 19:24:17',NULL,NULL);
INSERT INTO `users` VALUES (300,'Sakurai Takuya','takuyasak@mail.com','dBxMYz6G6l','sONJCVk7Dx','student','Sakurai Takuya','成为更好的自己，过更精彩的生活。',NULL,632,'2001-09-06 15:18:23',NULL,NULL);
INSERT INTO `users` VALUES (301,'Matsui Sakura','sam@icloud.com','jYBKHKF3b3','9hWw8oHVQh','student','Matsui Sakura','努力与坚持，终将迎来美好未来。',NULL,333,'2009-02-21 08:17:04',NULL,NULL);
INSERT INTO `users` VALUES (302,'Edwin Moreno','edwimoreno9@icloud.com','fvYl7laKSS','HApcfiq8wW','admin','Edwin Moreno','永远坚持自己的梦想，永不放弃。',NULL,878,'2002-05-22 11:32:56',NULL,NULL);
INSERT INTO `users` VALUES (303,'Meng Xiuying','xiuyingmeng@outlook.com','zaAYErK2Y7','XRtmEMlmmn','admin','Meng Xiuying','不负青春，不负自己。',NULL,407,'2016-11-03 07:32:32',NULL,NULL);
INSERT INTO `users` VALUES (304,'Sakamoto Hana','hanasa6@hotmail.com','fxLcEK2viI','bYLg02jWFX','teacher','Sakamoto Hana','不走捷径，才是成功的真正道路。',NULL,283,'2012-08-14 03:33:13',NULL,NULL);
INSERT INTO `users` VALUES (305,'Debbie Gutierrez','gutdebbie@outlook.com','cBsO3g5Rwu','fdwvi9iXu5','teacher','Debbie Gutierrez','人生没有彩排，每一刻都值得珍惜。',NULL,310,'2018-08-19 16:12:30',NULL,NULL);
INSERT INTO `users` VALUES (306,'Kaneko Daisuke','dkaneko@icloud.com','nnWV9PrITT','lHFYizixL6','admin','Kaneko Daisuke','不做旁观者，做自己人生的主角。',NULL,531,'2019-04-21 09:53:42',NULL,NULL);
INSERT INTO `users` VALUES (307,'Anthony Moreno','moreno10@outlook.com','1DmW82EbZu','ZnS87Zgefv','student','Anthony Moreno','别人说什么无关紧要，做自己喜欢的事最重要。',NULL,197,'2010-11-08 10:42:35',NULL,NULL);
INSERT INTO `users` VALUES (308,'Endo Ikki','endikk7@yahoo.com','3pImQzMRjv','99fVh9komA','student','Endo Ikki','生活是自己的，别让任何人左右。',NULL,823,'2004-11-19 19:13:30',NULL,NULL);
INSERT INTO `users` VALUES (309,'Fong Chun Yu','chunfon@icloud.com','jfsffIwQ3x','eOMkAcRvZt','student','Fong Chun Yu','做自己喜欢的事，过自己想过的生活。',NULL,925,'2006-09-30 02:24:03',NULL,NULL);
INSERT INTO `users` VALUES (310,'Kudo Kenta','kek229@outlook.com','4JUiSH43WM','P9oXGeBzcd','admin','Kudo Kenta','每一份坚持，都是对梦想的致敬。',NULL,240,'2012-11-06 17:05:42',NULL,NULL);
INSERT INTO `users` VALUES (311,'Kam On Kay','kamonkay1983@yahoo.com','DcdMzNDasy','KINUsD0tv0','admin','Kam On Kay','用心去生活，去爱，去奋斗。',NULL,839,'2007-09-16 03:35:18',NULL,NULL);
INSERT INTO `users` VALUES (312,'Peng Zitao','pengz2006@gmail.com','2ltocP5CeT','u6wf0HVxNM','teacher','Peng Zitao','别让别人左右你的心情。',NULL,391,'2019-08-14 16:48:13',NULL,NULL);
INSERT INTO `users` VALUES (313,'Au Wing Suen','wingsuenau@outlook.com','cyydwUc3qI','J9Ad4GCO36','teacher','Au Wing Suen','时间从不等待，努力让自己不后悔。',NULL,109,'2024-01-28 19:43:49',NULL,NULL);
INSERT INTO `users` VALUES (314,'Yang Ziyi','yang01@gmail.com','mP2gRyDCgq','TVQmtN6nnK','student','Yang Ziyi','不言放弃，坚韧前行。',NULL,778,'2023-11-19 06:50:04',NULL,NULL);
INSERT INTO `users` VALUES (315,'Abe Momoe','mabe7@gmail.com','7VdrrMeSrw','2PVkejN9H5','teacher','Abe Momoe','梦想就是努力的动力。',NULL,183,'2011-11-18 09:57:28',NULL,NULL);
INSERT INTO `users` VALUES (316,'Jiang Lu','lujiang@outlook.com','0NTx7Bxz01','CKfIQgedCn','student','Jiang Lu','做自己，世界会因你而不同。',NULL,36,'2006-05-23 11:04:24',NULL,NULL);
INSERT INTO `users` VALUES (317,'Tang Zitao','tanzi@mail.com','59Dvbq8Brl','iFWOeL755R','admin','Tang Zitao','每一份努力，都是成功的开始。',NULL,899,'2018-09-19 03:23:51',NULL,NULL);
INSERT INTO `users` VALUES (318,'Matsuda Kaito','kamatsu@gmail.com','OMIKmadVRk','d2mAIYEDJN','admin','Matsuda Kaito','不放弃任何可能，迎接每一个机会。',NULL,504,'2001-05-31 22:41:23',NULL,NULL);
INSERT INTO `users` VALUES (319,'Zhong Zhennan','zzhenn3@gmail.com','nsGaDCWMzZ','1vBYWDfS6g','admin','Zhong Zhennan','未来在前方，拼搏在当下。',NULL,256,'2010-08-15 20:40:46',NULL,NULL);
INSERT INTO `users` VALUES (320,'Francisco Hicks','hicfrancisco315@icloud.com','fC4lMOvEFS','jtBScICVD3','student','Francisco Hicks','人生如梦，唯奋斗不息。',NULL,264,'2012-08-27 14:58:03',NULL,NULL);
INSERT INTO `users` VALUES (321,'Kondo Daisuke','daiskondo902@yahoo.com','ongp10xDAv','9T5qb0SdGW','student','Kondo Daisuke','做自己喜欢的事，过自己的生活。',NULL,444,'2002-10-22 19:25:01',NULL,NULL);
INSERT INTO `users` VALUES (322,'Hashimoto Miu','hashmi44@icloud.com','7Y3b8fwUTu','IMds6W6r8K','admin','Hashimoto Miu','坚持到最后，就是胜利。',NULL,174,'2002-09-27 10:51:21',NULL,NULL);
INSERT INTO `users` VALUES (323,'Jia Zhiyuan','zhiyji1126@outlook.com','3zFVxO9V8T','TRGN1zJcnY','teacher','Jia Zhiyuan','世界因你而精彩，生活因你而有意义。',NULL,391,'2002-05-08 09:07:19',NULL,NULL);
INSERT INTO `users` VALUES (324,'Shibata Kaito','shibata809@outlook.com','TfPDpKphfg','kZQDODRwxr','teacher','Shibata Kaito','生活就是要充满希望。',NULL,692,'2018-12-27 13:00:08',NULL,NULL);
INSERT INTO `users` VALUES (325,'Mario Allen','allenmario827@icloud.com','mwrHOSit07','RMJT79Oyds','teacher','Mario Allen','让每一天都充满阳光和笑容。',NULL,186,'2020-12-04 02:51:21',NULL,NULL);
INSERT INTO `users` VALUES (326,'Lo Wing Sze','wingsze726@icloud.com','aZCTmpbaX6','kLUVq5bnHT','student','Lo Wing Sze','永远对自己有信心，未来一定会变好。',NULL,94,'2021-05-18 20:10:16',NULL,NULL);
INSERT INTO `users` VALUES (327,'Yao Jialun','yaoj04@yahoo.com','GH78NP1xx1','kErskzhjuk','teacher','Yao Jialun','每一次选择，都会成就不一样的人生。',NULL,86,'2021-01-13 20:26:44',NULL,NULL);
INSERT INTO `users` VALUES (328,'Chung Ming','chunming@mail.com','RgEu4guOtk','vDYwvwZskc','student','Chung Ming','勇敢追梦，改变自己的命运。',NULL,197,'2008-03-16 00:02:02',NULL,NULL);
INSERT INTO `users` VALUES (329,'Patrick Johnson','johnson413@gmail.com','61qxnTsUu0','WAxBiGbZ2u','admin','Patrick Johnson','用笑容面对一切，积极生活。',NULL,81,'2021-02-27 04:04:36',NULL,NULL);
INSERT INTO `users` VALUES (330,'Saito Aoi','saitaoi2@yahoo.com','zFkPdQlipA','lssZTeLwrj','student','Saito Aoi','做一个自由自在的人，勇敢追求自己的梦想。',NULL,683,'2013-03-15 23:38:25',NULL,NULL);
INSERT INTO `users` VALUES (331,'Clifford Hunt','chunt@gmail.com','PrOXFVe0Xc','LFjoENHuIA','teacher','Clifford Hunt','走自己的路，让别人去说。',NULL,583,'2023-06-29 04:51:17',NULL,NULL);
INSERT INTO `users` VALUES (332,'Sano Rena','renasano03@icloud.com','sgw1I8R0oV','0Cg2GXnIC7','teacher','Sano Rena','梦想可以改变一切，努力成就未来。',NULL,543,'2016-03-03 10:08:20',NULL,NULL);
INSERT INTO `users` VALUES (333,'Lau Tsz Hin','tszhinlau@gmail.com','QBFwh2hbys','Yxs0Ab5Dco','student','Lau Tsz Hin','青春是一场冒险，去勇敢追逐吧。',NULL,545,'2003-03-17 05:34:16',NULL,NULL);
INSERT INTO `users` VALUES (334,'Florence Garza','florencegar@icloud.com','q074MT4iFS','ZWQKzH9oez','teacher','Florence Garza','每一天都是崭新的希望。',NULL,729,'2001-02-09 07:40:36',NULL,NULL);
INSERT INTO `users` VALUES (335,'Elaine Rice','elric@yahoo.com','xA89EFQEh6','6MPMjBsjfG','admin','Elaine Rice','生活需要勇气，梦想需要坚持。',NULL,152,'2004-08-04 17:22:03',NULL,NULL);
INSERT INTO `users` VALUES (336,'Cho Wing Suen','chws81@mail.com','09gRZTVcAx','rIomBSqnpn','admin','Cho Wing Suen','做自己，过自己喜欢的生活。',NULL,643,'2000-08-29 02:00:43',NULL,NULL);
INSERT INTO `users` VALUES (337,'Jack Fisher','jacfi1219@icloud.com','LIlvmeeSzx','D7QcblgtWv','teacher','Jack Fisher','永远做最好的自己，成为更好的自己。',NULL,911,'2012-08-06 11:57:02',NULL,NULL);
INSERT INTO `users` VALUES (338,'Doris Dunn','doridunn@icloud.com','Wu0wPnDlTk','A9Shn0y2Sb','student','Doris Dunn','青春一去不复返，奋斗不止。',NULL,404,'2019-02-09 06:20:06',NULL,NULL);
INSERT INTO `users` VALUES (339,'Tin Tin Wing','tintw@hotmail.com','9pewSQpm73','wWLjXl65Ty','student','Tin Tin Wing','不惧未来，勇敢追逐梦想。',NULL,96,'2017-01-21 11:08:19',NULL,NULL);
INSERT INTO `users` VALUES (340,'Yau Hok Yau','yauhoky@icloud.com','mGY6sPOnxW','oSPQTHfSfW','admin','Yau Hok Yau','不负每一份坚持，成就最好的自己。',NULL,189,'2020-07-25 12:45:49',NULL,NULL);
INSERT INTO `users` VALUES (341,'Yuen Ho Yin','yuehy@icloud.com','FGQYx8IPKA','eIUXYTBy5w','student','Yuen Ho Yin','每一份努力，都是成功的铺路石。',NULL,220,'2008-02-18 15:49:50',NULL,NULL);
INSERT INTO `users` VALUES (342,'Cho Chung Yin','ccy@gmail.com','LDITO0ULRe','7NJQPh1PRW','teacher','Cho Chung Yin','未来掌握在自己手中。',NULL,677,'2020-05-04 21:58:04',NULL,NULL);
INSERT INTO `users` VALUES (343,'Yoshida Ryota','yoshidaryota@yahoo.com','5ych8KIm0r','ta98F43z7v','teacher','Yoshida Ryota','成功来自不懈的坚持。',NULL,602,'2006-10-15 04:51:45',NULL,NULL);
INSERT INTO `users` VALUES (344,'Taniguchi Mai','taniguchi7@gmail.com','7lVBsdWbHq','XJsP4YRJak','student','Taniguchi Mai','让生活的每一天都充满活力。',NULL,641,'2017-09-25 21:03:04',NULL,NULL);
INSERT INTO `users` VALUES (345,'Lu Yunxi','yl7@outlook.com','72c7Y5cnP5','TWlRELVkt2','admin','Lu Yunxi','每一步都走得坚定，每一个梦想都值得追求。',NULL,990,'2023-03-28 23:49:04',NULL,NULL);
INSERT INTO `users` VALUES (346,'Shi Zhiyuan','zhiyuans@outlook.com','JuTS10o6np','W2JVm13XwP','teacher','Shi Zhiyuan','梦想从未远离，努力让它触手可及。',NULL,238,'2021-11-04 16:58:32',NULL,NULL);
INSERT INTO `users` VALUES (347,'Zou Jiehong','jiehong8@outlook.com','oP2pIEeGZf','nH8qckTHln','teacher','Zou Jiehong','别让任何困难打败你，坚持到最后。',NULL,243,'2003-03-31 19:04:10',NULL,NULL);
INSERT INTO `users` VALUES (348,'Qin Lan','qinlan@icloud.com','SQI0kRn0m4','WL6TMjxJGj','admin','Qin Lan','每一天都是新的开始，活出自己的精彩。',NULL,737,'2000-02-11 22:11:33',NULL,NULL);
INSERT INTO `users` VALUES (349,'Dai Wai Lam','dai108@outlook.com','e1MCam8zjt','Sxyx3VjNz9','student','Dai Wai Lam','勇敢前行，路途永远都不怕远。',NULL,83,'2013-05-17 23:46:41',NULL,NULL);
INSERT INTO `users` VALUES (350,'Cho Ming','ming2@outlook.com','QfR1kduLNV','0WAU3UPRHf','admin','Cho Ming','让生活有意义，让自己有价值。',NULL,176,'2002-07-21 18:42:25',NULL,NULL);
INSERT INTO `users` VALUES (351,'Emma Butler','emmabu69@gmail.com','zXUEkgTetH','d0ylM1RXLX','student','Emma Butler','每一个微笑背后都是自信与勇气。',NULL,372,'2011-08-20 10:49:26',NULL,NULL);
INSERT INTO `users` VALUES (352,'Peng Xiaoming','pxiaoming@outlook.com','oy01RJ3TAO','lxT9SRwsXC','admin','Peng Xiaoming','做最好的自己，过最精彩的生活。',NULL,142,'2021-07-28 12:24:44',NULL,NULL);
INSERT INTO `users` VALUES (353,'Theodore Coleman','colemantheod5@icloud.com','7c5loUcBnc','RgvqVyIQyf','student','Theodore Coleman','生命的意义在于不断追求梦想。',NULL,206,'2001-02-11 20:39:23',NULL,NULL);
INSERT INTO `users` VALUES (354,'To Sum Wing','swto4@outlook.com','6K9DV8uiyE','fjuMOfLfCB','teacher','To Sum Wing','永远不要停止追逐梦想的步伐。',NULL,865,'2020-07-22 20:29:45',NULL,NULL);
INSERT INTO `users` VALUES (355,'Inoue Shino','shino63@outlook.com','wVBoxRbytp','yuzRNpIMph','student','Inoue Shino','每一天都值得全力以赴。',NULL,882,'2021-04-02 00:29:25',NULL,NULL);
INSERT INTO `users` VALUES (356,'Ichikawa Itsuki','ichikawaitsuki@gmail.com','pT2BC4Tpr6','4BZh3TBsHK','admin','Ichikawa Itsuki','走自己喜欢的路，做自己想做的人。',NULL,309,'2009-10-23 13:55:08',NULL,NULL);
INSERT INTO `users` VALUES (357,'Lok Lai Yan','lolaiya@outlook.com','qxXyhDuKTq','irsgQtzSML','teacher','Lok Lai Yan','时光匆匆，未来在脚下。',NULL,85,'2018-04-21 02:08:36',NULL,NULL);
INSERT INTO `users` VALUES (358,'Tse Wai Man','waimants@gmail.com','kWY4wgaysv','erMf382BDU','teacher','Tse Wai Man','永远不放弃自己，永远坚持梦想。',NULL,102,'2000-10-03 06:55:17',NULL,NULL);
INSERT INTO `users` VALUES (359,'Fujiwara Kenta','kefujiwara@outlook.com','02QUPUTg9i','ZG0u6KUWCB','student','Fujiwara Kenta','路有多远，心有多大，未来就有多精彩。',NULL,909,'2012-10-23 03:23:25',NULL,NULL);
INSERT INTO `users` VALUES (360,'Ichikawa Misaki','misaki2@hotmail.com','Wd1BRyRImY','ppAewmItjX','teacher','Ichikawa Misaki','永远勇敢，永远坚持，永远前行。',NULL,460,'2023-10-26 15:55:47',NULL,NULL);
INSERT INTO `users` VALUES (361,'Nakagawa Akina','akina2@gmail.com','Mi1nP5bJcG','RC9XUjUVVQ','admin','Nakagawa Akina','一切都从今天开始，未来更值得期待。',NULL,351,'2003-03-28 22:33:14',NULL,NULL);
INSERT INTO `users` VALUES (362,'Jerry Hughes','jerry59@mail.com','xKFF8yOgyp','dlJj2UX1cJ','admin','Jerry Hughes','梦想从未离开，一直在前方等着你。',NULL,789,'2001-05-13 10:26:43',NULL,NULL);
INSERT INTO `users` VALUES (363,'Gong Yunxi','yunxi18@mail.com','7WHAr6lRHm','dInBcrSWIe','admin','Gong Yunxi','人生没有捷径，只有不断努力。',NULL,795,'2022-07-26 15:13:28',NULL,NULL);
INSERT INTO `users` VALUES (364,'Ono Kazuma','kazuma215@icloud.com','1DFS0mbyS1','2e9voX5Aq5','admin','Ono Kazuma','每一次进步，都是一个新的开始。',NULL,641,'2005-02-05 23:12:15',NULL,NULL);
INSERT INTO `users` VALUES (365,'Kaneko Mitsuki','mkane68@gmail.com','CdHhaGd9Fx','pf9X9RFTgq','admin','Kaneko Mitsuki','相信自己，你就是最好的自己。',NULL,896,'2011-01-29 17:42:04',NULL,NULL);
INSERT INTO `users` VALUES (366,'Duan Lu','lud330@outlook.com','cV6fpG0hAV','kWmeJh0njG','teacher','Duan Lu','每一步都算数，终点就在前方。',NULL,546,'2022-01-30 22:28:24',NULL,NULL);
INSERT INTO `users` VALUES (367,'Kong Ka Man','kongkm1@outlook.com','zpfaUh3B95','k2ULOXkVO5','student','Kong Ka Man','生活需要一点点勇气，梦想需要满满的坚持。',NULL,481,'2021-04-20 18:34:10',NULL,NULL);
INSERT INTO `users` VALUES (368,'Ralph Harris','harrisralph@gmail.com','A04ijfTgxd','NhFRvXGntP','student','Ralph Harris','向阳而生，笑对人生。',NULL,743,'2004-10-24 22:24:36',NULL,NULL);
INSERT INTO `users` VALUES (369,'Tao Kar Yan','taokaryan6@mail.com','MpaQph4QJs','rXcdfqgTc7','teacher','Tao Kar Yan','路在脚下，梦在前方。',NULL,985,'2005-04-07 02:36:40',NULL,NULL);
INSERT INTO `users` VALUES (370,'Anne Evans','evansanne10@gmail.com','3c2tl4AkVd','q9BoneMTTJ','teacher','Anne Evans','让心自由，让未来不设限。',NULL,716,'2021-02-17 09:58:41',NULL,NULL);
INSERT INTO `users` VALUES (371,'Tang Rui','rut@yahoo.com','r8EHrmzE5T','DtPmRPi2aQ','admin','Tang Rui','永远怀抱希望，勇敢去追。',NULL,418,'2003-07-07 05:18:45',NULL,NULL);
INSERT INTO `users` VALUES (372,'Kojima Seiko','seiko511@gmail.com','GBMFPWDLJu','8mbaBOHZab','student','Kojima Seiko','坚持自我，才是最好的活法。',NULL,426,'2010-05-05 09:30:13',NULL,NULL);
INSERT INTO `users` VALUES (373,'Choi Wing Fat','choiwf@gmail.com','PEVHQP6jam','E5lLZOSUWP','student','Choi Wing Fat','梦想不会辜负每一个努力的灵魂。',NULL,61,'2001-01-09 20:49:03',NULL,NULL);
INSERT INTO `users` VALUES (374,'Yuen Ching Wan','yuecw1961@outlook.com','hgCs8BAWmq','hS2BbAU6ga','admin','Yuen Ching Wan','用心过每一天，追逐每一个梦想。',NULL,263,'2005-09-21 20:46:14',NULL,NULL);
INSERT INTO `users` VALUES (375,'Shi Lu','lushi@hotmail.com','tBsL9OEdfi','PPpe2ZLbqo','student','Shi Lu','做自己喜欢的事，爱自己想爱的人。',NULL,831,'2008-05-17 04:06:18',NULL,NULL);
INSERT INTO `users` VALUES (376,'Evelyn Washington','washington411@outlook.com','OwzJQWdOqa','C3VaCA14nu','student','Evelyn Washington','做最好的自己，生活更精彩。',NULL,828,'2002-04-25 04:36:44',NULL,NULL);
INSERT INTO `users` VALUES (377,'Yamazaki Momoe','myama@hotmail.com','feqGPnsJ28','YXvhxyOhtG','admin','Yamazaki Momoe','勇敢追逐，不惧未来。',NULL,421,'2020-10-05 08:16:48',NULL,NULL);
INSERT INTO `users` VALUES (378,'Nakayama Yuito','yuitonakayama6@gmail.com','s6Q1s3DRFC','MeE4j5QHHi','teacher','Nakayama Yuito','成长就是不断超越自我。',NULL,304,'2014-07-20 11:23:52',NULL,NULL);
INSERT INTO `users` VALUES (379,'Liang Jiehong','jiehongliang@mail.com','7iUfeAfp7G','g7SK65sim1','admin','Liang Jiehong','每一个微笑，都是对生活的赞美。',NULL,599,'2022-09-12 17:47:45',NULL,NULL);
INSERT INTO `users` VALUES (380,'Liu Jiehong','liuj7@gmail.com','50ch7ARmsM','Qvwj3kARU2','admin','Liu Jiehong','为梦想坚持，为未来拼搏。',NULL,677,'2004-02-27 23:57:35',NULL,NULL);
INSERT INTO `users` VALUES (381,'Feng Xiuying','fexi2@icloud.com','cQYXoRWSNC','QEbOi5z7Rx','admin','Feng Xiuying','不急不躁，活得从容自在。',NULL,770,'2010-06-14 01:40:19',NULL,NULL);
INSERT INTO `users` VALUES (382,'Takagi Mio','miotak5@icloud.com','fdtxNF0c6P','OGUHubTByi','teacher','Takagi Mio','世界如此美好，你我都值得拥有。',NULL,509,'2004-01-16 21:39:22',NULL,NULL);
INSERT INTO `users` VALUES (383,'Rhonda Gordon','rhonda2013@gmail.com','D1zO3cFcNR','2MZrMH1eHM','teacher','Rhonda Gordon','用勇气去迎接每一个挑战。',NULL,324,'2005-03-16 14:46:46',NULL,NULL);
INSERT INTO `users` VALUES (384,'Rhonda West','rhondawest@gmail.com','c1u6gSnbf4','pz8qYEDwvj','student','Rhonda West','做一个有温度的自己。',NULL,100,'2002-09-01 01:31:05',NULL,NULL);
INSERT INTO `users` VALUES (385,'Yao Zitao','yaozitao@outlook.com','XWi1BToVV3','PUiC84NwPP','teacher','Yao Zitao','每一天都是新的一天，充满希望。',NULL,627,'2019-04-05 01:59:05',NULL,NULL);
INSERT INTO `users` VALUES (386,'Daniel Salazar','salazar99@yahoo.com','SbkhzLtSHd','LH27wyndEI','student','Daniel Salazar','追逐梦想，不言放弃。',NULL,325,'2006-11-08 08:49:31',NULL,NULL);
INSERT INTO `users` VALUES (387,'Josephine Bell','bell315@mail.com','ZiLioYEFHw','FutAUb0eO7','teacher','Josephine Bell','世界如此之大，不去看看多可惜。',NULL,787,'2000-01-19 07:57:26',NULL,NULL);
INSERT INTO `users` VALUES (388,'Wada Minato','waminato@icloud.com','YIJwYBHZr5','EEgYyoMKcD','teacher','Wada Minato','勇敢追梦，走自己想走的路。',NULL,5,'2008-09-07 21:37:52',NULL,NULL);
INSERT INTO `users` VALUES (389,'Ogawa Daisuke','ogawad@hotmail.com','QIMjo396fR','HwiluxK3iJ','admin','Ogawa Daisuke','做自己的英雄，创造属于自己的未来。',NULL,981,'2019-11-25 00:48:01',NULL,NULL);
INSERT INTO `users` VALUES (390,'Curtis Grant','grantcurtis@outlook.com','1xQzqMDxda','wpNy9RtCNe','student','Curtis Grant','每一次奋斗，都是向目标靠近的一步。',NULL,22,'2005-03-29 00:24:41',NULL,NULL);
INSERT INTO `users` VALUES (391,'Kato Momoka','momokakato416@icloud.com','CbYgAVBtUn','vm1uTDEx4Q','teacher','Kato Momoka','不断追求，不断超越。',NULL,95,'2001-10-11 10:49:19',NULL,NULL);
INSERT INTO `users` VALUES (392,'April Young','april12@outlook.com','5zekiaiAtd','rQE9XWgS8k','teacher','April Young','梦想照进现实，未来在努力中成型。',NULL,548,'2012-12-22 18:23:19',NULL,NULL);
INSERT INTO `users` VALUES (393,'Zhang Zitao','zhangz6@gmail.com','s0tIm3QvjZ','HtCEThwKUy','student','Zhang Zitao','青春不再，奋斗未停。',NULL,724,'2019-04-20 18:42:56',NULL,NULL);
INSERT INTO `users` VALUES (394,'Cho Ho Yin','hoyin213@outlook.com','ob5sw67BPZ','WCTD5aPYiZ','student','Cho Ho Yin','未来在于今天的努力。',NULL,444,'2023-05-28 14:21:10',NULL,NULL);
INSERT INTO `users` VALUES (395,'Sandra Alvarez','sandralvar41@outlook.com','THeGXghx0i','zdlxbx4QwO','teacher','Sandra Alvarez','不忘初心，继续前行。',NULL,403,'2017-04-12 15:02:33',NULL,NULL);
INSERT INTO `users` VALUES (396,'Otsuka Yamato','yamatotsuka314@gmail.com','7vld3wQrTo','pDsCNnlVEd','student','Otsuka Yamato','永远做最真实的自己。',NULL,199,'2023-10-02 07:46:00',NULL,NULL);
INSERT INTO `users` VALUES (397,'Sit On Kay','sok@hotmail.com','mo0hmNqnOh','uC6zNzsiwK','student','Sit On Kay','为梦想努力，人生因此精彩。',NULL,343,'2003-05-02 22:13:28',NULL,NULL);
INSERT INTO `users` VALUES (398,'Wu Wai Han','wuwaihan8@icloud.com','Dmmb5nJ21m','4zNxlN63MI','student','Wu Wai Han','每一天都是新的起点，新的希望。',NULL,339,'2024-11-07 21:32:00',NULL,NULL);
INSERT INTO `users` VALUES (399,'Chad Ramirez','rachad@hotmail.com','ySw0G0zrJ4','IPMkmnJRJ1','admin','Chad Ramirez','不畏将来，不念过往。',NULL,633,'2020-06-01 05:21:01',NULL,NULL);
INSERT INTO `users` VALUES (400,'Zeng Lan','lanzen@outlook.com','sJKWij2VOo','ADzQjYkcsF','student','Zeng Lan','未来值得期待，生活从不言败。',NULL,639,'2024-07-17 13:10:41',NULL,NULL);
INSERT INTO `users` VALUES (401,'Wong Ho Yin','hywo@yahoo.com','CBlsrtOVzl','vV14AbstpM','admin','Wong Ho Yin','每一步都走得坚定，未来可期。',NULL,611,'2014-04-13 15:50:49',NULL,NULL);
INSERT INTO `users` VALUES (402,'Cheng Ho Yin','chenghoyin@outlook.com','IAV17WsvJp','noJaInTw12','teacher','Cheng Ho Yin','不为过去的自己后悔，只为未来的自己努力。',NULL,472,'2023-09-25 18:17:37',NULL,NULL);
INSERT INTO `users` VALUES (403,'Ho On Kay','hoonkay@gmail.com','ywlLCkTInv','8awEqcALPE','teacher','Ho On Kay','梦想太远，努力刚好。',NULL,367,'2011-03-26 08:47:22',NULL,NULL);
INSERT INTO `users` VALUES (404,'Chung Yu Ling','chuyl@icloud.com','ranGFW8jyW','87uxCQcRYj','admin','Chung Yu Ling','用微笑迎接每一天。',NULL,757,'2020-07-14 09:33:47',NULL,NULL);
INSERT INTO `users` VALUES (405,'Ng Ling Ling','linglingn18@outlook.com','3eYdjbDEfK','U7oiRM79J9','teacher','Ng Ling Ling','路上有风景，心里有梦想。',NULL,76,'2014-02-01 01:22:07',NULL,NULL);
INSERT INTO `users` VALUES (406,'Antonio Howard','antoniohowar1@icloud.com','tCcIFs5jP4','qEDHR3kLjT','admin','Antonio Howard','永远相信，最好的还在前方。',NULL,156,'2015-03-26 00:28:37',NULL,NULL);
INSERT INTO `users` VALUES (407,'Sakurai Kazuma','ks420@gmail.com','prytRAUdqH','YE8ehm2CNr','student','Sakurai Kazuma','你可以慢，但不能停。',NULL,240,'2019-06-03 12:14:09',NULL,NULL);
INSERT INTO `users` VALUES (408,'Sandra Richardson','richardsonsandra@mail.com','5TRCltjszP','S0Nm5TV7W1','admin','Sandra Richardson','只要心中有梦想，脚步永远不停。',NULL,386,'2020-12-26 07:40:53',NULL,NULL);
INSERT INTO `users` VALUES (409,'Kathryn Hughes','kathrynhughes@hotmail.com','4jx75p9426','nJgczySSrp','teacher','Kathryn Hughes','生活充满未知，但我愿迎接每一场挑战。',NULL,821,'2010-02-11 06:47:10',NULL,NULL);
INSERT INTO `users` VALUES (410,'Kam Kwok Ming','kmkam@hotmail.com','TlN8uqWEgH','r9e5TXf92g','student','Kam Kwok Ming','做一个快乐的人，活出自己最好的样子。',NULL,2,'2010-04-14 05:19:35',NULL,NULL);
INSERT INTO `users` VALUES (411,'Maeda Takuya','maedat1004@gmail.com','vbWfwksKSz','iI4MjY9UYq','admin','Maeda Takuya','不为他人活，只为自己梦想。',NULL,898,'2007-01-18 22:05:47',NULL,NULL);
INSERT INTO `users` VALUES (412,'Hasegawa Hikari','hikahaseg1@icloud.com','HkK6b4P4Iy','SHGd8jiIrZ','admin','Hasegawa Hikari','未来属于那些敢于改变的人。',NULL,660,'2014-10-23 12:40:22',NULL,NULL);
INSERT INTO `users` VALUES (413,'Tin Chun Yu','tinchunyu2004@icloud.com','HAM7kHaxGJ','URIj9csg7M','student','Tin Chun Yu','做一个有意义的灵魂。',NULL,541,'2012-08-15 16:48:31',NULL,NULL);
INSERT INTO `users` VALUES (414,'Zhu Jiehong','jiezhu@outlook.com','Eu0uu5ahsH','BcLvU7cpIA','teacher','Zhu Jiehong','别怕路远，终点会更美。',NULL,625,'2009-05-24 11:08:30',NULL,NULL);
INSERT INTO `users` VALUES (415,'Liu Shihan','lishih@hotmail.com','8P0nuEWQgF','Bfv4mJMSH9','student','Liu Shihan','让每一天都成为人生的新起点。',NULL,94,'2016-11-27 05:01:38',NULL,NULL);
INSERT INTO `users` VALUES (416,'Maeda Kazuma','kazumamaeda7@icloud.com','xw4qY21aft','L4ypSYpeVi','admin','Maeda Kazuma','心怀梦想，脚步坚定。',NULL,517,'2007-10-24 04:07:19',NULL,NULL);
INSERT INTO `users` VALUES (417,'Wang Anqi','anqiw820@gmail.com','Gwr8QB3tZF','cyvcXxhnrp','admin','Wang Anqi','走出舒适区，迎接不一样的未来。',NULL,963,'2022-04-14 00:42:58',NULL,NULL);
INSERT INTO `users` VALUES (418,'Ichikawa Akina','iakina53@mail.com','dtiKrSEMNL','sPsugt2NXE','teacher','Ichikawa Akina','努力的意义，是为了更好的明天。',NULL,53,'2006-07-24 23:09:31',NULL,NULL);
INSERT INTO `users` VALUES (419,'Su Jialun','su93@mail.com','9pe1aiBgli','KHX2ofEATG','student','Su Jialun','生活是一场冒险，我愿勇敢前行。',NULL,11,'2024-07-30 14:14:37',NULL,NULL);
INSERT INTO `users` VALUES (420,'Tse Chung Yin','chungyintse@outlook.com','ORqH0OFfBp','hHXn25FgFR','student','Tse Chung Yin','用一颗坚定的心，去追求梦想。',NULL,787,'2001-10-14 16:45:03',NULL,NULL);
INSERT INTO `users` VALUES (421,'Tam Wing Sze','tamws8@icloud.com','2Q1LJWtvY2','cpVOQQCYso','admin','Tam Wing Sze','每一步的努力，都是通往梦想的桥梁。',NULL,951,'2019-09-30 23:25:53',NULL,NULL);
INSERT INTO `users` VALUES (422,'Li Anqi','lian@icloud.com','wFHjLN4xoM','8sxuAPkIXL','admin','Li Anqi','不为别人而活，只为自己坚持。',NULL,810,'2019-07-15 06:26:29',NULL,NULL);
INSERT INTO `users` VALUES (423,'Shawn Stevens','steshaw@outlook.com','ciuLFCqpn9','gJkDhLQUMf','teacher','Shawn Stevens','拥抱未来，告别过去。',NULL,175,'2016-10-27 11:22:41',NULL,NULL);
INSERT INTO `users` VALUES (424,'Matsui Mitsuki','mitsukim@outlook.com','q8U3KR2Qae','xLIE8jVT5X','teacher','Matsui Mitsuki','让梦想成为前行的动力。',NULL,698,'2003-01-09 09:40:47',NULL,NULL);
INSERT INTO `users` VALUES (425,'Shen Jiehong','shen1986@gmail.com','eyCdIjVOv3','XIxOtqyslw','admin','Shen Jiehong','一颗勇敢的心，走向光明的未来。',NULL,61,'2017-05-14 21:18:12',NULL,NULL);
INSERT INTO `users` VALUES (426,'Kevin Ross','rok@gmail.com','kBAuRPNHDQ','VJqUiLMoQm','admin','Kevin Ross','活得简单，做自己最真实的样子。',NULL,361,'2014-02-07 11:28:06',NULL,NULL);
INSERT INTO `users` VALUES (427,'Fujiwara Airi','afujiw809@hotmail.com','6hToIkZ1PH','2AnM4EOxXN','admin','Fujiwara Airi','未来是由今天的努力决定的。',NULL,500,'2012-06-11 23:18:01',NULL,NULL);
INSERT INTO `users` VALUES (428,'Wada Momoka','mw1@yahoo.com','PChhrXOZXV','3Y1Q55ow4j','admin','Wada Momoka','从不惧怕艰难，因为我相信未来的自己。',NULL,823,'2017-03-07 12:17:50',NULL,NULL);
INSERT INTO `users` VALUES (429,'Hara Ren','ren2004@icloud.com','l9waF26Cbo','YOWfH0FWni','admin','Hara Ren','不再迷茫，向着梦想出发。',NULL,62,'2001-01-08 21:54:47',NULL,NULL);
INSERT INTO `users` VALUES (430,'Fong Ho Yin','hoyin1129@yahoo.com','DF8hIEP25d','rnrfSWWXDj','teacher','Fong Ho Yin','一直在前行，直到抵达梦想的彼岸。',NULL,208,'2000-01-08 17:01:03',NULL,NULL);
INSERT INTO `users` VALUES (431,'Ando Kaito','kaito3@gmail.com','bOejxG9K5v','9uOoLuFwfg','admin','Ando Kaito','用微笑战胜一切。',NULL,268,'2018-01-21 12:39:56',NULL,NULL);
INSERT INTO `users` VALUES (432,'Lisa Grant','grantlisa@icloud.com','iGvLqSMYfx','eRKqspZH4C','teacher','Lisa Grant','不管未来如何，我都能勇敢面对。',NULL,952,'2022-01-17 18:37:14',NULL,NULL);
INSERT INTO `users` VALUES (433,'Lam Siu Wai','siuwlam@hotmail.com','OYy6zEdw4l','LhuAvuJy3e','teacher','Lam Siu Wai','不忘初心，方得始终。',NULL,176,'2015-05-26 22:46:39',NULL,NULL);
INSERT INTO `users` VALUES (434,'Liu Xiaoming','xliu@outlook.com','RtTt2gSPYr','eqKrpp9C3R','teacher','Liu Xiaoming','每天进步一点点，梦想更进一步。',NULL,527,'2003-10-27 14:54:30',NULL,NULL);
INSERT INTO `users` VALUES (435,'Harry Roberts','roberha@mail.com','CAwqM1aW2h','9erG8RvNz8','student','Harry Roberts','不为明天忧虑，活在当下。',NULL,623,'2011-02-02 19:07:26',NULL,NULL);
INSERT INTO `users` VALUES (436,'Susan Jones','susanjones1@gmail.com','7BjiId38QW','N4NtkRq6Xg','student','Susan Jones','梦想永远不会放弃坚持的人。',NULL,66,'2006-02-19 07:46:02',NULL,NULL);
INSERT INTO `users` VALUES (437,'Tong Tsz Ching','tong54@icloud.com','Pfzwd0EF4M','aJWiITjxwA','teacher','Tong Tsz Ching','每天给自己一个微笑，迎接新的挑战。',NULL,312,'2016-07-07 05:46:54',NULL,NULL);
INSERT INTO `users` VALUES (438,'Jamie Murphy','jm2001@yahoo.com','M7PlDe2gwZ','0SuAaxr8kM','admin','Jamie Murphy','世界没有想象的那么大，梦想也没有那么远。',NULL,973,'2003-02-01 01:28:05',NULL,NULL);
INSERT INTO `users` VALUES (439,'Mok Wing Suen','wingsuenmo@icloud.com','QLuiVxEZhY','xBlYJgVqlt','student','Mok Wing Suen','走自己的路，让别人说去吧。',NULL,136,'2018-03-17 11:30:19',NULL,NULL);
INSERT INTO `users` VALUES (440,'Xiang Ziyi','zixian@mail.com','DdNavb5L5C','NpTlJStrIX','student','Xiang Ziyi','生活是一场马拉松，坚持就是胜利。',NULL,587,'2008-11-01 11:08:36',NULL,NULL);
INSERT INTO `users` VALUES (441,'Matthew Allen','amatthew3@hotmail.com','1QMUHxCibj','EUPJnCb2n6','teacher','Matthew Allen','走得慢不怕，最怕的是停下。',NULL,522,'2019-12-19 20:50:25',NULL,NULL);
INSERT INTO `users` VALUES (442,'Kojima Aoshi','aokojima@mail.com','5AUUZmQpql','WlKk1CijUp','teacher','Kojima Aoshi','永远相信自己，未来一定会更好。',NULL,428,'2010-06-21 22:34:32',NULL,NULL);
INSERT INTO `users` VALUES (443,'Tamura Miu','tammiu529@gmail.com','dhklRpLhBR','x8MeGxFRid','teacher','Tamura Miu','人生的意义就在于追逐与超越。',NULL,195,'2016-07-05 16:51:45',NULL,NULL);
INSERT INTO `users` VALUES (444,'Shimizu Momoka','sm228@mail.com','7yuqcpE8iS','q7kAKAAGYi','teacher','Shimizu Momoka','心若向阳，永不言败。',NULL,323,'2011-10-28 00:24:08',NULL,NULL);
INSERT INTO `users` VALUES (445,'Doris Nelson','nelsondoris@gmail.com','sNydyXERia','5V4v9qOBEz','teacher','Doris Nelson','勇敢去追，错过了不后悔。',NULL,6,'2024-09-15 14:10:55',NULL,NULL);
INSERT INTO `users` VALUES (446,'Crystal Salazar','salcrystal@yahoo.com','YvEStM0Sgk','jQgQ1zK8Qp','teacher','Crystal Salazar','不想当懒人，只想当梦想的追随者。',NULL,688,'2020-04-18 01:44:06',NULL,NULL);
INSERT INTO `users` VALUES (447,'Koon Tsz Ching','tckoo@outlook.com','TtrLIqBRVa','JKTw4CWNwO','teacher','Koon Tsz Ching','做自己喜欢做的事，过自己喜欢的生活。',NULL,938,'2019-10-30 20:27:16',NULL,NULL);
INSERT INTO `users` VALUES (448,'Lam Chi Ming','lcm41@icloud.com','sic1er35sH','gbrmzEAz1a','admin','Lam Chi Ming','生活不需要完美，只要真实。',NULL,750,'2008-03-20 01:32:31',NULL,NULL);
INSERT INTO `users` VALUES (449,'Dai Yunxi','dayun@outlook.com','vvojOC2YRI','QiPcWDej7n','student','Dai Yunxi','为自己努力，为梦想拼搏。',NULL,688,'2003-10-22 07:37:11',NULL,NULL);
INSERT INTO `users` VALUES (450,'Duan Jiehong','jdu@outlook.com','VVBzoydvw6','ifUHf72fbl','teacher','Duan Jiehong','每天都是全新的开始，充满无限可能。',NULL,171,'2008-02-04 13:05:13',NULL,NULL);
INSERT INTO `users` VALUES (451,'Jeremy Gardner','jga718@yahoo.com','tkmw5iekir','9HFnANPmji','admin','Jeremy Gardner','不为过去懊悔，向未来努力。',NULL,808,'2004-05-31 00:05:52',NULL,NULL);
INSERT INTO `users` VALUES (452,'Norma Jimenez','jimenez10@outlook.com','Jh4kfEzXgx','rdb9cfmwAU','student','Norma Jimenez','活在当下，尽力去爱。',NULL,925,'2013-03-26 12:32:43',NULL,NULL);
INSERT INTO `users` VALUES (453,'Sugawara Yuna','sugawarayu1@outlook.com','JIHytE7dB6','ukeVrYOZJO','student','Sugawara Yuna','把握现在，创造未来。',NULL,517,'2002-05-17 17:05:45',NULL,NULL);
INSERT INTO `users` VALUES (454,'Travis Gonzales','gonzatra40@gmail.com','S7U2tWrqj6','GbXTaTvc2M','admin','Travis Gonzales','做一个温暖的人，带给自己和别人阳光。',NULL,846,'2021-05-09 15:57:39',NULL,NULL);
INSERT INTO `users` VALUES (455,'Zheng Xiaoming','xiaoming3@yahoo.com','JlhmuBWwyG','dzF0MfSiDs','admin','Zheng Xiaoming','不为别人活，只为自己心中的梦想。',NULL,824,'2003-11-01 04:25:48',NULL,NULL);
INSERT INTO `users` VALUES (456,'Iwasaki Sara','iwasakisara@gmail.com','XPt8wJLt4s','JLLgXKALrP','student','Iwasaki Sara','青春不再，奋斗依旧。',NULL,889,'2002-03-10 02:37:14',NULL,NULL);
INSERT INTO `users` VALUES (457,'Sheh Chi Yuen','cys@hotmail.com','AOlEVDRGQ8','WbOkssL1fs','teacher','Sheh Chi Yuen','每一次失败，都是成功的积累。',NULL,750,'2015-03-30 03:28:21',NULL,NULL);
INSERT INTO `users` VALUES (458,'Qiu Lan','qiulan720@hotmail.com','DyUkiwEZJK','Soj3JVbRae','student','Qiu Lan','永远不放弃自己，相信未来会更好。',NULL,666,'2013-06-18 12:49:33',NULL,NULL);
INSERT INTO `users` VALUES (459,'Janice Richardson','richardsonjanic40@yahoo.com','j8R7hRHL1J','b9zYjVdAa4','student','Janice Richardson','用心去活，做自己喜欢的事。',NULL,932,'2021-02-13 03:07:58',NULL,NULL);
INSERT INTO `users` VALUES (460,'Nomura Sakura','sakura1966@hotmail.com','4ZGBm1qVxZ','q0LAurFwxJ','admin','Nomura Sakura','路虽远，行则必至。',NULL,152,'2013-12-04 07:02:21',NULL,NULL);
INSERT INTO `users` VALUES (461,'Qin Zhennan','zq87@gmail.com','sf51syT78m','6lifjyKOrR','teacher','Qin Zhennan','每一份坚持，都是成功的前提。',NULL,901,'2017-08-02 02:57:10',NULL,NULL);
INSERT INTO `users` VALUES (462,'Xie Shihan','xishiha@gmail.com','1rQfKI4vvC','2PQWYA6byQ','teacher','Xie Shihan','努力拼搏，梦想就在前方。',NULL,454,'2024-05-11 16:28:03',NULL,NULL);
INSERT INTO `users` VALUES (463,'Yamazaki Nanami','nanamiyamazaki726@gmail.com','1Wf8RgpC3J','aDqhlIm9oE','admin','Yamazaki Nanami','让梦想照进现实，迎接新的挑战。',NULL,33,'2004-01-06 05:27:09',NULL,NULL);
INSERT INTO `users` VALUES (464,'Tan Jiehong','tanjieh@icloud.com','zRJbGIajZ2','vobk6BimJX','teacher','Tan Jiehong','不为结果焦虑，享受过程。',NULL,127,'2002-12-21 06:38:12',NULL,NULL);
INSERT INTO `users` VALUES (465,'Peng Lan','lapen@gmail.com','GFVfHMOasK','hv9FixoC9r','teacher','Peng Lan','无论何时，勇敢向前。',NULL,428,'2000-12-22 01:48:04',NULL,NULL);
INSERT INTO `users` VALUES (466,'Ma Tin Wing','matw219@outlook.com','kWYFbPsOH9','HkPOuQA0OT','admin','Ma Tin Wing','用笑容面对世界，心中充满阳光。',NULL,590,'2002-01-20 15:04:11',NULL,NULL);
INSERT INTO `users` VALUES (467,'Yu Yuning','yuning6@icloud.com','5pEssWHJ3q','bksL4x0h0U','teacher','Yu Yuning','每天都是新的一天，充满新的希望。',NULL,50,'2014-12-25 14:31:43',NULL,NULL);
INSERT INTO `users` VALUES (468,'Lei Ziyi','ziyi130@hotmail.com','pbES60uEFq','HDrRka0IE2','student','Lei Ziyi','永远相信，自己有无限的可能。',NULL,569,'2014-02-12 07:50:44',NULL,NULL);
INSERT INTO `users` VALUES (469,'Chic Tsz Hin','chicth@mail.com','XMnVssZKDF','dCOOufMbtv','teacher','Chic Tsz Hin','生活的意义在于不断进步。',NULL,242,'2012-05-30 22:48:37',NULL,NULL);
INSERT INTO `users` VALUES (470,'Chad Young','youngchad@icloud.com','OXyjNCXe9i','0PITGMR4US','student','Chad Young','路的尽头是光明，继续向前。',NULL,435,'2006-06-02 11:18:23',NULL,NULL);
INSERT INTO `users` VALUES (471,'Yao Lan','yal@outlook.com','FDduoEv86K','0AbWCAxGJu','teacher','Yao Lan','成功在于坚持，幸福在于努力。',NULL,406,'2008-03-15 10:10:21',NULL,NULL);
INSERT INTO `users` VALUES (472,'Annie Perry','annie1962@icloud.com','1t1EzSYH51','uOp2xRjLxn','student','Annie Perry','人生最精彩的部分，往往是起点。',NULL,466,'2011-07-16 04:29:27',NULL,NULL);
INSERT INTO `users` VALUES (473,'Tang Yunxi','yunxitang6@gmail.com','7x4wnqFwof','Ddo6KlwxkI','admin','Tang Yunxi','每一次努力，都值得骄傲。',NULL,651,'2019-03-13 17:33:07',NULL,NULL);
INSERT INTO `users` VALUES (474,'Fan Anqi','fan8@gmail.com','6uBIFYK7yo','WM9IctgEa6','admin','Fan Anqi','用努力铺路，梦想照亮前行。',NULL,304,'2004-10-15 09:38:31',NULL,NULL);
INSERT INTO `users` VALUES (475,'Zou Rui','ruizou1971@hotmail.com','TVqFuNWoJa','GbTw3xJDPw','admin','Zou Rui','不后悔做过的事，不畏将来。',NULL,203,'2005-02-15 17:32:26',NULL,NULL);
INSERT INTO `users` VALUES (476,'Hsuan Ting Fung','tingfunghsuan@yahoo.com','4dgBr5PWfE','qbQ4krtmju','teacher','Hsuan Ting Fung','每天坚持一点点，未来就更精彩。',NULL,244,'2016-11-15 06:20:54',NULL,NULL);
INSERT INTO `users` VALUES (477,'Kato Aoi','katoaoi308@gmail.com','83Et8IMFnY','gLqTjntL25','student','Kato Aoi','做自己喜欢的事，过自己想过的生活。',NULL,215,'2000-02-03 06:35:04',NULL,NULL);
INSERT INTO `users` VALUES (478,'Lori Edwards','edwards2@hotmail.com','7wGwmym7no','bg52CulcL2','admin','Lori Edwards','勇敢追梦，坚持到底。',NULL,647,'2017-10-25 12:29:37',NULL,NULL);
INSERT INTO `users` VALUES (479,'Tam Chi Ming','tam1123@outlook.com','vjbA2uoJuJ','EwLpcKxcjo','admin','Tam Chi Ming','不怕风雨，只怕放弃。',NULL,705,'2004-02-14 04:09:02',NULL,NULL);
INSERT INTO `users` VALUES (480,'Diana Sanders','sanders1001@outlook.com','uMvRAHlpQl','HIhItoGGDt','teacher','Diana Sanders','世界因你而美丽，你因梦想而精彩。',NULL,65,'2020-02-22 08:31:33',NULL,NULL);
INSERT INTO `users` VALUES (481,'Ono Itsuki','itsukiono@gmail.com','v6ZvdgLUc3','sm6xtrtVLC','student','Ono Itsuki','每一次努力，都是成功的积累。',NULL,373,'2002-10-30 04:15:52',NULL,NULL);
INSERT INTO `users` VALUES (482,'Wanda Morris','wandamorris@mail.com','1RBFFq9RET','7Lwn8RpGSz','teacher','Wanda Morris','无论如何，不放弃任何一个追求。',NULL,474,'2014-12-12 04:59:08',NULL,NULL);
INSERT INTO `users` VALUES (483,'Mori Yota','yotamor@gmail.com','OrH9YsVi6J','DtlfvvxO4x','admin','Mori Yota','梦想不大，奋斗不止。',NULL,30,'2020-07-05 22:03:27',NULL,NULL);
INSERT INTO `users` VALUES (484,'Alice Thomas','athomas@gmail.com','jQq0o3j55T','KfAe7B9ruZ','teacher','Alice Thomas','不畏将来，不念过往。',NULL,358,'2022-02-17 09:26:08',NULL,NULL);
INSERT INTO `users` VALUES (485,'Goto Sara','sarag5@yahoo.com','F6qo6xwws8','5rF9HT31Mo','admin','Goto Sara','未来属于那些敢于追梦的人。',NULL,694,'2004-03-30 07:55:38',NULL,NULL);
INSERT INTO `users` VALUES (486,'Micheal Morgan','michealmor1999@icloud.com','cRabCAdcHa','K3Bl60HLQN','teacher','Micheal Morgan','青春就是勇敢追逐梦想的时光。',NULL,695,'2021-01-10 13:21:42',NULL,NULL);
INSERT INTO `users` VALUES (487,'Xiong Xiaoming','xiongx@hotmail.com','LTadrfyRj6','VjSDoe4Wck','student','Xiong Xiaoming','做自己想做的事，活成最好的自己。',NULL,322,'2021-10-24 13:17:59',NULL,NULL);
INSERT INTO `users` VALUES (488,'Ishikawa Momoe','ishikawamo10@outlook.com','dksJrbYqtH','08i30nQtN5','student','Ishikawa Momoe','不怕路远，只怕心不坚。',NULL,886,'2014-10-12 08:46:36',NULL,NULL);
INSERT INTO `users` VALUES (489,'Dawn Ramos','dawnram@yahoo.com','FwwzrUyVSw','2ti2FRjgnf','teacher','Dawn Ramos','未来的自己，一定会感谢现在拼搏的你。',NULL,720,'2008-08-06 08:39:18',NULL,NULL);
INSERT INTO `users` VALUES (490,'Joan Myers','mjoan@outlook.com','Pd1auUas40','ozHvDDFOL2','student','Joan Myers','没有最好的时光，只有努力的时光。',NULL,905,'2011-07-10 16:40:56',NULL,NULL);
INSERT INTO `users` VALUES (491,'Sato Akina','sak91@outlook.com','C1k63xZEBG','RZ1SXQrpYs','admin','Sato Akina','梦想从不嫌晚，努力才是关键。',NULL,765,'2003-03-03 18:02:10',NULL,NULL);
INSERT INTO `users` VALUES (492,'Dai Chi Ming','dchiming@gmail.com','RVXOqxyni5','wkMSfWKNPE','teacher','Dai Chi Ming','不为别的，只为更好的自己。',NULL,738,'2007-10-17 10:24:01',NULL,NULL);
INSERT INTO `users` VALUES (493,'Peter Reyes','reyes4@hotmail.com','EnMCji5qtd','gI6ncpxWyc','admin','Peter Reyes','每一天的坚持，都是未来的收获。',NULL,503,'2016-06-08 01:18:43',NULL,NULL);
INSERT INTO `users` VALUES (494,'Lui Ka Keung','kakeung8@yahoo.com','tNlzHvhdBo','zdB6zrDk7x','admin','Lui Ka Keung','用心活出每一天的精彩。',NULL,679,'2005-10-20 00:13:53',NULL,NULL);
INSERT INTO `users` VALUES (495,'Ng Ka Ming','ngkaming@hotmail.com','eYokHYRe9N','BqQDfBaEN1','admin','Ng Ka Ming','勇敢去追，错过了不后悔。',NULL,441,'2014-03-20 18:56:49',NULL,NULL);
INSERT INTO `users` VALUES (496,'Kyle Howard','kylehoward5@gmail.com','enIg6jDc7g','MiWYRVp7ws','student','Kyle Howard','生活从不缺少奇迹，只缺少发现的眼睛。',NULL,687,'2019-09-08 09:27:38',NULL,NULL);
INSERT INTO `users` VALUES (497,'Meng Yunxi','mengyunxi@mail.com','NGsBaGecwX','9s7crof0sS','student','Meng Yunxi','梦想不会辜负每一个坚持的人。',NULL,909,'2018-11-21 16:44:54',NULL,NULL);
INSERT INTO `users` VALUES (498,'Xiang Yuning','xiang109@hotmail.com','vx56SiOq3c','rYEg5NJhCX','student','Xiang Yuning','每一个今天，都是明天的回忆。',NULL,826,'2001-04-11 01:39:10',NULL,NULL);
INSERT INTO `users` VALUES (499,'Tanaka Hina','hina7@outlook.com','ckbSUjefhR','aWRu5js7Gs','student','Tanaka Hina','不负韶华，追求梦想。',NULL,282,'2000-10-21 01:28:59',NULL,NULL);
INSERT INTO `users` VALUES (500,'Yin Zhiyuan','yinzhiyuan@yahoo.com','GtiDfqWmC0','HaWToWiR2t','admin','Yin Zhiyuan','每一份努力，都会换来不一样的未来。',NULL,518,'2012-07-12 13:55:27',NULL,NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-18 21:43:58
