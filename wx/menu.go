package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) MenuCreate(data []byte) error {
	token, err := c.getTokenText()
	if err != nil {
		return err
	}
	res, err := http.Post(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", token), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
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
	logger.Println(result)
	return nil
}
