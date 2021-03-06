package fetch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestVersionFromStatic(t *testing.T) {
	minorVersions := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	testCases := []string{}
	for _, mv := range minorVersions {
		testCases = append(testCases, fmt.Sprintf("1.%d", mv))
	}

	for _, tc := range testCases {
		res := versionFromStatic(tc)
		if len(strings.Split(res, ".")) != 3 {
			t.Errorf("[from static %s] expected semver, got: %s\n", tc, res)
		}
	}
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	// This could actually use the req to form a response, but going with the
	// simplest/quickest test for now
	resBody := `{"data":{"repository":{"refs":{"nodes":[{"name":"go1.14.4"}]}}}}`
	res := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(resBody))}
	return res, nil
}

func TestVersionFromGitHub(t *testing.T) {
	// note: current MockClient always returns "1.14.4"
	Client = &MockClient{}

	testCases := []string{"1.13", "1"}

	for _, tc := range testCases {
		res := versionFromGitHub(tc, "12345")
		if len(strings.Split(res, ".")) != 3 {
			t.Errorf("[from github %s] expected semver, got: %s\n", tc, res)
		}
	}
}

func setupTestMod(t *testing.T, version string) string {
	testDir, err := createTmpDir(".test-tmp")
	check(err)
	modFileText := []byte("module testy\n\ngo " + version)
	testFileName := filepath.Join(testDir, "go.mod")
	err = ioutil.WriteFile(testFileName, modFileText, 0644)
	check(err)

	t.Cleanup(func() {
		os.RemoveAll(testDir)
	})
	return testDir
}

func setupEnvVar(t *testing.T, envVar string) {
	priorValue := os.Getenv(envVar)
	t.Cleanup(func() {
		os.Setenv(envVar, priorValue)
	})
}

func TestGetVersion(t *testing.T) {
	expect := "1.13"
	testDir := setupTestMod(t, expect)
	testFile := filepath.Join(testDir, "go.mod")

	version := getVersionFromModFile(testFile)
	if version != expect {
		t.Errorf("[version from go.mod] expected %s, got %s\n", expect, version)
	}
}

type LatestVersionTestCase struct {
	Name   string
	Key    string
	Short  string
	Expect string
}

func TestFetchLatestVersion(t *testing.T) {
	testCases := []LatestVersionTestCase{
		{
			Name:   "key env var clear",
			Key:    "",
			Short:  "1.14",
			Expect: "1.14.4"},
		{
			Name:   "key env var set",
			Key:    "12345",
			Short:  "1.14",
			Expect: "1.14.4"}}

	// note: current MockClient always returns "1.14.4"
	Client = &MockClient{}

	keyEnvVar := "GH_ACCESS_KEY"
	setupEnvVar(t, keyEnvVar)

	for _, tc := range testCases {
		os.Setenv(keyEnvVar, tc.Key)
		res := fetchLatestVersion(tc.Short)
		if res != tc.Expect {
			t.Errorf("[%s] expected %s, got %s\n", tc.Name, tc.Expect, res)
		}
	}
}

type ArchiveTestCase struct {
	Name    string
	Version string
	Os      string
	Arch    string
	Expect  string
}

func TestBuildArchiveName(t *testing.T) {
	testCases := []ArchiveTestCase{
		{
			Name:    "linux arm64",
			Version: "1.12.17",
			Os:      "linux",
			Arch:    "arm64",
			Expect:  "go1.12.17.linux-arm64.tar.gz"},
		{
			Name:    "windows 386",
			Version: "1.13.12",
			Os:      "windows",
			Arch:    "386",
			Expect:  "go1.13.12.windows-386.zip"},
		{
			Name:    "osx amd64",
			Version: "1.14.1",
			Os:      "darwin",
			Arch:    "amd64",
			Expect:  "go1.14.1.darwin-amd64.tar.gz"}}

	for _, tc := range testCases {
		res := buildArchiveName(tc.Version, tc.Os, tc.Arch)
		if res != tc.Expect {
			t.Errorf("[%s] expected %s, got %s\n", tc.Name, tc.Expect, res)
		}
	}
}

func TestBuildReport(t *testing.T) {
	keyEnvVar := "GH_ACCESS_KEY"
	setupEnvVar(t, keyEnvVar)
	os.Setenv(keyEnvVar, "")

	_version := "1.13"
	td := setupTestMod(t, _version)
	_fullversion := versionFromStatic(_version)

	_os := runtime.GOOS
	_arch := runtime.GOARCH

	_Archive := fmt.Sprintf("go%s.%s-%s", _fullversion, _os, _arch)
	if _os == "windows" {
		_Archive = _Archive + ".zip"
	} else {
		_Archive = _Archive + ".tar.gz"
	}

	r := BuildReport(td)
	if r.Os != _os {
		t.Errorf("[report] expected os %s, got %s\n", _os, r.Os)
	}
	if r.Arch != _arch {
		t.Errorf("[report] expected arch %s, got %s\n", _arch, r.Arch)
	}
	if r.Version != _fullversion {
		t.Errorf("[report] expected version %s, got %s\n", _fullversion, r.Version)
	}
	if r.Archive != _Archive {
		t.Errorf("[report] expected Archive %s, got %s\n", _Archive, r.Archive)
	}
}
