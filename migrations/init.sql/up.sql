CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    balance INT DEFAULT 1000,
    created_at DATE DEFAULT NOW()
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    reciever_id INT NOT NULL,
    amount INT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    item_name VARCHAR(50) NOT NULL,
    count INT
);