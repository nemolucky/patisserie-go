INSERT INTO categories (title , parent_id)
VALUES ($1, $2)
RETURNING id;