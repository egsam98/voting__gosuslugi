create table users (
    id bigserial primary key,
    passport varchar not null unique,
    fullname varchar not null,
    birth_date timestamp not null,
    death_date timestamp
);
