package hclconfig

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

type getterCall struct {
	src     string
	dest    string
	working string
}

func setupMockGetter(t *testing.T, err error) (Getter, *[]getterCall) {
	calls := &[]getterCall{}

	g := &GoGetter{
		get: func(src, dest, working string) error {
			*calls = append(*calls, getterCall{
				src:     src,
				dest:    dest,
				working: working,
			})

			return err
		},
	}

	return g, calls
}

func TestDoesNothingWhenFolderExistsAndIgnoreCacheFalse(t *testing.T) {
	dest := t.TempDir()
	downloadPath := path.Join(dest, "github.com_test")
	os.MkdirAll(downloadPath, os.ModePerm)

	g, calls := setupMockGetter(t, nil)

	_, err := g.Get("github.com/test", dest, false)
	require.NoError(t, err)

	require.Len(t, *calls, 0)
}

func TestCallsGetWhenFolderExistsAndIgnoreCacheTrue(t *testing.T) {
	dest := t.TempDir()

	g, calls := setupMockGetter(t, nil)

	_, err := g.Get("github.com/test", dest, true)
	require.NoError(t, err)

	require.Len(t, *calls, 1)
}

func TestCallsGetWithURLEncodedOutputFolder(t *testing.T) {
	g, calls := setupMockGetter(t, nil)

	_, err := g.Get("github.com/shipyard-run/hclconfig?ref=7271da1cd14778d3762304954d7061cc753da204", "/mycache", false)
	require.NoError(t, err)

	require.Len(t, *calls, 1)

	require.Equal(t, "/mycache/github.com_shipyard-run_hclconfig_ref=7271da1cd14778d3762304954d7061cc753da204", (*calls)[0].dest)
}