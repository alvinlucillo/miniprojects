DROP INDEX book_isbn13_idx;

ALTER TABLE book DROP CONSTRAINT IF EXISTS book_isbn13_unique;
ALTER TABLE book DROP CONSTRAINT IF EXISTS book_isbn10_unique;