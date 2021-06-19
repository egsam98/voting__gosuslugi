create table users (
    id bigserial primary key,
    passport varchar not null unique,
    fullname varchar not null,
    birth_date timestamp not null,
    death_date timestamp
);

-- name: FindByPassport :one
select * from users where passport = $1;

-- name: Create :exec
insert into users (passport, fullname, birth_date, death_date) values ($1, $2, $3, $4);
