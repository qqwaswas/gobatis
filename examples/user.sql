+-------+-------------------------------------------------------------------------------------+
| Table | Create Table                                                                        |
+-------+-------------------------------------------------------------------------------------+
| user  | CREATE TABLE `user` (                                                               |
|       |   `id` int(11) NOT NULL AUTO_INCREMENT,                                             |
|       |   `user_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,                      |
|       |   `age` smallint(6) NOT NULL,                                                       |
|       |   `addr` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,                      |
|       |   `passwd` varchar(60) COLLATE utf8mb4_unicode_ci DEFAULT NULL,                     |
|       |   `is_disable` tinyint(1) DEFAULT 0,                                                |
|       |   `money` float DEFAULT NULL,                                                       |
|       |   `total` decimal(10,2) DEFAULT NULL,                                               |
|       |   PRIMARY KEY (`id`)                                                                |
|       | ) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci |
+-------+-------------------------------------------------------------------------------------+

INSERT INTO gobatis.user (id, user_name, age, addr, passwd, is_disable, money, total) VALUES (1, 'sean', 32, 'shenzhen', '12232323', 1, 10.23, 1.33);
INSERT INTO gobatis.user (id, user_name, age, addr, passwd, is_disable, money, total) VALUES (2, 'juanmao', 30, 'hunan', 'sdfsdfsdfsd', 0, 11.56, 3.11);
INSERT INTO gobatis.user (id, user_name, age, addr, passwd, is_disable, money, total) VALUES (3, 'ouyang', 28, 'hunan', '12312312', 1, 18.73, 5.20);
