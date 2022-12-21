-- CREATE TABLE "smart_contracts" -------------------------------
CREATE TABLE IF NOT EXISTS `smart_contracts` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `marketplace_id` INT,
    `address` VARCHAR(120),
    `node_url` VARCHAR(250),
    `event` JSON,
    `resource_address` VARCHAR(120),
    `resource_node_url` VARCHAR(250),
    `resource_offset_stack` INT,
    PRIMARY KEY (`id`),
    CONSTRAINT `contract` UNIQUE(`address`)
    ) ENGINE = InnoDB;
-- ---------------------------------------------------------------

-- CREATE TABLE "marketplaces" -----------------------------------
CREATE TABLE IF NOT EXISTS `marketplaces` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `name` VARCHAR(120),
    PRIMARY KEY (`id`),
    CONSTRAINT `name` UNIQUE(`name`)
    ) ENGINE = InnoDB;
-- ----------------------------------------------------------------

-- CREATE TABLE "collections" ------------------------------------
CREATE TABLE IF NOT EXISTS `collections` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `marketplace_id` INT,
    `name` VARCHAR(120),
    `creator` VARCHAR(120),
    PRIMARY KEY (`id`),
    INDEX `marketplace` (`marketplace_id`),
    INDEX `name` (`name`)
    ) ENGINE = InnoDB;
-- ----------------------------------------------------------------

-- CREATE TABLE "collection_stats" ------------------------------------
CREATE TABLE IF NOT EXISTS `collection_stats` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `collection_id` INT,
    `total_volume` FLOAT,
    `total_buy` INT,
    `total_listing` INT,
    `floor_price` FLOAT,
    `created_at` DATETIME NULL DEFAULT NULL,
    `updated_at` DATETIME NULL DEFAULT NULL,
    `deleted_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `collection` (`collection_id`, `created_at`)
    ) ENGINE = InnoDB;
-- ----------------------------------------------------------------

-- CREATE TABLE "nfts" ------------------------------------
CREATE TABLE IF NOT EXISTS `nfts` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `marketplace_id` INT,
    `collection_id` INT,
    `owner_wallet_address` VARCHAR(120),
    `name` VARCHAR(120),
    `metadata` JSON,
    `listing_price` FLOAT,
    `listing_start_date` DATETIME NULL DEFAULT NULL,
    `listing_end_date` DATETIME NULL DEFAULT NULL,
    `status` VARCHAR(20),
    PRIMARY KEY (`id`),
    INDEX `nft` (`marketplace_id`, `collection_id`, `name`),
    INDEX `owner` (`owner_wallet_address`)
    ) ENGINE = InnoDB;
-- ----------------------------------------------------------------

-- CREATE TABLE "wallets" ------------------------------------
CREATE TABLE IF NOT EXISTS `wallets` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `address` VARCHAR(120),
    PRIMARY KEY (`id`),
    INDEX `address` (`address`)
    ) ENGINE = InnoDB;
-- ----------------------------------------------------------------

-- CREATE TABLE "transactions" ------------------------------------
CREATE TABLE IF NOT EXISTS `transactions` (
    `id` INT(0) UNSIGNED AUTO_INCREMENT NOT NULL,
    `tx_hash` VARCHAR(120),
    `tx_creation_number` VARCHAR(50),
    `tx_sequence_number` VARCHAR(50),
    `marketplace_id` INT,
    `collection_id` INT,
    `wallet_address` VARCHAR(120),
    `event` VARCHAR(50),
    `nft_id` INT,
    `from_address` VARCHAR(120),
    `to_address` VARCHAR(120),
    `price` FLOAT,
    `data` JSON,
    `created_at` DATETIME NULL DEFAULT NULL,
    `updated_at` DATETIME NULL DEFAULT NULL,
    `deleted_at` DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `unique_transaction` (`tx_hash`,`tx_creation_number`,`tx_sequence_number`),
    INDEX `marketplace` (`marketplace_id`,`collection_id`),
    INDEX `nft` (`nft_id`),
    INDEX `wallet` (`wallet_address`)
    ) ENGINE = InnoDB;
-- ----------------------------------------------------------------

insert into marketplaces (id, name)
values
    (1, 'BlueMove'),
    (2, 'Topaz');

insert into smart_contracts (id, marketplace_id, address, node_url, event, resource_address, resource_node_url, resource_offset_stack)
values (
           1,
           1,
           '0xd1fd99c1944b84d1670a2536417e997864ad12303d19eac725891691b04d614e',
           'https://fullnode.mainnet.aptoslabs.com',
           '{"0xd1fd99c1944b84d1670a2536417e997864ad12303d19eac725891691b04d614e::marketplaceV2::ListEvent":"listing","0xd1fd99c1944b84d1670a2536417e997864ad12303d19eac725891691b04d614e::marketplaceV2::BuyEvent":"buy","0xd1fd99c1944b84d1670a2536417e997864ad12303d19eac725891691b04d614e::marketplaceV2::DelistEvent":"delist","0xd1fd99c1944b84d1670a2536417e997864ad12303d19eac725891691b04d614e::marketplaceV2::ChangePriceEvent":"change_price","0xd1fd99c1944b84d1670a2536417e997864ad12303d19eac725891691b04d614e::offer_lib::AcceptOfferCollectionEvent":"accept_offer"}',
           '0x2a1f62a1663fc7e6c08753e8fc925fbcb946c4b80c5c95a95314a16bc3ac24bc',
           'https://indexer.mainnet.aptoslabs.com/v1/graphql',
           0
       );