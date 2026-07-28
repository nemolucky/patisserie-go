SELECT
    id,
    title,
    parent_id
FROM categories
ORDER BY parent_id NULLS FIRST, title;