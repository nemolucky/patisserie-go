CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

	first_name VARCHAR(100) NOT NULL,
	middle_name VARCHAR(100),
	last_name VARCHAR(100) NOT NULL,

	phone VARCHAR(20) UNIQUE NOT NULL
		CHECK (phone ~ '^\+[1-9]\d{7,14}$'),
	telegram VARCHAR(32) UNIQUE
		CHECK (telegram ~ ',^@[A-Za-z0-9_]{5,32}$'),
		
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);