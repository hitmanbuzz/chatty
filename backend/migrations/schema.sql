CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,  
    is_online BOOLEAN NOT NULL,
    last_online TIMESTAMPTZ NOT NULL,
    in_group BOOLEAN NOT NULL
);

CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(50) NOT NULL,
    users_id INTEGER[],
    msgs_id INTEGER[],
    create_time TIMESTAMPTZ NOT NULL,
    total_users INTEGER,
    owner_id INTEGER,
    CONSTRAINT group_owner_id
        FOREIGN KEY (owner_id)
        REFERENCES users(user_id)
);

CREATE TABLE messages (
    msg_id SERIAL PRIMARY KEY,
    reply_msg_id INTEGER,
    data TEXT,
    create_time TIMESTAMPTZ NOT NULL,
    user_id INTEGER,
    CONSTRAINT msg_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(user_id)
);
