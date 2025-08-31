ALTER TABLE `credentials`
    ADD COLUMN `active` BOOLEAN DEFAULT FALSE;

UPDATE `credentials`
    SET `active` = TRUE;
