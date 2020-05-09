package system

import (
	"fmt"

	"tino/util/config"
)

type Error interface {
	Error() string
	Map() Map
	ResultMap(args ...string) Map
}

type MessageLang struct {
	En string
	Zh string
}

type KError struct {
	Code    interface{}
	Message MessageLang `json:"message"`
}

func (e *KError) Error() string {
	return e.Message.En
}

func (e *KError) Map() Map {
	return Map{"code": e.Code, "message": Map{"zh": e.Message.Zh, "en": e.Message.En}}
}

func (e *KError) ResultMap(args ...string) Map {
	if args != nil {
		e.Message.En = args[0]
	}
	return Map{"status": e.Map()}
}

func GetKError(key string) *KError {
	lan := config.GetKeyLang(key)
	if lan == nil {
		logger.Println("Error: Can't find language resource,Must modify................................................................................")
		return &KError{Code: "ERROR_SYSTEM", Message: MessageLang{En: fmt.Sprintf("Language error :[%s]", key), Zh: fmt.Sprintf("多语言资源Error：[%d]", key)}}
	}
	return &KError{Code: lan.Code, Message: MessageLang{En: lan.En, Zh: lan.Zh}}
}

func ErrorResult(key string) Map {
	e := GetKError(key)
	return Map{"status": e.Map()}
}

func NewErrorResult(code string, msg string) Map {
	e := NewKError(code, msg)
	return Map{"status": e.Map()}
}

func NewKError(code string, msg string) *KError {
	return &KError{Code: code, Message: MessageLang{En: msg, Zh: msg}}
}

func SqlError(err error) *KError {
	return &KError{Code: DbError, Message: MessageLang{En: err.Error(), Zh: err.Error()}}
}
