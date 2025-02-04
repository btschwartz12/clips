clips:
	CGO_ENABLED=0 go build -o clips main.go

run-server: clips
	godotenv -f .env ./clips \
		--port 8000 \
		--var-dir var \
		--config-file config.yaml

clean:
	rm -f clips