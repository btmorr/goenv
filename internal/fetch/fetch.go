package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/mod/modfile"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Node struct {
	Name string `json:"name"`
}

type Refs struct {
	Nodes []Node `json:"nodes"`
}

type Repository struct {
	Refs Refs `json:"refs"`
}

type Data struct {
	Repository Repository `json:"repository"`
}

type QueryResult struct {
	Data Data `json:"data"`
}

// HTTPClient defines the client inerface for test mocks
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is an HTTPClient, defined here to enable mocking
var Client HTTPClient

// versionFromStatic looks up a short version string ("1.x") in the lookup
// table and returns a full version string ("1.x.x")
func versionFromStatic(version string) string {
	staticLookup := map[string]string{
		"1.14": "1.14.4",
		"1.13": "1.13.12",
		"1.12": "1.12.17",
		"1.11": "1.11.13",
		"1.10": "1.10.8",
		"1.9":  "1.9.7",
		"1.8":  "1.8.7",
		"1.7":  "1.7.6",
		"1.6":  "1.6.4",
		"1.5":  "1.5.4",
		"1.4":  "1.4.3",
		"1.3":  "1.3.3",
		"1.2":  "1.2.2",
		"1.1":  "1.1.2",
		"1.0":  "1.0.3"}
	return staticLookup[version]
}

// githubVersionRequest separates actual request logic from request
// construction and parsing, to make it easier to mock for testing
func githubVersionRequest(body string, accessKey string) []byte {
	req, err := http.NewRequest(
		"POST", "https://api.github.com/graphql", bytes.NewBuffer([]byte(body)))
	check(err)
	req.Header.Set("Authorization", "token "+accessKey)
	resp, err := Client.Do(req)
	check(err)
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	check(err)
	return resBody
}

// versionFromGitHub takes a short version ("1.x") and access key to use for
// querying GitHub's API, and returns the highest full version ("1.x.x") for
// the specified short version
func versionFromGitHub(version string, accessKey string) string {
	dat := `{"query": "query { repository(owner:\"golang\", name:\"go\") { refs(refPrefix: \"refs/tags/\", query: \"go%s\", last: 1) { nodes { name } } } }"}`
	reqBody := fmt.Sprintf(dat, version)

	resBody := githubVersionRequest(reqBody, accessKey)
	// fmt.Printf("%s\n", body)

	var r QueryResult
	err := json.Unmarshal(resBody, &r)
	check(err)
	return r.Data.Repository.Refs.Nodes[0].Name[2:]
}

// fetchLatestVersion translates a short version string into a full semver
// version string. If the GH_ACCESS_KEY environment variable is set, it uses
// GitHub's API (specifically the v4/GraphQL API) to determine the most
// recent version matching the short version provided. If the environment
// variable is not set, it uses a static lookup (currently in code)
func fetchLatestVersion(version string) string {
	var latestVersion string
	accessKey := os.Getenv("GH_ACCESS_KEY")
	if accessKey == "" {
		latestVersion = versionFromStatic(version)
	} else {
		latestVersion = versionFromGitHub(version, accessKey)
	}
	return latestVersion
}

// getVersionFromModFile reads the go.mod file in the working directory and
// returns the go version string
func getVersionFromModFile(filename string) string {
	_, err := os.Stat(filename)
	check(err)
	dat, err := ioutil.ReadFile(filename)
	check(err)
	modFile, err := modfile.Parse(filename, dat, nil)
	check(err)
	return modFile.Go.Version
}

// buildArchiveName constructs the name of the [tar|zip]file for go for the
// given version and the current OS and architecture
func buildArchiveName(version, os, arch string) string {
	extensions := map[string]string{
		"linux":   "tar.gz",
		"freebsd": "tar.gz",
		"darwin":  "tar.gz",
		"windows": "zip"}

	return fmt.Sprintf("go%s.%s-%s.%s", version, os, arch, extensions[os])
}

// ensureDirectory creates the directory if it does not exist (fail if path
// exists and is not a directory)
func ensureDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		var fileMode os.FileMode
		fileMode = os.ModeDir | 0775
		mdErr := os.MkdirAll(path, fileMode)
		return mdErr
	}
	if err != nil {
		return err
	}
	file, _ := os.Stat(path)
	if !file.IsDir() {
		return errors.New(path + " is not a directory")
	}
	return nil
}

// createTmpDir makes a directory named `dir` in the system temporary directory
// and returns the full path to the created directory
func createTmpDir(dir string) (string, error) {
	tmpDir := os.TempDir()
	dataDir := filepath.Join(tmpDir, dir)
	err := ensureDirectory(dataDir)
	return dataDir, err
}

// getArchive retreives the selected file from the Golang download server, and
// returns the path where the downloaded file is stored in a temp directory
func getArchive(name string) string {
	tempDir, err := createTmpDir(".goenv-tmp")
	check(err)

	prefix := "https://dl.google.com/go/"
	req, err := http.NewRequest("GET", prefix+name, nil)
	check(err)
	resp, err := Client.Do(req)
	check(err)
	defer resp.Body.Close()

	tempFile := filepath.Join(tempDir, name)
	out, err := os.Create(tempFile)
	check(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)
	return tempFile
}

func getShimPath() string {
	user, err := user.Current()
	check(err)
	homeDir := user.HomeDir
	return filepath.Join(homeDir, ".goenv-shims")
}

// installShim unzips the Archive into a hidden directory under the user home
// directory, and modifies the PATH environment variable in the current process
// so that this version of Go is found before the system version, then returns
// the shim path
func installShim(archive string, shimPath string) string {
	// note: move this to either a shell script or a different command
	err := ensureDirectory(shimPath)
	check(err)

	currentPath := os.Getenv("PATH")
	if !strings.Contains(currentPath, shimPath) {
		os.Setenv("PATH", fmt.Sprintf("%s:%s", shimPath, currentPath))
	}
	return shimPath
}

type Report struct {
	Os      string
	Arch    string
	Version string
	Archive string
}

// BuildReport is the primary runner for the application. It compiles all
// information needed to select a golang install package, based on current os,
// architecture, and the go version from go.mod. Then it downloads the selected
// file and installs it into a shim directory, modifies the PATH for the
// current process, and prints a message for the user to suggest a PATH export
// to add to the shell profile
func BuildReport(pwd string) Report {
	_os := runtime.GOOS
	_arch := runtime.GOARCH
	_version := getVersionFromModFile(filepath.Join(pwd, "go.mod"))
	if len(strings.Split(_version, ".")) < 3 {
		_version = fetchLatestVersion(_version)
	}
	fetchName := buildArchiveName(_version, _os, _arch)

	return Report{
		Os:      _os,
		Arch:    _arch,
		Version: _version,
		Archive: fetchName}
}

func FetchArchive(archiveName string) string {
	return getArchive(archiveName)
}

func init() {
	Client = &http.Client{}
}
