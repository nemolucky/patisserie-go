UPDATE products
SET (title, description, price) = ($2, $3, $4)
WHERE id = $1;