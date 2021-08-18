-- 0001.sql TEST

INSERT INTO sites (id, name)
VALUES (101, 'testsite'), (102, 'testsite2'), (103, 'testsite3');

INSERT INTO navigation_links (id, site_id, uri, label, position, sequence)
VALUES (104, 101, '/test1', 'test1', 'header', 1), (105, 101, '/test1', 'test1', 'header', 1),
(101, 103, '/test1', 'test1', 'header', 1), (102, 103, '/test2', 'test2', 'header', 1), (103, 103, '/test3', 'test3', 'header', 1);

INSERT INTO pages (id, site_id, uri_path, label)
VALUES (101, 101, '/home', 'Home'), (102, 101, '/createrows', 'Create');

INSERT INTO rows (id, sequence, page_id)
VALUES (101, 1, 101), (102, 2, 101), (103, 3, 101), (104, 1, 102), (105, 1, 102);

INSERT INTO row_titles (row_id, sequence, context)
VALUES (101, 1, 'test1'), (101, 2, 'test2'), (102, 1, 'test3'), (104, 1, 'deletable row');

INSERT INTO row_texts (row_id, sequence, context)
VALUES (101, 1, 'this is just a test text'), 
	(101, 2, 'and another one'),
	(103, 1, 'this is yet another text field');
