package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"golang.org/x/mod/modfile"
)

var staticLookup = map[string]string{
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

func versionFromStatic(version string) string {
	return staticLookup[version]
}

func versionFromGitHub(version string, accessKey string) string {
	dat := `{"query": "query { repository(owner:\"golang\", name:\"go\") { refs(refPrefix: \"refs/tags/\", query: \"go%s\", last: 1) { nodes { name } } } }"}`
	reqBody := fmt.Sprintf(dat, version)

	req, err := http.NewRequest(
		"POST", "https://api.github.com/graphql", bytes.NewBuffer([]byte(reqBody)))
	check(err)
	req.Header.Set("Authorization", "token "+accessKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	// fmt.Printf("%s\n", body)

	var r QueryResult
	err = json.Unmarshal(body, &r)
	check(err)
	return r.Data.Repository.Refs.Nodes[0].Name[2:]
}

func fetchLatestVersion(version string) string {
	var latestVersion string
	accessKey := os.Getenv("GH_ACCESS_KEY")
	if accessKey == "" {
		// if true {
		latestVersion = versionFromStatic(version)
		log.Println("Version from static lookup:", latestVersion)
	} else {
		latestVersion = versionFromGitHub(version, accessKey)
		log.Println("Version from GitHub tags:", latestVersion)
	}
	return latestVersion
}

func main() {
	_os := runtime.GOOS
	_arch := runtime.GOARCH
	log.Println("OS:                      ", _os)
	log.Println("Arch:                    ", _arch)

	filename := "go.mod"
	_, err := os.Stat(filename)
	check(err)
	dat, err := ioutil.ReadFile(filename)
	check(err)
	modFile, err := modfile.Parse(filename, dat, nil)
	check(err)
	_version := modFile.Go.Version
	log.Println("Raw Go version:          ", _version)
	if len(strings.Split(_version, ".")) < 3 {
		_version = fetchLatestVersion(_version)
	}

	extensions := map[string]string{
		"linux":   "tar.gz",
		"freebsd": "tar.gz",
		"darwin":  "tar.gz",
		"windows": "zip"}

	fetchName := fmt.Sprintf("go%s.%s-%s.%s", _version, _os, _arch, extensions[_os])
	log.Println("Fetching:                ", fetchName)
}
