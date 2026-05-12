.PHONY: build docker test clean

# Build Go binary
build:
	CGO_ENABLED=0 go build -o ascii-art-web .

# Build Docker image
docker: build
	sudo docker build -t ascii-art-app .

# Test Docker
test: docker
	sudo docker run ascii-art-app "Make Test" standard

# Run web server
run-web:
	sudo docker run -p 8080:8080 asciithe-art-app

# Clean up
clean:
	rm -f ascii-art-web
	sudo docker rm -f $$(sudo docker ps -aq) 2>/dev/null || true

# Quick development
dev:
	go run .