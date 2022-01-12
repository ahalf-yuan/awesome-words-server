package models

type IcibaReq struct {
	Word string `json:"word" from:"word"`
	Type string `json:"type"`
	Key  string `json:"key"`
}

type Exchange struct {
	WordPl    []string `json:"word_pl"`
	WordPast  []string `json:"word_past"`
	WordDone  []string `json:"word_done"`
	WordIng   []string `json:"word_ing"`
	WordThird []string `json:"word_third"`
	WordEr    []string `json:"word_er"`
	WordEst   []string `json:"word_est"`
}

type Part struct {
	Part  string   `json:"part"`
	Means []string `json:"means"`
}

type Symbol struct {
	PhEn    string `json:"ph_en"`
	PhAm    string `json:"ph_am"`
	PhOther string `json:"ph_other"`
	PhEnMp3 string `json:"ph_en_mp3"`
	PhAmMp3 string `json:"ph_am_mp3"`
	Parts   []Part `json:"parts"`
}

type IcibaResp struct {
	WordName string   `json:"word_name"`
	IsCRI    int      `json:"is_CRI"`
	Exchange Exchange `json:"exchange"`
	Symbols  []Symbol `json:"symbols"`
	Items    []string `json:"items"`
}
