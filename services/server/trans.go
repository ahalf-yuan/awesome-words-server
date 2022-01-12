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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonObj models.YoudaoResp
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}

// 如果这个接口调用失败，就用 sug 接口替代
// key 不一定稳定
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置参数
	params := req.URL.Query()
	params.Add("key", key)
	params.Add("type", reqType)
	params.Add("w", w)

	req.URL.RawQuery = params.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonObj models.IcibaResp
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}

func icibaSentence(ctx *gin.Context) {
	const host = "http://open.iciba.com/dsapi/"

	resp, err := http.Get(host)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonObj models.IcibaSentence
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}

// 金山 sug API
// https://dict-mobile.iciba.com/interface/index.php?c=word&nums=1&is_need_mean=1&word=test&m=getsuggest
func icibaSug(ctx *gin.Context) {
	const host = "https://dict-mobile.iciba.com/interface/index.php"
	queries := ctx.Request.URL.Query()

	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.URL.RawQuery = queries.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonObj models.IcibaSugResp
	json.Unmarshal(responseData, &jsonObj)

	ctx.JSON(http.StatusOK, errno.OK.WithData(jsonObj))
}
