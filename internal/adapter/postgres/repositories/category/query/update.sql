UPDATE categories 
SET title = $2
WHERE id = $1
RETURNING id, title, parent_id;
