package bot

import (
	"reflect"

	"github.com/johnnyeven/bot-sdk-go/bot/model"
)

// 技能基础类
type Bot struct {
	intentHandler map[string]func(bot *Bot, request *model.IntentRequest) // 针对intent requset不同intent的处理函数
	//eventHandler               map[string]func(bot *Bot, request *model.EventRequest)  // 针对事件的处理函数
	eventHandler               map[string]func(bot *Bot, request interface{}) // 针对事件的处理函数
	defaultEventHandler        func(bot *Bot, request interface{})
	launchRequestHandler       func(bot *Bot, request *model.LaunchRequest)       // 针对技能打开的处理函数
	sessionEndedRequestHandler func(bot *Bot, request *model.SessionEndedRequest) // 针对技能关闭的处理函数
	Request                    interface{}                                        // 对当前request的封装，需要在使用时断言，判断当前的类型
	Session                    *model.Session                                     // 对session的封装
	Response                   *model.Response                                    // 对技能返回的封装
}

// 创建常驻bot类，可维持在内存状态中, addhandler 和 addEventer事件可以缩减为一次
func NewBot() *Bot {
	return &Bot{
		intentHandler: make(map[string]func(bot *Bot, request *model.IntentRequest)),
		eventHandler:  make(map[string]func(bot *Bot, request interface{})),
	}
}

// 根据每个请求分别处理
func (b *Bot) Handler(request string) string {
	b.Request = model.NewRequest(request)
	b.Session = model.NewSession(model.GetSessionData(request))
	b.Response = model.NewResponse(b.Session, b.Request)

	b.dispatch()

	return b.Response.Build()
}

// 添加对intent的处理函数
func (b *Bot) AddIntentHandler(intentName string, fn func(bot *Bot, request *model.IntentRequest)) {
	if intentName != "" {
		b.intentHandler[intentName] = fn
	}
}

// 添加对事件的处理函数
func (b *Bot) AddEventListener(eventName string, fn func(bot *Bot, request interface{})) {
	if eventName != "" {
		b.eventHandler[eventName] = fn
	}
}

// 添加事件默认处理函数
// 比如，在播放视频时，技能会收到各种事件的上报，如果不想一一处理可以使用这个来添加处理
func (b *Bot) AddDefaultEventListener(fn func(bot *Bot, request interface{})) {
	b.defaultEventHandler = fn
}

// 打开技能时的处理
func (b *Bot) OnLaunchRequest(fn func(bot *Bot, request *model.LaunchRequest)) {
	b.launchRequestHandler = fn
}

// 技能关闭的处理，比如可以做一些清理的工作
// TIP: 根据协议，技能关闭返回的结果，DuerOS不会返回给用户。
func (b *Bot) OnSessionEndedRequest(fn func(bot *Bot, request *model.SessionEndedRequest)) {
	b.sessionEndedRequestHandler = fn
}

func (b *Bot) dispatch() {
	switch request := b.Request.(type) {
	case model.IntentRequest:
		b.processIntentHandler(request)
		return
	case model.LaunchRequest:
		b.processLaunchHandler(request)
		return
	case model.SessionEndedRequest:
		b.processSessionEndedHandler(request)
		return
	}
	b.processEventHandler(b.Request)
}

func (b *Bot) processLaunchHandler(request model.LaunchRequest) {
	if b.launchRequestHandler != nil {
		b.launchRequestHandler(b, &request)
	}
}

func (b *Bot) processSessionEndedHandler(request model.SessionEndedRequest) {
	if b.sessionEndedRequestHandler != nil {
		b.sessionEndedRequestHandler(b, &request)
	}
}

func (b *Bot) processIntentHandler(request model.IntentRequest) {
	intentName, _ := request.GetIntentName()
	fn, ok := b.intentHandler[intentName]

	if ok {
		fn(b, &request)
		return
	}
}

func (b *Bot) processEventHandler(req interface{}) {
	rVal := reflect.ValueOf(req)
	eventType := rVal.FieldByName("Type").Interface().(string)

	fn, ok := b.eventHandler[eventType]

	if ok {
		fn(b, req)
		return
	}

	if b.defaultEventHandler != nil {
		b.defaultEventHandler(b, req)
	}
}
