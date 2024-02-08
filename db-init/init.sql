-- public.customer определение

-- Drop table

-- DROP TABLE public.customer;

CREATE TABLE public.customer (
                                 id varchar NOT NULL, -- Уникальный идентификатор клиента
                                 "name" varchar NOT NULL, -- Имя пользователя
                                 is_active bool NULL,
                                 balance int4 NOT NULL, -- Балланс на счету
                                 "password" varchar NOT NULL, -- Хэшированный пароль
                                 email varchar NOT NULL, -- Почта пользователя
                                 CONSTRAINT customer_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.customer.id IS 'Уникальный идентификатор клиента';
COMMENT ON COLUMN public.customer."name" IS 'Имя пользователя';
COMMENT ON COLUMN public.customer.balance IS 'Балланс на счету';
COMMENT ON COLUMN public.customer."password" IS 'Хэшированный пароль';
COMMENT ON COLUMN public.customer.email IS 'Почта пользователя';

-- Permissions

ALTER TABLE public.customer OWNER TO postgres;
GRANT ALL ON TABLE public.customer TO postgres;


-- public.film определение

-- Drop table

-- DROP TABLE public.film;

CREATE TABLE public.film (
                             id varchar NOT NULL, -- Уникальный идентификатор фильма
                             day_price int4 NULL,
                             title varchar NULL, -- Название фильма
                             CONSTRAINT film_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.film.id IS 'Уникальный идентификатор фильма';
COMMENT ON COLUMN public.film.title IS 'Название фильма';

-- Permissions

ALTER TABLE public.film OWNER TO postgres;
GRANT ALL ON TABLE public.film TO postgres;


-- public.cassette определение

-- Drop table

-- DROP TABLE public.cassette;

CREATE TABLE public.cassette (
                                 id varchar NOT NULL, -- Уникальный идентификатор кассеты
                                 film_id varchar NOT NULL, -- Уникальный идентификатор фильма
                                 available bool NOT NULL DEFAULT false, -- Признак доступности кассеты к заказу
                                 CONSTRAINT cassette_pk PRIMARY KEY (id),
                                 CONSTRAINT cassette_fk FOREIGN KEY (film_id) REFERENCES public.film(id)
);

-- Column comments

COMMENT ON COLUMN public.cassette.id IS 'Уникальный идентификатор кассеты';
COMMENT ON COLUMN public.cassette.film_id IS 'Уникальный идентификатор фильма';
COMMENT ON COLUMN public.cassette.available IS 'Признак доступности кассеты к заказу';

-- Permissions

ALTER TABLE public.cassette OWNER TO postgres;
GRANT ALL ON TABLE public.cassette TO postgres;


-- public."order" определение

-- Drop table

-- DROP TABLE public."order";

CREATE TABLE public."order" (
                                id varchar NOT NULL, -- Уникальный идентификатор заказа
                                customer_id varchar NOT NULL, -- Уникальный идентификатор пользователя
                                "date" timetz NULL, -- Время создания заказа
                                CONSTRAINT order_pk PRIMARY KEY (id),
                                CONSTRAINT order_fk FOREIGN KEY (customer_id) REFERENCES public.customer(id)
);

-- Column comments

COMMENT ON COLUMN public."order".id IS 'Уникальный идентификатор заказа';
COMMENT ON COLUMN public."order".customer_id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN public."order"."date" IS 'Время создания заказа';

-- Permissions

ALTER TABLE public."order" OWNER TO postgres;
GRANT ALL ON TABLE public."order" TO postgres;


-- public.rent определение

-- Drop table

-- DROP TABLE public.rent;

CREATE TABLE public.rent (
                             id varchar NOT NULL, -- Уникальный идентификатор аренды
                             customer_id varchar NOT NULL, -- Уникальный идентификатор клиента
                             cassette_id varchar NOT NULL, -- Уникальный идентификатор кассеты
                             start_datetime timestamp NOT NULL, -- Время начала аренды
                             end_datetime timestamp NOT NULL, -- Время конца аренды
                             return_sign bool NOT NULL DEFAULT false, -- Признак возврата
                             CONSTRAINT rent_pk PRIMARY KEY (id),
                             CONSTRAINT rent_fk FOREIGN KEY (cassette_id) REFERENCES public.cassette(id),
                             CONSTRAINT rent_fk_1 FOREIGN KEY (customer_id) REFERENCES public.customer(id)
);

-- Column comments

COMMENT ON COLUMN public.rent.id IS 'Уникальный идентификатор аренды';
COMMENT ON COLUMN public.rent.customer_id IS 'Уникальный идентификатор клиента';
COMMENT ON COLUMN public.rent.cassette_id IS 'Уникальный идентификатор кассеты';
COMMENT ON COLUMN public.rent.start_datetime IS 'Время начала аренды';
COMMENT ON COLUMN public.rent.end_datetime IS 'Время конца аренды';
COMMENT ON COLUMN public.rent.return_sign IS 'Признак возврата';

-- Permissions

ALTER TABLE public.rent OWNER TO postgres;
GRANT ALL ON TABLE public.rent TO postgres;


-- public.cassette_in_order определение

-- Drop table

-- DROP TABLE public.cassette_in_order;

CREATE TABLE public.cassette_in_order (
                                          cassette_id varchar NOT NULL, -- Уникальный идентификатор кассеты
                                          order_id varchar NOT NULL, -- Уникальный идентификатор заказа
                                          rent_cost int4 NULL,
                                          CONSTRAINT cassette_in_order_pk PRIMARY KEY (cassette_id, order_id)
);

-- Column comments

COMMENT ON COLUMN public.cassette_in_order.cassette_id IS 'Уникальный идентификатор кассеты';
COMMENT ON COLUMN public.cassette_in_order.order_id IS 'Уникальный идентификатор заказа';

-- Permissions

ALTER TABLE public.cassette_in_order OWNER TO postgres;
GRANT ALL ON TABLE public.cassette_in_order TO postgres;

-- public.cassette_in_order внешние включи

ALTER TABLE public.cassette_in_order ADD CONSTRAINT cassette_in_order_fk FOREIGN KEY (cassette_id) REFERENCES public.cassette(id);
ALTER TABLE public.cassette_in_order ADD CONSTRAINT cassette_in_order_fk_1 FOREIGN KEY (order_id) REFERENCES public."order"(id);

-- Permissions

GRANT ALL ON SCHEMA public TO pg_database_owner;
GRANT USAGE ON SCHEMA public TO public;