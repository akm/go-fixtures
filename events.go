package fixtures

import "gorm.io/gorm"

type Handler = func()

type BaseEvent struct {
	handlers []Handler
}

func (h *BaseEvent) On(handler Handler) {
	h.handlers = append(h.handlers, handler)
}

type BeforeCreate struct {
	BaseEvent
}

func (h *BeforeCreate) BeforeCreate(tx *gorm.DB) (err error) {
	for _, handler := range h.handlers {
		handler()
	}
	return nil
}

type AfterCreateEvent struct {
	BaseEvent
}

func (h *AfterCreateEvent) AfterCreate(tx *gorm.DB) (err error) {
	for _, handler := range h.handlers {
		handler()
	}
	return nil
}
