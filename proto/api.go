package proto

type CommonResp struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

type ApplyMqReq struct {
	Ips []string `json:"ips"`
}

type ApplyMqResp struct {
	SendPort int `json:"send_port"`
	RevPort  int `json:"rev_port"`
}
