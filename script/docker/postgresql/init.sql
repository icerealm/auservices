CREATE USER docker WITH SUPERUSER PASSWORD 'docker';

CREATE DATABASE au OWNER docker ENCODING 'UTF-8';

CREATE TABLE category_message (
 id bigserial NOT NULL PRIMARY KEY,
 message_seq bigint NOT NULL,
 info json NOT NULL,
 status varchar(8) NOT NULL,
 create_dt timestamp default CURRENT_TIMESTAMP,
 rev_by varchar(64) NOT NULL
);

CREATE TABLE item_message (
 id bigserial NOT NULL PRIMARY KEY,
 message_seq bigint NOT NULL,
 info json NOT NULL,
 status varchar(8) NOT NULL,
 create_dt timestamp default CURRENT_TIMESTAMP,
 rev_by varchar(64) NOT NULL
);

CREATE TABLE category (
    id bigserial NOT NULL PRIMARY KEY,
    category_nm varchar(128) NOT NULL,
    category_desc varchar(512),
    category_type varchar(16) NOT NULL,
    user_id varchar(64) NOT NULL,
    rev_dt timestamp default CURRENT_TIMESTAMP,
    rev_by varchar(64) NOT NULL,
    UNIQUE(category_nm, user_id)
);

CREATE TABLE item_line (
    id bigserial NOT NULL PRIMARY KEY,
    item_line_nm varchar(128) NOT NULL,
    item_line_desc varchar(512),
    item_line_dt timestamp default current_date,
    item_value numeric(15,6) NOT NULL default 0.0,
    category_id bigint references category(ID) NOT NULL,
    user_id varchar(64) NOT NULL,
    rev_dt timestamp default CURRENT_TIMESTAMP,
    rev_by varchar(64) NOT NULL,
    UNIQUE(category_id, item_line_nm, user_id)
);

-- SQL FOR materialize view - showing total value of each category
SELECT c.id, c.category_nm, coalesce(v.total,0), c.category_type, c.user_id 
FROM 
	(select category_id, sum(item_value) total from item_line group by category_id) as v 
	RIGHT JOIN category c ON v.category_id = c.id;
