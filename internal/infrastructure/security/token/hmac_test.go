package token

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHMACHasher_Hash(t *testing.T) {
	hasher := NewHMACHasher("secret")

	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "пустой токен",
			token: "",
			want:  "f9e66e179b6747ae54108f82f8ade8b3c25d76fd30afde6c395822c530196169",
		},
		{
			name:  "обычный токен",
			token: "raw-token",
			want:  "f88f9081d43f807af1f49c0b2c17b6540e1b3b773c70c8133df3e024bd794f1a",
		},
		{
			name:  "другой токен",
			token: "another-token",
			want:  "cbed33de21d76f54e68479cbdad158441e50d1de76dacb71e68a0aa52c701ff0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := hasher.Hash(tt.token)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestHMACHasher_Generate(t *testing.T) {
	t.Parallel()

	hasher := NewHMACHasher("secret")

	t.Run("успешная генерация токена", func(t *testing.T) {
		t.Parallel()

		token, err := hasher.Generate()

		require.NoError(t, err)
		require.Len(t, token, 43)

		_, err = base64.RawURLEncoding.DecodeString(token)
		require.NoError(t, err)
	})

	t.Run("каждый вызов возвращает новый токен", func(t *testing.T) {
		t.Parallel()

		first, err := hasher.Generate()
		require.NoError(t, err)

		second, err := hasher.Generate()
		require.NoError(t, err)

		require.NotEqual(t, first, second)
	})
}
