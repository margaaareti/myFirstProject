CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    surname       varchar(255) not null,
    patronymic    varchar(255) not null,
    username      varchar(255) not null unique,
    password      varchar(255) not null,
    email         varchar(255) not null unique,
    verified      boolean default false,
    reg_date      timestamp with time zone default current_timestamp
);


CREATE TABLE students
(
    id serial not null unique,
    name    varchar(255) not null,
    surname varchar(255) not null,
    patronymic varchar(255),
    isu_number varchar(255) not null unique ,
    added_by int references users (id) on delete cascade not null,
    /*added_by varchar(255) not null,*/
    title varchar(255) not null,
    description varchar(255),
    reg_date      timestamp with time zone default current_timestamp
);



/*CREATE TABLE notes
(
    id serial not null unique,
    add_by int references users (id) on delete cascade not null,
    owner int references students (id) on delete cascade not null,
    title varchar(255) not null,
    description varchar(255)
);*/





