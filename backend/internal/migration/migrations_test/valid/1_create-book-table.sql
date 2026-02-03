CREATE TABLE book (
    id INTEGER PRIMARY KEY NOT NULL UNIQUE,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    published_year INTEGER,
    isbn TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL
);
