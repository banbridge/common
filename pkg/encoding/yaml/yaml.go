package yaml

import (
	"gopkg.in/yaml.v3"

	"github.com/banbridge/common/pkg/encoding"
)

// Name is the name registered for the yaml codec.
const Name = "yaml"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with yaml.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
