package tools

import (
	_ "golang.org/x/lint"
	_ "golang.org/x/tools/imports"
)

//go:generate go install golang.org/x/lint/golint
//go:generate go install golang.org/x/tools/cmd/goimports
