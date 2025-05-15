# Go Simple WhatsApp Gateway - Usage Guide

This guide explains how to run and use the Go Simple WhatsApp Gateway application.

## Running the Application

There are two ways to run the application:

### 1. Standard Mode

Run the application using the standard batch file:

```
D:\Dev\go-simple-whatsapp-gateway2\run.bat
```

### 2. Debug Mode

Run the application with additional debug logging:

```
D:\Dev\go-simple-whatsapp-gateway2\run_debug.bat
```

The debug mode provides more detailed logs which can be helpful for troubleshooting.

## Using the Application

Once the application is running, you can access the web UI and API endpoints:

### Web UI

1. Open your browser and navigate to: `http://localhost:8080`
2. You will be redirected to the dashboard

### Main Pages

- Dashboard: `http://localhost:8080/ui/dashboard`
- Client Management: `http://localhost:8080/ui/clients`
- QR Code Authentication: `http://localhost:8080/ui/qrcode/{client_id}`
- Send Message: `http://localhost:8080/ui/sendmessage/{client_id}`

### API Endpoints

The API requires authentication using the API key specified in your .env file (default: `changeme`).
You can pass the API key in one of two ways:

1. Via the `X-API-Key` header
2. Via the `api_key` query parameter

#### Main API Endpoints:

- List Clients: `GET /api/clients`
- Create Client: `POST /api/clients`
- Get Client Status: `GET /api/clients/{id}`
- Delete Client: `DELETE /api/clients/{id}`
- Generate QR Code: `GET /api/clients/{id}/qr`
- Send Message: `POST /api/clients/{id}/send`
- Logout Client: `POST /api/clients/{id}/logout`

## Creating and Connecting a WhatsApp Client

1. Go to the Clients page (`/ui/clients`)
2. Enter a unique Client ID (e.g., "personal", "work", etc.) and click "Create Client"
3. Once created, click on "QR Code" to generate a QR code for WhatsApp Web authentication
4. Open WhatsApp on your phone:
   - Go to Settings > Linked Devices
   - Tap on "Link a Device"
   - Scan the QR code displayed on your screen
5. After successful scanning, you'll be connected to WhatsApp and can send messages

## Troubleshooting

### Blank White Page

If you encounter a blank white page in the web UI, check:
1. Browser console for JavaScript errors
2. Application logs for any backend errors
3. Ensure templates and static files are accessible

### QR Code Generation Issues

If you have issues generating QR codes:
1. Check the application logs for specific errors
2. Ensure the client is not already logged in (you may need to logout first)
3. Try restarting the application
4. Check your browser console for any frontend JavaScript errors

### Connection Issues

If you're having trouble connecting to WhatsApp:
1. Make sure your phone has an active internet connection
2. Verify that WhatsApp is up to date on your phone
3. Try generating a new QR code
4. Check that your computer's time and date are accurate

## Additional Information

- The default API key is `changeme` (defined in .env file)
- WhatsApp client data is stored in the `whatsapp-data` directory
- Each client has its own SQLite database for session storage
