MIGRATE=migrate
DB_URL=postgres://postgres:nurlan050@localhost:5432/task_management_db?sslmode=disable

up:
	$(MIGRATE) -path migrations -database "$(DB_URL)" up

down:
	$(MIGRATE) -path migrations -database "$(DB_URL)" down

force:
	$(MIGRATE) -path migrations -database "$(DB_URL)" force 1