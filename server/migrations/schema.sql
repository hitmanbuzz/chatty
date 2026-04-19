CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL UNIQUE,  
    in_group BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(50) NOT NULL UNIQUE,
    total_users INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    owner_id INTEGER,
    CONSTRAINT group_owner_id
        FOREIGN KEY (owner_id)
        REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE group_members (
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
    reply_msg_id INTEGER REFERENCES messages(id) ON DELETE SET NULL,
    content VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id INTEGER,
    CONSTRAINT msg_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);
