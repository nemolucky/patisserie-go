SELECT
    p.id,
    p.title,
    p.description,
    p.price,
    p.is_active,

    COALESCE((
        SELECT jsonb_agg(
            jsonb_build_object(
                'id', c.id,
                'title', c.title,
                'parent_id', c.parent_id
            )
        )
        FROM product_categories pc
        JOIN categories c
            ON pc.category_id = c.id
        WHERE pc.product_id = p.id
    ), '[]'::json) AS categories,

    COALESCE((
        SELECT jsonb_agg(
            jsonb_build_object(
                'id', i.id,
                'key', i.key,
                'sort_order', pi.sort_order
            )
						ORDER BY pi.sort_order
        )
        FROM product_images pi
        JOIN images i
            ON pi.image_id = i.id
        WHERE pi.product_id = p.id
    ), '[]'::json) AS images,

    (
        SELECT jsonb_build_object(
            'calories', nf.calories,
            'protein', nf.protein,
            'fat', nf.fat,
            'carbohydrates', nf.carbohydrates,
            'weight', nf.weight
        )
        FROM nutrition_facts nf
        WHERE nf.product_id = p.id
    ) AS facts

FROM products p
WHERE p.id = $1;