CREATE TABLE IF NOT EXISTS user_roles (
    `user_id` INT UNSIGNED NOT NULL,
    `role_id` INT UNSIGNED NOT NULL,
    
    PRIMARY KEY (`user_id`, `role_id`),
    
    FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`role_id`) REFERENCES roles(`id`) ON DELETE CASCADE
);