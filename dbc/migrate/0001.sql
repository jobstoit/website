-- 0001.sql

CREATE TABLE admins (
	id SERIAL PRIMARY KEY,
	username VARCHAR NOT NULL UNIQUE,
	password VARCHAR NOT NULL
);

CREATE TABLE sites (
	id SERIAL PRIMARY KEY,
	created_by INT NOT NULL,
	name VARCHAR(50) NOT NULL UNIQUE,
	active_since TIMESTAMP DEFAULT NOW(),
	css_uri VARCHAR,
	title VARCHAR(50)
);

ALTER TABLE sites
ADD CONSTRAINT fk_sites_created_by
FOREIGN KEY (created_by) REFERENCES admins(id);

CREATE TABLE pages (
	id SERIAL PRIMARY KEY,
	site_id INT NOT NULL,
	uri_path VARCHAR NOT NULL,
	label VARCHAR(50) NOT NULL
);

ALTER TABLE pages
ADD CONSTRAINT fk_pages_site_id
FOREIGN KEY (site_id) REFERENCES sites(id);

CREATE TABLE rows (
	id SERIAL PRIMARY KEY,
	page_id INT NOT NULL,
	sequence INT NOT NULL
);

ALTER TABLE rows
ADD CONSTRAINT fk_rows_page_id
FOREIGN KEY (page_id) REFERENCES pages(id);

CREATE TABLE row_titles (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	context VARCHAR NOT NULL
);

ALTER TABLE row_titles
ADD CONSTRAINT fk_row_titles_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

CREATE TABLE row_texts (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	context TEXT NOT NULL
);

ALTER TABLE row_texts
ADD CONSTRAINT fk_row_texts_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

CREATE TYPE media_type AS ENUM ('video', 'audio', 'image');

CREATE TABLE row_media (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	uri VARCHAR NOT NULL,
	type media_type NOT NULL,
	alt VARCHAR(50)
);

ALTER TABLE row_media
ADD CONSTRAINT fk_row_media_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

CREATE TABLE row_buttons (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	label VARCHAR(50) NOT NULL,
	uri VARCHAR NOT NULL
);

ALTER TABLE row_buttons
ADD CONSTRAINT fk_row_buttons_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

CREATE TYPE nav_position AS ENUM ('header', 'footer', 'panel');

CREATE TABLE navigation_links (
	id SERIAL PRIMARY KEY,
	site_id INT NOT NULL,
	page_id INT NOT NULL,
	position nav_position NOT NULL,
	sequence INT NOT NULL
);

ALTER TABLE navigation_links
ADD CONSTRAINT fk_navigation_links_page_id
FOREIGN KEY (page_id) REFERENCES pages(id);

ALTER TABLE navigation_links
ADD CONSTRAINT fk_navigation_links_site_id
FOREIGN KEY (site_id) REFERENCES sites(id);
