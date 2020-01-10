package cache

//Reply ...
type Reply string

//Stored ..
const (
	StoredReply         Reply = "STORED"
	NotStoredReply            = "NOT_STORED"
	ErrReply                  = "ERROR"
	NotImplementedReply       = "NOT_IMPLEMENTED"
	ValueReply                = "VALUE"
	NotFoundReply             = "NOT_FOUND"
	DeletedReply              = "DELETED"
	ClientErrorReply          = "CLIENT_ERROR"
)
