package fromfile

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	// globalAssemblyVersionRegexp is used to find the argument `<AssemblyVersion>2016.7.00</AssemblyVersion>`
	csharpProjectVersionRegexp = regexp.MustCompile(`AssemblyVersion>(\d*|\.)*`)
)

type CsharpProjectVersionReader struct {
}

func (r CsharpProjectVersionReader) String() string {
	return "csharp-project"
}

func (r CsharpProjectVersionReader) SupportedFiles() []string {
	return []string{
		"\\.csproj$",
	}
}

func (r CsharpProjectVersionReader) ReadFileVersion(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	assemblyVersionLine := csharpProjectVersionRegexp.Find(content)
	if len(assemblyVersionLine) == 0 {
		return "", fmt.Errorf("AssemblyVersion not found in file %s", filePath)
	}

	version := string(VersionRegexp.Find(assemblyVersionLine))

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
