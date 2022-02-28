CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    username      varchar(255) not null unique,
    password      varchar(255) not null,
    email         varchar(255) not null unique,
    verified      boolean default false,
    reg_date      timestamp with time zone default current_timestamp
);
