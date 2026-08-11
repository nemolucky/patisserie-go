INSERT INTO users (
	first_name,
	middle_name,
	last_name,
	phone,
	telegram
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5)
RETURNING id;