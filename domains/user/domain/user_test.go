package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("Given valid username and nickname when create user then should return user", func(t *testing.T) {
		nickname := &Nickname{value: "testnick"}
		username := &Username{value: "testuser"}

		user := NewUser(nickname, username)

		assert.NotNil(t, user)
		assert.Equal(t, nickname, user.Nickname)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, ROLE_USER, user.Role)
	})
}
