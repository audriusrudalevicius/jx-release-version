package fromfile

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var (
	// globalAssemblyVersionRegexp is used to find the argument `AssemblyVersion("2016.7.00.00")`
	globalAssemblyVersionRegexp = regexp.MustCompile(`AssemblyVersion\("(\d*|\.)*"\)`)
	// VersionRegexp is used to find the version `2016.7.00.00`
	VersionRegexp = regexp.MustCompile(`([0-9]+|\.+)+`)
)

type AssemblyVersionReader struct {
}

func (r AssemblyVersionReader) String() string {
	return "csharp"
}

func (r AssemblyVersionReader) SupportedFiles() []string {
	return []string{
		"GlobalAssemblyInfo.cs",
	}
}

func (r AssemblyVersionReader) ReadFileVersion(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	assemblyVersionLine := globalAssemblyVersionRegexp.Find(content)
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
