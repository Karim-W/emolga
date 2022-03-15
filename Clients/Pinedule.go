package clients

import (
	"github.com/karim-w/emolga/utils/hermes"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PineduleClient struct {
	logger *zap.SugaredLogger
	hermes *hermes.HttpClient
}

func PineduleClientProvider(log *zap.SugaredLogger, hermes *hermes.HttpClient) *PineduleClient {
	return &PineduleClient{
		logger: log,
		hermes: hermes,
	}
}

var PineduleClientModule = fx.Option(fx.Provide(PineduleClientProvider))
