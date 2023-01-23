package maps

import (
	"github.com/halilylm/microservice/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopy(t *testing.T) {
	m := map[int]int{1: 2, 3: 4}
	c := Copy(m)
	test.AssertDeepEqual(t, m, c)
	c[1] = 3
	assert.Equal(t, m[1], 2)
	assert.Equal(t, c[1], 3)
}
