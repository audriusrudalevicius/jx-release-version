package fromfile

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestGlobalAssemblyVersionReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		filePath         string
		expected         string
		expectedErrorMsg string
	}{
		{
			name:     "standard file",
			filePath: "GlobalAssemblyInfo.cs",
			expected: "2016.7.0",
		},
		{
			name:             "file does not exists",
			filePath:         "does-not-exists.yaml",
			expectedErrorMsg: "open testdata/does-not-exists.yaml: no such file or directory",
		},
		{
			name:             "invalid file",
			filePath:         "Chart.yaml",
			expectedErrorMsg: "AssemblyVersion not found in file testdata/Chart.yaml",
		},
	}

	reader := AssemblyVersionReader{}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := reader.ReadFileVersion(filepath.Join("testdata", test.filePath))
			if len(test.expectedErrorMsg) > 0 {
				require.EqualError(t, err, test.expectedErrorMsg)
				assert.Empty(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
