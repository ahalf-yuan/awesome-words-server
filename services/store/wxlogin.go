package store

type DatabaseQueryReq struct {
	Env   string `json:"env"`
	Query string `json:"query"`
}

type WxcodeReq struct {
	Page  string `json:"page"`
	Width int64  `json:"width"`
	Scene string `json:"scene"`
}
