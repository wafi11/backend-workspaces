create table users (
    id serial primary key,
    username varchar(100),
    email varchar(200),
    password_hash  text,
    phone_number VARCHAR(200),
    is_deleted BOOLEAN DEFAULT false,
    pasword varchar(200),
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE UNIQUE INDEX idx_users_username on users(username)  where is_deleted = false;
CREATE UNIQUE INDEX idx_users_email on users(email)  where is_deleted = false;
CREATE UNIQUE INDEX idx_users_phone_number on users(phone_number)  where is_deleted = false;


