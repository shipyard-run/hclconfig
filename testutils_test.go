package hclconfig

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestProcessesTypes(t *testing.T) {
	vars := map[string]cty.Value{}
	vars["string"] = cty.StringVal("abc")
	vars["number"] = cty.NumberIntVal(23)
	vars["bool"] = cty.BoolVal(true)
	vars["array"] = cty.ListVal(
		[]cty.Value{
			cty.StringVal("abc"),
			cty.StringVal("123"),
		})

	vars["map"] = cty.MapVal(map[string]cty.Value{
		"foo": cty.StringVal("abc"),
	})

	output := ParseVars(vars)

	require.Equal(t, "abc", output["string"])

	num, _ := output["number"].(*big.Float).Int64()
	require.Equal(t, int64(23), num)

	require.True(t, output["bool"].(bool))

	require.Equal(t, "abc", output["array"].([]interface{})[0])
	require.Equal(t, "123", output["array"].([]interface{})[1])

	require.Equal(t, "abc", output["map"].(map[string]interface{})["foo"])
}

// CreateConfigFromStrings is a test helper function that
// parses the given contents strings as HCL and returns a Shipyard Config
func CreateConfigFromStrings(t *testing.T, contents ...string) (*Config, string) {
	//dir := CreateTestFiles(t, contents...)

	//c := resources.NewConfig()
	//err := ParseFolder(dir, c, false, "", false, []string{}, nil, "")
	//require.NoError(t, err)

	//return c, dir

	return nil, ""
}

// ParseVars converts a map[string]cty.Value into map[string]interface
// where the interface are generic go types like string, number, bool, slice, map
func ParseVars(value map[string]cty.Value) map[string]interface{} {
	vars := map[string]interface{}{}

	for k, v := range value {
		vars[k] = castVar(v)
	}

	return vars
}

func castVar(v cty.Value) interface{} {
	if v.Type() == cty.String {
		return v.AsString()
	} else if v.Type() == cty.Bool {
		return v.True()
	} else if v.Type() == cty.Number {
		return v.AsBigFloat()
	} else if v.Type().IsObjectType() || v.Type().IsMapType() {
		return ParseVars(v.AsValueMap())
	} else if v.Type().IsTupleType() || v.Type().IsListType() {
		i := v.ElementIterator()
		vars := []interface{}{}
		for {
			if !i.Next() {
				// cant iterate
				break
			}

			_, value := i.Element()
			vars = append(vars, castVar(value))
		}

		return vars
	}

	return nil
}

// createsTestFiles creates a temporary directory and
// stores temp files into it
// returns directory containing files
// cleanup function
// usage:
// d, cleanup := createTestFiles(t, `cluster "abc" {}`, `docs "bcdf" {}`)
// defer cleanup()
func CreateTestFiles(t *testing.T, contents ...string) string {
	dir := createTempDirectory(t)

	for _, x := range contents {
		createNamedFile(t, dir, "*.hcl", x)
	}

	t.Cleanup(func() {
		removeTestFiles(t, dir)
	})

	return dir
}

// createTestFile creates a hcl file from the given contents
func CreateTestFile(t *testing.T, contents string) string {
	dir := createTempDirectory(t)

	t.Cleanup(func() {
		removeTestFiles(t, dir)
	})

	return createNamedFile(t, dir, "*.hcl", contents)
}

// create a temporary directory
func createTempDirectory(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err)
	}

	return dir
}

func createNamedFile(t *testing.T, dir, name, contents string) string {
	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		t.Fatalf("Error creating temp file %s", err)
	}
	defer f.Close()

	if _, err := f.WriteString(contents); err != nil {
		t.Fatalf("Error writing temp file contents: %s", err)
	}

	return f.Name()
}

// remove test files cleans up any temporary files created
// with createTestFile
func removeTestFiles(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("Unable to remove temporary files %s", err)
	}
}
