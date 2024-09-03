package postgres

import (
	"database/sql"
	"goP2Pbackend/internal/domain"
)

type artboardRepository struct {
	db *sql.DB
}

func NewArtboardRepository(db *sql.DB) domain.ArtboardRepository {
	return &artboardRepository{db: db}
}

func (r *artboardRepository) Create(artboard *domain.Artboard) error {
	query := `INSERT INTO artboards (id, name, owner_id, created_at, updated_at, shareable_id, is_read_only) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, artboard.ID, artboard.Name, artboard.OwnerID, artboard.CreatedAt, artboard.UpdatedAt, artboard.ShareableID, artboard.IsReadOnly)
	return err
}

func (r *artboardRepository) GetByID(id string) (*domain.Artboard, error) {
	query := `SELECT id, name, owner_id, created_at, updated_at, shareable_id, is_read_only FROM artboards WHERE id = $1`
	var artboard domain.Artboard
	err := r.db.QueryRow(query, id).Scan(&artboard.ID, &artboard.Name, &artboard.OwnerID, &artboard.CreatedAt, &artboard.UpdatedAt, &artboard.ShareableID, &artboard.IsReadOnly)
	if err != nil {
		return nil, err
	}
	return &artboard, nil
}

func (r *artboardRepository) GetByOwnerID(ownerID string) ([]*domain.Artboard, error) {
	query := `SELECT id, name, owner_id, created_at, updated_at, shareable_id, is_read_only FROM artboards WHERE owner_id = $1`
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artboards []*domain.Artboard
	for rows.Next() {
		var artboard domain.Artboard
		err := rows.Scan(&artboard.ID, &artboard.Name, &artboard.OwnerID, &artboard.CreatedAt, &artboard.UpdatedAt, &artboard.ShareableID, &artboard.IsReadOnly)
		if err != nil {
			return nil, err
		}
		artboards = append(artboards, &artboard)
	}
	return artboards, nil
}

func (r *artboardRepository) Update(artboard *domain.Artboard) error {
	query := `UPDATE artboards SET name = $2, updated_at = $3, shareable_id = $4, is_read_only = $5 WHERE id = $1`
	_, err := r.db.Exec(query, artboard.ID, artboard.Name, artboard.UpdatedAt, artboard.ShareableID, artboard.IsReadOnly)
	return err
}

func (r *artboardRepository) Delete(id string) error {
	query := `DELETE FROM artboards WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
