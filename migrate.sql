CREATE TABLE `expenses` (
  `id` int(11) NOT NULL PRIMARY key AUTO_INCREMENT,
  `amount` int(11) NOT NULL,
  `category_id` int(11) NOT NULL,
  `comment` text NOT NULL,
  `date` timestamp NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `categories` (
  `id` int(11) NOT NULL PRIMARY key AUTO_INCREMENT,
  `name` text NOT NULL,
  `deleted` tinyint(1) NOT NULL DEFAULT 0,
  `created` timestamp NOT NULL DEFAULT current_timestamp(),
  UNIQUE KEY `uniq_name` (`name`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE USER 'expenses'@'localhost' IDENTIFIED BY 'expenses';
GRANT ALL PRIVILEGES ON expenses.* TO 'expenses'@'localhost';
