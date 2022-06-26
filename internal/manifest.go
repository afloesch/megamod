package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

const ESRBERating ESRBString = "E"
const ESRBEPlusRating ESRBString = "E10+"
const ESRBTeenRating ESRBString = "T"
const ESRBMatureRating ESRBString = "M"
const ESRBAdultRating ESRBString = "AO"

type ESRBString string

type Game struct {
	Executable string `json:"executable" yaml:"executable"`
	Version    string `json:"version" yaml:"version"`
}

type AgeRating struct {
	ESRB ESRBString `json:"esrb" yaml:"esrb"`
}

type Manifest struct {
	Dependency map[string]string `json:"dependency" yaml:"dependency"`
	Directory  string            `json:"directory" yaml:"directory"`
	Game       Game              `json:"game" yaml:"game"`
	Name       string            `json:"name" yaml:"name"`
	AgeRating  AgeRating         `json:"rating" yaml:"rating"`
	Schema     string            `json:"schema" yaml:"schema"`
	URL        string            `json:"url" yaml:"url"`
	Version    string            `json:"version" yaml:"version"`

	repo string
}

func (m *Manifest) FetchDepsManifest() ([]*Manifest, error) {
	var deps []*Manifest

	for _, v := range m.Dependency {
		d, err := FetchManifest(v)
		if err != nil {
			return nil, err
		}
		deps = append(deps, d)
	}

	return deps, nil
}

func FetchManifest(repo string) (*Manifest, error) {
	resp, err := resty.New().R().
		Get(fmt.Sprintf("https://raw.githubusercontent.com/%s", repo))
	if err != nil {
		return nil, err
	}

	return ParseYAMLManifest(resp.Body())
}

func ParseManifestFromFile(fpath string) (*Manifest, error) {
	data, err := os.ReadFile(filepath.Clean(fpath))
	if err != nil {
		return nil, err
	}

	man, err := ParseYAMLManifest(data)
	if err != nil {
		return nil, err
	}

	return man, nil
}

func ParseYAMLManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func ParseJSONManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	err := json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}
