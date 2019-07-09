package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// 设置所属行业
func (c *Client) SetIndustry(xs2 []string) error {
	token, err := c.getTokenText()
	if err != nil {
		return err
	}
	data := fmt.Sprintf(`{"industry_id1":"%s","industry_id2":"%s"}`, xs2[0], xs2[1])
	res, err := http.Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token=%s", token), "application/json", strings.NewReader(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	result := struct {
		Errcode int
		Errmsg  string
	}{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}
	if result.Errcode != 0 {
		return fmt.Errorf("errcode:%d,errmsg:%s", result.Errcode, result.Errmsg)
	}
	return nil
}

// 获取设置的行业信息
type class struct {
	First_class  string
	Second_class string
}
type industry struct {
	Primary_industry   class
	Secondary_industry class
}

func (i *industry) String() string {
	return fmt.Sprintf("主营行业:%s,%s;副营行业:%s,%s", i.Primary_industry.First_class, i.Primary_industry.Second_class, i.Secondary_industry.First_class, i.Secondary_industry.Second_class)
}

func (c *Client) GetIndustry() (string, error) {
	token, err := c.getTokenText()
	if err != nil {
		return "", err
	}
	res, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=%s", token))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	defer res.Body.Close()
	in := industry{}
	if err := json.NewDecoder(res.Body).Decode(&in); err != nil {
		return "", err
	}
	return in.String(), nil
}

// 获取模板列表
type MsgTpls struct {
	Template_list []struct {
		Template_id, Title, Content string
	}
}

func (c *Client) GetAllMsgTpl() (*MsgTpls, error) {
	t, err := c.getTokenText()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	res, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=%s", t))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	tpls := MsgTpls{}
	if err := json.NewDecoder(res.Body).Decode(&tpls); err != nil {
		return nil, err
	}
	return &tpls, err
}

// 删除模板
func (c *Client) DelMsgTpl(msgTplId string) error {
	t, err := c.getTokenText()
	if err != nil {
	}
	data := fmt.Sprintf(`{"template_id":"%s"}`, msgTplId)
	res, err := http.Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/del_private_template?access_token=%s", t), "application/json", strings.NewReader(data))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	m := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return err
	}
	if m["errcode"].(float64) != 0 {
		return fmt.Errorf("errcode:%f,errmsg:%s", m["errcode"].(float64), m["errmsg"])
	}
	return nil
}

// 发送模板消息
type TplMsgData struct {
	Touser      string `json:"touser"`      //接收者openid
	Template_id string `json:"template_id"` //	模板ID
	Url         string `json:"url"`         //	模板跳转链接（海外帐号没有跳转能力）
	Miniprogram struct {
		Appid    string `json:"appid"`
		Pagepath string `json:"pagepath"`
	} `json:"miniprogram"` //	跳小程序所需数据，不需跳小程序可不用传该数据
	Data map[string]map[string]string `json:"data"` //模板数据
}

func (c *Client) TplMsgSend(tplMsgData *TplMsgData) (int, error) {
	t, err := c.getTokenText()
	if err != nil {
	}
	data, err := json.Marshal(tplMsgData)
	if err != nil {
		return 0, err
	}
	res, err := http.Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", t), "application/json", bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error in status code:%d", res.StatusCode)
	}
	m := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return 0, err
	}
	if m["errcode"].(float64) != 0 {
		return 0, fmt.Errorf("errcode:%f,errmsg:%s", m["errcode"].(float64), m["errmsg"])
	}
	return int(m["msgid"].(float64)), nil
}

// 事件推送TODO
// 在模版消息发送任务完成后，微信服务器会将是否送达成功作为通知，发送到开发者中心中填写的服务器配置地址中。
