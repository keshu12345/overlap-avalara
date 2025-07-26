package constants

type Code string

func (c Code) String() string {
	return string(c)
}

const (
	BadRequest          Code = "BAD_REQUEST"
	NotFound            Code = "NOT_FOUND"
	RequestNotValid     Code = "REQUEST_NOT_VALID"
	RequestInvalid      Code = "REQUEST_INVALID"
	UnmarshalError      Code = "UNMARSHAl_ERROR"
	MarshalError        Code = "MARSHAL_ERR"
	ParseIntError       Code = "PARSE_INT_ERROR"
	DataNotFoundDbError Code = "DATA_NOT_FOUND_DB_ERROR"
	GoroutineError      Code = "GOROUTINE_ERROR"
	ParseFilesError     Code = "PARSE_FILES_ERROR"
	NotFoundMapError    Code = "NOT_FOUND_MAP_ERROR"
	UrlError            Code = "URL_ERROR"
	StatusUnauthorized  Code = "UNAUTHORIZED_ERROR"
)

type Filename string

const (
	OverlapFile Filename = "overlap"
)

func (fn Filename) String() string {
	return string(fn)
}
