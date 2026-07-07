package projects

import (
	"database/sql"
)

// Project struct holds project metadata.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SeedProjects seeds default projects into SQLite.
func SeedProjects(db *sql.DB) error {
	defaultProjects := []Project{
		{
			ID:          "koalasync",
			Name:        "KoalaSync",
			Description: "Visual translation portal for KoalaSync extension.",
		},
	}

	for _, p := range defaultProjects {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)", p.ID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			_, err = db.Exec(
				"INSERT INTO projects (id, name, description) VALUES (?, ?, ?)",
				p.ID, p.Name, p.Description,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
