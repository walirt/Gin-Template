package constant

type ResponseError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

var (
	OK                = ResponseError{ErrorCode: 0, ErrorMsg: "ok"}
	ParameterAbnormal = ResponseError{ErrorCode: 1, ErrorMsg: "parameter abnormal"}
	TokenAbnormal     = ResponseError{ErrorCode: 2, ErrorMsg: "token abnormal"}
	VersionAbnormal   = ResponseError{ErrorCode: 3, ErrorMsg: "version abnormal"}
	OtherError        = ResponseError{ErrorCode: 4, ErrorMsg: "other error, %v"}
	UnknownError      = ResponseError{ErrorCode: 5, ErrorMsg: "unknown error"}

	UsernameOrPasswordWrong               = ResponseError{ErrorCode: 100, ErrorMsg: "user name or password wrong"}
	ConfigurationNotNeedUpdate            = ResponseError{ErrorCode: 101, ErrorMsg: "no need to update configuration"}
	GuidDoesNotExist                      = ResponseError{ErrorCode: 102, ErrorMsg: "the requested guid does not exist"}
	ModeValueAbnormal                     = ResponseError{ErrorCode: 103, ErrorMsg: "mode value abnormal"}
	StrategyIDDoesNotExist                = ResponseError{ErrorCode: 104, ErrorMsg: "the strategy id does not exist"}
	BeginTimeLaterThanEndTime             = ResponseError{ErrorCode: 105, ErrorMsg: "begin-time is later than end-time"}
	TimeFormatIncorrect                   = ResponseError{ErrorCode: 106, ErrorMsg: "incorrect time format"}
	SerialNumDoesNotExist                 = ResponseError{ErrorCode: 107, ErrorMsg: "the serial num does not exist"}
	MsgTypeDoesNotExist                   = ResponseError{ErrorCode: 108, ErrorMsg: "the msg type does not exist"}
	RepeatAlarmOperation                  = ResponseError{ErrorCode: 109, ErrorMsg: "repeat alarm operation"}
	ValueAbnormal                         = ResponseError{ErrorCode: 110, ErrorMsg: "value abnormal"}
	ControlFailure                        = ResponseError{ErrorCode: 111, ErrorMsg: "control failure"}
	SpecifiedControlPointIsUncontrollable = ResponseError{ErrorCode: 112, ErrorMsg: "the specified control point is uncontrollable"}
)

type Response struct {
	ResponseError
	Data interface{} `json:"data"`
}

type ResponsePointsPart struct {
	Guid      string `json:"guid"`
	Status    int    `json:"status"`
	Value     string `json:"value"`
	Timestamp int    `json:"timestamp"`
}

type ResponseDevicesPart struct {
	Guid   string               `json:"guid"`
	Status int                  `json:"status"`
	Points []ResponsePointsPart `json:"points"`
}
