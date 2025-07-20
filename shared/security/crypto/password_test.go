package crypto

import (
	"testing"
)

func TestArgon2PasswordHasher_HashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		password := "testPassword123!"

		// When
		hash, err := hasher.HashPassword(password)

		// Then
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		if hash == "" {
			t.Error("Hash should not be empty")
		}
		if len(hash) == 0 {
			t.Error("Hash length should be greater than 0")
		}
	})

	t.Run("should generate different hashes for same password", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		password := "testPassword123!"

		// When
		hash1, err1 := hasher.HashPassword(password)
		hash2, err2 := hasher.HashPassword(password)

		// Then
		if err1 != nil {
			t.Fatalf("Failed to hash password first time: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("Failed to hash password second time: %v", err2)
		}
		if hash1 == hash2 {
			t.Error("Two hashes of the same password should be different due to salt")
		}
	})
}

func TestArgon2PasswordHasher_VerifyPassword(t *testing.T) {
	t.Run("should verify correct password successfully", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		password := "testPassword123!"
		hash, _ := hasher.HashPassword(password)

		// When
		err := hasher.VerifyPassword(hash, password)

		// Then
		if err != nil {
			t.Errorf("Password verification should succeed: %v", err)
		}
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		password := "testPassword123!"
		wrongPassword := "wrongPassword"
		hash, _ := hasher.HashPassword(password)

		// When
		err := hasher.VerifyPassword(hash, wrongPassword)

		// Then
		if err == nil {
			t.Error("Password verification should fail for wrong password")
		}
	})

	t.Run("should fail with invalid hash format", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		invalidHash := "invalid_hash"
		password := "password"

		// When
		err := hasher.VerifyPassword(invalidHash, password)

		// Then
		if err == nil {
			t.Error("Should fail with invalid hash format")
		}
	})

	t.Run("should fail with empty hash", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		emptyHash := ""
		password := "password"

		// When
		err := hasher.VerifyPassword(emptyHash, password)

		// Then
		if err == nil {
			t.Error("Should fail with empty hash")
		}
	})

	t.Run("should fail with empty password", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		password := "testPassword123!"
		emptyPassword := ""
		hash, _ := hasher.HashPassword(password)

		// When
		err := hasher.VerifyPassword(hash, emptyPassword)

		// Then
		if err == nil {
			t.Error("Should fail with empty password")
		}
	})
}

func TestArgon2PasswordHasher_EdgeCases(t *testing.T) {
	t.Run("should handle long password", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		longPassword := "a" + string(make([]byte, 1000)) + "very_long_password_123!"

		// When
		hash, err := hasher.HashPassword(longPassword)

		// Then
		if err != nil {
			t.Fatalf("Failed to hash long password: %v", err)
		}
		if hash == "" {
			t.Error("Hash should not be empty for long password")
		}
	})

	t.Run("should handle special characters", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		specialPassword := "!@#$%^&*()_+{}|:<>?[]\\;'\",./"

		// When
		hash, err := hasher.HashPassword(specialPassword)

		// Then
		if err != nil {
			t.Fatalf("Failed to hash password with special characters: %v", err)
		}
		if hash == "" {
			t.Error("Hash should not be empty for special characters")
		}

		// Verify the password
		err = hasher.VerifyPassword(hash, specialPassword)
		if err != nil {
			t.Errorf("Password verification should succeed for special characters: %v", err)
		}
	})

	t.Run("should handle unicode characters", func(t *testing.T) {
		// Given
		hasher := NewArgon2PasswordHasher()
		unicodePassword := "비밀번호123!🔒"

		// When
		hash, err := hasher.HashPassword(unicodePassword)

		// Then
		if err != nil {
			t.Fatalf("Failed to hash unicode password: %v", err)
		}
		if hash == "" {
			t.Error("Hash should not be empty for unicode password")
		}

		// Verify the password
		err = hasher.VerifyPassword(hash, unicodePassword)
		if err != nil {
			t.Errorf("Password verification should succeed for unicode characters: %v", err)
		}
	})
}

func BenchmarkArgon2PasswordHasher_HashPassword(b *testing.B) {
	// Given
	hasher := NewArgon2PasswordHasher()
	password := "testPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// When
		_, err := hasher.HashPassword(password)

		// Then
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArgon2PasswordHasher_VerifyPassword(b *testing.B) {
	// Given
	hasher := NewArgon2PasswordHasher()
	password := "testPassword123!"
	hash, err := hasher.HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// When
		err := hasher.VerifyPassword(hash, password)

		// Then
		if err != nil {
			b.Fatal(err)
		}
	}
}
