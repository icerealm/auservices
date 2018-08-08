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
SELECT c.id, c.category_nm, coalesce(v.total,0) total, c.category_type, c.user_id, to_char(now(), 'MM/YYYY') month_year
FROM 
	(select category_id, sum(item_value) total 
	 from item_line where item_line_dt between (date_trunc('month', now())::date) and 
												 (date_trunc('month', now()::date) + interval '1 month' - interval '1 day')::date
												group by category_id) as v 
	RIGHT JOIN category c ON v.category_id = c.id;

-- SQL - showing items of the month
select * from item_line where item_line_dt between (date_trunc('month', now())::date) and now();

-- SQL - showing items of the week
select * from item_line where item_line_dt between date_trunc('week', now())::date and (date_trunc('week', now()) + '6 days') ::date;



