package spider

import (
	"net/url"
)

// Target 目标网站
type Target struct {
	url         *url.URL
	interpreter *Interpreter
}

func (t *Target) Interpreter() *Interpreter {
	return t.interpreter
}

func (t *Target) SetInterpreter(interpreter *Interpreter) {
	t.interpreter = interpreter
}

func (t *Target) Url() *url.URL {
	return t.url
}

func (t *Target) SetUrl(url *url.URL) {
	t.url = url
}
