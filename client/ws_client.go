package client

import (
	"fmt"
	"time"

	"github.com/portto/solana-go-sdk/rpc"
	"golang.org/x/net/websocket"
)

type WebsocketClient struct {
	endpoint string
	origin   string
}

func NewWebsocketClient(endpoint string, origin string) WebsocketClient {
	return WebsocketClient{
		endpoint: endpoint,
		origin:   origin,
	}
}

func (c *WebsocketClient) ConfirmTransaction(signature string, timeout time.Duration) error {
	ws, err := websocket.Dial(c.endpoint, "", c.origin)
	if err != nil {
		return err
	}

	cfg := rpc.SignatureSubscribeConfig{Commitment: "finalized"}

	var res rpc.SignatureSubscribeResponse
	err = websocket.JSON.Send(ws, []interface{}{signature, cfg})
	if err := checkRpcResult(res.GeneralResponse, err); err != nil {
		return err
	}

	var notif rpc.SignatureNotificationCall
	ws.SetDeadline(time.Now().Add(timeout))
	if err := websocket.JSON.Receive(ws, &notif); err != nil {
		return err
	}

	if notif.Params.Result.Value.Error != "" {
		return fmt.Errorf(notif.Params.Result.Value.Error)
	}

	if notif.Params.Subscription != res.Result {
		return fmt.Errorf("Subscription id mismatch. Call: %d, result: %d", notif.Params.Subscription, res.Result)
	}

	ws.Close()

	return nil
}
