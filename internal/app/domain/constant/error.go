package constant

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	ErrSys001 = Error{
		Code:    "SYS_001",
		Message: "Internal server error",
	}
)
