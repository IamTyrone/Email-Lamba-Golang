# Email Lambda Service

A Go-based AWS Lambda function that sends emails using Amazon SES (Simple Email Service). This function is designed to be triggered via AWS API Gateway.

## Features

- Sends emails via Amazon SES.
- Supports both HTML and Text content.
- Secures the endpoint using a static API Key.
- lightweight Go runtime.

## Prerequisites

- Go 1.23 or later
- AWS Account with SES setup
- Verified Sender Identity (Email or Domain) in AWS SES

## Configuration

The function requires the following environment variable to be set in the Lambda configuration:

| Variable | Description |
|----------|-------------|
| `STATIC_API_KEY` | A secret key used to authenticate requests. |

## Building and Deployment

To build the function for AWS Lambda (Linux/amd64):

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap
```

Upload `function.zip` to your AWS Lambda function. Ensure the handler is set to `bootstrap` (if using `provided.al2` or `provided.al2023` runtime) or match your specific runtime settings.

## API Documentation

### Endpoint

The Lambda function is expected to be invoked via an API Gateway Proxy integration.

### Authentication

All requests must include the following header:

- **Header:** `x-api-key`
- **Value:** Your configured `STATIC_API_KEY`

### Request Format

**Method:** `POST`

**Content-Type:** `application/json`

**Body:**

```json
{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "Email Subject",
  "body": "<h1>Hello</h1><p>This is the email body.</p>",
  "isHtml": true
}
```

| Field | Type | Description |
|-------|------|-------------|
| `from` | string | The sender's email address (must be verified in SES). |
| `to` | string | The recipient's email address. |
| `subject` | string | The subject line of the email. |
| `body` | string | The content of the email. |
| `isHtml` | boolean | Set to `true` for HTML content, `false` for plain text. |

### Responses

#### Success (200 OK)

```json
{
  "message": "Email sent successfully"
}
```

#### Unauthorized (401 Unauthorized)

Returned if the `x-api-key` is missing or invalid.

```json
{
  "message": "Unauthorized"
}
```

#### Bad Request (400 Bad Request)

Returned if the JSON body is invalid or malformed.

```json
{
  "message": "Invalid JSON body"
}
```

#### Internal Server Error (500 Internal Server Error)

Returned if AWS SES fails to send the email.

```json
{
  "message": "Failed to send email: <error_details>"
}
```
