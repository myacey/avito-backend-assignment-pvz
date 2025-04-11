CREATE TYPE role_enum AS ENUM ('employee', 'moderator');

CREATE TABLE IF NOT EXISTS users (
    "id" UUID PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL CHECK("email" LIKE '%@%'),
    "password" varchar NOT NULL,
    "role" role_enum NOT NULL
);

CREATE TYPE city_enum AS ENUM('Москва', 'СПб', 'Казань');

CREATE TABLE IF NOT EXISTS pvz (
    "id" UUID PRIMARY KEY,
    "registration_date" TIMESTAMPTZ NOT NULL,
    "city" city_enum NOT NULL
);

CREATE TYPE status_enum AS ENUM ('in_progress', 'close');

CREATE TABLE IF NOT EXISTS receptions (
    "id" UUID PRIMARY KEY,
    "date_time" TIMESTAMPTZ NOT NULL,
    "pvz_id" UUID REFERENCES pvz ("id") NOT NULL,
    "status" status_enum NOT NULL DEFAULT('in_progress')
);
CREATE INDEX ON receptions ("pvz_id", "status");

CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE IF NOT EXISTS products (
    "id" UUID PRIMARY KEY,
    "date_time" TIMESTAMPTZ NOT NULL DEFAULT(NOW()),
    "type" product_type NOT NULL,
    "reception_id" UUID REFERENCES receptions ("id") NOT NULL
);
