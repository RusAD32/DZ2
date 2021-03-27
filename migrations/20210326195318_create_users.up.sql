CREATE TABLE users (
    login varchar not null ,
    password varchar not null
);

CREATE TABLE cars (
    mark varchar not null primary key,
    max_speed integer not null,
    distance integer not null,
    handler varchar not null,
    stock varchar not null
);