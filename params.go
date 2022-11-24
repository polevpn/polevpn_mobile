package polevpnmobile

type errorCode struct {
	Code int
	Msg  string
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
