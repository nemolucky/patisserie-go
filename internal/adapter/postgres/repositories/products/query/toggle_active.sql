UPDATE products
SET is_active = NOT is_active
WHERE id = $1
RETURNING is_active;