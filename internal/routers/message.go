package routers

type Operation string

const (
	OperationCreate            Operation = "I want to create a DataverseOS app."
	OperationUpdate            Operation = "I want to update a DataverseOS app."
	OperationSetExternalModels Operation = "I want to set external models to a DataverseOS app."
)

type CreateMessage struct {
	Operation Operation `yaml:"Operation"`
	Slug      string    `yaml:"Slug"`
	Ceramic   *string   `yaml:"Ceramic"`
	Models    []string  `yaml:"Models"`
}

type SetExternalModelsMessage struct {
	Operation Operation `yaml:"Operation"`
	Slug      string    `yaml:"Slug"`
	Schema    string    `yaml:"Schema"`
}

type ResponseNonce[T any] struct {
	Message string
	Data    T
	Nonce   string
}
