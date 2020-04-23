package errkit

import "strconv"

type LogicCode struct {
	Code int    `json:"code" desc:"业务错误码"`
	Msg  string `json:"msg" desc:"错误描述"`
}

func New(Code int, Msg string) error {
	return &LogicCode{Code: Code, Msg: Msg}
}

func (code *LogicCode) Error() string {
	return "[" + strconv.Itoa(code.Code) + "]<" + code.Msg + ">"
}
