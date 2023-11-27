package message

type Package struct {
	ServerID uint32
	Module   string
	Body     any
	TTL      int32
}

type BroadcastPackage struct {
	ServerType string
	Module     string
	Body       any
}

type RPCRequest struct {
	Caller   string
	ServerID uint32
	Module   string
	Req      any
	Resp     any
	Cb       func(resp any, err error)
	Err      error
}

type RPCPackage struct {
	Req  any
	Resp chan any
	Err  chan error
}
