package fromfile

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var (
	ErrFileHasNoVersion = errors.New("the file has no version")
)

type Strategy struct {
	Dir      string
	FilePath string
}

func (s Strategy) ReadVersion() (*semver.Version, error) {
	var (
		dir = s.Dir
		err error
	)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	var (
		reader    FileVersionReader
		filePaths []string
	)
	if len(s.FilePath) > 0 {
		filePath := filepath.Join(dir, s.FilePath)
		filePaths = append(filePaths, filePath)
		reader, err = s.getReader()
	} else {
		reader, filePaths, err = s.autoDetect(dir)
	}
	if err != nil {
		return nil, err
	}

	var version string
	for _, filePath := range filePaths {
		log.Logger().Debugf("Reading version from file %s using reader %s", filePath, reader)
		version, err = reader.ReadFileVersion(filePath)
		if errors.Is(err, ErrFileHasNoVersion) {
			log.Logger().Debugf("File %s has no version", filePath)
			continue
		}
		if err != nil {
			return nil, err
		}
		if version != "" {
			break
		}
	}

	if version == "" {
		return nil, fmt.Errorf("could not find version from %s using reader %s", filePaths, reader)
	}

	log.Logger().Debugf("Found version %s", version)
	return semver.NewVersion(version)
}

func (s Strategy) BumpVersion(_ semver.Version) (*semver.Version, error) {
	return s.ReadVersion()
}

func (s Strategy) autoDetect(dir string) (FileVersionReader, []string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("could not list files, in directory %s", dir)
	}

	for _, reader := range fileVersionReaders {
		var filePaths []string
		for _, fileNamePattern := range reader.SupportedFiles() {
			var fileNameExpr = regexp.MustCompile(fileNamePattern)
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				filePath := filepath.Join(dir, file.Name())
				if fileNameExpr.FindString(file.Name()) != "" {
					log.Logger().Debugf("Adding file %s as a candidate to read version using %s reader", filePath, reader.String())
					filePaths = append(filePaths, filePath)
				}
			}
		}
		if len(filePaths) > 0 {
			return reader, filePaths, nil
		}
	}

	return nil, nil, fmt.Errorf("could not find a file to read version from, in directory %s", dir)
}

func (s Strategy) getReader() (FileVersionReader, error) {
	for _, reader := range fileVersionReaders {
		for _, filePattern := range reader.SupportedFiles() {
			var pattenMatcher = regexp.MustCompile(filePattern)
			if pattenMatcher.FindString(path.Base(s.FilePath)) != "" {
				return reader, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find a file version reader for %s", s.FilePath)
}

type FileVersionReader interface {
	ReadFileVersion(filePath string) (string, error)
	SupportedFiles() []string
	String() string
}

// fileVersionReaders is an ordered list of all readers to try
// when auto-detecting the file to use
var fileVersionReaders = []FileVersionReader{
	HelmChartVersionReader{},
	MakefileVersionReader{},
	AutomakeVersionReader{},
	CMakeVersionReader{},
	PythonVersionReader{},
	MavenPOMVersionReader{},
	JsPackageVersionReader{},
	GradleVersionReader{},
	AssemblyVersionReader{},
	CsharpProjectVersionReader{},
}
