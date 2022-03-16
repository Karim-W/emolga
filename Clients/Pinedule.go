package clients

import (
	"github.com/karim-w/emolga/Clients/clientmodels"
	"github.com/karim-w/emolga/utils/hermes"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PineduleClient struct {
	baseUrl string
	logger  *zap.SugaredLogger
	hermes  *hermes.HttpClient
}

func (p *PineduleClient) FetchConfrenceIdFromHearingId(hearingId string, TransactionId string) (string, error) {
	p.logger.Infow("FetchConfrenceIdFromHearingId", "hearingId:", hearingId, "TransactionId:", TransactionId)
	url := p.baseUrl + "/api/hearings/" + hearingId + "/session?testIfActive=true"
	cid := clientmodels.ConferenceId{}
	head := make(map[string]string)
	head["TransactionId"] = TransactionId
	if statusCode, ok, err := p.hermes.Get(url, head).Result(&cid); !ok {
		p.logger.Errorw("FetchConfrenceIdFromHearingId", "hearingId:", hearingId, "TransactionId:", TransactionId, "statusCode:", statusCode, "err:", err)
		return "", err
	} else {
		return cid.ConferenceId, nil
	}

}

func PineduleClientProvider(log *zap.SugaredLogger, hermes *hermes.HttpClient) *PineduleClient {
	return &PineduleClient{
		baseUrl: "https://pndl-lasc-cxstg.azurewebsites.net",
		logger:  log,
		hermes:  hermes,
	}
}

var PineduleClientModule = fx.Option(fx.Provide(PineduleClientProvider))
