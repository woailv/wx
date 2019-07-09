// 通过网页授权获取微信用户基本信息
package wx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 获取token(获取用户信息的之前要获取的token)
type profileToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// RefreshToken string `json:"refresh_token"`
	Openid  string `json:"openid"`
	Scope   string `json:"scope"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (c *Client) getProfileToken(code string) (*profileToken, error) {
	res, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", c.appID, c.appSecret, code))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	pt := new(profileToken)
	err = json.NewDecoder(res.Body).Decode(pt)
	if pt.Errcode != 0 {
		return nil, fmt.Errorf("errcode:%d,errmsg:%s", pt.Errcode, pt.Errmsg)
	}
	return pt, err
}

// 使用token 获取用户信息
type Profile struct {
	Openid     string   `bson:"-" json:"openid"`              //
	Nickname   string   `bson:"nickname" json:"nickname"`     //
	Sex        int      `bson:"sex" json:"sex"`               //
	Province   string   `bson:"province" json:"province"`     //
	City       string   `bson:"city" json:"city"`             //
	Country    string   `bson:"country" json:"country"`       //
	HeadImgUrl string   `bson:"headimgurl" json:"headimgurl"` //
	Privilege  []string `bson:"privilege,omitempty" json:"privilege"`
	Unionid    string   `bson:"unionid,omitempty" json:"unionid"`
}

func (c *Client) GetProfile(code string) (*Profile, error) {
	pt, err := c.getProfileToken(code)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", pt.AccessToken, pt.Openid))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	profile := new(Profile)
	err = json.NewDecoder(res.Body).Decode(profile)
	return profile, err
}
