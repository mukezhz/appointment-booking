-- Update users table to use full_name instead of first_name/last_name
ALTER TABLE `users`
  -- Create new columns
  ADD COLUMN `full_name` varchar(255) NULL,
  ADD COLUMN `full_name_ja` varchar(255) NULL;

-- Copy existing name data to new columns
UPDATE `users` 
SET `full_name` = CONCAT_WS(' ', `first_name`, `last_name`),
    `full_name_ja` = CONCAT_WS(' ', `first_name_ja`, `last_name_ja`);

-- Drop old columns
ALTER TABLE `users`
  DROP COLUMN `first_name`,
  DROP COLUMN `last_name`,
  DROP COLUMN `first_name_ja`,
  DROP COLUMN `last_name_ja`;
