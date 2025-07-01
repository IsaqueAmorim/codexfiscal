package service

import (
	"encoding/json"
	"fmt"

	"github.com/IsaqueAmorim/codexfiscal/internal/model"
	"github.com/IsaqueAmorim/codexfiscal/internal/repository"
	"github.com/IsaqueAmorim/codexfiscal/pkg/utils"
)

type NCMService interface {
	CreateNCM(ncm *model.NCM) error
	UpdateNCM(ncm *model.NCM) error
	DeleteNCM(id string) error
	GetNCMByCode(code string) (*model.NCM, error)
	GetNCMByID(id string) (*model.NCM, error)
	GetNCMByText(text string) (*model.NCM, error)
	GetNCMSByCodes(codes []string) ([]*model.NCM, error)
	GetNCMSByText(text string) ([]*model.NCM, error)
	BulkInsertNCMs(ncms []*model.NCM) error
}

func NewNCMService(repo repository.NCMRepository) NCMService {
	return &ncmService{
		repo: repo,
	}
}

type ncmService struct {
	repo repository.NCMRepository
}

func (s *ncmService) CreateNCM(ncm *model.NCM) error {
	if ncm == nil {
		return fmt.Errorf("ncm cannot be nil")
	}

	if utils.IsNullOrWhiteSpace(ncm.Code) {
		return fmt.Errorf("code cannot be null or empty")
	}

	if utils.IsNullOrWhiteSpace(ncm.Description) {
		return fmt.Errorf("description cannot be null or empty")
	}

	ncm.CodeNoSimbols = utils.RemoveSymbols(ncm.Code)

	return s.repo.CreateNCM(ncm)
}

func (s *ncmService) UpdateNCM(ncm *model.NCM) error {
	if ncm == nil {
		return fmt.Errorf("ncm cannot be nil")
	}

	if utils.IsNullOrWhiteSpace(ncm.ID) {
		return fmt.Errorf("id cannot be null or empty")
	}

	if utils.IsNullOrWhiteSpace(ncm.Code) {
		return fmt.Errorf("code cannot be null or empty")
	}

	if utils.IsNullOrWhiteSpace(ncm.Description) {
		return fmt.Errorf("description cannot be null or empty")
	}

	ncm.CodeNoSimbols = utils.RemoveSymbols(ncm.Code)

	return s.repo.UpdateNCM(ncm)
}

func (s *ncmService) DeleteNCM(id string) error {
	if utils.IsNullOrWhiteSpace(id) {
		return fmt.Errorf("id cannot be null or empty")
	}

	err := s.repo.DeleteNCM(id)
	if err != nil {
		return fmt.Errorf("error deleting ncm with id %s: %w", id, err)
	}

	return nil
}

func (s *ncmService) GetNCMByCode(code string) (*model.NCM, error) {
	if utils.IsNullOrWhiteSpace(code) {
		return nil, fmt.Errorf("code cannot be null or empty")
	}

	ncm, err := s.repo.GetNCMByCode(utils.RemoveSymbols(code))

	if err != nil {
		return nil, fmt.Errorf("ncm with code %s not found", code)
	}

	return ncm, nil
}

func (s *ncmService) GetNCMByID(id string) (*model.NCM, error) {
	if utils.IsNullOrWhiteSpace(id) {
		return nil, fmt.Errorf("id cannot be null or empty")
	}

	ncm, err := s.repo.GetNCMByID(id)

	if err != nil {
		return nil, fmt.Errorf("ncm with id %s not found", id)
	}

	return ncm, nil
}

func (s *ncmService) GetNCMByText(text string) (*model.NCM, error) {
	if utils.IsNullOrWhiteSpace(text) {
		return nil, fmt.Errorf("text cannot be null or empty")
	}

	ncm, err := s.repo.GetNCMByText(text)

	if err != nil {
		return nil, fmt.Errorf("ncm with text %s not found", text)
	}

	return ncm, nil
}

func (s *ncmService) GetNCMSByCodes(codes []string) ([]*model.NCM, error) {
	if len(codes) == 0 {
		return nil, fmt.Errorf("codes cannot be empty")
	}

	ncms, err := s.repo.GetNCMSByCodes(codes)
	if err != nil {
		return nil, fmt.Errorf("ncm with codes %v not found", codes)
	}

	return ncms, nil
}

func (s *ncmService) GetNCMSByText(text string) ([]*model.NCM, error) {
	if utils.IsNullOrWhiteSpace(text) {
		return nil, fmt.Errorf("text cannot be null or empty")
	}

	ncms, err := s.repo.GetNCMSByText(text)
	if err != nil {
		return nil, fmt.Errorf("ncm with text %s not found", text)
	}

	return ncms, nil
}

func (s *ncmService) BulkInsertNCMs(ncms []*model.NCM) error {
	if len(ncms) == 0 {
		return fmt.Errorf("ncms cannot be empty")
	}

	s.repo.BulkInsertNCMs(ncms)

	if err := s.repo.BulkInsertNCMs(ncms); err != nil {
		return fmt.Errorf("error bulk inserting ncms: %w", err)
	}

	return nil
}

func (s *ncmService) ImportFromJson(jsonContent string) ([]*model.NCM, error) {
	if utils.IsNullOrWhiteSpace(jsonContent) {
		return nil, fmt.Errorf("json cannot be null or empty")
	}

	var ncms []*model.NCM
	err := json.Unmarshal([]byte(jsonContent), &ncms)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}

	//TODO: Lógica para atualizar ou inserir os NCMs no repositório
	for _, ncm := range ncms {
		if err := s.CreateNCM(ncm); err != nil {
			return nil, fmt.Errorf("error creating ncm from json: %w", err)
		}
	}
	return ncms, nil
}
