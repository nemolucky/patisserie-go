UPDATE nutrition_facts
SET (calories, protein, fat, carbohydrates, weight) = ($2, $3, $4, $5, $6)
WHERE product_id = $1;