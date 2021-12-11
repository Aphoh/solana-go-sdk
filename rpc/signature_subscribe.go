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

type SignatureSubscribeRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Id      uint64        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewSignatureSubscribeRequest(signature string, cfg SignatureSubscribeConfig) SignatureSubscribeRequest {
	return SignatureSubscribeRequest{
		JsonRPC: "2.0",
		Id:      1,
		Method:  "signatureSubscribe",
		Params:  []interface{}{signature, cfg},
	}
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
