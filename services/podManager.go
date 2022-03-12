package services

import (
	"github.com/karim-w/emolga/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PodManager struct {
	logger *zap.SugaredLogger
	pods   []models.Pod
}

func (p *PodManager) AddPod(pod models.Pod) {
	p.pods = append(p.pods, pod)
}

func (p *PodManager) RemovePod(podId string) {
	for i, pod := range p.pods {
		if pod.PodName == podId {
			p.pods = append(p.pods[:i], p.pods[i+1:]...)
			return
		}
	}
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\::::: DI :::::::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
func PodManagerProvider(log *zap.SugaredLogger) *PodManager {
	return &PodManager{
		logger: log,
		pods:   []models.Pod{},
	}
}

var PodManagerModule = fx.Provide(PodManagerProvider)
