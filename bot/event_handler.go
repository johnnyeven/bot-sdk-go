package bot

import (
	"github.com/johnnyeven/bot-sdk-go/bot/model"
)

// ListTemplate 列表选择事件
// Display.ButtonClicked 事件
// https://dueros.baidu.com/didp/doc/dueros-bot-platform/dbp-custom/display-template_markdown#Display.ElementSelected%E4%BA%8B%E4%BB%B6
func (b *Bot) OnDisplayElementSelected(fn func(bot *Bot, request *model.EventRequest)) {
	b.AddEventListener(
		"Display.ElementSelected", func(bot *Bot, request interface{}) {
			req := request.(model.EventRequest)
			fn(bot, &req)
		},
	)
}

// Display.ButtonClicked 事件
// ```javascript
//{
//	"type": "Display.ButtonClicked",
//	"requestId": "{{STRING}}",
//	"timestamp": "{{STRING}}",
//	"token": "{{STRING}}",
//	"buttonType": "{{ENUM}}"
//}
// ```

// LinkAccountSucceeded 事件
// ```javascript
//
//	{
//	   "type": "Connections.Response",
//	   "name": "LinkAccountSucceeded",
//	   "requestId": "{{STRING}}",
//	   "timestamp": {{INT32}},
//	   "token": "{{STRING}}"
//	}
//
// ```
func (b *Bot) OnLinkAccountSuccessed(fn func(bot *Bot, request *model.EventRequest)) {
	b.AddEventListener(
		"LinkAccountSucceeded", func(bot *Bot, request interface{}) {
			req := request.(model.EventRequest)
			fn(bot, &req)
		},
	)
}

// Screen.LinkClicked事件
// https://dueros.baidu.com/didp/doc/dueros-bot-platform/dbp-custom/cards_markdown#Screen.LinkClicked%E4%BA%8B%E4%BB%B6
//
//	{
//	   "type": "Screen.LinkClicked",
//	   "url": "{{STRING}}",
//	   "requestId": "{{STRING}}",
//	   "timestamp": {{INT32}}
//	   "token": "{{STRING}}"
//	}
func (b *Bot) OnScreenLinkClicked(fn func(bot *Bot, request *model.EventRequest)) {
	b.AddEventListener(
		"Screen.LinkClicked", func(bot *Bot, request interface{}) {
			req := request.(model.EventRequest)
			fn(bot, &req)
		},
	)
}
