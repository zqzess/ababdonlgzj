package types

type Conf struct {
	Account        string `json:"account,omitempty"`
	Password       string `json:"password,omitempty"`
	Port           int    `json:"port"`
	Welcome        string `json:"welcome"`
	UploadSize     int64  `json:"uploadSize"`
	Tmp            string `json:"tmp,omitempty"`
	AutoRemoveTime int    `json:"autoRemoveTime"`
}

type JSResp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    any    `json:"data,omitempty"`
}

type Response struct {
	Flag string      `json:"flag"`
	Data interface{} `json:"data,omitempty"`
}
