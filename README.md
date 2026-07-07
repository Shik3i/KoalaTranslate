# KoalaTranslate

> **Important**: KoalaTranslate currently starts as a hardcoded KoalaStuff translation portal. It is not yet a public platform for arbitrary projects.

KoalaTranslate is a self-hosted visual translation portal for KoalaStuff projects, starting with KoalaSync. It enables non-technical translators to contribute localization changes directly in-context by clicking on simulated extension screens instead of editing raw JSON files or managing GitHub workflows manually.

---

## Tech Stack
* **Frontend**: SvelteKit + TypeScript, styled with Vanilla CSS (self-hosted layout assets, local system fonts, dark mode).
* **Backend**: Go with HTTP standard library routing and CORS support.
* **Database**: SQLite database using CGO-free `modernc.org/sqlite` drivers.
* **Packaging**: Single container multi-stage build, serving embedded SvelteKit assets directly from Go using `go:embed`.

---

## Non-Goals & Scope (Phase 1)
* **No open registration / registration UI**: Admin account seeded via environment variables.
* **No runtime AI**: No automated translation or screen-building.
* **No SaaS features**: No project creation/management UI.
* **No direct commits to main**: No direct write access to upstream project main branches.

---

## Environment Variables
The application reads the following values on startup:
* `ADMIN_USERNAME`: Admin login user (default: `admin`)
* `ADMIN_PASSWORD`: Admin login password (default: `adminpassword`)
* `SESSION_SECRET`: Secure cookie secret key
* `DATABASE_PATH`: SQLite database file path (default: `koalatranslate.db`)
* `PORT`: Server listening port (default: `8080`)
* `GITHUB_TOKEN`: Bot account Personal Access Token (optional)
* `GITHUB_BOT_OWNER`: Bot account owner username holding the fork (optional)
* `GITHUB_BOT_FORK_REPO`: Bot fork repository name (optional)

---

## Quick Start (Docker Compose)

To start the unified production application, run:
```bash
docker compose up --build
```
The application will listen on [http://localhost:8080](http://localhost:8080). 

Sign in with your seeded `ADMIN_USERNAME` and `ADMIN_PASSWORD` credentials.

---

## Development Setup

If you want to run frontend and backend separately in development mode:

1. **Start the Go Backend**:
   ```bash
   cd server
   go run cmd/server/main.go
   ```
2. **Start the SvelteKit Dev Server**:
   ```bash
   cd web
   npm run dev
   ```
   Vite will run on [http://localhost:5173](http://localhost:5173) and proxy all `/api/*` endpoints to the Go server on port 8080.

---

## Translation Workflow

### 1. Import translations
On the dashboard, click **Import** to open the drawer:
* **Manual upload**: Select and upload a `messages.json` file from your disk.
* **GitHub URL import**: Paste a github.com blob URL or raw.githubusercontent.com URL pointing to a `messages.json` dictionary. Only GitHub domains are allowed, and files are limited to 10MB.

### 2. Edit translations
Select a target language and pick one of the 7 available screens. Hovering over a text field in the simulator displays a dashed border. Click the text box to activate the editor sidebar where you can modify the translation, view descriptions, and check placeholder validation warning reports (e.g. checks if `{seconds}` is present in target).

### 3. Export translations
Click **Export** to save your work:
* **Download**: Directly download the updated `messages.json` file.
* **GitHub Pull Request**: 
  1. KoalaTranslate will verify the bot account has a fork of the target repo.
  2. It will create a branch in the bot's fork (e.g., `koalatranslate-update-de-abcdef`).
  3. It will commit the updated file to that fork branch.
  4. It will open a Pull Request from `bot-fork:branch` to `upstream:main`.
  *Note: KoalaTranslate never writes directly to the upstream/main project.*

---

## Adding More Hardcoded Projects

The code is structured with a modular architecture to easily support additional KoalaStuff projects:

### 1. Backend Structure (`server/`)
* To add project metadata, update [projects.go](file:///Users/koala/Documents/Workspaces/KoalaTranslate/server/internal/projects/projects.go) to seed new entries in the SQLite `projects` table.
* Keep the API controllers in [httpapi.go](file:///Users/koala/Documents/Workspaces/KoalaTranslate/server/internal/httpapi/httpapi.go) generic so they handle translations and locales dynamically based on `:projectId`.

### 2. Frontend Structure (`web/`)
* Recreate new visual screens under `web/src/lib/projects/{newProjectID}/screens/`.
* Declare the screens list and dynamically load the respective project screens in `web/src/routes/+page.svelte`.
