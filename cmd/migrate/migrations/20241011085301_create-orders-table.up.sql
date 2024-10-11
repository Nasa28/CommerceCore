CREATE TABLE IF NOT EXISTS orders
(
    `id` INT UNSIGNED AUTO_INCREMENT,
    `userId` INT UNSIGNED NOT NULL,
    `totalAmount` DECIMAL(10, 2) NOT NULL,
    `status` ENUM('pending', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    `shippingAddress` TEXT NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`userId`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
