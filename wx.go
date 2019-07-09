// 微信公众号常用接口
package wx

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"sort"
)

var logger = log.New(os.Stdout, "[wx]", log.LstdFlags|log.Lshortfile)

// 用户参数配置
// tokenCacheFunc action:get,set; set values[0] = value ,values[1] = expiry
type Client struct {
	appID          string
	appSecret      string
	tokenCacheFunc func(action, key string, values ...interface{}) (string, error)
	logger         *log.Logger
}

func New(appid, appsecret string) *Client {
	return &Client{appID: appid, appSecret: appsecret, logger: logger}
}

func (c *Client) SetLogger(logger *log.Logger) {
	c.logger = logger
}

func (c *Client) SetTokenCacheFunc(f func(action, key string, values ...interface{}) (string, error)) {
	c.tokenCacheFunc = f
}

// 签名验证(确定是否为微信服务器发过来的消息),token:公众号页面配置的token
func SignVerify(token, signature, timestamp, nonce string) bool {
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, signature := range strs {
		str += signature
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil)) == signature
}
