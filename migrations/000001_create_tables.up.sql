CREATE TYPE role AS ENUM ('employee', 'moderator');

CREATE TABLE IF NOT EXISTS users (
    "id" UUID PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL CHECK("email" LIKE '%@%'),
    "role" role NOT NULL
);

CREATE TYPE city AS ENUM('Москва', 'Санкт-Петербург', 'Казань');

CREATE TABLE IF NOT EXISTS pvz (
    "id" UUID PRIMARY KEY,
    "registration_date" TIMESTAMP NOT NULL,
    "city" city NOT NULL
);

CREATE TYPE status AS ENUM ('in_progress', 'close');

CREATE TABLE IF NOT EXISTS reception (
    "id" UUID PRIMARY KEY,
    "date_time" TIMESTAMP NOT NULL,
    "pvz_id" UUID REFERENCES pvz ("id"),
    "status" status NOT NULL DEFAULT('in_progress')
);

CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TABLE IF NOT EXISTS product (
    "id" UUID PRIMARY KEY,
    "date_time" TIMESTAMP NOT NULL DEFAULT(NOW()),
    "type" product_type NOT NULL,
    "reception_id" UUID REFERENCES Reception ("id")
);

CREATE TABLE IF NOT EXISTS reception_product (
    "reception_id" UUID NOT NULL REFERENCES reception ("id"),
    "product_id" UUID NOT NULL REFERENCES product ("id"),
    PRIMARY KEY ("reception_id", "product_id")
);
