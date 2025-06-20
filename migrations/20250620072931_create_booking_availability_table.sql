-- Create "availabilities" table
CREATE TABLE `availabilities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `uuid` binary(16) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `weekday` varchar(10) NOT NULL,
  `start_time` varchar(5) NOT NULL,
  `end_time` varchar(5) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_availabilities_deleted_at` (`deleted_at`),
  INDEX `idx_availabilities_user_id` (`user_id`),
  INDEX `idx_availabilities_uuid` (`uuid`),
  UNIQUE INDEX `uni_availabilities_uuid` (`uuid`),
  CONSTRAINT `fk_availabilities_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "bookings" table
CREATE TABLE `bookings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `uuid` binary(16) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `guest_name` varchar(255) NOT NULL,
  `guest_email` varchar(255) NOT NULL,
  `date` datetime(3) NOT NULL,
  `start_time` varchar(5) NOT NULL,
  `end_time` varchar(5) NOT NULL,
  `status` varchar(20) NULL DEFAULT "confirmed",
  PRIMARY KEY (`id`),
  INDEX `idx_bookings_deleted_at` (`deleted_at`),
  INDEX `idx_bookings_user_id` (`user_id`),
  INDEX `idx_bookings_uuid` (`uuid`),
  UNIQUE INDEX `uni_bookings_uuid` (`uuid`),
  CONSTRAINT `fk_bookings_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
