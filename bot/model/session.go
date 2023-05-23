package model

import (
	"github.com/johnnyeven/bot-sdk-go/bot/data"
)

type Session struct {
	data data.Session
}

func NewSession(data data.Session) *Session {
	if data.Attributes == nil {
		data.Attributes = make(map[string]string)
	}

	return &Session{
		data: data,
	}
}

// 当前session是否是新的
func (s *Session) IsNew() bool {
	return s.data.New
}

// 获取session id
func (s *Session) GetId() string {
	return s.data.SessionId
}

// 获取session中对应字段的值
func (s *Session) GetAttribute(key string) string {
	value, ok := s.data.Attributes[key]
	if ok {
		return value
	}
	return ""
}

// 设置session中对应字段的值
func (s *Session) SetAttribute(key, value string) {
	if key == "" {
		return
	}

	s.data.Attributes[key] = value
}

func (s *Session) GetData() data.Session {
	return s.data
}
