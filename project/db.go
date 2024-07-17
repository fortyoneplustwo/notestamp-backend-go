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
func (p *ProjectDB) Add(uid int, proj Metadata) error {
	stmt, err := p.db.Prepare(`
    INSERT INTO projects (id, format, title, label, src, mime_type, media_hash, user_id) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  `)
	if err != nil {
		return err
	}

	pid := strconv.Itoa(uid) + "/" + proj.Title

	_, err = stmt.Exec(pid, proj.Format, proj.Title, proj.Label, proj.Src, proj.Mimetype, proj.MediaHash, uid)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectDB) Get(uid int, title string) (Metadata, error) {
	pid := strconv.Itoa(uid) + "/" + title

	project := Metadata{}
	row := p.db.QueryRow("SELECT title, label, format, src, mime_type, media_hash FROM projects WHERE id = $1", pid)
	if err := row.Scan(
		&project.Title,
		&project.Label,
		&project.Format,
		&project.Src,
		&project.Mimetype,
    &project.MediaHash,
	); err != nil {
		return Metadata{}, err
	}

	return project, nil
}

func(p *ProjectDB) FindMediaDup(uid int, title string) (bool, error) {
	pid := strconv.Itoa(uid) + "/" + title
  row := p.db.QueryRow("SELECT media_hash FROM projects WHERE pid = $1", pid)
  var mediaHash string
    if err := row.Scan(&mediaHash); err != nil {
    return false, err
  }
  
  return true, nil
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

func (p *ProjectDB) Remove(uid int, title string) (Metadata, error) {
	var project Metadata
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
