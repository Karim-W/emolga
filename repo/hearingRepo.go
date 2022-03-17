package repo

import (
	"encoding/json"
	"sync"

	"github.com/karim-w/emolga/clients"
	"github.com/karim-w/emolga/common"
	"github.com/karim-w/emolga/helpers/redishelper"
	"github.com/karim-w/emolga/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HearingRepo struct {
	logger  *zap.SugaredLogger
	client  *redishelper.RedisManager
	pndl    *clients.PineduleClient
	manager *common.UserManager
}

func HearingRepoProvider(logger *zap.SugaredLogger, client *redishelper.RedisManager, p *clients.PineduleClient) *HearingRepo {
	return &HearingRepo{
		logger: logger,
		client: client,
		pndl:   p,
	}
}

func (h *HearingRepo) GetUsersInHearing(Hearing string, tid string) (*[]string, error) {
	return h.manager.GetUsersInSet(Hearing, tid)
}

func (h *HearingRepo) ExpandUserDetails(usersList *[]string) (*[]models.RedisUserEntry, error) {
	return h.manager.ExpandUserDetails(usersList)
}

func (h *HearingRepo) MapUserByState(users *[]models.RedisUserEntry) *map[string][]models.RedisUserEntry {
	return h.manager.MapUserByState(users)
}

func (h *HearingRepo) AddPstnUser(pUser models.PstnUser, tid string) error {
	wg := sync.WaitGroup{}
	wg.Add(len(pUser.HearingIds) + 1)
	go func() {
		defer wg.Done()
		h.addPstnUserOnRedis(pUser, tid)
	}()
	for _, hearingId := range pUser.HearingIds {
		go func(hearingId string) {
			defer wg.Done()
			h.addPstnUserToHearingSet("Pstn-User-"+pUser.Email, hearingId)
			h.addPstnUserToSessionSet("Pstn-User-"+pUser.Email, hearingId, tid)
		}(hearingId)
	}
	wg.Wait()
	return nil
}

//\\//\\//\\//\\//\\//\\//\\//\\::::: PRIVATE ::::://\\//\\//\\//\\//\\//\\//\\//\\
func (h *HearingRepo) addPstnUserOnRedis(pUser models.PstnUser, tid string) {
	if stringfiedText, err := json.Marshal(pUser); err == nil {
		h.client.AddKeyValuePair("Pstn-User-"+pUser.Email, string(stringfiedText))
	} else {
		h.logger.Errorw("AddPSTNUser", "err:", err)
	}
}
func (h *HearingRepo) addPstnUserToSessionSet(user string, hearingId string, tid string) {
	if sid, err := h.pndl.FetchConfrenceIdFromHearingId(hearingId, tid); err == nil {
		h.client.AddToSet(sid, user)
	} else {
		h.logger.Errorw("AddPSTNUser", "err:", err)
	}
}

func (h *HearingRepo) addPstnUserToHearingSet(key string, hearingId string) {
	h.client.AddToSet(hearingId, key)
}

var HearingRepoModule = fx.Option(fx.Provide(HearingRepoProvider))
