CREATE DATABASE goapi;

CREATE TABLE product (
    id serial primary key,
    name varchar(255),
    description varchar(255),
    price float,
    available_flag bool
);
