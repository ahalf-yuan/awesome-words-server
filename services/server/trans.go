package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"wordshub/services/common/errno"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

// https://aidemo.youdao.com/trans
func youdao(ctx *gin.Context) {
	youdaoReq := ctx.MustGet(gin.BindKey).(*store.YoudaoReq)

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

	var jsonObj store.YoudaoResp
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}
