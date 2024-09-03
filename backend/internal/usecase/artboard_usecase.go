package usecase

import (
	"goP2Pbackend/internal/domain"

	"github.com/google/uuid"
)

type artboardUsecase struct {
	artboardRepo    domain.ArtboardRepository
	artboardStorage domain.ArtboardStorage
}

func NewArtboardUsecase(ar domain.ArtboardRepository, as domain.ArtboardStorage) domain.ArtboardUsecase {
	return &artboardUsecase{
		artboardRepo:    ar,
		artboardStorage: as,
	}
}

func (a *artboardUsecase) Create(artboard *domain.Artboard) error {
	artboard.ID = uuid.New().String()
	artboard.ShareableID = uuid.New().String()
	return a.artboardRepo.Create(artboard)
}

func (a *artboardUsecase) GetByID(id string) (*domain.Artboard, error) {
	return a.artboardRepo.GetByID(id)
}

func (a *artboardUsecase) GetByOwnerID(ownerID string) ([]*domain.Artboard, error) {
	return a.artboardRepo.GetByOwnerID(ownerID)
}

func (a *artboardUsecase) Update(artboard *domain.Artboard) error {
	return a.artboardRepo.Update(artboard)
}

func (a *artboardUsecase) Delete(id string) error {
	return a.artboardRepo.Delete(id)
}

func (a *artboardUsecase) GenerateShareableLink(artboardID string, isReadOnly bool) (string, error) {
	artboard, err := a.artboardRepo.GetByID(artboardID)
	if err != nil {
		return "", err
	}

	artboard.IsReadOnly = isReadOnly
	artboard.ShareableID = uuid.New().String()

	err = a.artboardRepo.Update(artboard)
	if err != nil {
		return "", err
	}

	return artboard.ShareableID, nil
}

func (a *artboardUsecase) SaveArtboardData(artboardID string, data []byte) error {
	return a.artboardStorage.Save(artboardID, data)
}

func (a *artboardUsecase) LoadArtboardData(artboardID string) ([]byte, error) {
	return a.artboardStorage.Load(artboardID)
}
