app:
	CGO_ENABLED=0 go build -o app main.go

run-server: app
	godotenv -f .env ./app \
		--port 8000 \
		--var-dir var \
		--media-dir /Volumes/SAUL

clean:
	rm -f app