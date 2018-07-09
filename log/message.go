package log

import (
	"fmt"
)

type message struct {
	values []interface{}
}

func newMessage(v ...interface{}) fmt.Stringer {
	message := new(message)
	message.values = v
	return message
}

func (this *message) String() string {
	return fmt.Sprint(this.values...)
}

// format
type formattedMessage struct {
	format string
	values []interface{}
}

func newFormattedMessage(format string, v ...interface{}) *formattedMessage {
	message := new(formattedMessage)
	message.format = format
	message.values = v
	return message
}

func (this *formattedMessage) String() string {
	return fmt.Sprintf(this.format, this.values...)
}
