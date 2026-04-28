-- +goose Up
Create table if not exists products (
    id serial primary key,
    name text not null,
    description text not null,
    price integer not null check (price >= 0),
    quantity integer not null default 0 check (quantity >= 0),
    created_at timestamptz not null default now()
);

-- +goose Down
DROP TABLE IF EXISTS products;
