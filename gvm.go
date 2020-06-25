package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

// buildTarfileName constructs the name of the tarfile for go for the given
// version and the current OS and architecture
func buildTarfileName(version, os, arch string) string {
	extensions := map[string]string{
		"linux":   "tar.gz",
		"freebsd": "tar.gz",
		"darwin":  "tar.gz",
		"windows": "zip"}

	return fmt.Sprintf("go%s.%s-%s.%s", version, os, arch, extensions[os])
}

type Report struct {
	Os      string
	Arch    string
	Version string
	Tarfile string
}

// buildReport compiles all information needed to select a golang install
// package, based on current os, architecture, and the go version from go.mod
func buildReport(pwd string) Report {
	_os := runtime.GOOS
	_arch := runtime.GOARCH
	_version := getVersionFromModFile(filepath.Join(pwd, "go.mod"))
	if len(strings.Split(_version, ".")) < 3 {
		_version = fetchLatestVersion(_version)
	}
	fetchName := buildTarfileName(_version, _os, _arch)

	return Report{
		Os:      _os,
		Arch:    _arch,
		Version: _version,
		Tarfile: fetchName}
}

func init() {
	Client = &http.Client{}
}

func main() {
	fmt.Printf("%+v\n", buildReport("."))
}
