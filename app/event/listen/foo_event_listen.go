package listen

import (
	"fmt"
	event2 "gorm-gen-skeleton/app/event/entity"
	"gorm-gen-skeleton/internal/event"
)

type FooListen struct {
}

func (*FooListen) Listen() []event.EventInterface {
	return []event.EventInterface{
		&event2.FooEvent{},
	}
}

func (*FooListen) Process(data any) {
	fmt.Printf("%v --> %s \n", data, "exec FooListen.Process")
}
