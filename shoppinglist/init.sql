-- init.sql
CREATE TABLE IF NOT EXISTS shopping_items (
    id SERIAL PRIMARY KEY,
    shopping_item VARCHAR(255) NOT NULL,
    shopping_amount INT NOT NULL
);
