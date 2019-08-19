# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.25)
# Database: go_test
# Generation Time: 2019-08-19 07:51:26 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table a
# ------------------------------------------------------------

DROP TABLE IF EXISTS `a`;

CREATE TABLE `a` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(11) DEFAULT NULL,
  `create` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `a` WRITE;
/*!40000 ALTER TABLE `a` DISABLE KEYS */;

INSERT INTO `a` (`id`, `name`, `create`)
VALUES
	(1,'a',NULL),
	(2,'b',NULL),
	(3,'c',NULL);

/*!40000 ALTER TABLE `a` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table class
# ------------------------------------------------------------

DROP TABLE IF EXISTS `class`;

CREATE TABLE `class` (
  `id` int(11) unsigned NOT NULL,
  `name` varchar(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `class` WRITE;
/*!40000 ALTER TABLE `class` DISABLE KEYS */;

INSERT INTO `class` (`id`, `name`)
VALUES
	(1,'tset'),
	(2,'张三');

/*!40000 ALTER TABLE `class` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table student
# ------------------------------------------------------------

DROP TABLE IF EXISTS `student`;

CREATE TABLE `student` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `no` varchar(20) DEFAULT NULL,
  `class_id` int(11) DEFAULT NULL,
  `create_date` timestamp NULL DEFAULT NULL,
  `is_reading` int(11) DEFAULT NULL,
  `High` double DEFAULT NULL,
  `Weight` decimal(10,0) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `student` WRITE;
/*!40000 ALTER TABLE `student` DISABLE KEYS */;

INSERT INTO `student` (`id`, `name`, `address`, `no`, `class_id`, `create_date`, `is_reading`, `High`, `Weight`)
VALUES
	(83,'张三','中国','00000000',1,'2019-08-15 11:28:11',1,NULL,NULL);

/*!40000 ALTER TABLE `student` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table teacher
# ------------------------------------------------------------

DROP TABLE IF EXISTS `teacher`;

CREATE TABLE `teacher` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `no` varchar(20) DEFAULT NULL,
  `class_id` int(11) DEFAULT NULL,
  `create_date` timestamp NULL DEFAULT NULL,
  `is_reading` int(11) DEFAULT NULL,
  `High` double DEFAULT NULL,
  `Weight` decimal(10,0) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `teacher` WRITE;
/*!40000 ALTER TABLE `teacher` DISABLE KEYS */;

INSERT INTO `teacher` (`id`, `name`, `address`, `no`, `class_id`, `create_date`, `is_reading`, `High`, `Weight`)
VALUES
	(84,'zzq',NULL,NULL,1,'2019-08-19 14:39:50',1,160.256,NULL),
	(85,'zzq',NULL,NULL,1,NULL,1,160.256,NULL),
	(86,'zzq',NULL,NULL,NULL,'2019-08-19 14:41:32',1,NULL,NULL),
	(87,'zzq',NULL,NULL,NULL,'2019-08-19 15:50:13',1,160.256,NULL);

/*!40000 ALTER TABLE `teacher` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
