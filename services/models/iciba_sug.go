package models

// c=word&m=getsuggest&nums=10&is_need_mean=1&word=test
type IcibaSugReq struct {
	Catalog    string `json:"catalog" from:"c"`
	Method     string `json:"method" from:"m"`
	Word       string `json:"word" from:"word"`
	Nums       string `json:"nums"`
	IsNeedMean string `json:"is_need_mean"`
}

type Mean struct {
	Part  string   `json:"part"`
	Means []string `json:"means"`
}

type SugItem struct {
	Key        string `json:"key"`
	Paraphrase string `json:"paraphrase"`
	Value      int    `json:"value"`
	Means      []Mean `json:"means"`
}

type IcibaSugResp struct {
	Message []SugItem `json:"message"`
	Status  int       `json:"status"`
}
