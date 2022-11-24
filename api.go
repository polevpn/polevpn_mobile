package polevpnmobile

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func AddAccessServer(reqs string) string {

	var req reqAddAccessServer
	var err error

	err = json.Unmarshal([]byte(reqs), &req)

	resp := &respAddAccessServer{errorCode: errorCode{Code: 0, Msg: "ok"}}

	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}

	err = addAccessServer(req.accessServer)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

func UpdateAccessServer(reqs string) string {

	var req reqUpdateAccessServer
	var err error

	err = json.Unmarshal([]byte(reqs), &req)

	resp := &respUpdateAccessServer{errorCode: errorCode{Code: 0, Msg: "ok"}}

	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}

	err = updateAccessServer(req.accessServer)

	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

func DeleteAccessServer(reqs string) string {

	var req reqDeleteAccessServer
	var err error
	err = json.Unmarshal([]byte(reqs), &req)

	resp := &respDeleteAccessServer{errorCode: errorCode{Code: 0, Msg: "ok"}}

	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}

	err = deleteAccessServer(req.ID)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

func GetAllAccessServer(reqs string) string {

	var req reqGetAllAccessServer
	var err error

	err = json.Unmarshal([]byte(reqs), &req)

	resp := &respGetAllAccessServer{errorCode: errorCode{Code: 0, Msg: "ok"}}

	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}

	servers, err := getAllAccessServer()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		data, _ := json.Marshal(resp)
		return string(data)
	}
	resp.Servers = servers
	data, _ := json.Marshal(resp)
	return string(data)
}

func GetAllLogs() string {

	logFilePath := plog.GetLogPath() + string(os.PathSeparator) + getAppName() + "-" + getTimeNowDate() + ".log"

	logData, err := ioutil.ReadFile(logFilePath)

	if err != nil {
		plog.Error(err)
		return ""
	}

	return string(logData)
}
