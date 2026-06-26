package protocol

const ControlPort uint16 = 6472

type ClientHello struct {
	RequiredPort uint16 `json:"ClientHello"`
}

type ServerHello struct {
	ActualPort uint16 `json:"ServerHello"`
}
