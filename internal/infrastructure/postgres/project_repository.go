package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/tomtom2k/kairo-anchor-server/internal/domain/project"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *project.Project) error {
	tasksJSON, err := json.Marshal(p.Tasks)
	if err != nil {
		return err
	}

	documentsJSON, err := json.Marshal(p.Documents)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO projects (user_id, name, description, status, progress, start_date, end_date, tasks, documents, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		p.UserID, p.Name, p.Description, p.Status, p.Progress,
		p.StartDate, p.EndDate, tasksJSON, documentsJSON,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProjectRepository) Update(ctx context.Context, p *project.Project) error {
	tasksJSON, err := json.Marshal(p.Tasks)
	if err != nil {
		return err
	}

	documentsJSON, err := json.Marshal(p.Documents)
	if err != nil {
		return err
	}

	query := `
		UPDATE projects
		SET name = $1, description = $2, status = $3, progress = $4,
		    start_date = $5, end_date = $6, tasks = $7, documents = $8, updated_at = NOW()
		WHERE id = $9 AND user_id = $10
		RETURNING updated_at
	`

	result := r.db.QueryRowContext(ctx, query,
		p.Name, p.Description, p.Status, p.Progress,
		p.StartDate, p.EndDate, tasksJSON, documentsJSON,
		p.ID, p.UserID,
	)

	err = result.Scan(&p.UpdatedAt)
	if err == sql.ErrNoRows {
		return errors.New("project not found or access denied")
	}
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("project not found or access denied")
	}

	return nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*project.Project, error) {
	var p project.Project
	var tasksJSON, documentsJSON []byte

	query := `
		SELECT id, user_id, name, description, status, progress,
		       start_date, end_date, tasks, documents, created_at, updated_at
		FROM projects
		WHERE id = $1 AND user_id = $2
	`

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Description, &p.Status, &p.Progress,
		&p.StartDate, &p.EndDate, &tasksJSON, &documentsJSON,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tasksJSON, &p.Tasks); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(documentsJSON, &p.Documents); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]project.Project, error) {
	query := `
		SELECT id, user_id, name, description, status, progress,
		       start_date, end_date, tasks, documents, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []project.Project
	for rows.Next() {
		var p project.Project
		var tasksJSON, documentsJSON []byte

		err := rows.Scan(
			&p.ID, &p.UserID, &p.Name, &p.Description, &p.Status, &p.Progress,
			&p.StartDate, &p.EndDate, &tasksJSON, &documentsJSON,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tasksJSON, &p.Tasks); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(documentsJSON, &p.Documents); err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
