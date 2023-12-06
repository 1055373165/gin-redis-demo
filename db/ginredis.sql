CREATE DATABASE IF NOT EXISTS `ginredis`
USE `ginredis`;

--
-- Table structure for table `comment`
--
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `comment_id` bigint unsigned NOT NULL,
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `post_id` bigint NOT NULL,
    `author_id` bigint NOT NULL,
    `parent_id` bigint NOT NULL,
    `status` tinyint unsigned NOT NULL DEFAULT '1',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP，
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_comment_id` (`comment_id`),
    KEY `idx_author_id` (`authour_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `comment`
--
LOCK TABLES `comment` WRITE;
UNLOCK TABLES;

--
-- Table structure for table `community`
--
DROP TABLE IF EXISTS `community`;
CREATE TABLE `community` (
    `id` int NOT NULL AUTO_INCREMENT,
    `community_id` int unsigned NOT NULL,
    `community_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `introduction` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_community_id` (`community_id`),
    UNIQUE KEY `idx_community_id` (`community_name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARACTER=utf8mb4 COLLATE=utf8mb4_general_ci; 

--
-- Dumping data for table `community`
--
LOCK TABLES `community` WRITE;
INSERT INTO `community` VALUES (1, 1, 'Go Topic', 'Golang', '2006-11-01 00:10:10', '2023-05-05 00:29:54'),(2,2,'Leetcode Topic','？？？？','2020-01-01 00:00:00','2023-05-05 00:28:54'),(3,3,'C++ Topic','C ++ ','2023-05-05 00:31:41','2023-05-05 00:31:41');
UNLOCK TABLES;

--
-- Table structure for table `post`
--
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `post_id` bigint NOT NULL COMMENT '帖子 id',
    `title` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
    `content` varchar(8192) CHARACTER  SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
    `author_id` bigint NOT NULL COMMENT '作者的用户 id',
    `community_id` bigint NOT NULL COMMENT '作者所属社区',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '帖子状态',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id` (`post_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARACTER=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `post`
--
LOCK TABLES `post` WRITE;
INSERT INTO `post` VALUES (1,63542771746603009,'holy shit ','holy shit !',63542625180844033,2,1,'2023-04-23 08:40:43','2023-04-23 08:40:43'),(2,63543686792740865,'123','321',63542625180844033,1,1,'2023-04-23 08:49:48','2023-04-23 08:49:48'),(3,63547823483781121,'1','2',63542625180844033,5,1,'2023-04-23 09:30:54','2023-05-04 07:44:57'),(4,63940350946836481,'gogogo','shit shit shit',63542625180844033,2,1,'2023-04-26 02:30:19','2023-04-26 02:30:19'),(5,65132075690229761,'5.4','5.4',63542625180844033,1,1,'2023-05-04 07:49:02','2023-05-04 07:49:02'),(6,65132733440983041,'66','66',63542625180844033,1,1,'2023-05-04 07:55:34','2023-05-04 07:55:34'),(7,65132756559986689,'66','66',63542625180844033,1,1,'2023-05-04 07:55:48','2023-05-04 07:55:48'),(8,65132782313013249,'66','666',63542696282685441,1,1,'2023-05-04 07:56:03','2023-05-04 07:56:03'),(9,65132859806973953,'44','44',63542625180844033,1,1,'2023-05-04 07:56:49','2023-05-04 07:56:49'),(10,65132899669639169,'11','11',63542625180844033,1,1,'2023-05-04 07:57:13','2023-05-04 07:57:13'),(11,65153146984333313,'test-test','test-test',65094992003072001,1,1,'2023-05-04 11:18:21','2023-05-04 11:18:21'),(12,65231969750876161,'test2023-test2023','test2023-test2023',65094992003072001,2,1,'2023-05-05 00:21:23','2023-05-05 00:21:23'),(13,65232007466057729,'test2023-test2023','test2023-test2023',65094992003072001,2,1,'2023-05-05 00:21:46','2023-05-05 00:21:46'),(14,65232924257026049,'2023-5-5','2023-5-5',63542625180844033,2,1,'2023-05-05 00:30:52','2023-05-05 00:30:52');
UNLOCK TABLES;

--
-- Table structure for table `user`
--
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `user_id` bigint NOT NULL,
    `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `gender` tinyint NOT NULL DEFAULT '0',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE,
    UNIQUE KEY `idx_email` (`email`) USING BTREE
)ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `user`
--
LOCK TABLES `user` WRITE;
INSERT INTO `user` VALUES (1,63542625180844033,'feng','31323334f5d77a10ae47e3738837865e6a831793','feng@126.com',0,'2023-04-23 08:39:16','2023-06-05 01:01:25'),(2,63542696282685441,'dev','31323334f5d77a10ae47e3738837865e6a831793','dev@126.com',0,'2023-04-23 08:39:58','2023-06-05 01:01:25'),(3,63948201509519361,'liang','31323334f5d77a10ae47e3738837865e6a831793','liang@126.com',0,'2023-04-26 03:48:18','2023-06-05 01:01:25'),(4,64078576667852801,'holy','31323334f5d77a10ae47e3738837865e6a831793','holy',1,'2023-04-27 01:23:28','2023-04-27 01:23:28'),(5,65094992003072001,'feng121','3132333435368cab425cf81f49fef958042bfb16029b','123@123.com',1,'2023-05-04 01:40:38','2023-05-04 01:40:38'),(6,65095113201680385,'feng1211','3132333435368cab425cf81f49fef958042bfb16029b','123@123.com',1,'2023-05-04 01:41:50','2023-05-04 01:41:50'),(7,65097239747362817,'feng12112','3132333435368cab425cf81f49fef958042bfb16029b','555@11.com',1,'2023-05-04 02:02:58','2023-06-06 11:48:26');
UNLOCK TABLES;
-- Dump completed on 2023-06-30 17:55:45
