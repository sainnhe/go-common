package constant

import "time"

const (
	// LogInitTimeFormat x is the time format used when initializing the logger.
	LogInitTimeFormat = time.StampMilli

	// LogMsgInitNil is the log message for initialization error due to nil field.
	LogMsgInitNil = "Init failed. Detected nil dependency."
	// LogMsgInitSuccess is the message for successful initialization.
	LogMsgInitSuccess = "Init success."

	// LogAttrError is the field name for error.
	LogAttrError = "error"
	// LogAttrAPI is the field name for API.
	LogAttrAPI = "api"
	// LogAttrMethod is the field name for method.
	LogAttrMethod = "method"
	// LogAttrDO is the field name for Data Object.
	LogAttrDO = "do"
	// LogAttrSQL is the field name for SQL.
	LogAttrSQL = "sql"
	// LogAttrHandler is the field name for handler.
	LogAttrHandler = "handler"
	// LogAttrWantMethod is the field name for desired method.
	LogAttrWantMethod = "want"
	// LogAttrGetMethod is the field name for obtained method.
	LogAttrGetMethod = "get"
)

type logAttrCtxKeyType struct{}

// LogAttrCtxKey is the context key whose value should be logged.
var LogAttrCtxKey = logAttrCtxKeyType{}
