package golden

import (
	"fmt"
)

// interceptor need to intercept the output of logs from testify.assert.
type interceptor string

func (i *interceptor) Errorf(format string, args ...interface{}) {
	*i += interceptor(fmt.Sprintf(format, args...))
}

func (i interceptor) String() string {
	return string(i)
}