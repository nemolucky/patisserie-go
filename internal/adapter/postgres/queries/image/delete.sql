DELETE FROM images 
WHERE id = ANY($1);