package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"wordshub/services/common/errno"
	"wordshub/services/models"

	"github.com/gin-gonic/gin"
)

// https://aidemo.youdao.com/trans
func youdao(ctx *gin.Context) {
	youdaoReq := ctx.MustGet(gin.BindKey).(*models.YoudaoReq)

	reqMap := map[string][]string{}
	reqMap["q"] = []string{youdaoReq.Word}

	resp, err := http.PostForm("http://aidemo.youdao.com/trans", reqMap)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// log.Fatal(err)
	}

	var jsonObj models.YoudaoResp
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}

func iciba(ctx *gin.Context) {
	const host = "https://dict-co.iciba.com/api/dictionary.php"
	const key = "1F9CA812CB18FFDFC95FC17E9C57A5E1"
	const reqType = "json"

	values := ctx.Request.URL.Query()
	w := values.Get("w")
	if w == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errno.ErrParam)
		return
	}

	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		// log.Fatal(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置参数
	params := req.URL.Query()
	params.Add("key", key)
	params.Add("type", reqType)
	params.Add("w", w)

	req.URL.RawQuery = params.Encode()
	resp, _ := http.DefaultClient.Do(req)
	responseData, _ := ioutil.ReadAll(resp.Body)

	var jsonObj models.IcibaResp
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}
