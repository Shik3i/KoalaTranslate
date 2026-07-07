package httpapi

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"koalatranslate/internal/auth"
	"koalatranslate/internal/github"
	"koalatranslate/internal/i18n"
	"koalatranslate/internal/projects"
)

// Server struct handles dependencies.
type Server struct {
	DB               *sql.DB
	SessionSecret    string // Used as cookie config check or simple verification if needed
	GitHubToken      string
	GitHubBotOwner   string
	GitHubBotForkRepo string
}

// NewServer creates a new API server.
func NewServer(db *sql.DB) *Server {
	return &Server{
		DB:                 db,
		SessionSecret:      getEnv("SESSION_SECRET", "default_secret_32_bytes_long_required"),
		GitHubToken:        os.Getenv("GITHUB_TOKEN"),
		GitHubBotOwner:     os.Getenv("GITHUB_BOT_OWNER"),
		GitHubBotForkRepo:  os.Getenv("GITHUB_BOT_FORK_REPO"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Response helpers
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// SessionContextKey is the custom type for request context keys.
type SessionContextKey string
const UserContextKey SessionContextKey = "user"

// AuthMiddleware protects routes from unauthenticated access.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health check is public
		if r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			writeError(w, http.StatusUnauthorized, "missing session token")
			return
		}

		userID, username, err := auth.ValidateSession(s.DB, cookie.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired session: "+err.Error())
			return
		}

		// Keep username in headers or extend context if needed
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
		r.Header.Set("X-Username", username)

		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes registers all API routes on the serve mux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)

	// Protected routes (we'll wrap these in s.AuthMiddleware if registering globally,
	// or we can wrap a sub-mux. For standard lib, we can wrap individual handlers
	// or use a middleware chain. Let's write a simple wrapper helper).
	wrap := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s.AuthMiddleware(http.HandlerFunc(h)).ServeHTTP(w, r)
		}
	}

	mux.HandleFunc("GET /api/auth/me", wrap(s.handleMe))
	mux.HandleFunc("GET /api/projects", wrap(s.handleGetProjects))
	mux.HandleFunc("GET /api/projects/{projectId}", wrap(s.handleGetProjectByID))
	mux.HandleFunc("POST /api/projects/{projectId}/import/upload", wrap(s.handleImportUpload))
	mux.HandleFunc("POST /api/projects/{projectId}/import/github", wrap(s.handleImportGitHub))
	mux.HandleFunc("GET /api/projects/{projectId}/locales", wrap(s.handleGetLocales))
	mux.HandleFunc("GET /api/projects/{projectId}/translations", wrap(s.handleGetTranslations))
	mux.HandleFunc("PATCH /api/projects/{projectId}/translations/{key}", wrap(s.handlePatchTranslation))
	mux.HandleFunc("GET /api/projects/{projectId}/export", wrap(s.handleExport))
	mux.HandleFunc("POST /api/projects/{projectId}/export/github-pr", wrap(s.handleExportGitHubPR))
}

// GET /api/health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// POST /api/auth/login
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var userID int
	var passwordHash string
	err := s.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", req.Username).Scan(&userID, &passwordHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	if !auth.CheckPasswordHash(req.Password, passwordHash) {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// 24 hour session duration
	token, expiresAt, err := auth.CreateSession(s.DB, userID, 24)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session: "+err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false, // Set true in production if running HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"username": req.Username,
		"status":   "logged_in",
	})
}

// POST /api/auth/logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		_ = auth.DeleteSession(s.DB, cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		MaxAge:   -1,
	})

	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

// GET /api/auth/me
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Username")
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"username":      username,
		"authenticated": true,
	})
}

// GET /api/projects
func (s *Server) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("SELECT id, name, description FROM projects")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database query failed: "+err.Error())
		return
	}
	defer rows.Close()

	var list []projects.Project
	for rows.Next() {
		var p projects.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			writeError(w, http.StatusInternalServerError, "database scan failed: "+err.Error())
			return
		}
		list = append(list, p)
	}

	writeJSON(w, http.StatusOK, list)
}

// GET /api/projects/{projectId}
func (s *Server) handleGetProjectByID(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var p projects.Project
	err := s.DB.QueryRow("SELECT id, name, description FROM projects WHERE id = ?", projectID).Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "project not found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// Helper to store parsed dictionary into SQLite
func (s *Server) saveImportedTranslations(projectID, localeCode, localeName string, isSource bool, dict i18n.MessagesDict, username string, sourceType string, sourceURL string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Ensure locale is registered
	_, err = tx.Exec(
		`INSERT INTO locales (project_id, code, name, is_source) 
		 VALUES (?, ?, ?, ?) 
		 ON CONFLICT(project_id, code) DO UPDATE SET name=excluded.name, is_source=excluded.is_source`,
		projectID, localeCode, localeName, isSource,
	)
	if err != nil {
		return err
	}

	// 2. Insert keys and values
	for keyName, entry := range dict {
		// Insert translation key if not exists
		_, err = tx.Exec(
			`INSERT INTO translation_keys (project_id, key_name, description) 
			 VALUES (?, ?, ?) 
			 ON CONFLICT(project_id, key_name) DO UPDATE SET description=coalesce(excluded.description, translation_keys.description)`,
			projectID, keyName, entry.Description,
		)
		if err != nil {
			return err
		}

		var keyID int
		err = tx.QueryRow("SELECT id FROM translation_keys WHERE project_id = ? AND key_name = ?", projectID, keyName).Scan(&keyID)
		if err != nil {
			return err
		}

		// Insert/Update value
		// Since it's imported, it resets the dirty flag
		_, err = tx.Exec(
			`INSERT INTO translation_values (key_id, locale_code, value, is_dirty, last_edited_by, updated_at) 
			 VALUES (?, ?, ?, 0, ?, CURRENT_TIMESTAMP) 
			 ON CONFLICT(key_id, locale_code) DO UPDATE SET value=excluded.value, is_dirty=0, last_edited_by=excluded.last_edited_by, updated_at=CURRENT_TIMESTAMP`,
			keyID, localeCode, entry.Message, username,
		)
		if err != nil {
			return err
		}
	}

	// 3. Log import
	_, err = tx.Exec(
		`INSERT INTO imports (project_id, locale_code, source_type, source_url, imported_by) 
		 VALUES (?, ?, ?, ?, ?)`,
		projectID, localeCode, sourceType, sourceURL, username,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// POST /api/projects/{projectId}/import/upload
func (s *Server) handleImportUpload(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	username := r.Header.Get("X-Username")

	// Limit multipart memory parsing (e.g. 10MB)
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing translation file")
		return
	}
	defer file.Close()

	localeCode := r.FormValue("locale")
	localeName := r.FormValue("localeName")
	isSourceStr := r.FormValue("isSource")
	if localeCode == "" || localeName == "" {
		writeError(w, http.StatusBadRequest, "locale and localeName fields are required")
		return
	}

	isSource := isSourceStr == "true"

	// Read and parse file
	// Limit read to 10MB + 1 byte
	limitReader := io.LimitReader(file, 10*1024*1024+1)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read file: "+err.Error())
		return
	}

	if len(data) > 10*1024*1024 {
		writeError(w, http.StatusBadRequest, "file exceeds size limit of 10MB")
		return
	}

	dict, err := i18n.ParseAndValidate(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse messages.json: "+err.Error())
		return
	}

	err = s.saveImportedTranslations(projectID, localeCode, localeName, isSource, dict, username, "manual", "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save translations: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "success",
		"keysCount":   len(dict),
		"locale":      localeCode,
		"localeName":  localeName,
	})
}

// POST /api/projects/{projectId}/import/github
func (s *Server) handleImportGitHub(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	username := r.Header.Get("X-Username")

	var req struct {
		URL        string `json:"url"`
		Locale     string `json:"locale"`
		LocaleName string `json:"localeName"`
		IsSource   bool   `json:"isSource"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" || req.Locale == "" || req.LocaleName == "" {
		writeError(w, http.StatusBadRequest, "url, locale, and localeName are required fields")
		return
	}

	// Fetch file content
	client := &http.Client{Timeout: 15 * time.Second}
	data, err := github.FetchGitHubFile(client, req.URL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "GitHub download failed: "+err.Error())
		return
	}

	dict, err := i18n.ParseAndValidate(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to parse messages.json from GitHub: "+err.Error())
		return
	}

	err = s.saveImportedTranslations(projectID, req.Locale, req.LocaleName, req.IsSource, dict, username, "github", req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save translations: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "success",
		"keysCount":   len(dict),
		"locale":      req.Locale,
		"localeName":  req.LocaleName,
	})
}

// GET /api/projects/{projectId}/locales
func (s *Server) handleGetLocales(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	rows, err := s.DB.Query("SELECT code, name, is_source FROM locales WHERE project_id = ?", projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	type Locale struct {
		Code     string `json:"code"`
		Name     string `json:"name"`
		IsSource bool   `json:"isSource"`
	}

	var list []Locale
	for rows.Next() {
		var l Locale
		if err := rows.Scan(&l.Code, &l.Name, &l.IsSource); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		list = append(list, l)
	}

	writeJSON(w, http.StatusOK, list)
}

// TranslationListItem represents a translation value together with its English reference and details.
type TranslationListItem struct {
	Key         string   `json:"key"`
	Description string   `json:"description,omitempty"`
	SourceValue string   `json:"sourceValue"`
	TargetValue string   `json:"targetValue"`
	IsDirty     bool     `json:"isDirty"`
	Warnings    []string `json:"warnings,omitempty"`
}

// GET /api/projects/{projectId}/translations
func (s *Server) handleGetTranslations(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	locale := r.URL.Query().Get("locale")

	if locale == "" {
		writeError(w, http.StatusBadRequest, "locale parameter is required")
		return
	}

	// 1. Fetch the source (is_source = 1) locale code
	var sourceLocale string
	err := s.DB.QueryRow("SELECT code FROM locales WHERE project_id = ? AND is_source = 1 LIMIT 1", projectID).Scan(&sourceLocale)
	if err != nil {
		// Fallback to "en" if not found
		sourceLocale = "en"
	}

	// 2. Fetch all keys for the project
	rows, err := s.DB.Query(
		`SELECT tk.id, tk.key_name, tk.description 
		 FROM translation_keys tk 
		 WHERE tk.project_id = ?`,
		projectID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	type keyInfo struct {
		id   int
		name string
		desc string
	}
	var keys []keyInfo
	for rows.Next() {
		var k keyInfo
		if err := rows.Scan(&k.id, &k.name, &k.desc); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		keys = append(keys, k)
	}

	// 3. For each key, query the source value and target value
	list := []TranslationListItem{}
	for _, k := range keys {
		var sourceVal string
		_ = s.DB.QueryRow("SELECT value FROM translation_values WHERE key_id = ? AND locale_code = ?", k.id, sourceLocale).Scan(&sourceVal)

		var targetVal string
		var isDirty bool
		err = s.DB.QueryRow("SELECT value, is_dirty FROM translation_values WHERE key_id = ? AND locale_code = ?", k.id, locale).Scan(&targetVal, &isDirty)
		if err != nil {
			// No translation yet
			targetVal = ""
			isDirty = false
		}

		warnings := i18n.CheckPlaceholderWarnings(sourceVal, targetVal)

		list = append(list, TranslationListItem{
			Key:         k.name,
			Description: k.desc,
			SourceValue: sourceVal,
			TargetValue: targetVal,
			IsDirty:     isDirty,
			Warnings:    warnings,
		})
	}

	writeJSON(w, http.StatusOK, list)
}

// PATCH /api/projects/{projectId}/translations/{key}
func (s *Server) handlePatchTranslation(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	keyName := r.PathValue("key")
	username := r.Header.Get("X-Username")

	var req struct {
		Locale string `json:"locale"`
		Value  string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Locale == "" {
		writeError(w, http.StatusBadRequest, "locale field is required")
		return
	}

	// Find the key_id
	var keyID int
	err := s.DB.QueryRow("SELECT id FROM translation_keys WHERE project_id = ? AND key_name = ?", projectID, keyName).Scan(&keyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "key not found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Upsert the value and mark is_dirty = 1
	_, err = s.DB.Exec(
		`INSERT INTO translation_values (key_id, locale_code, value, is_dirty, last_edited_by, updated_at) 
		 VALUES (?, ?, ?, 1, ?, CURRENT_TIMESTAMP) 
		 ON CONFLICT(key_id, locale_code) DO UPDATE SET value=excluded.value, is_dirty=1, last_edited_by=excluded.last_edited_by, updated_at=CURRENT_TIMESTAMP`,
		keyID, req.Locale, req.Value, username,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update database: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

// GET /api/projects/{projectId}/export
func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	locale := r.URL.Query().Get("locale")

	if locale == "" {
		writeError(w, http.StatusBadRequest, "locale parameter is required")
		return
	}

	// Fetch all translation values for this locale
	rows, err := s.DB.Query(
		`SELECT tk.key_name, tk.description, tv.value 
		 FROM translation_values tv 
		 JOIN translation_keys tk ON tv.key_id = tk.id 
		 WHERE tk.project_id = ? AND tv.locale_code = ?`,
		projectID, locale,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	dict := make(i18n.MessagesDict)
	for rows.Next() {
		var key, desc, val string
		if err := rows.Scan(&key, &desc, &val); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dict[key] = i18n.MessageEntry{
			Message:     val,
			Description: desc,
		}
	}

	// Log export audit trail
	username := r.Header.Get("X-Username")
	_, _ = s.DB.Exec(
		"INSERT INTO exports (project_id, locale_code, export_type, exported_by) VALUES (?, ?, ?, ?)",
		projectID, locale, "download", username,
	)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=messages.json"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Build JSON with indentation to match standard format
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(dict)
}

// POST /api/projects/{projectId}/export/github-pr
func (s *Server) handleExportGitHubPR(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	username := r.Header.Get("X-Username")

	var req struct {
		Locale         string `json:"locale"`
		UpstreamOwner  string `json:"upstreamOwner"`
		UpstreamRepo   string `json:"upstreamRepo"`
		UpstreamBranch string `json:"upstreamBranch"`
		FilePath       string `json:"filePath"` // e.g. _locales/de/messages.json
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Locale == "" || req.UpstreamOwner == "" || req.UpstreamRepo == "" || req.UpstreamBranch == "" || req.FilePath == "" {
		writeError(w, http.StatusBadRequest, "locale, upstreamOwner, upstreamRepo, upstreamBranch, and filePath are required")
		return
	}

	if s.GitHubToken == "" || s.GitHubBotOwner == "" || s.GitHubBotForkRepo == "" {
		writeError(w, http.StatusServiceUnavailable, "GitHub bot fork credentials are not configured on this server")
		return
	}

	// 1. Fetch translation entries for export
	rows, err := s.DB.Query(
		`SELECT tk.key_name, tk.description, tv.value 
		 FROM translation_values tv 
		 JOIN translation_keys tk ON tv.key_id = tk.id 
		 WHERE tk.project_id = ? AND tv.locale_code = ?`,
		projectID, req.Locale,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	dict := make(i18n.MessagesDict)
	for rows.Next() {
		var key, desc, val string
		if err := rows.Scan(&key, &desc, &val); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dict[key] = i18n.MessageEntry{
			Message:     val,
			Description: desc,
		}
	}

	// 2. Format JSON
	jsonData, err := json.MarshalIndent(dict, "", "  ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to format export JSON: "+err.Error())
		return
	}

	// 3. Base64 encode JSON
	base64Content := base64.StdEncoding.EncodeToString(jsonData)

	// 4. Create PR via bot fork model
	client := &http.Client{Timeout: 30 * time.Second}
	commitMsg := fmt.Sprintf("Update %s translations via KoalaTranslate", req.Locale)
	prTitle := fmt.Sprintf("Update %s translations via KoalaTranslate", req.Locale)
	
	// Count dirty keys
	var changedKeys int
	_ = s.DB.QueryRow(
		`SELECT count(*) FROM translation_values tv 
		 JOIN translation_keys tk ON tv.key_id = tk.id 
		 WHERE tk.project_id = ? AND tv.locale_code = ? AND tv.is_dirty = 1`,
		projectID, req.Locale,
	).Scan(&changedKeys)

	// Fetch source locale
	var sourceLocale string
	_ = s.DB.QueryRow("SELECT code FROM locales WHERE project_id = ? AND is_source = 1 LIMIT 1", projectID).Scan(&sourceLocale)
	if sourceLocale == "" {
		sourceLocale = "en"
	}

	prBody := fmt.Sprintf(
		"This PR updates %s translations for KoalaSync.\n\n"+
			"Generated by KoalaTranslate.\n\n"+
			"Summary:\n"+
			"* Project: KoalaSync\n"+
			"* Locale: %s\n"+
			"* File: %s\n"+
			"* Changed keys: %d\n"+
			"* Source locale: %s\n\n"+
			"Please review the changed translations before merging.",
		req.Locale, req.Locale, req.FilePath, changedKeys, sourceLocale,
	)

	prNumber, prURL, err := github.CreateForkPR(
		client,
		s.GitHubToken,
		s.GitHubBotOwner,
		s.GitHubBotForkRepo,
		req.UpstreamOwner,
		req.UpstreamRepo,
		req.UpstreamBranch,
		req.FilePath,
		commitMsg,
		prTitle,
		prBody,
		base64Content,
	)

	if err != nil {
		writeError(w, http.StatusBadGateway, "GitHub PR creation failed: "+err.Error())
		return
	}

	// 5. Reset dirty status after successful commit/PR creation
	_, err = s.DB.Exec(
		`UPDATE translation_values 
		 SET is_dirty = 0 
		 WHERE locale_code = ? AND key_id IN (
			 SELECT id FROM translation_keys WHERE project_id = ?
		 )`,
		req.Locale, projectID,
	)
	if err != nil {
		// Log warning but don't fail request since PR was already opened
		fmt.Fprintf(os.Stderr, "failed to reset is_dirty flag: %v\n", err)
	}

	// 6. Log PR export
	_, _ = s.DB.Exec(
		`INSERT INTO exports (project_id, locale_code, export_type, pr_number, pr_url, exported_by) 
		 VALUES (?, ?, ?, ?, ?, ?)`,
		projectID, req.Locale, "github_pr", prNumber, prURL, username,
	)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "success",
		"prNumber":  prNumber,
		"prUrl":     prURL,
		"branch":    fmt.Sprintf("koalatranslate-update-%s", req.Locale),
		"keysCount": changedKeys,
	})
}
