package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	testDSN := "init_database_test.db"

	t.Run("should return error if failed to connect to database", func(t *testing.T) {
		_, err := InitDatabase("//error")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrFailedToConnectDatabase)
	})

	t.Run("should return database connection", func(t *testing.T) {
		db, err := InitDatabase(testDSN)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NotNil(t, db.conn)
		assert.NotNil(t, db.Models)
		assert.NotNil(t, db.Models.Users)
	})
}
