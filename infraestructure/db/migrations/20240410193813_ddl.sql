-- +goose Up
-- +goose StatementBegin
create table merchant (
    id uuid primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    password varchar(255) not null,
    create_at timestamp not null default (now()),
    updated_at timestamp not null default (now())
);

comment on table merchant is 'Storing merchant data';

create unique index merchant_email_uindex on merchant (email);

create table tx
(
    id uuid primary key,
    merchant_id uuid not null,
    amount float default 0 not null,
    tx_type bigserial not null,
    create_at timestamp not null default (now()),
    updated_at timestamp not null default (now())
);

comment on table tx is 'Storing transactions data';

alter table tx add foreign key (merchant_id) references merchant (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
