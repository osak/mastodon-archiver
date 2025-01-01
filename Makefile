download:
	go build -o bin/download cmd/download/main.go

migration_up:
	migrate -path db/migrations -database "mysql://mastodon_archiver:test@tcp(127.0.0.1)/db" up
