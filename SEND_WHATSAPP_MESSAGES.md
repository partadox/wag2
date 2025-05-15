# Sending WhatsApp Messages Guide

This guide explains how to properly send WhatsApp messages using the Go Simple WhatsApp Gateway.

## Phone Number Format

When sending WhatsApp messages, it's important to format the recipient's phone number correctly. The following formats are accepted:

1. **Standard International Format**: 
   `628123456789` or `+628123456789`
   - Include the country code (e.g., 62 for Indonesia, 1 for US/Canada)
   - Do not include spaces, dashes, or special characters
   - The plus sign (+) is optional and will be automatically removed

2. **WhatsApp JID Format**:
   `628123456789@s.whatsapp.net`
   - This is the internal WhatsApp format
   - The application automatically adds the "@s.whatsapp.net" suffix if not provided

## Common Errors and Solutions

### "Invalid recipient: not a user JID"

This error occurs when the phone number is not properly formatted for WhatsApp. To fix:

1. Ensure you're using the full international format including country code
2. Remove spaces, dashes and special characters
3. Try adding the "@s.whatsapp.net" suffix manually

### "Not connected" or "Not logged in"

This means the WhatsApp client isn't properly connected. To fix:

1. Go to the client's page and check its status
2. If disconnected, click on the QR Code button to authenticate
3. Scan the QR code with your WhatsApp mobile app
4. Try sending the message again after connecting

### "Failed to send message"

This generic error can occur for various reasons:

1. Check if the recipient's phone number exists on WhatsApp
2. Verify your internet connection
3. Make sure your WhatsApp session is still valid
4. Try reconnecting using the QR code

## Testing Your Connection

Before sending real messages, you can test your connection by:

1. Sending a message to your own phone number
2. Using the "Status" API to check if your client is properly connected
3. Looking at the client's status on the dashboard or client details page

## API Endpoints

If you're using the API directly:

```
POST /api/clients/{client_id}/send
{
  "recipient": "628123456789",
  "message": "Your message here"
}
```

Include the API key in the `X-API-Key` header or as the `api_key` query parameter.
