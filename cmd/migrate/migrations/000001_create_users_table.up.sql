create table if not exists users (
    id serial primary key,
    name text not null,
    email text not null unique,
    password text not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);