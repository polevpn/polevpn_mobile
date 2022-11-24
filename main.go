package polevpnmobile

import (
	"sync"

	"github.com/polevpn/anyvalue"
	"github.com/polevpn/elog"
	core "github.com/polevpn/polevpn_core"
)

const (
	POLEVPN_MOBILE_INIT     = 0
	POLEVPN_MOBILE_STARTED  = 1
	POLEVPN_MOBILE_STOPPED  = 2
	POLEVPN_MOBILE_STARTING = 3
	POLEVPN_MOBILE_STOPPING = 4
)

var plog *elog.EasyLogger

type PoleVPNEventHandler interface {
	OnStartedEvent()
	OnStoppedEvent()
	OnErrorEvent(errtype string, errmsg string)
	OnAllocEvent(ip string, dns string)
	OnReconnectingEvent()
	OnReconnectedEvent()
}

type PoleVPNLogHandler interface {
	OnWrite(data string)
	OnFlush()
}

type PoleVPN struct {
	handler PoleVPNEventHandler
	client  *core.PoleVpnClient
	mutex   *sync.Mutex
	state   int
	ip      string
}

func init() {

	plog = elog.GetLogger()
	core.SetLogger(plog)
	defer plog.Flush()
}

func SetLogPath(path string) {
	plog.SetLogPath(path)
}

func SetLogLevel(level string) {
	plog.SetLogLevel(level)
}

func NewPoleVPN() *PoleVPN {
	return &PoleVPN{mutex: &sync.Mutex{}, state: POLEVPN_MOBILE_INIT}
}

func (pvm *PoleVPN) eventHandler(event int, client *core.PoleVpnClient, av *anyvalue.AnyValue) {

	switch event {
	case core.CLIENT_EVENT_ADDRESS_ALLOCED:
		{
			if pvm.handler != nil {
				pvm.ip = av.Get("ip").AsStr()

				pvm.handler.OnAllocEvent(av.Get("ip").AsStr(), av.Get("dns").AsStr())
			}
		}
	case core.CLIENT_EVENT_STOPPED:
		{
			pvm.state = POLEVPN_MOBILE_STOPPED
			if pvm.handler != nil {
				pvm.handler.OnStoppedEvent()
			}
		}
	case core.CLIENT_EVENT_RECONNECTED:
		if pvm.handler != nil {
			pvm.handler.OnReconnectedEvent()
		}
	case core.CLIENT_EVENT_RECONNECTING:
		if pvm.handler != nil {
			pvm.handler.OnReconnectingEvent()
		}
	case core.CLIENT_EVENT_STARTED:
		pvm.state = POLEVPN_MOBILE_STARTED
		if pvm.handler != nil {
			pvm.handler.OnStartedEvent()
		}
	case core.CLIENT_EVENT_ERROR:
		if pvm.handler != nil {
			pvm.handler.OnErrorEvent(av.Get("type").AsStr(), av.Get("error").AsStr())
		}
	default:
		plog.Error("invalid evnet=", event)
	}

}

func (pvm *PoleVPN) Attach(fd int) {
	if pvm.state != POLEVPN_MOBILE_STARTED {
		return
	}
	tundevice := core.AttachTunDevice(fd)
	pvm.client.AttachTunDevice(tundevice)
}

func (pvm *PoleVPN) Start(endpoint string, user string, pwd string, sni string) {

	pvm.mutex.Lock()
	defer pvm.mutex.Unlock()
	if pvm.state != POLEVPN_MOBILE_INIT && pvm.state != POLEVPN_MOBILE_STOPPED {
		return
	}

	client, err := core.NewPoleVpnClient()
	if err != nil {
		if pvm.handler != nil {
			pvm.handler.OnErrorEvent("start", err.Error())
		}
		return
	}

	pvm.client = client
	pvm.state = POLEVPN_MOBILE_STARTING
	pvm.client.SetEventHandler(pvm.eventHandler)
	go pvm.client.Start(endpoint, user, pwd, sni, true)
}

func (pvm *PoleVPN) Stop() {
	pvm.mutex.Lock()
	defer pvm.mutex.Unlock()
	if pvm.state == POLEVPN_MOBILE_STARTED || pvm.state == POLEVPN_MOBILE_STARTING {
		pvm.state = POLEVPN_MOBILE_STOPPING
		go pvm.client.Stop()
	}

}

func (pvm *PoleVPN) GetUpBytes() int64 {

	if pvm.client == nil {
		return 0
	}

	up, _ := pvm.client.GetUpDownBytes()
	return int64(up)
}

func (pvm *PoleVPN) GetDownBytes() int64 {

	if pvm.client == nil {
		return 0
	}

	_, down := pvm.client.GetUpDownBytes()
	return int64(down)
}

func (pvm *PoleVPN) GetRemoteIP() string {
	return pvm.client.GetRemoteIP()
}

func (pvm *PoleVPN) GetLocalIP() string {
	return pvm.ip
}

func (pvm *PoleVPN) GetState() int {
	return pvm.state
}

func (pvm *PoleVPN) CloseConnect(flag bool) {

	if pvm.state == POLEVPN_MOBILE_STARTED {
		pvm.client.CloseConnect(flag)
	}

}

func (pvm *PoleVPN) SetEventHandler(handler PoleVPNEventHandler) {
	pvm.handler = handler
}
