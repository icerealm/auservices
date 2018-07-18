CREATE USER docker WITH SUPERUSER PASSWORD 'docker';

CREATE DATABASE au OWNER docker ENCODING 'UTF-8';

CREATE TABLE category_message (
 ID bigserial NOT NULL PRIMARY KEY,
 message_seq bigint NOT NULL,
 info json NOT NULL,
 status varchar(8) NOT NULL,
 create_dt timestamp default current_date
);

CREATE TABLE item_message (
 ID bigserial NOT NULL PRIMARY KEY,
 message_seq bigint NOT NULL,
 info json NOT NULL,
 status varchar(8) NOT NULL,
 create_dt timestamp default current_date
);

CREATE TABLE category (
    ID bigserial NOT NULL PRIMARY KEY,
    category_nm varchar(128) NOT NULL,
    category_desc varchar(512),
    category_type varchar(16) NOT NULL,
    rev_dt timestamp default current_date,
    rev_by varchar(64) NOT NULL,
    UNIQUE(category_nm, rev_by)
);

CREATE TABLE item_line (
    ID bigserial NOT NULL PRIMARY KEY,
    item_line_nm varchar(128) NOT NULL,
    item_line_desc varchar(512),
    item_line_dt timestamp default current_date,
    item_value numeric(15,6) NOT NULL default 0.0,
    category_id bigint references category(ID) NOT NULL,
    rev_dt timestamp default current_date
);
