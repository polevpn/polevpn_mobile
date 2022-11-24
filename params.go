package polevpnmobile

type errorCode struct {
	Code int
	Msg  string
}

type reqConnectAccessServer struct {
	accessServer
}

type respConnectAccessServer struct {
	errorCode
}

type reqStopAccessServer struct {
}

type respStopAccessServer struct {
	errorCode
}

type reqAddAccessServer struct {
	accessServer
}

type respAddAccessServer struct {
	errorCode
}

type reqGetAllAccessServer struct {
}

type respGetAllAccessServer struct {
	Servers []accessServer
	errorCode
}

type reqUpdateAccessServer struct {
	accessServer
}

type respUpdateAccessServer struct {
	errorCode
}

type reqDeleteAccessServer struct {
	ID uint
}

type respDeleteAccessServer struct {
	errorCode
}

type reqGetAllLogs struct {
}

type respGetAllLogs struct {
	errorCode
}

type reqGetUpDownBytes struct {
}

type respGetUpDownBytes struct {
	errorCode
}

type reqGetVersion struct {
}

type respGetVersion struct {
	errorCode
	Version string
}
