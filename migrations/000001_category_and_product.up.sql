CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products
(
    id          UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    price       BIGINT       NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE categories
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    title      VARCHAR(255) NOT NULL UNIQUE,
    parent_id UUID REFERENCES categories (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE images
(
    id         UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    key        TEXT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

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
    sort_order INTEGER   NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE nutrition_facts
(
    product_id    UUID PRIMARY KEY REFERENCES products (id) ON DELETE CASCADE,
    calories      INTEGER       NOT NULL,
    protein       DECIMAL(6, 2) NOT NULL,
    fat           DECIMAL(6, 2) NOT NULL,
    carbohydrates DECIMAL(6, 2) NOT NULL,
    weight        DECIMAL(6, 2) NOT NULL,
    created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_categories_category_id ON product_categories (category_id);

CREATE INDEX idx_product_images_product_id ON product_images (product_id);