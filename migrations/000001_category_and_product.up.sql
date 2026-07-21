CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products
(
    id          UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    price       BIGINT       NOT NULL CHECK (price > 0),
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE categories
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    title      VARCHAR(255) NOT NULL UNIQUE,
    parent_id UUID REFERENCES categories (id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE images
(
    id         UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    key        TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE product_categories
(
    product_id  UUID      NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    category_id UUID      NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (product_id, category_id)
);

CREATE TABLE product_images
(
    product_id UUID      NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    image_id UUID      NOT NULL REFERENCES images (id) ON DELETE CASCADE,
    sort_order INTEGER   NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (product_id, image_id),

    ADD CONSTRAINT unq_img_order UNIQUE(product_id, sort_order)
);

CREATE TABLE nutrition_facts
(
    product_id    UUID PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
    calories      INTEGER       NOT NULL CHECK (calories >= 0),
    protein       DECIMAL(6, 2) NOT NULL CHECK (protein >= 0),
    fat           DECIMAL(6, 2) NOT NULL CHECK (fat >= 0),
    carbohydrates DECIMAL(6, 2) NOT NULL CHECK (carbohydrates >= 0),
    weight        DECIMAL(6, 2) NOT NULL CHECK (weight > 0),
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);