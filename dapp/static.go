package dapp

import (
	_ "embed"
)

var (
	//go:embed content_folder.graphql
	ContentFloder []byte
	//go:embed index_file.graphql
	IndexFile []byte
	//go:embed index_folder.graphql
	IndexFolder []byte
)
