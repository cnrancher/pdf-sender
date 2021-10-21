package apis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeader(t *testing.T) {
	totest := tableWithHeader(0, 0, []string{"test1", "test2"})
	target := map[string]string{
		"A1": "test1",
		"B1": "test2",
	}
	assert.Equal(t, target, totest)

	target["A2"] = "data1"
	target["B2"] = "data2"

	setRow(totest, 0, 1, []string{"data1", "data2"})
	assert.Equal(t, target, totest)
}
