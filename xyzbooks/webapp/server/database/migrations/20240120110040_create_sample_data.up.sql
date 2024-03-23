-- Generate UUIDs for the new author and book
DO $$
DECLARE
    author_hartse UUID;
    author_templer UUID;
    author_duras UUID;
    author_amis UUID;
    author_flagg UUID;
    author_paglia UUID;
    author_rilke UUID;

    publisher_paste UUID;
    publisher_weekly UUID;
    publisher_graywolf UUID;
    publisher_mcsweeny UUID;

    book_elf UUID;
    book_cosmo UUID;
    book_essex UUID;
    book_mister UUID;
    book_underwater UUID;
BEGIN 
    -- Create authors
    author_hartse := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, last_name) VALUES (author_hartse, 'John', 'Hartse');
    author_templer := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, middle_name, last_name) VALUES (author_templer, 'Hannah', 'P', 'Templer');
    author_duras := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, middle_name, last_name) VALUES (author_duras, 'Margueriete', 'Z', 'Duras');
    author_amis := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, last_name) VALUES (author_amis, 'Kingsley', 'Amis');
    author_flagg := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, middle_name, last_name) VALUES (author_flagg, 'Fanny', 'Peters', 'Flagg');
    author_paglia := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, middle_name, last_name) VALUES (author_paglia, 'Camille', 'Byron', 'Paglia');
    author_rilke := uuid_generate_v4(); 
    INSERT INTO author (id, first_name, middle_name, last_name) VALUES (author_rilke, 'Rainer', 'Steel', 'Rilke');

    -- Create publishers
    publisher_paste := uuid_generate_v4();
    INSERT INTO publisher (id, name) VALUES (publisher_paste, 'Paste Magazine');
    publisher_weekly:= uuid_generate_v4();
    INSERT INTO publisher (id, name) VALUES (publisher_weekly, 'Publishers Weekly');
    publisher_graywolf:= uuid_generate_v4();
    INSERT INTO publisher (id, name) VALUES (publisher_graywolf, 'Graywolf Press');
    publisher_mcsweeny:= uuid_generate_v4();
    INSERT INTO publisher (id, name) VALUES (publisher_mcsweeny, 'McSweeney''s');

    -- create books
    book_elf := uuid_generate_v4();
    INSERT INTO book (id, title, isbn10, isbn13, publication_year, publisher_id, edition, price) VALUES (book_elf, 'American Elf', '1891830856', '9781891830853', 2004, publisher_paste, 'Book 2', 1000.00);
    book_cosmo := uuid_generate_v4();
    INSERT INTO book (id, title, isbn10, isbn13, publication_year, publisher_id, edition, price) VALUES (book_cosmo, 'Cosmoknights', '1603094547', '9781603094542', 2019, publisher_weekly, 'Book 1', 2000.00);
    book_essex := uuid_generate_v4();
    INSERT INTO book (id, title, isbn10, isbn13, publication_year, publisher_id, price) VALUES (book_essex, 'Essex County', '160309038X', '9781603090384', 1990, publisher_graywolf, 500.00);
    book_mister := uuid_generate_v4();
    INSERT INTO book (id, title, isbn10, isbn13, publication_year, publisher_id, edition, price) VALUES (book_mister, 'Hey, Mister (Vol 1)', '1891830023', '9781891830020', 2000, publisher_graywolf, 'After School
Special', 1200.00);
    book_underwater :=  uuid_generate_v4();
    INSERT INTO book (id, title, isbn10, isbn13, publication_year, publisher_id, price) VALUES (book_underwater, 'The Underwater Welder', '1603093982', '9781603093989', 2022, publisher_mcsweeny, 3000.00);

    -- create authors_books
    INSERT INTO authors_books (author_id, book_id) VALUES (author_hartse, book_elf);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_templer, book_elf);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_duras, book_elf);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_amis, book_cosmo);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_amis, book_essex);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_templer, book_mister);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_flagg, book_mister);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_paglia, book_mister);
    INSERT INTO authors_books (author_id, book_id) VALUES (author_rilke, book_underwater);
END $$;