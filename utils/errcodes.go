package errcodes

type ErrCode int

const (
	OK            ErrCode = 0
	JSONParseErr  ErrCode = 1
	UnkownCommand ErrCode = 2
)
