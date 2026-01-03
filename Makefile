include .env

create_migrations: 
	sqlx migrate add -r init