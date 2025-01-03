package polevpnmobile

import (
	"runtime/debug"
	"strings"
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
	OnAllocEvent(ip string, dns string, routes string)
	OnReconnectingEvent()
	OnReconnectedEvent()
}

type PoleVPNLogHandler interface {
	OnWrite(data string)
}

type PoleVPN struct {
	handler PoleVPNEventHandler
	client  core.PoleVpnClient
	mutex   *sync.Mutex
	state   int
	ip      string
	dns     string
	routes  string
}

type LogCallback struct {
	logHandler PoleVPNLogHandler
}

func (lc *LogCallback) Write(data []byte) (int, error) {
	if lc.logHandler != nil {
		lc.logHandler.OnWrite(string(data))
	}
	return len(data), nil
}

func (lc *LogCallback) SetLogHandler(handler PoleVPNLogHandler) {
	lc.logHandler = handler
}

var logCallback LogCallback

func init() {

	plog = elog.GetLogger()
	plog.SetCallback(&logCallback)
	core.SetLogger(plog)
	defer plog.Flush()
}

func SetLogPath(path string) {
	plog.SetLogPath(path)
}

func SetLogLevel(level string) {
	plog.SetLogLevel(level)
}

func SetLogHandler(handler PoleVPNLogHandler) {
	logCallback.SetLogHandler(handler)
}

func NewPoleVPN() *PoleVPN {
	debug.SetMemoryLimit(1024 * 1024 * 10)
	return &PoleVPN{mutex: &sync.Mutex{}, state: POLEVPN_MOBILE_INIT}
}

func (pvm *PoleVPN) eventHandler(event int, client core.PoleVpnClient, av *anyvalue.AnyValue) {

	switch event {
	case core.CLIENT_EVENT_ADDRESS_ALLOCED:
		{
			if pvm.handler != nil {
				pvm.ip = av.Get("ip").AsStr()
				pvm.dns = av.Get("dns").AsStr()
				routes, _ := av.Get("route").EncodeJson()
				pvm.routes = string(routes)

				pvm.handler.OnAllocEvent(av.Get("ip").AsStr(), av.Get("dns").AsStr(), pvm.routes)
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

func (pvm *PoleVPN) AttachIos(fd int) {
	if pvm.state != POLEVPN_MOBILE_STARTED {
		return
	}
	tundevice := core.AttachTunDeviceIos(fd)
	pvm.client.AttachTunDevice(tundevice)
}

func (pvm *PoleVPN) Start(endpoint string, user string, pwd string, sni string, skipSSLVerify bool, deviceType string, deviceId string) {

	pvm.mutex.Lock()
	defer pvm.mutex.Unlock()
	if pvm.state != POLEVPN_MOBILE_INIT && pvm.state != POLEVPN_MOBILE_STOPPED {
		return
	}
	var err error
	var client core.PoleVpnClient

	if strings.HasPrefix(endpoint, "proxy://") {
		client, err = core.NewPoleVpnClientProxy()
	} else {
		client, err = core.NewPoleVpnClientVLAN()
	}

	if err != nil {
		if pvm.handler != nil {
			pvm.handler.OnErrorEvent("start", err.Error())
		}
		return
	}

	pvm.client = client
	pvm.state = POLEVPN_MOBILE_STARTING
	pvm.client.SetEventHandler(pvm.eventHandler)
	go pvm.client.Start(endpoint, user, pwd, sni, skipSSLVerify, deviceType, deviceId)
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

func (pvm *PoleVPN) GetRoutes() string {
	return pvm.routes
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
