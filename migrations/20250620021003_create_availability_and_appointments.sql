-- Modify "availabilities" table
ALTER TABLE `availabilities` DROP INDEX `idx_availabilities_uuid`;
-- Modify "availabilities" table
ALTER TABLE `availabilities` RENAME COLUMN `resource_id` TO `doctor_id`, DROP COLUMN `recur_rule`, ADD COLUMN `day_of_week` tinyint NULL, ADD COLUMN `start_date` date NULL, ADD COLUMN `end_date` date NULL, ADD COLUMN `slot_duration` bigint NOT NULL, DROP INDEX `idx_availabilities_end_time`, DROP INDEX `idx_availabilities_resource_id`, DROP INDEX `idx_availabilities_start_time`, ADD UNIQUE INDEX `idx_availabilities_uuid` (`uuid`), DROP INDEX `uni_availabilities_uuid`, ADD INDEX `idx_availabilities_day_of_week` (`day_of_week`), ADD INDEX `idx_availabilities_doctor_id` (`doctor_id`), ADD INDEX `idx_availabilities_start_date` (`start_date`);
-- Create "appointments" table
CREATE TABLE `appointments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `uuid` binary(16) NOT NULL,
  `doctor_id` binary(16) NOT NULL,
  `patient_id` binary(16) NOT NULL,
  `availability_id` bigint unsigned NOT NULL,
  `appointment_date` date NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `status` varchar(20) NULL DEFAULT "scheduled",
  `cancellation_reason` text NULL,
  `lock_key` varchar(255) NULL,
  `lock_expiry` datetime(3) NULL,
  `cooldown_until` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_appointments_appointment_date` (`appointment_date`),
  INDEX `idx_appointments_deleted_at` (`deleted_at`),
  INDEX `idx_appointments_doctor_id` (`doctor_id`),
  INDEX `idx_appointments_patient_id` (`patient_id`),
  INDEX `idx_appointments_status` (`status`),
  UNIQUE INDEX `idx_appointments_uuid` (`uuid`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Drop "bookings" table
DROP TABLE `bookings`;
-- Drop "organizations" table
DROP TABLE `organizations`;
-- Drop "resources" table
DROP TABLE `resources`;
