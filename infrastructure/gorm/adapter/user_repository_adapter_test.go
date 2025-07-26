package adapter

import (
	"context"
	"github.com/pdh9523/go-hexarch/domains/user/domain"
	userError "github.com/pdh9523/go-hexarch/domains/user/domain/errorCode"
	"github.com/pdh9523/go-hexarch/infrastructure/gorm/mapper"
	"github.com/pdh9523/go-hexarch/infrastructure/gorm/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)
	return db
}

func createTestDomainUser(nickname, username string) *domain.User {
	user, _ := domain.NewUser(nickname, username, "password")
	return user
}

func createTestModelUser(db *gorm.DB, nickname, username string) *model.User {
	user := &model.User{
		Username: username,
		Nickname: nickname,
		Role:     string(domain.ROLE_USER),
	}
	db.Create(user)
	return user
}

func TestUserRepositoryAdapter_Save(t *testing.T) {
	t.Run("Given valid domain user when save then sohuld create model user successfully", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		nickname := "test"
		username := "test"

		domainUser := createTestDomainUser(nickname, username)
		modelUser, err := repo.Save(context.Background(), domainUser)

		assert.NoError(t, err)
		var count int64
		db.Model(&model.User{}).Where("id = ?", modelUser.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Given duplicate nickname when save then should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)

		nickname := "test"
		username1 := "test1"
		username2 := "test2"

		_ = createTestModelUser(db, nickname, username1)

		domainUser2 := createTestDomainUser(nickname, username2)
		modelUser2, err := repo.Save(context.Background(), domainUser2)

		assert.Error(t, err)
		assert.Nil(t, modelUser2)
	})

	t.Run("Given duplicate username when save then should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)

		nickname1 := "test1"
		nickname2 := "test2"
		username := "test"

		_ = createTestModelUser(db, nickname1, username)
		modelUser2 := createTestDomainUser(nickname2, username)
		modelUser2, err := repo.Save(context.Background(), modelUser2)

		assert.Error(t, err)
		assert.Nil(t, modelUser2)
	})

	t.Run("Given two of valid domain users when save then their ID should be different", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)

		nickname1 := "test1"
		username1 := "test1"
		nickname2 := "test2"
		username2 := "test2"

		domainUser1 := createTestDomainUser(nickname1, username1)
		modelUser1, err := repo.Save(context.Background(), domainUser1)

		assert.NoError(t, err)
		var count int64
		db.Model(&model.User{}).Where("id = ?", modelUser1.ID).Count(&count)
		assert.Equal(t, int64(1), count)

		domainUser2 := createTestDomainUser(nickname2, username2)
		modelUser2, err := repo.Save(context.Background(), domainUser2)

		assert.NoError(t, err)
		db.Model(&model.User{}).Where("id = ?", modelUser2.ID).Count(&count)
		assert.Equal(t, int64(1), count)

		assert.NotEqual(t, modelUser1.ID, modelUser2.ID)
	})
}

func TestUserRepositoryAdapter_FindByID(t *testing.T) {
	t.Run("Given existing user when find by ID then should return user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)

		username := "test"
		nickname := "test"

		userModel := createTestModelUser(db, nickname, username)
		foundUser, err := repo.FindByID(context.Background(), userModel.ID)

		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, userModel.ID, foundUser.ID)
		assert.Equal(t, userModel.Username, foundUser.Username)
		assert.Equal(t, userModel.Nickname, foundUser.Nickname)
	})

	t.Run("Given non-existing user when find by ID then should return nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)

		foundUser, err := repo.FindByID(context.Background(), "non-existing")

		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("Given blank id when find by ID then should return nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")

		foundUser, err := repo.FindByID(context.Background(), "")

		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("Given SQL Injection attempt when find by ID then should return nil and prevent injection", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)

		user1 := createTestModelUser(db, "test1", "test1")
		_ = createTestModelUser(db, "test2", "test2")

		injectionAttempts := []string{
			"' OR '1'='1",
			"'; DROP TABLE users; --",
			"' UNION SELECT * FROM users --",
			"1' OR 1=1 --",
			"'; DELETE FROM users WHERE '1'='1",
		}

		for _, injection := range injectionAttempts {
			foundUser, err := repo.FindByID(context.Background(), injection)
			assert.NoError(t, err)
			assert.Nil(t, foundUser)
		}

		var count int64
		db.Model(&model.User{}).Count(&count)
		assert.Equal(t, int64(2), count)

		validUser, err := repo.FindByID(context.Background(), user1.ID)
		assert.NoError(t, err)
		assert.NotNil(t, validUser)
		assert.Equal(t, user1.ID, validUser.ID)
	})
}

func TestUserRepositoryAdapter_FindByNickname(t *testing.T) {
	t.Run("Given existing user when find by Nickname then should return user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		user1 := createTestModelUser(db, "test", "test")

		foundUser, err := repo.FindByNickname(context.Background(), user1.Nickname)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user1.ID, foundUser.ID)
	})

	t.Run("Given non-existing user when find by Nickname then should return nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")

		foundUser, err := repo.FindByNickname(context.Background(), "other")
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("Given blank nickname when find by Nickname then should return nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		foundUser, err := repo.FindByNickname(context.Background(), "")
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserRepositoryAdapter_FindByUsername(t *testing.T) {
	t.Run("Given existing user when find by Username then should return user", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		user1 := createTestModelUser(db, "test", "test")
		foundUser, err := repo.FindByUsername(context.Background(), user1.Username)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user1.ID, foundUser.ID)
	})

	t.Run("Given non-existing user when find by Username then should return nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		foundUser, err := repo.FindByUsername(context.Background(), "other")
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})

	t.Run("Given blank username when find by Username then should return nil", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		foundUser, err := repo.FindByUsername(context.Background(), "")
		assert.NoError(t, err)
		assert.Nil(t, foundUser)
	})
}

func TestUserRepositoryAdapter_ExistsByNickname(t *testing.T) {
	t.Run("Given existing nickname when check nickname existence then should return true", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		isFound, err := repo.ExistsByNickname(context.Background(), "test")

		assert.NoError(t, err)
		assert.True(t, isFound)
	})

	t.Run("Given non-existing nickname when check nickname then should return false", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		isFound, err := repo.ExistsByNickname(context.Background(), "other")
		assert.NoError(t, err)
		assert.False(t, isFound)
	})
}

func TestUserRepositoryAdapter_ExistsByUsername(t *testing.T) {
	t.Run("Given existing user when check username existence then should return true", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		isFound, err := repo.ExistsByUsername(context.Background(), "test")
		assert.NoError(t, err)
		assert.True(t, isFound)
	})

	t.Run("Given non-existing user when check username existence then should return false", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		isFound, err := repo.ExistsByUsername(context.Background(), "other")
		assert.NoError(t, err)
		assert.False(t, isFound)
	})
}

func TestUserRepositoryAdapter_Update(t *testing.T) {
	t.Run("Given existing user when update then should change existing user's data", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		userModel := createTestModelUser(db, "test", "test")

		userDomain1, err := mapper.ToDomainUser(userModel)
		assert.NoError(t, err)
		userDomain1.Username = "test1"

		userDomain2, err := repo.Update(context.Background(), userDomain1)

		assert.NoError(t, err)
		assert.NotNil(t, userDomain2)
		assert.Equal(t, userDomain1.ID, userDomain2.ID)
		assert.Equal(t, userDomain1.Nickname, userDomain2.Nickname)
		assert.Equal(t, userDomain1.Role, userDomain2.Role)
		assert.Equal(t, userDomain1.Username, userDomain2.Username)
	})

	t.Run("Given non-existing user when update then should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")
		userDomain := createTestDomainUser("test1", "test1")

		updatedUser, err := repo.Update(context.Background(), userDomain)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.ErrorIs(t, err, userError.ErrUserNotFound)
	})
}

func TestUserRepositoryAdapter_DeleteByID(t *testing.T) {
	t.Run("Given existing user when delete then should delete existing user's data", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		userModel := createTestModelUser(db, "test", "test")

		err := repo.DeleteByID(context.Background(), userModel.ID)
		assert.NoError(t, err)

		var count int64
		db.WithContext(context.Background()).Model(&model.User{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Given non-existing user when delete then should return error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewUserRepositoryAdapter(db)
		_ = createTestModelUser(db, "test", "test")

		err := repo.DeleteByID(context.Background(), "other")
		assert.Error(t, err)
		assert.ErrorIs(t, err, userError.ErrUserNotFound)
	})
}
