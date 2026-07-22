SELECT
		i.id,
		i.key
FROM images i
LEFT JOIN product_images pi
		ON pi.image_id = i.id
WHERE pi.image_id IS NULL;