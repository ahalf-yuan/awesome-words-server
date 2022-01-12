package models

type IcibaSentence struct {
	Caption  string `json:"caption"`
	Dateline string `json:"dateline"`
	Content  string `json:"content"`
	Note     string `json:"note"`

	FenxiangImg string `json:"fenxiang_img"`
	Picture     string `json:"picture"`
	Picture2    string `json:"picture2"`
	Picture3    string `json:"picture3"`
	Picture4    string `json:"picture4"`

	Sid  string `json:"sid"`
	Tts  string `json:"tts"`
	Love string `json:"love"`
	SpPv string `json:"sp_pv"`
	SPv  string `json:"s_pv"`

	Tags        []string `json:"tags"`
	Translation string   `json:"translation"`
}
