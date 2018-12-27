package errcodes

type ErrCode int

const (
	OK                ErrCode = 0
	JSONParseErr      ErrCode = 1
	UnkownCommand     ErrCode = 2
	AlreadyInside     ErrCode = 3
	UnfulfilledReqs   ErrCode = 4
	CannotInTown      ErrCode = 5
	DoesNotExist      ErrCode = 6
	OutOfResource     ErrCode = 7
	NewTurn           ErrCode = 8
	MonstersStillLive ErrCode = 9
)
