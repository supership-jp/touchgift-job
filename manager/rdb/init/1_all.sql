-- mysqldump: [Warning] Using a password on the command line interface can be insecure.
-- MySQL dump 10.13  Distrib 8.0.34, for Linux (x86_64)
--
-- Host: localhost    Database: retail
-- ------------------------------------------------------
-- Server version	8.0.34

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
-- mysqldump: Error: 'Access denied; you need (at least one of) the PROCESS privilege(s) for this operation' when trying to dump tablespaces

--
-- Table structure for table `banner`
--

DROP TABLE IF EXISTS `banner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `banner` (
  `id` int NOT NULL AUTO_INCREMENT,
  `img_url` varchar(255) NOT NULL COMMENT 'クリエイティブ画像URL',
  `xid` varchar(32) NOT NULL COMMENT 'img_url内に含むユニーク文字列',
  `height` int NOT NULL,
  `width` int NOT NULL,
  `extension` enum('image/png','image/jpeg') NOT NULL COMMENT '拡張子のタイプ。pngの場合image/png、jpegの場合image/jpeg',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_beaf7c371ddbfd3aec513671e7` (`xid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `coupon`
--

DROP TABLE IF EXISTS `coupon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `coupon` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `name` varchar(255) NOT NULL COMMENT 'クーポン名',
  `status` enum('0','1') NOT NULL DEFAULT '0' COMMENT '審査ステータス。0: 審査前(下書き), 1: OK',
  `code` varchar(255) DEFAULT NULL COMMENT 'クーポンコード',
  `img_url` varchar(255) NOT NULL COMMENT 'クーポン画像URL',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `xid` varchar(32) NOT NULL COMMENT 'img_url内に含むユニーク文字列',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_04007043d9cb7e062723c13181` (`xid`),
  UNIQUE KEY `IDX_72e3a38e2d1c160d3dba89b8ab` (`organization_code`,`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `creative`
--

DROP TABLE IF EXISTS `creative`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `creative` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `name` varchar(255) NOT NULL COMMENT 'クリエイティブ名',
  `status` enum('0','1','2','3','4') NOT NULL DEFAULT '0' COMMENT '審査ステータス。0: 審査前(下書き), 1: 審査中, 2: 審査OK(公開中), 3: NG, 4: 停止中',
  `click_url` varchar(255) DEFAULT NULL COMMENT '遷移先URL',
  `creative_type` enum('banner','video') NOT NULL COMMENT 'クリエイティブ種別',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `creative_banner`
--

DROP TABLE IF EXISTS `creative_banner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `creative_banner` (
  `id` int NOT NULL AUTO_INCREMENT,
  `creative_id` int NOT NULL COMMENT 'クリエイティブID',
  `banner_id` int NOT NULL COMMENT 'バナーID',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `REL_71154171edaeca0de953b088f4` (`creative_id`),
  UNIQUE KEY `REL_075f52f6fa73f1d964d76e4f84` (`banner_id`),
  CONSTRAINT `FK_075f52f6fa73f1d964d76e4f847` FOREIGN KEY (`banner_id`) REFERENCES `banner` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_71154171edaeca0de953b088f4f` FOREIGN KEY (`creative_id`) REFERENCES `creative` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `creative_video`
--

DROP TABLE IF EXISTS `creative_video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `creative_video` (
  `id` int NOT NULL AUTO_INCREMENT,
  `creative_id` int NOT NULL COMMENT 'クリエイティブID',
  `video_id` int NOT NULL COMMENT 'ビデオID',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `REL_bb754ed848cc7ee4b17a341ce1` (`creative_id`),
  UNIQUE KEY `REL_4f31b5c357a04ff557e987d945` (`video_id`),
  CONSTRAINT `FK_4f31b5c357a04ff557e987d9453` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_bb754ed848cc7ee4b17a341ce11` FOREIGN KEY (`creative_id`) REFERENCES `creative` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gimmick`
--

DROP TABLE IF EXISTS `gimmick`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gimmick` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `name` varchar(255) NOT NULL COMMENT 'ギミック名',
  `img_url` varchar(255) NOT NULL COMMENT 'ギミック画像URL',
  `status` enum('0','1','2','3','4') NOT NULL DEFAULT '0' COMMENT '審査ステータス。0: 審査前(下書き), 1: 審査中, 2: 審査OK(公開中), 3: NG, 4: 停止中',
  `xid` varchar(32) NOT NULL COMMENT 'S3へのアップロード用識別子',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_4da2b8bfc4b59dc0ec2468b3b0` (`xid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `location`
--

DROP TABLE IF EXISTS `location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `location` (
  `id` int NOT NULL AUTO_INCREMENT,
  `origin_id` varchar(255) NOT NULL,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `name` varchar(255) DEFAULT NULL COMMENT '管理名',
  `width` int DEFAULT NULL COMMENT '枠の横幅',
  `height` int DEFAULT NULL COMMENT '枠の縦幅',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_66440ed6cb1acd0240ba6cc947` (`origin_id`,`organization_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `migrations`
--

DROP TABLE IF EXISTS `migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `migrations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `timestamp` bigint NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `store`
--

DROP TABLE IF EXISTS `store`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `store` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `store_id` varchar(255) NOT NULL COMMENT '店舗ID',
  `name` varchar(255) NOT NULL COMMENT '店舗名',
  `post_code` varchar(8) DEFAULT NULL COMMENT '郵便番号',
  `province_code` varchar(2) DEFAULT NULL COMMENT '都道府県コード',
  `address` varchar(255) DEFAULT NULL COMMENT '住所（漢字）',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_f5393e62c0378af2a73664e520` (`store_id`,`organization_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `store_update_history`
--

DROP TABLE IF EXISTS `store_update_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `store_update_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `file_name` varchar(255) NOT NULL COMMENT 'ファイル名',
  `xid` varchar(32) NOT NULL COMMENT 'S3アップロード用識別子',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_81b6186ac397a87cede61f5648` (`xid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `touch_point`
--

DROP TABLE IF EXISTS `touch_point`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `touch_point` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `point_unique_id` varchar(30) NOT NULL COMMENT 'NFC: NFCID, QR: 印字管理ID',
  `print_management_id` varchar(30) NOT NULL COMMENT '印字管理ID',
  `store_id` varchar(255) NOT NULL COMMENT '店舗ID',
  `type` varchar(255) NOT NULL COMMENT '接触機器の種別',
  `name` varchar(255) DEFAULT NULL COMMENT 'NFC名',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_b8f52033d364a02db8ac40b40b` (`point_unique_id`),
  UNIQUE KEY `IDX_521605b4828c5311c1921631fd` (`print_management_id`),
  KEY `FK_48e60d26ce1179f9021d5d916ae` (`store_id`),
  CONSTRAINT `FK_48e60d26ce1179f9021d5d916ae` FOREIGN KEY (`store_id`) REFERENCES `store` (`store_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `touch_point_update_history`
--

DROP TABLE IF EXISTS `touch_point_update_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `touch_point_update_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_code` varchar(255) NOT NULL COMMENT '組織コード',
  `file_name` varchar(255) NOT NULL COMMENT 'ファイル名',
  `xid` varchar(32) NOT NULL COMMENT 'S3アップロード用識別子',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_41aefe2e9fd6daa7f3eb7a2675` (`xid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `video`
--

DROP TABLE IF EXISTS `video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `video` (
  `id` int NOT NULL AUTO_INCREMENT,
  `video_url` varchar(255) NOT NULL COMMENT 'クリエイティブ動画URL',
  `endcard_url` varchar(255) NOT NULL COMMENT 'エンドカード画像URL',
  `video_xid` varchar(32) NOT NULL COMMENT 'video_url内に含むユニーク文字列',
  `endcard_xid` varchar(32) NOT NULL COMMENT 'endcard_url内に含むユニーク文字列',
  `height` int NOT NULL,
  `width` int NOT NULL,
  `extension` enum('video/quicktime','video/mp4') NOT NULL,
  `endcard_height` int NOT NULL,
  `endcard_width` int NOT NULL,
  `endcard_extension` enum('image/png','image/jpeg') NOT NULL COMMENT 'エンドカードの拡張子。pngの場合image/png、jpegの場合image/jpeg',
  `created_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT 'レコードが作成された日時',
  `updated_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT 'レコードが更新された日時',
  `last_updated_by` int NOT NULL COMMENT '最終更新者のユーザID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_e66ade2997af86b2270b537132` (`video_xid`),
  UNIQUE KEY `IDX_747a5c13abca59b4e7e1b130da` (`endcard_xid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-06-09 12:05:06
