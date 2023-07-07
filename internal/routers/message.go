package routers

type DeployMessage struct {
	CeramicURL string  `yaml:"CeramicUrl"`
	Models     []Model `yaml:"Models"`
}

type Model struct {
	Schema         string
	IsPublicDomain bool
	Encryptable    []string
}

type ModelResult struct {
	StreamID  string
	ModelName string
	Schema    string
}

type ResponseNonce[T any] struct {
	Message string
	Data    T
	Nonce   string
}
