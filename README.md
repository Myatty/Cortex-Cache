# Cortex Cache — Project Documentation

## Overview

Cortex Cache is a Go web application for creating, viewing, and managing short notes called "snippets" with user accounts, sessions, CSRF protection, and server-side rendered HTML templates.  
Main app entry: `main.go`. The app uses dependency injection via an application struct to pass shared services (logger, DB models, session manager, template cache) to handlers.

## Quick Links

- Application entry: `main.go`
- Router and middleware: `routes.go` 
- Request handlers: `handlers.go`
- Template loading + helper: `templates.go`
- Helpers + rendering: `helpers.go`
- Middleware: `middleware.go`
- Models (DB access): `snippets.go`, `users.go`
- Custom defined Errors: `errors.go`
- Validator: `validator.go`
- Embedded UI files: `efs.go` and templates/static files in `html` and `static`
- Module and deps: `go.mod`

## Technologies and Libraries

- Go 1.24 (toolchain pinned in `go.mod`)
- HTTP routing: [github.com/julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- Middleware chaining: [github.com/justinas/alice](https://github.com/justinas/alice)
- Sessions: [github.com/alexedwards/scs/v2](https://github.com/alexedwards/scs/v2) and MySQL store [github.com/alexedwards/scs/mysqlstore](https://github.com/alexedwards/scs/mysqlstore)
- CSRF protection: [github.com/justinas/nosurf](https://github.com/justinas/nosurf)
- Form decoding: [github.com/go-playground/form/v4](https://github.com/go-playground/form)
- MySQL driver: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- Password hashing: [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- Templates: `html/template`, embedded via Go's embed FS (see `ui/efs.go`)

## High-Level Architecture

### `main.go` initializes global services:
- Loggers
- DB pool via `openDB` (`cmd/web/main.go#openDB`)
- Template cache via `newTemplateCache` (`cmd/web/templates.go`)
- Form decoder (`go-playground/form`)
- SCS session manager with MySQL store
- Constructs the application struct which holds these dependencies and is used as the receiver for handlers and middleware.

### HTTP server
Uses a configured `http.Server` instance in `main.go` so timeouts, TLS, and `ErrorLog` are set explicitly.

### Routing and Middleware
- `routes()` (`cmd/web/routes.go`) builds routes with `httprouter` and wraps handlers with middleware chains produced by `alice`.
- The chain includes CSRF (`noSurf`), session loading (`app.sessionManager.LoadAndSave`), and authentication (`app.authenticate`), plus standard middleware like `secureHeaders`, `logRequest`, and `recoverPanic`.
- Handlers are methods on the `application` type so they can access shared dependencies (see `handlers.go`).
- Templates are embedded at build time and parsed into a cache map used for rendering (`cmd/web/templates.go` and `ui.Files` from `ui/efs.go`).

## Key Implementation Details

### Application Struct (Dependency Container)
Defined in `main.go`.  
Example:  
`application` holds `infoLog`, `errorLog`, `snippets`, `users`, `templateCache`, `formDecoder`, `sessionManager`.  
Purpose: keep handler methods simple and testable by providing a single place for app-wide dependencies.

### Routing
- Router: `httprouter` in `routes.go`.
- Dynamic routes use `httprouter` params (e.g., `/snippet/view/:id`) and handlers extract params via `httprouter.ParamsFromContext(r.Context())`.
- Static files served from the embedded FS:  
  `router.Handler(http.MethodGet, "/static/*filepath", fileServer)` using `http.FS(ui.Files)`.

### Middleware
- **Secure headers:** `secureHeaders` — sets headers like CSP, X-Frame-Options.
- **CSRF protection:** `noSurf` uses `nosurf` and sets Secure/HttpOnly cookie.
- **Sessions:** `scs` session manager (MySQL store) in `main`.
- **Authentication middleware:**
  - `authenticate` checks session `authenticatedUserID`, sets a context key when authenticated.
  - `requireAuthentication` redirects unauthenticated users to login and disables caching for protected pages.

### Sessions and CSRF
- **Sessions:** `scs v2` with MySQL store (configured in `cmd/web/main.go`). Lifetime set to 12 hours.
- **CSRF:** `nosurf` token is injected into templates via `templateData.CSRFToken` in `helpers.go`.

### Templates
- Templates and static assets are embedded with Go's `embed`: `efs.go`.
- Template cache is built in `templates.go` by parsing base, partials, and the page templates. Templates are stored in a map keyed by filename.
- **Rendering pipeline:**
  `app.render` in `helpers.go` retrieves the parsed template from the cache and executes base with the provided `templateData`.

### Form Handling & Validation
- Form decoding uses `go-playground/form` in `helpers.go` via `decodePostForm`.
- Validation helpers provided in `validator.go` (e.g., `NotBlank`, `MaxChars`, `Matches`, `PermittedValue`) and a `Validator` type that collects `FieldErrors` and `NonFieldErrors`.
- Forms embed `validator.Validator` to reuse its methods (`CheckField`, `AddFieldError`, `Valid`), e.g., `snippetCreateForm` in `handlers.go`.

### Models / Database Access

#### Snippets (`snippets.go`)
- `Insert(title, content, expires int)` → returns inserted ID  
- `Get(id)` → retrieves snippet or returns `ErrNoRecord`  
- `Latest()` → returns up to 10 latest non-expired snippets

#### Users (`users.go`)
- `Insert(name, email, password)` — hashes password with bcrypt and checks for duplicate email using MySQL error inspection  
- `Authenticate(email, password)` — checks credentials and returns user ID  
- `Exists(id)` — checks if user exists  
- Models return application-level errors from `errors.go` (`ErrNoRecord`, `ErrInvalidCredentials`, `ErrDuplicateEmail`).

---

## Database Schema (Inferred)

### Snippets Table (used by `snippets` model)
- **Columns:**
  - `id` (PK, auto-increment)
  - `title` (varchar)
  - `content` (text)
  - `created` (timestamp)
  - `expires` (timestamp)

### Users Table (used by `users` model)
- **Columns:**
  - `id` (PK)
  - `name`
  - `email` (unique constraint; code checks `users_uc_email` index)
  - `hashed_password`
  - `created` (UTC_TIMESTAMP())

### Example `CREATE TABLE` SQL (adapt to your MySQL flavor)
```sql

CREATE TABLE snippets (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATETIME NOT NULL,
  expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL,
  created DATETIME NOT NULL
);

CREATE TABLE sessions (
  token CHAR(43) PRIMARY KEY,
  data BLOB NOT NULL,
  expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

```

## Session Table

### Purpose
The session table stores server-side session data for the `scs` session manager (configured in `cmd/web/main.go`). It holds serialized session state (e.g., flash messages, authenticated user ID) so handlers and middleware can read/write session values across requests.

### Why It’s Used
- Persist sessions across server restarts.
- Share sessions across multiple app instances (horizontal scaling).
- Keep sensitive data off client cookies (cookies only store session token).
- Allow server-side session expiry and revocation.

### Typical Contents
- **session token / id** (string) — references a row in the session table via the client cookie.
- **serialized session data** (BLOB) — includes values like flash messages, `authenticatedUserID`.
- **expiry timestamp** — when the session becomes invalid.
- (optional) **created/updated timestamps** or **user id** for auditing.

### Lifecycle
1. Client holds a cookie with the session token.
2. On each request, `scs` looks up the token in the session table, loads data into memory, and exposes `Get/Put` methods.
3. When the app writes session changes, `scs` updates the row; expired rows are cleaned up.

---

## Data Flow Example: Create Snippet

1. **Client** sends a `POST` request with `title`, `content`, and `expires` fields to `/snippet/create`.
2. **Route** points to `app.snippetCreatePost` (protected by `requireAuthentication` middleware).
3. **Handler** calls `app.decodePostForm` to decode POST values into `snippetCreateForm`.
4. **Validator** functions check input fields and add field errors if validation fails.  
   - If invalid, the form is re-rendered with form field errors.
5. **On success**, `app.snippets.Insert` stores the snippet in the database.
6. **Handler** sets a flash message via `app.sessionManager.Put(r.Context(), "flash", "...")` and redirects to the snippet’s view page.
7. **Flash message** is read in `templateData` via `app.sessionManager.PopString` and rendered on the next page.

## Running and Configuration

### Default Flags

- `-addr` — default `:4000` (listen address), see `main.go`
- `-dsn` — default `web:Lucifer@/cortexCache?parseTime=true` (MySQL DSN)

### Build and Run

- NOTE: Please Install Go and configure underlying MySQL Database on your machine before Building and running the Program

```bash

# 1. Create MYSQL Database
CREATE DATABASE cortexCache;

# 2. Import the database (OR CREATE ON YOUR OWN USING DATABASE SCHEMA ABOVE)
mysql -u [username] -p cortexCache < cortexCache.sql

# OPTIONAL: Update the default DSN in `main.go`
dsn := flag.String("dsn", "[usernameUsedByApp]:Password@/cortexCache?parseTime=true", "MySQL data source name")
  
# 3. Generate self-signed TLS Cert
mkdir tls && cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost


# 4. Compile the go program(from repo root)
go build -o cortexcache ./cmd/web
# run the program
./cortexcache -addr=":4000" -dsn="username:password@tcp(localhost:3306)/cortexCache?parseTime=true"

# OR run directly using
go run ./cmd/web -addr=":4000" -dsn="username:password@tcp(localhost:3306)/cortexCache?parseTime=true"
```

### TLS

- ListenAndServeTLS is used in main. The TLS certificate and key are expected at cert.pem and key.pem (the tls directory is gitignored). For local development you can replace call with srv.ListenAndServe() or provide self-signed cert/key.

- ## Security Considerations & Further Improvements

- **Password hashing:** Uses `bcrypt` which costs (`12`).
- **CSRF protection:** Enabled via `nosurf` and `CSRFToken` is injected into templates.
- **Sessions:** Stored in database (`scs mysqlstore`). Ensure session table is created according to `scs/mysqlstore` docs.
- **TLS:** Required by `main.go` — ensure certificates are present under tls/ for testing and production.
- **Static file server:** Uses embedded FS, which is safe. Be careful with serving files from arbitrary paths if that changes.
- **Further Improvements:** Considering rate limiting and stricter Content Security Policy (CSP) if exposing to the public internet. Will be containerized using Docker.

## Testing and Extension Points

- Handlers and middleware are methods on `application`, making unit injection straightforward. Can create tests that initialize an application with mocked DB/models and an in-memory session manager.
- Template rendering uses an in-memory cache that makes template tests reproducible.
- To add features: add models under `models`, handlers in `handlers.go`, and routes in `routes.go`.
- Further Testing will be done and will be updating this section once I've done it.(Will be doing Unit testing, End to End testing and Integration Testing.)
- This app will be containerized using Docker.

---

## Files of Interest (Quick Reference)

- `main.go` — app bootstrap and server  
- `routes.go` — routing & handler wiring  
- `handlers.go` — application handlers for snippets and user account flows  
- `helpers.go` — rendering, `decodePostForm`, error helpers  
- `middleware.go` — security & request middleware  
- `templates.go` — template functions and cache  
- `internal/models/*.go` — DB access and domain errors  
- `validator.go` — form validation helpers  
- `ui/html/*.tmpl.html` and `ui/static/*` — frontend templates and assets  
- `efs.go` — embedding templates and static assets
