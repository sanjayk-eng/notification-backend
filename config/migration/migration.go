package migration

func CreateTable() string {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		phone VARCHAR(15) UNIQUE NOT NULL,
		country_code VARCHAR(5) NOT NULL DEFAULT '91',
		fullName VARCHAR(100) NOT NULL DEFAULT 'username',
		profilePhoto TEXT DEFAULT 'https://i.pravatar.cc/150?img=3',
		status VARCHAR(255) DEFAULT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	return query
}
