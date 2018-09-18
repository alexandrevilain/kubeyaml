package kubernetes

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Input is the top level YAML document the system receives.
type Input struct {
	Kind       string
	APIVersion string
	Data       map[string]interface{}
}

// Loader defines a struct that can read data from a stream into an internal type.
type Loader struct{}

// NewLoader returns a Loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load reads the input and returns the internal type representing the top level document
// that is properly cleaned.
func (l *Loader) Load(reader io.Reader) (*Input, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read incoming reader: %v", err)
	}

	incoming := map[string]interface{}{}
	if err := yaml.Unmarshal(b, incoming); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml with error %v", err)
	}

	val, ok := incoming["apiVersion"]
	if !ok {
		return nil, NewRequiredKeyNotFoundError("apiVersion")
	}

	apiVersion, ok := val.(string)
	if !ok {
		return nil, NewUnknownTypeError(val)
	}

	val, ok = incoming["kind"]
	if !ok {
		return nil, NewRequiredKeyNotFoundError("kind")
	}

	kind, ok := val.(string)
	if !ok {
		return nil, NewUnknownTypeError(val)
	}

	delete(incoming, "apiVersion")
	delete(incoming, "kind")
	return &Input{
		APIVersion: apiVersion,
		Kind:       kind,
		Data:       incoming,
	}, nil
}
