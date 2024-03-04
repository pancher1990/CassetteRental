-- +goose Up
-- +goose StatementBegin
create table public.customer (
    id          serial      primary key,
    created_at  timestamptz not null default current_timestamp,
    name        text        not null,
    is_active   boolean     not null default true,
    balance     integer     not null default 0,
    password    text        not null,
    email       text        not null unique
);

comment on column public.customer.id is 'Уникальный идентификатор';
comment on column public.customer.name is 'Имя';
comment on column public.customer.is_active is 'Признак активности';
comment on column public.customer.balance is 'Балланс на счету';
comment on column public.customer.password is 'Хэшированный пароль';
comment on column public.customer.email is 'Почта';

alter table public.customer owner to postgres;
grant all on table public.customer to postgres;

create table public.film (
    id          serial      primary key,
    created_at  timestamptz not null default current_timestamp,
    price       integer     not null,
    title       text        not null unique
);

comment on column public.film.id is 'Уникальный идентификатор';
comment on column public.film.price is 'Стоимость суточной аренды';
comment on column public.film.title is 'Название';

alter table public.film owner to postgres;
grant all on table public.film to postgres;

create table public.cassette (
    id              serial   primary key,
    film_id         int      not null references public.film(id),
    is_available    boolean  not null default false
);

comment on column public.cassette.id is 'Уникальный идентификатор кассеты';
comment on column public.cassette.film_id is 'Уникальный идентификатор фильма';
comment on column public.cassette.is_available is 'Признак доступности кассеты к заказу';

alter table public.cassette owner to postgres;
grant all on table public.cassette to postgres;

create table public.order (
    id              serial      primary key,
    customer_id     int         not null references public.customer(id),
    created_at      timestamptz not null default current_timestamp,
    return_deadline timestamptz not null,
    is_active       boolean     not null default true
);

comment on column public.order.id is 'Уникальный идентификатор';
comment on column public.order.customer_id is 'Уникальный идентификатор пользователя';
comment on column public.order.created_at is 'Время создания заказа';
comment on column public.order.return_deadline is 'Крайний срок возврата';
comment on column public.order.is_active is 'Признак активности';

alter table public.order owner to postgres;
grant all on table public.order to postgres;

create table public.order_cassette (
    id          serial  primary key,
    order_id    int     not null references public.order(id),
    cassette_id int     not null references public.cassette(id)
);

comment on column public.order_cassette.id is 'Уникальный идентификатор';
comment on column public.order_cassette.cassette_id is 'Уникальный идентификатор кассеты';
comment on column public.order_cassette.order_id is 'Уникальный идентификатор заказа';

alter table public.order_cassette owner to postgres;
grant all on table public.order_cassette to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table public.order_cassette;
drop table public.order;
drop table public.cassette;
drop table public.film;
drop table public.customer;
-- +goose StatementEnd
