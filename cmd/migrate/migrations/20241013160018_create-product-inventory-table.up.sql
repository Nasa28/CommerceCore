CREATE TABLE IF NOT EXISTS product_inventory (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `product_id` INT UNSIGNED NOT NULL,
    `quantity_available` INT NOT NULL CHECK (quantity_available >= 0),
    `stock` INT NOT NULL CHECK (stock >= 0),

    PRIMARY KEY (`id`),

    FOREIGN KEY (`product_id`) REFERENCES products(`id`)
);