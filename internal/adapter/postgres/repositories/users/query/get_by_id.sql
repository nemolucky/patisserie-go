SELECT id, 
			 first_name, 
			 middle_name, 
			 last_name, 
			 phone, 
			 telegram 
FROM users 
WHERE id = $1;