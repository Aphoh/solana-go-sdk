package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func noTestSub(t *testing.T) {
	c := NewWebsocketClient("ws://api.devnet.solana.com", "http://localhost")
	err := c.ConfirmTransaction("N5zzJTsrjJnkr6n8s8MSBju5qgmQ4kSYFjmuZZn1xsgtY8PVMeeLETkokqjHzB7XgXi4Y22ruXLVtFmkcAi47QD", 500*time.Millisecond)
	assert.NoError(t, err)
}
