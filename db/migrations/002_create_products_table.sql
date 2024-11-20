CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    product_name VARCHAR(100),
    product_description TEXT,
    product_images TEXT[],
    compressed_product_images TEXT[],
    product_price DECIMAL
);

