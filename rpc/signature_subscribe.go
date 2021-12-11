package rpc

import (
	"net/rpc"
)

type SignatureSubscribeResponse struct {
	GeneralResponse
	Result uint32 `json:"result"`
}

type SignatureNotificationCall struct {
	JsonRPC string                      `json:"jsonrpc"`
	Method  string                      `json:"method"`
	Params  SignatureNotificationParams `json:"params"`
}

type SignatureNotificationResult struct {
	Context Context                    `json:"context"`
	Value   SignatureNotificationValue `json:"value"`
}

type SignatureNotificationParams struct {
	Result       SignatureNotificationResult `json:"result"`
	Subscription uint32                      `json:"subscription"`
}

type SignatureNotificationValue struct {
	Error string `json:"err"`
}

type SignatureSubscribeConfig struct {
	Commitment Commitment `json:"commitment,omitempty"`
}

func (c *RpcClient) SignatureSubscribe(client *rpc.Client, signature string) (SignatureSubscribeResponse, error) {
	var reply SignatureSubscribeResponse
	err := client.Call("signatureSubscribe", signature, &reply)
	return reply, err
}

func (c *RpcClient) SignatureSubscribeWithConfig(client *rpc.Client, signature string, cfg SignatureSubscribeConfig) (SignatureSubscribeResponse, error) {
	var reply SignatureSubscribeResponse
	err := client.Call("signatureSubscribe", []interface{}{signature, cfg}, &reply)
	return reply, err
}
