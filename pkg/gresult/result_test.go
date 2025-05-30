package gresult

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {

	re := OK("123")

	data, _ := json.Marshal(re)

	t.Log(string(data))

	h := R[string]{}

	_ = json.Unmarshal(data, &h)
	t.Log(h)

	t.Log(h.typ())
}
