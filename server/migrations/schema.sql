CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,  
    is_online BOOLEAN NOT NULL,
    in_group BOOLEAN NOT NULL
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(50) NOT NULL,
    users_id INTEGER[],
    msgs_id INTEGER[],
    total_users INTEGER,
    owner_id INTEGER,
    CONSTRAINT group_owner_id
        FOREIGN KEY (owner_id)
        REFERENCES users(id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    reply_msg_id INTEGER,
    content VARCHAR(100),
    user_id INTEGER,
    CONSTRAINT msg_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);
