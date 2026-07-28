INSERT INTO product_categories (product_id, category_id)
SELECT $1, unnest($2::uuid[]);