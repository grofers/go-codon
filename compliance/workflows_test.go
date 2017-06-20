package main

import (
	"testing"
	"github.com/grofers/go-codon/testing/workflows"
	"github.com/stretchr/testify/assert"
)

func TestDirect(t *testing.T) {
	var_map := map[string]interface{} {}
    result_i := workflows.GetDirect(&var_map)
    result, ok := result_i.(map[string]interface{})
    if !ok {
    	t.FailNow()
    }
    assert.Equal(t, result["body"], "OK", "Response body not OK")
    assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}
