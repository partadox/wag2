# Deployment Guide

This document explains how to deploy the Go Simple WhatsApp Gateway on various platforms.

## Deploy on VPS

### Using Docker (Recommended)

1. Clone the repository on your VPS:
   ```bash
   git clone https://github.com/yourusername/go-simple-whatsapp-gateway2.git
   cd go-simple-whatsapp-gateway2
   ```

2. Update the API key in the docker-compose.yml file (replace "changeme" with your desired API key)

3. Build and start the container:
   ```bash
   docker-compose up -d
   ```

4. The application will be available at http://your-server-ip:8080

### Using Docker Manually

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-simple-whatsapp-gateway2.git
   cd go-simple-whatsapp-gateway2
   ```

2. Build the Docker image:
   ```bash
   docker build -t whatsapp-gateway .
   ```

3. Run the container:
   ```bash
   docker run -d \
     --name whatsapp-gateway \
     -p 8080:8080 \
     -e API_KEY=your_api_key \
     -v whatsapp-data:/app/whatsapp-data \
     whatsapp-gateway
   ```

4. The application will be available at http://your-server-ip:8080

### Without Docker

1. Install Go 1.20 or higher on your VPS

2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-simple-whatsapp-gateway2.git
   cd go-simple-whatsapp-gateway2
   ```

3. Build the application:
   ```bash
   go build
   ```

4. Create an .env file with your configuration:
   ```
   LISTEN_ADDR=:8080
   API_KEY=your_api_key
   WHATSAPP_DATA_DIR=./whatsapp-data
   ```

5. Run the application:
   ```bash
   ./go-simple-whatsapp-gateway2
   ```

6. For production use, consider using a process manager like systemd or supervisor to keep the application running.

## Deploy on Coolify

1. Set up a Coolify instance following their documentation

2. Create a new service and select "Git Repository"

3. Connect your fork of this repository

4. Configure the build settings:
   - **Build Pack**: Docker
   - **Environment Variables**:
     - API_KEY: your_api_key
     - LISTEN_ADDR: :8080
     - WHATSAPP_DATA_DIR: /app/whatsapp-data
     - GIN_MODE: release

5. Deploy the application

6. Coolify will use the Dockerfile in the .coolify directory, build the Docker image, and run it automatically.

## Troubleshooting

### Template Not Found Error

If you see an error like `html/template: pattern matches no files: templates/*`, it means the application can't find the template files. This may happen if:

1. The templates folder is not in the same directory as the executable
2. Permissions issues prevent the application from reading the templates directory
3. The path is hardcoded to a Windows-style path (D:/ etc.) which doesn't exist on Linux

Solution:
- Make sure the templates directory is in the correct location
- Check file permissions
- Use relative paths like `./templates/*` instead of absolute paths

### Database Errors

If you encounter database errors, check:
1. The whatsapp-data directory exists and is writable
2. The WHATSAPP_DATA_DIR environment variable is set correctly
3. If using Docker, ensure you're using a volume to persist data

### Connection Issues

If clients fail to connect to WhatsApp, check:
1. Your server has internet access
2. Firewall rules are not blocking outgoing connections
3. The whatsmeow library is up to date (check UPDATE_WHATSMEOW.md)