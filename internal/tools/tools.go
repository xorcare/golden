// Copyright (c) 2019-2024 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tools
// +build tools

package tools_test

import (
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
)
