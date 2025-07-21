package domain

import (
	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNickname(t *testing.T) {
	t.Run("Given valid nickname when create new nickname then should return true", func(t *testing.T) {
		testNickname := "testNickname"
		nickname, err := NewNickname(testNickname)

		assert.Nil(t, err)
		assert.NotNil(t, nickname)
		assert.Equal(t, nickname.value, testNickname)
	})

	t.Run("Given short nickname when create new nickname then should return false", func(t *testing.T) {
		testNickname := "a"
		nickname, err := NewNickname(testNickname)
		assert.NotNil(t, err)
		assert.Nil(t, nickname)
		assert.ErrorIs(t, err, errorCode.ErrNicknameTooShort)
	})

	t.Run("Given long nickname when create new nickname then should return false", func(t *testing.T) {
		testNickname := "testestestestestestestestesetestestestestestestestestestestestestestest"
		nickname, err := NewNickname(testNickname)
		assert.NotNil(t, err)
		assert.Nil(t, nickname)
		assert.ErrorIs(t, err, errorCode.ErrNicknameTooLong)
	})

	t.Run("Given nickname with invalid characters when create new nickname then should return false", func(t *testing.T) {
		testNickname := "test!"
		nickname, err := NewNickname(testNickname)
		assert.NotNil(t, err)
		assert.Nil(t, nickname)
		assert.ErrorIs(t, err, errorCode.ErrNicknameInvalidCharacters)
	})
}

func TestNickname_ToString(t *testing.T) {
	t.Run("Given valid nickname when create new nickname then should have a string representation of nickname", func(t *testing.T) {
		testNickname := "test"
		nickname, _ := NewNickname(testNickname)

		assert.Equal(t, testNickname, nickname.ToString())
	})

	t.Run("Given valid nickname when create new nickname then should distinguish other nickname", func(t *testing.T) {
		testNickname := "test"
		nickname, _ := NewNickname(testNickname)

		otherNickname := "otherTest"
		assert.NotEqual(t, nickname.ToString(), otherNickname)
	})
}

func TestNickname_Equals(t *testing.T) {
	t.Run("Given same nickname when create two of same nickname then should have equal string representation of nickname", func(t *testing.T) {
		testNickname := "test"
		nickname1, _ := NewNickname(testNickname)
		nickname2, _ := NewNickname(testNickname)

		assert.True(t, nickname1.Equals(nickname2))
	})

	t.Run("Given different nickname when create two of distinguishable nickname then should not have equal string representation of nickname", func(t *testing.T) {
		testNickname1 := "test"
		nickname1, _ := NewNickname(testNickname1)

		testNickname2 := "tset"
		nickname2, _ := NewNickname(testNickname2)

		assert.False(t, nickname1.Equals(nickname2))
	})
}

func TestNickname_IsEmpty(t *testing.T) {
	t.Run("Given empty nickname when check nickname empty then should return true", func(t *testing.T) {
		testNickname := &Nickname{value: ""}
		assert.True(t, testNickname.IsEmpty())
	})

	t.Run("Given nickname not empty when check nickname empty then should return false", func(t *testing.T) {
		testNickname := &Nickname{value: "test"}
		assert.False(t, testNickname.IsEmpty())
	})
}
