package wxutils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

// 刷新和获取最新 access_token

type AccessTokenReq struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type AccTokenCache struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	RecordTime  int64  `json:"create_time"`
}

const TOKEN_FILE_NAME = "mp_token_info.json"

// 获取小程序全局 token
func GetAccessToken(accessTokenReq *AccessTokenReq) (string, error) {
	// 优化点：文件缓存太原始了
	fileContent := readFile(TOKEN_FILE_NAME)

	if fileContent == nil {
		// get and write
		// 不存在
		accessToken, err := reqAndWriteAccToken(accessTokenReq)
		return accessToken, err
	}

	jsonData := AccTokenCache{}
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		// handle error: access_token 解析失败
		accessToken, _ := reqAndWriteAccToken(accessTokenReq)
		return accessToken, nil
	}

	accessToken := jsonData.AccessToken
	// compare expires_time (unit second)
	expiresIn := jsonData.ExpiresIn
	recordTime := jsonData.RecordTime
	timeNowStamp := time.Now().Unix()

	if accessToken == "" || timeNowStamp-recordTime > expiresIn-600 {
		// 过期了
		accessToken, _ := reqAndWriteAccToken(accessTokenReq)
		return accessToken, nil
	}

	return accessToken, nil
}

func reqAndWriteAccToken(accessTokenReq *AccessTokenReq) (string, error) {

	jMap, err := requestToken(accessTokenReq.AppId, accessTokenReq.AppSecret)
	accessToken, _ := jMap["access_token"].(string)
	expiresIn, _ := jMap["expires_in"].(int64)
	if err != nil {
		return accessToken, err
	}

	content := AccTokenCache{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
		RecordTime:  time.Now().Unix(),
	}

	// trans struct to []byte
	contentStr := EncodeToBytes(&content)
	writeFile(TOKEN_FILE_NAME, contentStr)

	return accessToken, nil
}

//获取wx_AccessToken 拼接get请求 解析返回json结果 返回 AccessToken和err
func requestToken(appid, secret string) (map[string]interface{}, error) {
	u, err := url.Parse("https://api.weixin.qq.com/cgi-bin/token")
	if err != nil {
		log.Fatal()
	}
	paras := &url.Values{}
	//设置请求参数
	paras.Set("appid", appid)
	paras.Set("secret", secret)
	paras.Set("grant_type", "client_credential")
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("request token err :" + err.Error())
	}

	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New("request token response json parse err :" + err.Error())
	}
	if jMap["errcode"] == nil || jMap["errcode"] == 0 {
		// accessToken, _ := jMap["access_token"].(string)
		return jMap, nil
	} else {
		//返回错误信息
		errcode := jMap["errcode"].(string)
		errmsg := jMap["errmsg"].(string)
		err = errors.New(errcode + ":" + errmsg)
		return nil, err
	}
}

func readFile(filepath string) []byte {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		// log.Fatal(err)
	}
	return content
}

func writeFile(filepath string, message []byte) {
	// 0644 拥有读写权限
	err := ioutil.WriteFile(filepath, message, 0644)
	if err != nil {
		// log.Fatal(err)
		log.Error().Err(err).Msg("Error write file")
	}
}

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		// log.Fatal(err)
		// log error 数据处理失败
	}
	return buf.Bytes()
}
