package system

/***************************************************************
 * 通用错误码
 *
 * 只有通用的，常用的错误码放在这里，每个接口私有的错误码，放在接口内部
 ***************************************************************/

const (
	//	SqlErrorID    = 9
	//	DelErrorID    = 7
	//	ModifyErrorID = 8

	SystemError           = "SystemError"
	DbError               = "ERROR_DB_OPERATION"
	ParamError            = "ParamError"
	AccountNullError      = "AccountNullError"
	AccountHaveExistError = "AccountHaveExistError"
	//	NoDataError           = "NoDataError"

	FileNoExistError  = "FileNoExistError"
	PasswordNullError = "PasswordNullError"
	PasswordError     = "PasswordError"
	OldPasswordError  = "OldPasswordError"

	PasswordNoRegular    = "PasswordNoRegular"
	MailNoRegularError   = "MailNoRegularError"
	MobileNoRegularError = "MobileNoRegularError"
	VCodeError           = "VCodeError"

	AccountNoExistError = "AccountNoExistError"
	TokenError          = "TokenError"
	//	UserNoExistError       = "UserNoExistError"
	//	AccountOrPasswordError = "AccountOrPasswordError"
	AccountHaveBindError = "AccountHaveBindError"
	AuthenticError       = "AuthenticError"

	HaveDeleteError = "HaveDeleteError"
	BindRemoveError = "BindRemoveError"
	FrequentlyError = "FrequentlyError"
	//	UserNoFindError = "UserNoFindError"
	MacRepeatError = "MacRepeatError"
	MailNoExist    = "MailNoExist"
	StockNoExist   = "StockNoExist"
)

var (
	SuccessStatus        = &KError{Code: 0, Message: MessageLang{En: "Successful", Zh: "成功"}}
	ParameterErrorStatus = &KError{Code: "ERROR_PARAMS", Message: MessageLang{En: "Param Error", Zh: "参数解析错误"}}
)
