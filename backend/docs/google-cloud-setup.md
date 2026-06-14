# Google Cloud Setup (Drive + YouTube)

This guide walks through obtaining Google OAuth credentials for **Google Drive import** and **YouTube publishing** in PostEaze. You can reuse the same Firebase/Google Cloud project — Firebase projects are GCP projects.

## Prerequisites

- A Google account with access to [Google Cloud Console](https://console.cloud.google.com/)
- PostEaze backend and frontend running locally (or production domains ready for redirect URIs)

## 1. Create or select a project

1. Open [Google Cloud Console](https://console.cloud.google.com/)
2. Use the project selector to create a new project or select your existing Firebase project
3. Note the **Project ID** for reference

## 2. Enable APIs

Go to **APIs & Services → Library** and enable:

| API | Used for |
|-----|----------|
| **Google Drive API** | Browse files, download media, list revisions |
| **YouTube Data API v3** | Channel info, resumable video uploads |

## 3. Configure OAuth consent screen

1. Go to **APIs & Services → OAuth consent screen**
2. Choose **External** user type (unless you use Google Workspace internal-only)
3. Fill in app name, support email, and developer contact
4. Add scopes:
   - `https://www.googleapis.com/auth/drive.readonly`
   - `https://www.googleapis.com/auth/youtube.upload`
   - `https://www.googleapis.com/auth/youtube.readonly`
   - `https://www.googleapis.com/auth/userinfo.email`
5. While in **Testing** mode, add team Google accounts under **Test users**

> **Restricted scope note:** `drive.readonly` is a restricted scope. Verification is required only when you move to **Production** and want public access. Testing mode works for listed test users without full verification.

## 4. Create OAuth 2.0 Client ID

1. Go to **APIs & Services → Credentials**
2. **Create Credentials → OAuth client ID**
3. Application type: **Web application**
4. **Authorized JavaScript origins** (examples):
   - `http://localhost:5173`
   - `https://your-production-domain.com`
5. **Authorized redirect URIs** (must match frontend popup callbacks):
   - `http://localhost:5173/oauth/google/drive/callback`
   - `http://localhost:5173/oauth/google/youtube/callback`
   - Production equivalents on your domain

Copy the **Client ID** and **Client secret**.

## 5. Map credentials to environment variables

### Backend (`backend/.env` or deployment secrets)

| Variable | Value |
|----------|--------|
| `GOOGLE_CLIENT_ID` | OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | OAuth client secret |
| `GOOGLE_OAUTH_REDIRECT_URI` | Backend token exchange redirect (often same as one frontend callback URI) |
| `API_PUBLIC_BASE_URL` | Public API base used for signed media stream URLs (e.g. `https://api.posteaze.com/api`) |
| `ENCRYPTION_KEY` | Already required — used to encrypt stored tokens and sign stream URLs |

### Frontend (`frontend/.env`)

| Variable | Value |
|----------|--------|
| `VITE_GOOGLE_CLIENT_ID` | Same OAuth client ID |
| `VITE_GOOGLE_DRIVE_REDIRECT_URI` | `http://localhost:5173/oauth/google/drive/callback` |
| `VITE_GOOGLE_YOUTUBE_REDIRECT_URI` | `http://localhost:5173/oauth/google/youtube/callback` |

Restart backend and frontend after changing env vars.

## 6. Quotas and limits

### YouTube Data API v3

- Default quota: **10,000 units/day** per project
- `videos.insert` (resumable upload) costs about **1,600 units** per upload → roughly **6 uploads/day** at default quota
- Request a quota increase in Cloud Console if you need more production volume

### Google Drive API

- Per-user rate limits apply; large file streaming uses range requests and does not copy files to PostEaze blob storage when over 50MB

### PostEaze media limits

| Path | Limit |
|------|--------|
| Small import (≤ 50MB) | Copied to Vercel Blob at import |
| Large video (> 50MB) | Stays in Drive; served via signed stream proxy |
| Instagram image | ≤ 8MB |
| Instagram feed video | ≤ 100MB |
| Instagram Reels | ≤ 1GB (future) |
| YouTube | Video only; streams Drive → YouTube in chunks |

## 7. Token behavior

- Use `access_type=offline` and `prompt=consent` on the OAuth authorize URL so Google returns a **refresh token**
- Refresh tokens are stored encrypted in `user_integrations` (Drive) and `channel_tokens` (YouTube)
- In **Testing** mode, refresh tokens may expire after **7 days**; move the OAuth app to **Production** for long-lived refresh tokens for all users

## 8. Verify the integration

1. **Drive:** Media Workspace → Connect Google Drive → Import a file
2. **Large video:** Import a >50MB video; confirm `storage_provider` is `google_drive` and a `stream_url` appears on the version
3. **YouTube:** Channels → YouTube → Connect channel → Schedule a video post to YouTube
4. **Stream proxy:** Open the signed `stream_url` in a browser or `curl -I` and confirm `Content-Type` and `Accept-Ranges: bytes`

## Related docs

- [Firebase Authentication](./firebase-authentication.md) — login (separate from Drive/YouTube OAuth)
- Migration `010_google_drive` — `user_integrations`, Drive columns on media tables
