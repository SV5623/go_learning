create table if not exists events (
    id serial primary key,
    owner_id integer not null /*references users(id) on delete cascade*/,
    title text not null,
    description text,
    location text,
    start_time timestamp with time zone not null,
    end_time timestamp with time zone not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
    foreign key (owner_id) references users(id) on delete cascade
);