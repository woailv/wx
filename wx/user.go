package wx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// http请求方式: GET（请使用https协议）
// https://api.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&next_openid=NEXT_OPENID
/*
{
    "total":2,
    "count":2,
    "data":{
    "openid":["OPENID1","OPENID2"]},
    "next_openid":"NEXT_OPENID"
}*/
type User struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		Openid []string `json:"openid"`
	} `json:"data"`
	Next_openid string `json:"next_openid"`

	Errcode int
	Errmsg  string
}

func (c *Client) UserGet(nextOpenid ...string) (*User, error) {
	token, err := c.getTokenText()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s", token)
	if len(nextOpenid) > 1 {
		return nil, fmt.Errorf("nextOpenid 错误")
	}
	if len(nextOpenid) == 1 {
		url = url + "&next_openid=" + nextOpenid[0]
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	u := new(User)
	if err = json.NewDecoder(res.Body).Decode(u); err != nil {
		return nil, err
	}
	if u.Errcode != 0 {
		return nil, fmt.Errorf("errcode:%d,errmsg:%s", u.Errcode, u.Errmsg)
	}
	return u, nil
}

// http请求方式: GET
// https://api.weixin.qq.com/cgi-bin/user/info?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
/*
{
    "subscribe": 1,
    "openid": "o6_bmjrPTlm6_2sgVt7hMZOPfL2M",
    "nickname": "Band",
    "sex": 1,
    "language": "zh_CN",
    "city": "广州",
    "province": "广东",
    "country": "中国",
    "headimgurl":"http://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
    "subscribe_time": 1382694957,
    "unionid": " o6_bmasdasdsad6_2sgVt7hMZOPfL"
    "remark": "",
    "groupid": 0,
    "tagid_list":[128,2],
    "subscribe_scene": "ADD_SCENE_QR_CODE",
    "qr_scene": 98765,
    "qr_scene_str": ""
}*/
type UserInfo struct { //带扩展
	Openid   string
	Nickname string

	Errcode int
	Errmsg  string
}

func (c *Client) UserInfo(openid string) (*UserInfo, error) {
	token, err := c.getTokenText()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", token, openid)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	ui := new(UserInfo)
	if err = json.NewDecoder(res.Body).Decode(ui); err != nil {
		return nil, err
	}
	if ui.Errcode != 0 {
		return nil, fmt.Errorf("errcode:%d,errmsg:%s", ui.Errcode, ui.Errmsg)
	}
	return ui, nil
}
