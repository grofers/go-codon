package main

import (
	"testing"
	"github.com/grofers/go-codon/testing/workflows"
	"github.com/stretchr/testify/assert"
)

func deepcopyMap(src map[string]interface{}) map[string]interface{} {
	dst := map[string]interface{} {}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func TestDirect(t *testing.T) {
	var_map := map[string]interface{} {}
	var_map_c := deepcopyMap(var_map)
	result_i := workflows.GetDirect(&var_map_c)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], "OK", "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestSimple(t *testing.T) {
	var_map := map[string]interface{} {
		"string_val": "test_val",
		"boolean_val": true,
		"int_val": int64(10),
		"float_val": float64(1.11),
	}
	var_map_c := deepcopyMap(var_map)
	result_i := workflows.GetSimple(&var_map_c)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], var_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestNoInput(t *testing.T) {
	var_map := map[string]interface{} {}
	var_map_c := deepcopyMap(var_map)
	result_i := workflows.NoInput(&var_map_c)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], var_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestChain(t *testing.T) {
	assert_map := map[string]interface{} {
		"first_val": "success-first",
		"second_val": "success-second",
	}
	var_map := map[string]interface{} {}
	result_i := workflows.GetChain(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestChainError(t *testing.T) {
	assert_map := map[string]interface{} {
		"first_val": "error-first",
		"second_val": "error-second",
	}
	var_map := map[string]interface{} {}
	result_i := workflows.GetChainError(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestChainComplete(t *testing.T) {
	assert_map := map[string]interface{} {
		"first_val": "success-first",
		"second_val": "success-second",
		"third_val": "complete-first",
		"fourth_val": "complete-fourth",
	}
	var_map := map[string]interface{} {}
	result_i := workflows.GetChainComplete(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestChainErrorComplete(t *testing.T) {
	assert_map := map[string]interface{} {
		"first_val": "error-first",
		"second_val": "error-second",
		"third_val": "complete-first",
		"fourth_val": "complete-fourth",
	}
	var_map := map[string]interface{} {}
	result_i := workflows.GetChainErrorComplete(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestPublishSerial(t *testing.T) {
	assert_map := map[string]interface{} {
		"retval": float64(9),
	}
	var_map := map[string]interface{} {}
	result_i := workflows.PublishSerial(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestStartConcurrency(t *testing.T) {
	assert_map := map[string]interface{} {
		"val1": 2,
		"val2": 1,
	}
	var_map := map[string]interface{} {
		"val1": 1,
		"val2": 2,
	}
	result_i := workflows.StartTest(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestComplexOutput(t *testing.T) {
	assert_map := map[string]interface{} {
		"dict": map[string]interface{} {
			"a": "a",
			"b": float64(1),
			"c": true,
			"d": "true",
		},
		"list": []interface{} {float64(10),float64(20),float64(30)},
	}
	var_map := deepcopyMap(assert_map)
	result_i := workflows.ComplexOutput(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestNoPublish(t *testing.T) {
	assert_map := map[string]interface{} {
		"val": 1,
	}
	var_map := deepcopyMap(assert_map)
	result_i := workflows.NoPublish(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}

func TestRecursion(t *testing.T) {
	assert_map := map[string]interface{} {
		"val": float64(10),
	}
	var_map := map[string]interface{} {
		"val": 1,
	}
	result_i := workflows.Recursion(&var_map)
	result, ok := result_i.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	assert.Equal(t, result["body"], assert_map, "Response body not OK")
	assert.Equal(t, result["status_code"], float64(200), "Response code not 200")
}
