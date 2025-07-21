package domain

import (
	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUsername(t *testing.T) {
	t.Run("Given valid username when create username then should return user", func(t *testing.T) {
		testUsername := "testusername"
		username, err := NewUsername(testUsername)

		assert.NoError(t, err)
		assert.NotNil(t, username)
	})

	t.Run("Given short username when create username then should return error", func(t *testing.T) {
		testUsername := "1"
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameTooShort)
	})

	t.Run("Given long username when create username then should return error", func(t *testing.T) {
		testUsername := "1234567890123456789012345678901234576890"
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameTooLong)
	})

	t.Run("Given username with no length should return false", func(t *testing.T) {
		testUsername := ""
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameTooShort)
	})

	t.Run("Given username with invalid Character should return false", func(t *testing.T) {
		testUsername := "%test%"
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameInvalidCharacters)
	})

	t.Run("Given username with blank character should return false", func(t *testing.T) {
		testUsername := "  "
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameInvalidCharacters)
	})

	t.Run("Given username with Korean should return false", func(t *testing.T) {
		testUsername := "테스트유저"
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameInvalidCharacters)
	})

	t.Run("Given upper case character should return false", func(t *testing.T) {
		testUsername := "TESTUSER"
		username, err := NewUsername(testUsername)
		assert.Nil(t, username)
		assert.ErrorIs(t, err, errorCode.ErrUsernameInvalidCharacters)
	})
}

func TestUsername_ToString(t *testing.T) {
	t.Run("Given username when check value with toString should return username", func(t *testing.T) {
		testUsername := "test"
		username := &Username{value: testUsername}
		compare := "test1"

		assert.Equal(t, username.ToString(), testUsername)
		assert.NotEqual(t, username.ToString(), compare)
	})
}

func TestUsername_Equals(t *testing.T) {
	t.Run("Given two of same username when compare should return true", func(t *testing.T) {
		testUsername := "testusername"
		username1, err := NewUsername(testUsername)
		assert.NoError(t, err)
		username2, err := NewUsername(testUsername)
		assert.NoError(t, err)
		assert.Equal(t, true, username1.Equals(username2))
	})

	t.Run("Given two of different username when compare should return false", func(t *testing.T) {
		testUsername1 := &Username{value: "test1"}
		testUsername2 := &Username{value: "test2"}

		assert.NotEqual(t, testUsername1, testUsername2)
		assert.False(t, testUsername1.Equals(testUsername2))
	})
}
