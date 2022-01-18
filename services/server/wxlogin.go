package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wordshub/services/common/errno"
	"wordshub/services/common/wxutils"
	"wordshub/services/conf"
	"wordshub/services/store"

	"github.com/gin-gonic/gin"
)

func weappText(ctx *gin.Context, conf conf.Config) {
	ctx.JSON(http.StatusOK, errno.OK.WithData(gin.H{"name": conf.AppId, "secret": conf.AppSecret}))
}

// query user info from weapp database
func queryUserInfo(ctx *gin.Context, conf conf.Config) {
	// 1.get uuid from query
	// 2.get user info in weapp-db by uuid
	query := ctx.Request.URL.Query()
	uuid := query.Get("uuid")

	accessTokenReq := wxutils.AccessTokenReq{
		AppId:     conf.AppId,
		AppSecret: conf.AppSecret,
	}
	accessToken, err := wxutils.GetAccessToken(&accessTokenReq)
	if err != nil {
		// handle error
		// something wrong with accessToken
		return
	}

	databaseQueryReq := store.DatabaseQueryReq{
		Env:   conf.AppCloudEnv,
		Query: "db.collection(\"uuids\").where({uuid:" + uuid + "}).get()",
	}

	body, _ := json.Marshal(databaseQueryReq)

	const weappHost = "https://api.weixin.qq.com/tcb/databasequery"
	url := weappHost + "?access_token=" + accessToken
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//Failed to read response.
			panic(err)
		}

		//Convert bytes to String and print
		jsonStr := string(body)
		// TODO: do login and set jwt to PC
		fmt.Println("Response: ", jsonStr)

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp.Status)
	}
}

// 生成小程序码
func getWxCode(ctx *gin.Context, conf conf.Config) {
	const wxPagePath = "pages/scan_login/index"

	accessTokenReq := wxutils.AccessTokenReq{
		AppId:     conf.AppId,
		AppSecret: conf.AppSecret,
	}
	//
	accessToken, err := wxutils.GetAccessToken(&accessTokenReq)
	if err != nil {
		// handle error
		// something wrong with accessToken
		return
	}

	query := ctx.Request.URL.Query()
	uuid := query.Get("uuid")
	scene := uuid
	databasequeryReq := store.WxcodeReq{
		Page:  wxPagePath,
		Width: 230,
		Scene: scene,
	}

	body, _ := json.Marshal(databasequeryReq)

	const wxCodeUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
	url := wxCodeUrl + "?access_token=" + accessToken
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		// handle error
	}
	// mimeType
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	ctx.Data(200, "image/png", pix)
}
