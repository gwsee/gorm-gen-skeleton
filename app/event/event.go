package event

import (
	"gorm-gen-skeleton/app/event/listen"
	"gorm-gen-skeleton/internal/variable"
)

type Event struct {
}

func (*Event) Init() error {
	err := variable.Event.Register(&listen.FooListen{})
	if err != nil {
		return err
	}
	return nil
}
