// 微信消息处理
package wx

import (
	"encoding/xml"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

/*
微信消息处理
*/
//消息类型
const (
	MT_EVENT = iota + 1
	MT_TEXT
	MT_IMAGE
	MT_VOICE
	MT_VIDEO
	MT_SHORTVIDEO
	MT_LOCATION
	MT_LINK
)

// 消息类型
type MT struct {
	XML     xml.Name `xml:"xml"`
	MsgType string   `xml:"MsgType"`
}

func getMT(data []byte) int {
	v := new(MT)
	if err := xml.Unmarshal(data, v); err != nil {
		log.Panicln(err)
	}
	switch v.MsgType {
	case "event":
		return MT_EVENT
	case "text":
		return MT_TEXT
	case "image":
		return MT_IMAGE
	case "voice":
		return MT_VOICE
	case "video":
		return MT_VIDEO
	case "shortvideo":
		return MT_SHORTVIDEO
	case "location":
		return MT_LOCATION
	case "link":
		return MT_LINK
	}
	panic("不能识别的微信事件推送")
}

const (
	E_SUBSCRIBE   = iota + 1 //关注
	E_UNSUBSCRIBE            //取消关注
	E_SCAN                   //用户已关注时的事件推送
	E_LOCATION               //上报地理位置事件
	E_CLICK                  //自定义菜单事件
	E_VIEW                   //点击菜单跳转链接时的事件推送
)

type ET struct {
	XML    xml.Name `xml:"xml"`
	Eevent string   `xml:"Event"`
}

// 获取事件类型
func getET(data []byte) int {
	v := new(ET)
	if err := xml.Unmarshal(data, v); err != nil {
		log.Panicln(err)
	}
	switch v.Eevent {
	case "subscribe":
		return E_SUBSCRIBE
	case "unsubscribe":
		return E_UNSUBSCRIBE
	case "SCAN":
		return E_SCAN
	case "LOCATION":
		return E_LOCATION
	case "CLICK":
		return E_CLICK
	case "VIEW":
		return E_VIEW
	}
	return -1 //暂未识别
}

type GetMsgBaser interface {
	getMsgBase() *msgBase
}

type msgBase struct {
	XML          xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MT
}

func (m msgBase) getMsgBase() *msgBase {
	return &m
}

type eventBase struct {
	msgBase
	Event string `xml:"Event" text:"事件类型"`
}

// 1. 用户未关注时，进行关注后的事件推送,2. 用户已关注时的事件推送
type SubEvent struct {
	eventBase
	EventKey string `xml:"EventKey" text:"事件KEY值，qrscene_为前缀，后面为二维码的参数值"`
	Ticket   string `xml:"Ticket" text:"二维码的ticket，可用来换取二维码图片"`
}

// 获取int类型的从场景值
func (se *SubEvent) GetQRsceneInt() (int, error) {
	xs := strings.Split(se.EventKey, "_")
	if len(xs) == 2 {
		v, err := strconv.Atoi(xs[1])
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	return 0, errors.New("扫描了不到场景值的二维码")
}

func (se *SubEvent) MustGetQRsceneInt() int {
	v, _ := se.GetQRsceneInt()
	return v
}

func getSubEvent(data []byte) *SubEvent {
	v := new(SubEvent)
	if err := xml.Unmarshal(data, v); err != nil {
		log.Panicln(err)
	}
	return v
}

// 取消关注事件
type UnSubEvent struct {
	eventBase
}

func getUnSubEvent(data []byte) *UnSubEvent {
	v := new(UnSubEvent)
	if err := xml.Unmarshal(data, v); err != nil {
		log.Panicln(err)
	}
	return v
}

type TextMsg struct {
	msgBase
	Content string `xml:"Content"`
	MsgId   string `xml:"MsgId"`
}

func getTextMsg(data []byte) *TextMsg {
	v := new(TextMsg)
	if err := xml.Unmarshal(data, v); err != nil {
		log.Panicln(err)
	}
	return v
}

/*
解析微信消息类型和消息体
返回消息类型和消息体
*/
func GetMsg(data []byte) (interface{}, int, int) {
	switch getMT(data) {
	case MT_TEXT:
		return getTextMsg(data), MT_TEXT, 0
	case MT_EVENT:
		switch getET(data) {
		case E_SUBSCRIBE:
			return getSubEvent(data), MT_EVENT, E_SUBSCRIBE
		case E_UNSUBSCRIBE:
			return getUnSubEvent(data), MT_EVENT, E_UNSUBSCRIBE
		default:
			// log.Printf("待开发的事件类型:%s", string(data))
			return nil, 0, 0
		}
	default:
		// log.Printf("待开发的消息类型:%s", string(data))
		return nil, 0, 0
	}
}

/*响应消息数据格式*/
type CDATA struct {
	Text string `xml:",innerxml"`
}

func newCDATA(v string) CDATA {
	return CDATA{"<![CDATA[" + v + "]]>"}
}

type replyBase struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   CDATA
	MsgType      CDATA
}

/*响应消息*文本格式*/
type textReply struct {
	replyBase
	Content CDATA
}

func ReplyText(mber GetMsgBaser, text string) *textReply {
	mb := mber.getMsgBase()
	tr := new(textReply)
	tr.ToUserName = newCDATA(mb.FromUserName)
	tr.FromUserName = newCDATA(mb.ToUserName)
	tr.CreateTime = newCDATA(strconv.FormatInt(time.Now().Unix(), 10))
	tr.MsgType = newCDATA("text")
	tr.Content = newCDATA(text)
	return tr
}
