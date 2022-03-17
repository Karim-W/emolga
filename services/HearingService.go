package services

import (
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/repo"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HearingService struct {
	logger *zap.SugaredLogger
	repo   *repo.HearingRepo
}

func (h *HearingService) GetUsersInHearing(hearing string, tid string) (*[]string, error) {
	return h.repo.GetUsersInHearing(hearing, tid)
}

func (h *HearingService) GetExpandedUsersInHearing(hId string, tid string) (*[]models.RedisUserEntry, error) {
	if usersList, err := h.repo.GetUsersInHearing(hId, tid); err == nil {
		return h.repo.ExpandUserDetails(usersList)
	} else {
		return nil, err
	}
}

func (h *HearingService) GetUsersMappedByState(hearing string, tid string) (*map[string][]models.RedisUserEntry, error) {
	if usersList, err := h.GetExpandedUsersInHearing(hearing, tid); err == nil {
		return h.repo.MapUserByState(usersList), nil
	} else {
		return nil, err
	}
}

func (h *HearingService) AddPSTNUser(pUser models.PstnUser, tid string) error {

	return h.repo.AddPstnUser(pUser, tid)
}

func HearingServiceProvider(log *zap.SugaredLogger, r *repo.HearingRepo) *HearingService {
	return &HearingService{
		logger: log,
		repo:   r,
	}
}

var HearingServiceModule = fx.Option(fx.Provide(HearingServiceProvider))
