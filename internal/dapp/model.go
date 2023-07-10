package dapp

type DeployMessage struct {
	Models []StreamModel `yaml:"models"`
}

type StreamModel struct {
	Schema         string   `yaml:"schema"`
	IsPublicDomain bool     `yaml:"isPublicDomain"`
	Encryptable    []string `yaml:"encryptable"`
}

type ModelResult struct {
	StreamID string
	Schema   string
}
