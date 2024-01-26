// Copyright (c) 2019-2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

var _ assert.TestingT = new(interceptor)
var _ fmt.Stringer = interceptor("")

// interceptor need to intercept the output of logs from testify.assert.
type interceptor string

func (i *interceptor) Errorf(format string, args ...interface{}) {
	*i += interceptor(fmt.Sprintf(format, args...))
}

func (i interceptor) String() string {
	return string(i)
}
