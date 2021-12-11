package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSub(t *testing.T) {
	c := NewWebsocketClient("ws://api.devnet.solana.com", "http://localhost")
	err := c.ConfirmTransaction("4YtTstrrxufdZoWoRHYzULzcWut7LrAyYxzeFo3Hj9foM7j5v6Zw4h9yCNxzVyuMhs1W5ivJMJELpLKo8ZQh3uPA", 20000*time.Millisecond)
	assert.NoError(t, err)
}
