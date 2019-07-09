package wx

import (
	"dog/util/log"
	"encoding/json"
	"fmt"
	"net/http"
)

// 获取token(发信息的)->生产环境建议将token存入缓存中
type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

func (c *Client) getToken() (*token, error) {
	res, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", c.appID, c.appSecret))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	t := new(token)
	err = json.NewDecoder(res.Body).Decode(t)
	if t.Errcode != 0 {
		return nil, fmt.Errorf("errcode:%d,errmsg:%s", t.Errcode, t.Errmsg)
	}
	log.Warning.Println(t, err)
	return t, err
}

func (c *Client) getTokenText() (string, error) {
	if c.tokenCacheFunc == nil {
		t, err := c.getToken()
		if err != nil {
			return "", err
		}
		return t.AccessToken, nil
	}
	if t, err := c.tokenCacheFunc("get", c.appID+c.appSecret); err == nil {
		return t, nil
	}
	t, err := c.getToken()
	if err != nil {
		return "", err
	}
	if _, err := c.tokenCacheFunc("set", c.appID+c.appSecret, t.AccessToken, t.ExpiresIn); err != nil {
		logger.Printf("warning in token cache set:%s", err.Error())
	}
	return t.AccessToken, nil
}
