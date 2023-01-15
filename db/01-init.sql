CREATE SEQUENCE IF NOT EXISTS account_id;

CREATE TABLE
    "accounts" (
        "id" int4 NOT NULL DEFAULT nextval('account_id':: regclass),
        "balance" float8 NOT NULL DEFAULT 0,
        PRIMARY KEY ("id")
    );

CREATE TABLE
    IF NOT EXISTS "cloud_pockets" (
        id SERIAL PRIMARY KEY,
        name TEXT,
        category TEXT,
        currency TEXT,
        balance float8 NOT NULL DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS "transaction" (
    id serial NOT NULL,
    source_cloud_pocket_id int8 NULL,
    destination_cloud_pocket_id int8 NULL,
    amount float8 NULL,
    datetime timestamp NULL,
    description text NULL,
    status text NULL,
    CONSTRAINT expenses_pkey PRIMARY KEY (id)
);