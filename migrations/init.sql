CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    balance INT DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    sender_name VARCHAR(50) NOT NULL,
    receiver_id INT NOT NULL,
    receiver_name VARCHAR(50) NOT NULL,
    amount INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_name VARCHAR(50) NOT NULL,
    count INT NOT NULL DEFAULT 1,
    CONSTRAINT unique_user_item UNIQUE (user_id, item_name)
);

-- Индекс для username в users (hash-индекс)
CREATE INDEX IF NOT EXISTS idx_users_username_hash ON users USING HASH (username);

-- Индекс для sender_id и receiver_id в transactions (hash-индексы)
CREATE INDEX IF NOT EXISTS idx_transactions_sender_id_hash ON transactions USING HASH (sender_id);
CREATE INDEX IF NOT EXISTS idx_transactions_receiver_id_hash ON transactions USING HASH (receiver_id);

-- Индекс для user_id в inventory (hash-индекс)
CREATE INDEX IF NOT EXISTS idx_inventory_user_id_hash ON inventory USING HASH (user_id);
