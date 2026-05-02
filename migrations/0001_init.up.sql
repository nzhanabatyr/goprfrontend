CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       password TEXT NOT NULL
);

CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       author_id INT
);

CREATE TABLE favorite_books (
                                id SERIAL PRIMARY KEY,
                                user_id INT,
                                book_id INT,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);