-- Create "todos" table
CREATE TABLE `todos` (
  `id` binary(16) NOT NULL,
  `title` longtext NOT NULL,
  `description` longtext NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;