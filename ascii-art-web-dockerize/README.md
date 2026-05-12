# ASCII Art Web - Dockerized 🐳

A Go web application that generates ASCII art from text, now containerized with Docker.

## Features

- 🎨 Generate ASCII art from any text
- 🌐 Web interface with multiple banner styles
- 💻 Command-line interface
- 🐳 Docker support for easy deployment

## Quick Start

### Using Docker (Recommended)


# Clone the repository
git clone https://github.com/yourusername/ascii-art-dockerize
cd ascii-art-dockerize

# Build the Docker image
CGO_ENABLED=0 go build -o ascii-art-web .
sudo docker build -t ascii-art-app .

# Run CLI mode
sudo docker run ascii-art-app "Hello World" standard

# Run web mode
sudo docker run -p 8080:8080 ascii-art-app


Local Development
bash
# Run without Docker
go run .              # Web mode
go run . "Hello" standard  # CLI mode

# Build binary
go build -o ascii-art-web .
./ascii-art-web



  ## Usage


# CLI Mode
docker run ascii-art-app "Your text here" [banner]

# Banner options: standard, shadow, thinkertoy
docker run ascii-art-app "Hello" shadow


# Web Mode
docker run -p 8080:8080 ascii-art-app
# -Open http://localhost:8080


# Docker Commands
Command	Description
sudo docker build -t ascii-art-app .	        -Build image
sudo docker run ascii-art-app "text" standard	   -CLI mode
sudo docker run -p 8080:8080 ascii-art-app	     -Web mode
sudo docker run -d -p 8080:8080 --name myapp ascii-art-app	                  -Run in background
sudo docker ps	                 -List running containers
sudo docker stop myapp	         -Stop container
sudo docker rm myapp	           -Remove container
sudo docker rmi ascii-art-app	   -Remove image


# Project Structure
ascii-art-dockerize/
├── Dockerfile           # Docker configuration
├── Makefile            # Build automation
├── main.go             # Entry point
├── handlers/           # HTTP handlers
├── templates/          # HTML templates
└── README.md           # This file


# Dockerfile
FROM alpine:latest
COPY ascii-art-web /app/
COPY templates /app/templates
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["./ascii-art-web"]


# Makefile Commands
make build    # Build Go binary
make docker   # Build Docker image
make test     # Test Docker image
make run-web  # Run web server
make clean    # Clean artifacts


# Requirements
Go 1.24+
Docker
Linux/Unix environment

# Banner Styles
standard - Default ASCII art style
shadow - Shadow effect style
thinkertoy - ThinkerToy style



## Troubleshooting 


#Docker daemon not running:
sudo systemctl start docker
sudo systemctl enable docker


#Permission denied:
sudo usermod -aG docker $USER
newgrp docker  # or logout/login


#Binary not found in container:
bash
CGO_ENABLED=0 go build -o ascii-art-web .
sudo docker build -t ascii-art-app .


# Authors
- Clinton Odhiambo
- Brown Ajuoga
- Evans Juma