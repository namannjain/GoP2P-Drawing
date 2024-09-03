package domain

import "time"

type Artboard struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ShareableID string    `json:"shareable_id"`
	IsReadOnly  bool      `json:"is_read_only"`
}

type ArtboardRepository interface {
	Create(artboard *Artboard) error
	GetByID(id string) (*Artboard, error)
	GetByOwnerID(ownerID string) ([]*Artboard, error)
	Update(artboard *Artboard) error
	Delete(id string) error
}

type ArtboardStorage interface {
	Save(artboardID string, data []byte) error
	Load(artboardID string) ([]byte, error)
}

type ArtboardUsecase interface {
	Create(artboard *Artboard) error
	GetByID(id string) (*Artboard, error)
	GetByOwnerID(ownerID string) ([]*Artboard, error)
	Update(artboard *Artboard) error
	Delete(id string) error
	GenerateShareableLink(artboardID string, isReadOnly bool) (string, error)
	SaveArtboardData(artboardID string, data []byte) error
	LoadArtboardData(artboardID string) ([]byte, error)
}
