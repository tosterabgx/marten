package protocol

const ControlPort uint16 = 6472
const DefaultServerAddr = "marten.tosterabgx.me"

type ClientHello struct {
	RequiredPort uint16 `json:"ClientHello"`
}

type ServerHello struct {
	ActualPort uint16 `json:"ServerHello"`
}
