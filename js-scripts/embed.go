package jsscripts

import (
	_ "embed"
)

var (
	//go:embed dist/check.js
	CheckSyntax string
	//go:embed dist/deploy.js
	DeployModel string
	//go:embed dist/admin-access.js
	AdminAccess string
	//go:embed dist/indexed-models.js
	IndexedModels string
)
