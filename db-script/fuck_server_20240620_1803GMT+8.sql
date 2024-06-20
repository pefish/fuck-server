#
# SQL Export
# Created by Querious (303000)
# Created: June 20, 2024 at 18:04:00 GMT+8
# Encoding: Unicode (UTF-8)
#


SET @ORIG_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;

SET @ORIG_UNIQUE_CHECKS = @@UNIQUE_CHECKS;
SET UNIQUE_CHECKS = 0;

SET @ORIG_TIME_ZONE = @@TIME_ZONE;
SET TIME_ZONE = '+00:00';

SET @ORIG_SQL_MODE = @@SQL_MODE;
SET SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO';



CREATE TABLE `program` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(225) NOT NULL,
  `content` text NOT NULL,
  `logs` text,
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0 waiting to execute, 1 executing, 2 waiting to exit, 3 exited with error, 4 exited',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3 COMMENT='程序';




LOCK TABLES `program` WRITE;
INSERT INTO `program` (`id`, `name`, `content`, `logs`, `status`, `created_at`, `updated_at`) VALUES 
	(1,'test','function main () {\n    console.log(123)\n\n    console.log("a")\n    console.log("fsb")\n    console.log("gw54w4")\n    console.log("drtgbr")\n}\n\nmain()\n\n','123\na\nfsb\ngw54w4\ndrtgbr\n',4,'2024-06-20 09:07:33','2024-06-20 09:37:17');
UNLOCK TABLES;






SET FOREIGN_KEY_CHECKS = @ORIG_FOREIGN_KEY_CHECKS;

SET UNIQUE_CHECKS = @ORIG_UNIQUE_CHECKS;

SET @ORIG_TIME_ZONE = @@TIME_ZONE;
SET TIME_ZONE = @ORIG_TIME_ZONE;

SET SQL_MODE = @ORIG_SQL_MODE;



# Export Finished: June 20, 2024 at 18:04:00 GMT+8

