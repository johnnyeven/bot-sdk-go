package model

import (
	"encoding/json"

	"github.com/johnnyeven/bot-sdk-go/bot/data"
	"github.com/johnnyeven/bot-sdk-go/bot/util"
)

type Response struct {
	session *Session
	request interface{}
	data    map[string]interface{}
}

func NewResponse(session *Session, request interface{}) *Response {
	d := make(map[string]interface{})
	return &Response{
		data:    d,
		session: session,
		request: request,
	}
}

/**
 * 询问用户时，返回的speech.
 * 此时设备的麦克风会进入收音状态，比如设备灯光亮起
 * TIP: 一般技能要完成一项任务，还缺少一些信息，主动发起对用户的询问的时候使用
 */
func (r *Response) Ask(speech string) *Response {
	r.Tell(speech)
	r.HoldOn()
	return r
}

func (r *Response) AskSlot(speech string, slot string) *Response {
	r.Ask(speech)

	request, ok := r.request.(IntentRequest)
	if ok {
		request.Dialog.ElicitSlot(slot)
	}
	return r
}

/**
 * 回复用户，返回的speech
 */
func (r *Response) Tell(speech string) *Response {
	r.data["outputSpeech"] = util.FormatSpeech(speech)
	return r
}

/**
 * 回复用户，返回的speech
 */
func (r *Response) Reprompt(speech string) *Response {
	r.data["reprompt"] = map[string]interface{}{
		"outputSpeech": util.FormatSpeech(speech),
	}
	return r
}

/**
 * 返回卡片.
 * 针对有屏幕的设备，比如: 电视、show，可以呈现更多丰富的信息给用户
 * 卡片协议参考：TODO
 */
func (r *Response) DisplayCard(card interface{}) *Response {
	r.data["card"] = card

	return r
}

/**
 * 返回指令. 比如，返回音频播放指令，使设备开始播放音频
 * TIP: 可以同时返回多个指令，设备按返回顺序执行这些指令，指令协议参考TODO
 */
func (r *Response) Command(directive interface{}) *Response {
	_, ok := r.data["directives"]
	if !ok {
		r.data["directives"] = make([]interface{}, 0)
	}

	directives, ok := r.data["directives"].([]interface{})
	directives = append(directives, directive)

	r.data["directives"] = directives

	return r
}

/**
 * 保持会话.
 * 此时设备的麦克风会自动开启监听用户说话
 */
func (r *Response) HoldOn() *Response {
	r.data["shouldEndSession"] = false
	return r
}

/**
 * 保持会话.
 * 关闭麦克风
 */
func (r *Response) CloseMicrophone() *Response {
	r.data["expectSpeech"] = true
	return r
}

func (r *Response) Build() string {
	//session
	attributes := r.session.GetData().Attributes

	ret := map[string]interface{}{
		"version":  "2.0",
		"session":  data.SessionResponse{Attributes: attributes},
		"response": r.data,
	}

	//intent request
	request, ok := r.request.(IntentRequest)
	if ok {
		ret["context"] = data.ContextResponse{Intent: request.Dialog.Intents[0].GetData()}

		directive := request.Dialog.GetDirective()
		if directive != nil {
			r.Command(directive)
		}
	}

	response, _ := json.Marshal(ret)

	return string(response)
}

func (r *Response) GetData() map[string]interface{} {
	return r.data
}
