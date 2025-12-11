package jwt

import (
	"encoding/json"
	"testing"
)

func TestGenerateTokenPair(t *testing.T) {
	pair, err := GenerateTokenPair(123)
	if err != nil {
		t.Error(err)
	}

	token, _ := json.Marshal(pair)
	t.Log(string(token))
}
