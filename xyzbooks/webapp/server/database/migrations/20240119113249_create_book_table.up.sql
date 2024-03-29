CREATE TABLE book (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(100) NOT NULL,
    isbn10 VARCHAR(10) DEFAULT '',
    isbn13 VARCHAR(13) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    publication_year INTEGER NOT NULL,
    publisher_id UUID NOT NULL CONSTRAINT fk_publisher_id REFERENCES publisher(id),
    edition VARCHAR(100) DEFAULT '',
    image_url VARCHAR(1000) DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);