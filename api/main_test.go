package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/util"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestConfig() *util.Config {
	return &util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := NewTestConfig()

	server, err := NewServer(config, store)

	require.NoError(t, err)

	return server
}
