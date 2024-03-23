CREATE INDEX book_isbn13_idx ON book (isbn13);

ALTER TABLE book ADD CONSTRAINT book_isbn13_unique UNIQUE (isbn13);
ALTER TABLE book ADD CONSTRAINT book_isbn10_unique UNIQUE (isbn10);