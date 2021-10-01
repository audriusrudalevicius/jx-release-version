package fromfile

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type VersionFileReader struct {
}

func (r VersionFileReader) String() string {
	return "gradle"
}

func (r VersionFileReader) SupportedFiles() []string {
	return []string{
		"^VERSION$",
	}
}

func (r VersionFileReader) ReadFileVersion(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	version := string(VersionRegexp.Find(content))

	parts := strings.Split(version, ".")
	if len(parts) < 3 {
		return "", fmt.Errorf("version value not found in file %s", filePath)
	}

	majorVer, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return "", err
	}
	minorVer, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return "", err
	}
	patchVer, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d", majorVer, minorVer, patchVer), nil
}
