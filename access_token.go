package dingtalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fromiuan/dingtalk/lib"
)

var (
	defaultExpires = 7200
)

//ResAccessToken struct
type ResAccessToken struct {
	AccessToken string `json:"access_token"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

//GetAccessToken 获取access_token
func (c *Client) GetAccessToken() (accessToken string, err error) {
	c.Tlock.Lock()
	defer c.Tlock.Unlock()

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", c.AppID)
	val := c.Cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	var resAccessToken ResAccessToken
	resAccessToken, err = c.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从服务器获取token
func (c *Client) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?appkey=%s&appsecret=%s", gettoken, c.AppKey, c.AppSecret)
	var body []byte
	body, err = lib.Get(url).AsBytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		return errors.New(resAccessToken.ErrMsg)
	}

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", c.AppID)
	err = c.Cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(defaultExpires)*time.Second)
	return
}
