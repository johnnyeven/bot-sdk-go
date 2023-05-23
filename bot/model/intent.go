package model

import (
	"github.com/johnnyeven/bot-sdk-go/bot/data"
)

type Intent struct {
	data               data.Intent
	Name               string
	ConfirmationStatus string
}

func NewIntent(intent data.Intent) *Intent {
	return &Intent{
		data: intent,
		Name: intent.Name,

		ConfirmationStatus: intent.ConfirmationStatus,
	}
}

// 根据槽位名获取槽位
func (i *Intent) GetSlot(name string) *data.Slot {
	slot, ok := i.data.Slots[name]
	if ok {
		return &slot
	}
	return nil
}

// 根据槽位名获取槽位对应的值
func (i *Intent) GetSlotValue(name string) string {
	slot := i.GetSlot(name)
	if slot != nil {
		return slot.Value
	}
	return ""
}

// 根据槽位名获取槽位对应的状态
func (i *Intent) GetSlotStatus(name string) string {
	slot := i.GetSlot(name)
	if slot != nil {
		return slot.ConfirmationStatus
	}
	return ""

}

// 设置槽位的值
func (i *Intent) SetSlotValue(name string, value string) bool {
	slot := i.GetSlot(name)
	if slot != nil {
		slot.Value = value
		return true
	} else {
		// TODO new a slot
	}
	return false
}

func (i *Intent) GetData() data.Intent {
	return i.data
}
