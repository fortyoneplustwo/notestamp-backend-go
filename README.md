# notestamp-backend-go

Experimental backend API for [notestamp.com](https://notestamp.com) implemented for educational purposes. Written in Go.

## Stack

- **Go** with the standard library `net/http` (Go 1.22+ method/pattern-based routing, no framework)
- **Firestore** for user accounts, revoked-token records, and committed note metadata
- **Redis** for staging in-progress saves before they're committed
- **Google Cloud Storage** for note content and media files
- **JWT** (access + refresh tokens, HMAC-signed) for auth, delivered via `HttpOnly` cookies
- **Argon2id** for password hashing

## Architecture

The app is split into small, independently testable packages, each exposing a store type behind an interface:

```
auth/        password hashing, JWT issuing/verification, revoked-token store
handler/     HTTP handlers, each depending on narrow interfaces (defined in handler/interfaces)
media/       GCS-backed store for media file objects
metadata/    Firestore-backed store for note metadata + a Redis-backed staging area
middleware/  authentication middleware (verifies/refreshes JWTs, injects user into context)
notes/       GCS-backed store for note content objects
user/        Firestore-backed user store + request-context helpers
cleanup/     best-effort rollback logic for partially-failed saves
```

Handlers depend on small interfaces rather than concrete store types (see `handler/interfaces`), so each handler can be unit tested against mocks instead of real Firestore/Redis/GCS clients.

### Save flow: stage, upload, commit

Saving a note involves multiple failure-prone steps (writing metadata, uploading media, uploading note content) potentially against different backends. Instead of doing all of this in one request, the API splits it into three steps:

1. **Stage** (`POST /save-with-media` or `/save-without-media`) — validates and writes the note's metadata to Redis with a TTL, and returns pre-built upload paths. The client uploads the actual note/media content directly using those paths.
2. **Update** (`PUT /update-notes`) — the client can update note content directly at any point.
3. **Commit** (`POST /commit`) — verifies that the media/notes actually landed in storage, then promotes the staged metadata into permanent Firestore storage and clears the Redis entry.

If verification at commit time fails, `cleanup.FailedSave` is fired to asynchronously remove whatever was already written (staged metadata, media, notes) so nothing is left orphaned, and logs any cleanup failures for later reconciliation rather than failing loudly.

The upload/verification checks in the commit step run concurrently (fan-out over goroutines, fan-in over a channel) since they hit independent backends.

### Auth

- Access + refresh JWT pair issued on login/registration, stored as `HttpOnly`, `Secure`, `SameSite=Strict` cookies.
- The `Authenticate` middleware verifies the access token on every protected request; if it's expired (but otherwise valid) and its refresh token isn't revoked, it transparently issues a new pair.
- Logout/deregistration record the token in a revoked-token collection so a stolen token can't be replayed after logout, even if it hasn't expired yet.

## API

| Method | Path | Description |
|---|---|---|
| POST | `/register` | Create a user |
| POST | `/login` | Authenticate, receive token cookies |
| DELETE | `/logout` | Revoke current token, clear cookies |
| DELETE | `/deregister` | Delete user and all associated data |
| GET | `/api/list` | List a user's saved note metadata |
| POST | `/api/save-with-media` | Stage metadata for a note with an attached media file, get upload paths |
| POST | `/api/save-without-media` | Stage metadata for a note with only an external source (e.g. YouTube link) |
| PUT | `/api/update-notes` | Overwrite note content |
| POST | `/api/commit` | Verify uploads and promote staged metadata to permanent storage |
| DELETE | `/api/remove` | Delete a note and its associated media/content |

Routes under `/api/*` require a valid session and go through the auth middleware.

## Testing

Each store and handler has unit tests against mock implementations of its dependencies (see the `mock.go` files throughout), so the suite runs without needing live Firestore/Redis/GCS connections.

```
go test ./...
```

## Running locally

Requires a `.env` file in the project root directory with:

```
PORT=
REDIS_ADDR=
FIREBASE_CONF=       # path to a Firebase service account JSON
USER_COLLECTION=
REVOKED_COLLECTION=
METADATA_COLLECTION=
NOTES_BUCKET=
MEDIA_BUCKET=
JWT_SIGNING_KEY=
JWT_ACC_EXP_IN=       # hours
JWT_REF_EXP_IN=       # hours
```

```
go run ./cmd
```
