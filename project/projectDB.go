package project

import (
	"database/sql"
	"strconv"
)

type ProjectDB struct {
	db *sql.DB
}

func NewProjectDB(db *sql.DB) *ProjectDB {
	return &ProjectDB{db: db}
}

// Implement ProjectStore interface
func (p *ProjectDB) Add(uid int, proj Project) error {
	stmt, err := p.db.Prepare(`
    INSERT INTO projects (id, format, title, label, src, mime_type, user_id) 
    VALUES ($1, $2, $3, $4, $5, $6, $7)
  `)
	if err != nil {
		return err
	}

	pid := strconv.Itoa(uid) + "/" + proj.Title

	_, err = stmt.Exec(pid, proj.Format, proj.Title, proj.Label, proj.Src, proj.Mimetype, uid)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectDB) Get(uid int, title string) (Project, error) {
	pid := strconv.Itoa(uid) + "/" + title

	project := Project{}
	row := p.db.QueryRow("SELECT title, label, format, src, mime_type FROM projects WHERE id = $1", pid)
	if err := row.Scan(
		&project.Title,
		&project.Label,
		&project.Format,
		&project.Src,
		&project.Mimetype,
	); err != nil {
		return Project{}, err
	}

	return project, nil
}

func (p *ProjectDB) List(uid int) ([]string, error) {
	var projects []string

	stmt, err := p.db.Prepare("SELECT title FROM projects WHERE user_id = $1")
	if err != nil {
		return projects, err
	}

	rows, err := stmt.Query(uid)
	if err != nil {
		return projects, err
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return projects, err
		}
		projects = append(projects, title)
	}

	return projects, nil
}

func (p *ProjectDB) Remove(uid int, title string) (Project, error) {
	var project Project
	project, err := p.Get(uid, title)
	if err != nil {
		return project, err
	}

	stmt, err := p.db.Prepare("DELETE FROM projects WHERE id = $1")
	if err != nil {
		return project, err
	}

	pid := strconv.Itoa(uid) + "/" + title

	_, err = stmt.Exec(pid)
	if err != nil {
		return project, err
	}

	return project, nil
}
