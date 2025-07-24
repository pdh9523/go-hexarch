package domain

import (
	"github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("Given valid username and nickname when create user then should return user", func(t *testing.T) {
		nickname := "testnick"
		username := "testuser"

		user, err := NewUser(nickname, username)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, nickname, user.Nickname)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, ROLE_USER, user.Role)
	})

	t.Run("Given invalid username when create user then should return error", func(t *testing.T) {
		nickname := "testnick"
		username := "TESTUSER"

		user, err := NewUser(nickname, username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errorCode.ErrUsernameInvalidCharacters)
	})

	t.Run("Given invalid nickname when create user then should return error", func(t *testing.T) {
		nickname := "test!"
		username := "testuser"

		user, err := NewUser(nickname, username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errorCode.ErrNicknameInvalidCharacters)
	})
}

// Username 검증 테스트
func TestValidateUsername(t *testing.T) {
	t.Run("Given valid username when validate then should return no error", func(t *testing.T) {
		validUsernames := []string{
			"testuser",
			"test123",
			"user1",
			"abc123def",
			"test",
		}

		for _, username := range validUsernames {
			err := validateUsername(username)
			assert.NoError(t, err, "Username should be valid: %s", username)
		}
	})

	t.Run("Given too short username when validate then should return error", func(t *testing.T) {
		shortUsernames := []string{
			"",
			"a",
		}

		for _, username := range shortUsernames {
			err := validateUsername(username)
			assert.ErrorIs(t, err, errorCode.ErrUsernameTooShort, "Username should be too short: %s", username)
		}
	})

	t.Run("Given too long username when validate then should return error", func(t *testing.T) {
		longUsername := strings.Repeat("a", 17) // 17자
		err := validateUsername(longUsername)
		assert.ErrorIs(t, err, errorCode.ErrUsernameTooLong)

		veryLongUsername := strings.Repeat("a", 40) // 40자
		err = validateUsername(veryLongUsername)
		assert.ErrorIs(t, err, errorCode.ErrUsernameTooLong)
	})

	t.Run("Given username with invalid characters when validate then should return error", func(t *testing.T) {
		invalidUsernames := []string{
			"TESTUSER",  // 대문자
			"test!",     // 특수문자
			"test user", // 공백
			"테스트유저",     // 한글
			"test@user", // @ 문자
			"test-user", // 하이픈
			"test_user", // 언더스코어
			"test.user", // 점
			"%test%",    // 퍼센트
			"  ",        // 공백만
		}

		for _, username := range invalidUsernames {
			err := validateUsername(username)
			assert.ErrorIs(t, err, errorCode.ErrUsernameInvalidCharacters, "Username should be invalid: %s", username)
		}
	})
}

// Nickname 검증 테스트
func TestValidateNickname(t *testing.T) {
	t.Run("Given valid nickname when validate then should return no error", func(t *testing.T) {
		validNicknames := []string{
			"testnick",
			"TestNick",
			"테스트닉네임",
			"test123",
			"TEST",
			"테스트",
			"Nick123",
			"한글English123",
		}

		for _, nickname := range validNicknames {
			err := ValidateNickname(nickname)
			assert.NoError(t, err, "Nickname should be valid: %s", nickname)
		}
	})

	t.Run("Given too short nickname when validate then should return error", func(t *testing.T) {
		shortNicknames := []string{
			"",
			"a",
			"한",
		}

		for _, nickname := range shortNicknames {
			err := ValidateNickname(nickname)
			assert.ErrorIs(t, err, errorCode.ErrNicknameTooShort, "Nickname should be too short: %s", nickname)
		}
	})

	t.Run("Given too long nickname when validate then should return error", func(t *testing.T) {
		longNickname := strings.Repeat("a", 21) // 21자
		err := ValidateNickname(longNickname)
		assert.ErrorIs(t, err, errorCode.ErrNicknameTooLong)

		longKoreanNickname := strings.Repeat("한", 21) // 21자 한글
		err = ValidateNickname(longKoreanNickname)
		assert.ErrorIs(t, err, errorCode.ErrNicknameTooLong)
	})

	t.Run("Given nickname with invalid characters when validate then should return error", func(t *testing.T) {
		invalidNicknames := []string{
			"test!",     // 특수문자
			"test@nick", // @ 문자
			"test nick", // 공백
			"test-nick", // 하이픈
			"test_nick", // 언더스코어
			"test.nick", // 점
			"test#nick", // 해시
			"test$nick", // 달러
			"test%nick", // 퍼센트
			"test&nick", // 앰퍼샌드
		}

		for _, nickname := range invalidNicknames {
			err := ValidateNickname(nickname)
			assert.ErrorIs(t, err, errorCode.ErrNicknameInvalidCharacters, "Nickname should be invalid: %s", nickname)
		}
	})

	t.Run("Given UTF-8 nickname when validate then should count correctly", func(t *testing.T) {
		koreanNickname := strings.Repeat("한", 20)
		err := ValidateNickname(koreanNickname)
		assert.NoError(t, err)

		shortKoreanNickname := "한글"
		err = ValidateNickname(shortKoreanNickname)
		assert.NoError(t, err)

		mixedNickname := "한글Test123" // 9자
		err = ValidateNickname(mixedNickname)
		assert.NoError(t, err)
	})
}
