-- MariaDB dump 10.19  Distrib 10.11.6-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: db    Database: mariadb
-- ------------------------------------------------------
-- Server version	10.11.9-MariaDB-ubu2204

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (id, username, email, password_hash, created_at, updated_at) VALUES
(UUID(), 'foobar', 'jonesrussell42@gmail.com', '$2a$10$7U0oMJZ0qtKcrJPI0otrXOTczXRfHdYD64JZ6oB2QTluNMSF9zmE6', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

SET @user_id = (SELECT id FROM users WHERE username = 'foobar');

LOCK TABLES `campaigns` WRITE;
/*!40000 ALTER TABLE `campaigns` DISABLE KEYS */;
INSERT INTO `campaigns` (id, name, description, template, owner_id, created_at, updated_at) VALUES
(UUID(), 'Unmarked Burials', 'Urge for increased funding...', '... template content ...', @user_id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(UUID(), 'Climate Action Now', 'Advocate for stronger climate policies...', '... template content ...', @user_id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
/*!40000 ALTER TABLE `campaigns` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-10-12 20:34:17
