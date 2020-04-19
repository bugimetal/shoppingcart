-- +goose Up

CREATE TABLE IF NOT EXISTS `shoppingcart` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (`id`),
    INDEX `user_id_cart_id` (`user_id`, `id`)
)
DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci
ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `shoppingcart_item` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `shoppingcart_id` BIGINT NOT NULL,
    `product_id` BIGINT NOT NULL,
    `quantity` UNSIGNED INT NOT NULL DEFAULT 1,
    `price` BIGINT NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (`id`),
    INDEX `shoppingcart_id_item_id` (`shoppingcart_id`, `id`),
    FOREIGN KEY (`shoppingcart_id`) REFERENCES shoppingcart(id)
)
DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci
ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS `shoppingcart_item`;
DROP TABLE IF EXISTS `shoppingcart`;
