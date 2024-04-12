-- +goose Up
-- +goose StatementBegin
alter table public.customer
    add column is_admin bool not null default false;

create table session (
    id          serial      primary key,
    token       text        not null unique,
    customer_id int         not null references public.customer (id),
    created_at  timestamptz not null default current_timestamp,
    expire_at   timestamptz not null default current_timestamp + interval '1 month'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table session;

alter table public.customer
    drop column is_admin;
-- +goose StatementEnd
