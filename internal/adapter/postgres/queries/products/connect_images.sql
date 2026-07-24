INSERT INTO product_images (
    product_id,
    image_id,
    sort_order
)
SELECT
    $1,
    image_id,
    sort_order
FROM unnest(
    $2::uuid[],
    $3::int[]
) AS t(image_id, sort_order);