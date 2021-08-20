-- 0001.sql

CREATE TYPE row_style
AS ENUM('hero', 'c1', 'c2', 'c3');

CREATE TYPE media_type
AS ENUM ('video', 'audio', 'image');

CREATE TYPE nav_position
AS ENUM ('header', 'footer', 'panel');

CREATE TABLE sites (
	id SERIAL PRIMARY KEY,
	created_by VARCHAR,
	name VARCHAR(50) NOT NULL UNIQUE,
	active_since TIMESTAMP DEFAULT NOW(),
	css_uri VARCHAR,
	title VARCHAR(50)
);

CREATE TABLE pages (
	id SERIAL PRIMARY KEY,
	site_id INT NOT NULL,
	uri_path VARCHAR NOT NULL,
	label VARCHAR(50) NOT NULL
);

CREATE TABLE rows (
	id SERIAL PRIMARY KEY,
	page_id INT NOT NULL,
	sequence INT NOT NULL,
	style row_style NOT NULL
);

CREATE TABLE row_titles (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	context VARCHAR NOT NULL
);

CREATE TABLE row_texts (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	context TEXT NOT NULL
);

CREATE TABLE row_media (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	uri VARCHAR NOT NULL,
	type media_type NOT NULL,
	alt VARCHAR(50)
);

CREATE TABLE row_buttons (
	id SERIAL PRIMARY KEY,
	row_id INT NOT NULL,
	sequence INT NOT NULL,
	label VARCHAR(50) NOT NULL,
	uri VARCHAR NOT NULL
);

CREATE TABLE navigation_links (
	id SERIAL PRIMARY KEY,
	site_id INT NOT NULL,
	uri VARCHAR NOT NULL,
	label VARCHAR(50) NOT NULL,
	position nav_position NOT NULL,
	sequence INT NOT NULL
);

CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	password_hash VARCHAR NOT NULL
);

CREATE TABLE sessions (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	token VARCHAR UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	expires_at TIMESTAMP NOT NULL
);

ALTER TABLE pages
ADD CONSTRAINT fk_pages_site_id
FOREIGN KEY (site_id) REFERENCES sites(id);

ALTER TABLE rows
ADD CONSTRAINT fk_rows_page_id
FOREIGN KEY (page_id) REFERENCES pages(id);

ALTER TABLE row_titles
ADD CONSTRAINT fk_row_titles_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

ALTER TABLE row_texts
ADD CONSTRAINT fk_row_texts_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

ALTER TABLE row_media
ADD CONSTRAINT fk_row_media_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

ALTER TABLE row_buttons
ADD CONSTRAINT fk_row_buttons_row_id
FOREIGN KEY (row_id) REFERENCES rows(id);

ALTER TABLE navigation_links
ADD CONSTRAINT fk_navigation_links_site_id
FOREIGN KEY (site_id) REFERENCES sites(id);

ALTER TABLE sessions
ADD CONSTRAINT fk_sessions_user_id
FOREIGN KEY (user_id) REFERENCES users(id);
