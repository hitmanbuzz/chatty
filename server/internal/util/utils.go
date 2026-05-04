package util

const MAX_USER_BYTE = 128
const MAX_USER_LEN = 24

const MAX_MSG_BYTE = 1024
const MAX_MSG_LEN = 100

const DAY_SECS = 86400

type MessageStat int

const (
	EMPTY_MSG MessageStat = iota
	BAD_MSG
	GOOD_MSG
)

// This types if for checking during websocket unlike MessageStat which is after this
type MessageType int

const (
	DISCONNECT MessageType = iota
	ABRUPT_DISCONNECT
	ABNORMAL_DISCONNECT
	FAIL_JSON_PARSE
	EXCEED_MAX_MSG_BYTE
	NIL
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)
