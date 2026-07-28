SELECT
    p.id,
    p.title,
    p.price,
    pi.url,
    pi.sort_order,
    COALESCE(
        array_agg(c.title)
            FILTER (WHERE c.id IS NOT NULL),
        '{}'
    ) AS categories
FROM products p
LEFT JOIN product_images pi
    ON p.id = pi.product_id AND pi.sort_order = 1
LEFT JOIN product_categories pc
    ON p.id = pc.product_id
LEFT JOIN categories c
    ON c.id = pc.category_id
WHERE p.is_active = TRUE
AND ($1 = '' OR p.title ILIKE '%' || $1 || '%')
GROUP BY
    p.id,
    pi.url,
    pi.sort_order
ORDER BY p.title
LIMIT $2
OFFSET $3;