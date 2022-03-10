package services

import (
	"github.com/karim-w/emolga/models"
	"github.com/karim-w/emolga/modules/scclient"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PodManager struct {
	logger *zap.SugaredLogger
	pods   []models.Pod
}

func (p *PodManager) AddPod(pod models.Pod) {
	client := scclient.New("ws://" + pod.PodIp + "/socketcluster/")
	/*
			{
		            username: "kyle",
		            password: "thebestpassword",
		            client_name: this.client_name,
		          }
	*/
	loginJson := make(map[string]string, 3)
	loginJson["username"] = "kyle"
	loginJson["password"] = "thebestpassword"
	loginJson["client_name"] = "emolga"
	// client.SetBasicListener(p.onConnect, p.onConnectError, p.onDisconnect)
	//client.SetAuthenticationListener(p.onSetAuthentication, p.onAuthentication)
	go func() {
		client.Connect()
		client.Emit("login", loginJson)
		//  func(eventName string, error interface{}, data interface{}) {
		// 	if error == nil {
		// 		// fmt.Println("Got ack for emit event with data ", data, " and error ", error)
		// 		pod.Client = client
		// 		p.pods = append(p.pods, pod)
		// 	}
		// })
	}()
}

func (p *PodManager) RemovePod(podId string) {
	for i, pod := range p.pods {
		if pod.PodName == podId {
			p.pods = append(p.pods[:i], p.pods[i+1:]...)
			return
		}
	}
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\::::: SocketClusterHandlers :::::::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
func (p *PodManager) onConnect(scc scclient.Client) {
	p.logger.Info("Connected to socketcluster server")
}

func (p *PodManager) onDisconnect(scc scclient.Client, err error) {
	p.logger.Info("Disconnected from socketcluster server")
	p.logger.Error(err)
}

func (p *PodManager) onConnectError(scc scclient.Client, err error) {
	p.logger.Error(err)
}

func (p *PodManager) onSetAuthentication(scc scclient.Client, token string) {
	p.logger.Info("Client Received Auth token :	", token)
}

func (p *PodManager) onAuthentication(scc scclient.Client, isAuthenticated bool) {
	p.logger.Info("Client Authentication Status :	", isAuthenticated)
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\::::: DI :::::::://\\//\\//\\//\\//\\//\\//\\//\\//\\//\\
func PodManagerProvider(log *zap.SugaredLogger) *PodManager {
	return &PodManager{
		logger: log,
		pods:   []models.Pod{},
	}
}

var PodManagerModule = fx.Provide(PodManagerProvider)
