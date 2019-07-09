package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 生成带场景参数的二维码

// 临时二维码请求说明
// http请求方式: POST
// URL: https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=TOKEN
// POST数据格式：json
// POST数据例子：{"expire_seconds": 604800, "action_name": "QR_SCENE", "action_info": {"scene": {"scene_id": 123}}}
// 或者也可以使用以下POST数据创建字符串形式的二维码参数：
// {"expire_seconds": 604800, "action_name": "QR_STR_SCENE", "action_info": {"scene": {"scene_str": "test"}}}

// 永久二维码请求说明
// http请求方式: POST
// URL: https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=TOKEN
// POST数据格式：json
// POST数据例子：{"action_name": "QR_LIMIT_SCENE", "action_info": {"scene": {"scene_id": 123}}}
// 或者也可以使用以下POST数据创建字符串形式的二维码参数：
//              {"action_name": "QR_LIMIT_STR_SCENE", "action_info": {"scene": {"scene_str": "test"}}}

// 二维码类型，QR_SCENE为临时的整型参数值，QR_STR_SCENE为临时的字符串参数值，QR_LIMIT_SCENE为永久的整型参数值，QR_LIMIT_STR_SCENE为永久的字符串参数值
type QRCode struct {
	Expire_seconds int                               `json:"expire_seconds,omitempty"` //临时二维码过期时间
	Action_name    string                            `json:"action_name"`
	Action_info    map[string]map[string]interface{} `json:"action_info"`
}

// {"ticket":"gQH47joAAAAAAAAAASxodHRwOi8vd2VpeGluLnFxLmNvbS9xL2taZ2Z3TVRtNzJXV1Brb3ZhYmJJAAIEZ23sUwMEmm
// 3sUw==","expire_seconds":60,"url":"http://weixin.qq.com/q/kZgfwMTm72WWPkovabbI"}
type QRCodeRes struct {
	Ticket         string
	Expire_seconds int
	Url            string
	Errcode        int
	Errmsg         string
}

func (c *Client) QRCodeGen() {
	token, err := c.getTokenText()
	if err != nil {

	}
	http.Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", token), "application/json", nil)
}

// 永久二维码int参数生成
func (c *Client) QRCodeForEverGen(sceneId int) (*QRCodeRes, error) {
	token, err := c.getTokenText()
	if err != nil {
		return nil, err
	}
	qrc := &QRCode{
		Action_name: "QR_LIMIT_SCENE",
		Action_info: map[string]map[string]interface{}{"scene": map[string]interface{}{"scene_id": sceneId}},
	}
	data, _ := json.Marshal(qrc)
	res, err := http.Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", token), "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	qrcRes := new(QRCodeRes)
	if err = json.NewDecoder(res.Body).Decode(qrcRes); err != nil {
		return nil, err
	}
	if qrcRes.Errcode != 0 {
		return nil, fmt.Errorf("errcode:%d,errmsg:%s", qrcRes.Errcode, qrcRes.Errmsg)
	}
	return qrcRes, nil
}

// HTTP GET请求（请使用https协议）https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=TICKET
// 提醒：TICKET记得进行UrlEncode
// 获取图片字节数据
func (c *Client) QRCodeGet(ticket string) ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", ticket))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}

func QRCodeUrl(ticket string) string {
	return fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", url.QueryEscape(ticket))
}
