# Webhooks Package

The Webhooks package handles incoming webhook events from external providers (currently Instagram/Meta). It provides endpoints for verification and event processing.

## Endpoints

### 1. Instagram Webhook Verification
- **URL**: `/api/v1/webhooks/instagram`
- **Method**: `GET`
- **Purpose**: Handles Meta's webhook verification challenge.
- **Parameters**:
  - `hub.mode`: Must be `subscribe`.
  - `hub.verify_token`: Must match `INSTAGRAM_WEBHOOK_VERIFY_TOKEN` in env.
  - `hub.challenge`: Random string sent by Meta.
- **Response**: Returns the `hub.challenge` string if verification succeeds.

### 2. Instagram Webhook Event Receiver
- **URL**: `/api/v1/webhooks/instagram`
- **Method**: `POST`
- **Purpose**: Receives and processes webhook events.
- **Security**: Verifies `X-Hub-Signature-256` header using `INSTAGRAM_APP_SECRET`.
- **Process**:
  1.  Verifies the request signature.
  2.  Parses the JSON payload.
  3.  Iterates through entries and changes.
  4.  Enqueues an asynchronous task (via Asynq) for each event.
  5.  Returns `200 OK` ("EVENT_RECEIVED") immediately to acknowledge receipt.

## Supported Events

| Event Type | Task Type | Description |
| :--- | :--- | :--- |
| `comments` | `instagram:comment` | New comments on media. |
| `mentions` | `instagram:mention` | User mentions in comments or captions. |
| `story_insights` | `instagram:story_insight` | Insights data for stories (impressions, reach, etc.). |

## Configuration

Required environment variables:
- `INSTAGRAM_WEBHOOK_VERIFY_TOKEN`: Token for verification challenge.
- `INSTAGRAM_APP_SECRET`: App Secret for signature verification.
