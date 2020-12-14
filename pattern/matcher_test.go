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
	{
		rs, err := pattern.StringEval(`CtxJSONInt("roles.0.id")`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "1", rs)
	}
	{
		rs, err := pattern.StringEval(`CtxJSONInt("age")`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "24", rs)
	}
}

func TestJSON(t *testing.T) {
	sample := `[{"ID": 123, "Name": "管宜尧"},{"ID": 124, "Name": "李逍遥"}]`
	{
		rs, err := pattern.StringEval(`len(CtxJSONArray(""))`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "2", rs)
	}
	{
		rs, err := pattern.StringEval(`CtxJSONArray(".ID")`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "[123 124]", rs)
	}
	{
		rs, err := pattern.StringEval(`map(CtxJSONStrArray(".Name"), {# + '*'})`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "[管宜尧* 李逍遥*]", rs)
	}
	{
		rs, err := pattern.StringEval(`any(CtxJSONIntArray(".ID"), {# == 124})`, sample)
		assert.NoError(t, err)
		assert.Equal(t, "true", rs)
	}
}
