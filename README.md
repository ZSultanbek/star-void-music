# 🎧 Star Void Music — Full Project Specification

## 📌 Overview

**Star Void Music** is a self-hosted music streaming service built with **Golang (backend)** and **React (frontend)**. The system allows users to stream music from a personal server (home PC), manage their library, and interact with playlists.

This project is designed to:

* Teach backend architecture in Go
* Implement real-world audio streaming
* Build a modern frontend with React
* Deploy a local server accessible via the internet

---

# 🧱 Tech Stack

## Backend

* Language: Go
* Framework: Gin
* Authentication: JWT
* Database: PostgreSQL
* Query Tool: sqlc (preferred) or pgx
* File Storage: Local filesystem (`/storage/music`)

## Frontend

* React (Vite or Next.js optional)
* Fetch/Axios for API calls
* Native HTML5 `<audio>` player

## DevOps / Infrastructure

* Docker + Docker Compose
* Cloudflare Tunnel (cloudflared) for external access

---

# 🏗️ System Architecture

```
[ Client (Browser / Phone) ]
            ↓
     React Frontend (SPA)
            ↓ HTTP API
      Go Backend (Gin)
            ↓
   Service Layer (Logic)
            ↓
 Repository Layer (DB)
            ↓
PostgreSQL + File Storage
```

---

# 🗂️ Project Structure (Monorepo)

```
star-void-music/
│
├── backend/
│   ├── cmd/
│   │   └── server/main.go
│   │
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── models/
│   │   └── streaming/
│   │
│   ├── pkg/
│   │   ├── jwt/
│   │   └── hash/
│   │
│   ├── db/
│   │   ├── migrations/
│   │   ├── queries/
│   │   └── sqlc/     # Generated Go code from sqlc
│   │   
│   │
│   ├── storage/music/
│   ├── Dockerfile
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── services/   # API calls
│   │   └── App.jsx
│   │
│   ├── public/
│   ├── package.json
│   └── vite.config.js
│
├── docker-compose.yml
└── project.md
```

### ✅ Backend Initialization Scope (First Delivery)

Initialize backend with a clean, scalable skeleton:

- Create base module and folders:
  - `cmd/server`
  - `internal/config`
  - `internal/handler`
  - `internal/service`
  - `internal/repository`
  - `internal/middleware`
  - `internal/models`
  - `internal/streaming`
  - `pkg/jwt`
  - `pkg/hash`
- Bootstrap Gin server in `cmd/server/main.go`
- Add `GET /health` endpoint returning basic service status
- Wire router setup and dependency injection entrypoint (minimal, ready for expansion)
- Keep implementation simple, readable, and production-oriented

---

# 🧠 Core Features

## 1. Authentication

* Register
* Login
* JWT-based authentication
* Roles:

  * `user`
  * `admin`

---

## 2. Songs

* Admin uploads songs
* Metadata stored in DB
* File stored locally
* Users can stream songs

---

## 3. User Library

* Users can add/remove songs
* Personal collection of songs

---

## 4. Streaming

* HTTP streaming via Range Requests
* Supports:

  * Seek
  * Pause/Resume
* Implemented via `http.ServeContent`

---

## 5. Playlists (Future)

* Create playlists
* Add/remove songs
* Ordered songs list

---

# 🗄️ Database Schema

## users

```
  id (UUID)
  email (string)
  password_hash (string)
  role (string) // user | admin
  created_at
```

## artists
```
  id (UUID, PK)
  name (string)
  slug (string, UNIQUE) -- e.g. "star-void"
  created_at (timestamp)
```

## albums
```
  id (UUID, PK)
  title (string)
  artist_id (UUID, FK -> artists.id)
  cover_image_url (string)
  release_date (date)
  created_at (timestamp)
```

## songs
```
  id (UUID, PK)
  title (string)
  album_id (UUID, FK -> albums.id)
  filepath (string)
  duration (int)
  uploaded_by (UUID, FK -> users.id)
  created_at (timestamp)
```

## user_library
```
  user_id (UUID, FK -> users.id)
  song_id (UUID, FK -> songs.id)
  added_at (timestamp)

  PRIMARY KEY (user_id, song_id)
```

## playlists
```
  id (UUID)
  user_id
  name
  created_at
```

## playlist_songs
```
  playlist_id
  song_id
  position
```


---
# REST API structure

## users
```
    POST   /api/users
    GET    /api/users
    GET    /api/users/:id
    GET    /api/users?email=ex@mail.com
    PATCH  /api/users/:id
    DELETE /api/users/:id
```

## songs
```
    POST   /api/songs
    GET    /api/songs/:id
    GET    /api/songs/:id/stream
    PATCH  /api/songs/:id
    DELETE /api/songs/:id
```

## user_library
```
    POST   /api/me/library (add song)
    GET    /api/me/library (list songs)
    DELETE /api/me/library/:song_id
```

## artists
```
    POST   /api/artists
    GET    /api/artists
    GET    /api/artists/:id
    GET    /api/artists/:id/albums
    GET    /api/artists/:id/songs
    PATCH  /api/artists/:id
```

## albums
```
    POST   /api/albums
    GET    /api/albums
    GET    /api/albums/:id
    GET    /api/albums/:id/songs
    PATCH  /api/albums/:id
```
---

---

# 🔗 Relationships

* User ↔ Songs (many-to-many via user_library)
* Album → Songs (one-to-many)
* User → Playlists (one-to-many)
* Playlist ↔ Songs (many-to-many with order)

---

# 🔁 User Flow

## 🧑 Regular User

### 1. Register / Login

* Sends credentials → receives JWT

### 2. Browse Songs

* Fetch `/api/songs`

### 3. Stream Song

* Click play → frontend loads:

```
GET /api/songs/{id}/stream
```

### 4. Add to Library

```
POST /api/me/library
```

### 5. View Library

```
GET /api/me/library
```

---

## 🛠️ Admin

### Upload Song

* Upload via multipart form:

```
POST /api/admin/songs
```

Backend:

* Saves file to `/storage/music`
* Stores metadata in DB

---

# 🔊 Streaming Design

* Endpoint:

```
GET /api/songs/{id}/stream
```

* Backend:

  * Validate JWT
  * Locate file
  * Use `http.ServeContent`

* Frontend:

```html
<audio controls src="/api/songs/{id}/stream"></audio>
```

---

# 🔐 Security

* JWT for authentication
* Middleware to protect routes
* Role-based access for admin endpoints
* Validate file uploads (only mp3 initially)

---

# 🐳 Docker Setup

## docker-compose.yml

```
version: '3.9'

services:
  db:
    image: postgres:15
    container_name: svm-db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: starvoid
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    container_name: svm-backend
    ports:
      - "8080:8080"
    depends_on:
      - db

  frontend:
    build: ./frontend
    container_name: svm-frontend
    ports:
      - "3000:3000"
```

---

# 🌐 External Access (Cloudflare Tunnel)

Run:

```
cloudflared tunnel --url http://localhost:8080
```

This exposes your backend securely to your phone.

---

# 🚀 Development Phases

## Phase 1 (MVP)

* Auth
* Upload songs
* Stream songs
* Basic UI
* **Backend skeleton initialized (Gin bootstrap + health endpoint + clean folder layout)**

## Phase 2

* User library
* Search

## Phase 3

* Playlists
* Downloads

---

# 🧠 Design Principles

* Keep it simple first
* Clean architecture (handler → service → repository)
* Avoid overengineering
* Build streaming early
* Focus on working MVP

---


**End of Specification**
