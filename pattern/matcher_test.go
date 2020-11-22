package pattern_test

import (
	"testing"

	"github.com/mylxsw/adanos-scheduler/pattern"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	sample := `{"name": "Tom", "age": 24}`
	{
		rs, err := pattern.StringMatch("true", sample)
		assert.NoError(t, err)
		assert.True(t, rs)
	}
	{
		rs, err := pattern.StringMatch(`Int(JQ(".age")) > 20`, sample)
		assert.NoError(t, err)
		assert.True(t, rs)
	}
	{
		rs, err := pattern.StringMatch(`Int(JQ(".age")) > 25`, sample)
		assert.NoError(t, err)
		assert.False(t, rs)
	}
}

func TestEval(t *testing.T) {
	sample := `{"name": "Tom", "age": 24, "roles": [{"id": 1, "name": "admin"},{"id":2, "name":"editor"}]}`
	{
		rs, err := pattern.StringEval(`JQ(".age")`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "24", rs)
	}
	{
		rs, err := pattern.StringEval(`JQ(".roles[0].name")`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "admin", rs)
	}
}
