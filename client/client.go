package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

type Client struct {
	RpcClient rpc.RpcClient
}

func NewClient(endpoint string) *Client {
	return &Client{rpc.NewRpcClient(endpoint)}
}

// GetBalance fetch users lamports(SOL) balance
func (c *Client) GetBalance(ctx context.Context, base58Addr string) (uint64, error) {
	res, err := c.RpcClient.GetBalance(ctx, base58Addr)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result.Value, nil
}

// GetBalance fetch users lamports(SOL) balance with specific commitment
func (c *Client) GetBalanceWithConfig(ctx context.Context, base58Addr string, cfg rpc.GetBalanceConfig) (uint64, error) {
	res, err := c.RpcClient.GetBalanceWithConfig(ctx, base58Addr, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result.Value, nil
}

// GetTokenAccountBalance returns the token balance of an SPL Token account
func (c *Client) GetTokenAccountBalance(ctx context.Context, base58Addr string) (uint64, uint8, error) {
	res, err := c.RpcClient.GetTokenAccountBalance(ctx, base58Addr)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, 0, err
	}
	balance, err := strconv.ParseUint(res.Result.Value.Amount, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to cast token amount, err: %v", err)
	}
	return balance, res.Result.Value.Decimals, nil
}

// GetTokenAccountBalance returns the token balance of an SPL Token account
func (c *Client) GetTokenAccountBalanceWithConfig(ctx context.Context, base58Addr string, cfg rpc.GetTokenAccountBalanceConfig) (uint64, uint8, error) {
	res, err := c.RpcClient.GetTokenAccountBalanceWithConfig(ctx, base58Addr, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, 0, err
	}
	balance, err := strconv.ParseUint(res.Result.Value.Amount, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to cast token amount, err: %v", err)
	}
	return balance, res.Result.Value.Decimals, nil
}

// GetTokenSupply returns the total supply of an SPL Token type.
func (c *Client) GetTokenSupply(ctx context.Context, mintAddr string) (uint64, uint8, error) {
	res, err := c.RpcClient.GetTokenSupply(ctx, mintAddr)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, 0, err
	}
	balance, err := strconv.ParseUint(res.Result.Value.Amount, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to cast token amount, err: %v", err)
	}
	return balance, res.Result.Value.Decimals, nil
}

// GetTokenSupply returns the total supply of an SPL Token type.
func (c *Client) GetTokenSupplyWithConfig(ctx context.Context, mintAddr string, cfg rpc.GetTokenSupplyConfig) (uint64, uint8, error) {
	res, err := c.RpcClient.GetTokenSupplyWithConfig(ctx, mintAddr, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, 0, err
	}
	balance, err := strconv.ParseUint(res.Result.Value.Amount, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to cast token amount, err: %v", err)
	}
	return balance, res.Result.Value.Decimals, nil
}

type AccountInfo struct {
	Lamports  uint64
	Owner     string
	Excutable bool
	RentEpoch uint64
	Data      []byte
}

// GetAccountInfo return account's info
func (c *Client) GetAccountInfo(ctx context.Context, base58Addr string) (AccountInfo, error) {
	return c.processGetAccountInfo(c.RpcClient.GetAccountInfoWithConfig(ctx, base58Addr, rpc.GetAccountInfoConfig{
		Encoding: rpc.GetAccountInfoConfigEncodingBase64,
	}))
}

type GetAccountInfoConfig struct {
	Commitment rpc.Commitment
	DataSlice  *rpc.GetAccountInfoConfigDataSlice
}

// GetAccountInfoWithConfig return account's info
func (c *Client) GetAccountInfoWithConfig(ctx context.Context, base58Addr string, cfg GetAccountInfoConfig) (AccountInfo, error) {
	return c.processGetAccountInfo(c.RpcClient.GetAccountInfoWithConfig(ctx, base58Addr, rpc.GetAccountInfoConfig{
		Encoding:   rpc.GetAccountInfoConfigEncodingBase64,
		Commitment: cfg.Commitment,
		DataSlice:  cfg.DataSlice,
	}))
}

func (c *Client) processGetAccountInfo(res rpc.GetAccountInfoResponse, err error) (AccountInfo, error) {
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return AccountInfo{}, err
	}
	return c.rpcAccountInfoToClientAccountInfo(res.Result.Value)
}

func (c *Client) rpcAccountInfoToClientAccountInfo(v rpc.GetAccountInfoResultValue) (AccountInfo, error) {
	if v == (rpc.GetAccountInfoResultValue{}) {
		return AccountInfo{}, nil
	}

	data, ok := v.Data.([]interface{})
	if !ok {
		return AccountInfo{}, fmt.Errorf("failed to cast raw response to []interface{}")
	}
	if data[1] != string(rpc.GetAccountInfoConfigEncodingBase64) {
		return AccountInfo{}, fmt.Errorf("encoding mistmatch")
	}
	rawData, err := base64.StdEncoding.DecodeString(data[0].(string))
	if err != nil {
		return AccountInfo{}, fmt.Errorf("failed to base64 decode data")
	}
	return AccountInfo{
		Lamports:  v.Lamports,
		Owner:     v.Owner,
		Excutable: v.Excutable,
		RentEpoch: v.RentEpoch,
		Data:      rawData,
	}, nil
}

type GetMultipleAccountsConfig struct {
	Commitment rpc.Commitment
	DataSlice  *rpc.GetMultipleAccountsConfigDataSlice
}

// GetMultipleAccounts returns multiple accounts info
func (c *Client) GetMultipleAccounts(ctx context.Context, base58Addrs []string) ([]AccountInfo, error) {
	return c.processGetMultipleAccounts(c.RpcClient.GetMultipleAccountsWithConfig(ctx, base58Addrs, rpc.GetMultipleAccountsConfig{
		Encoding: rpc.GetMultipleAccountsConfigEncodingBase64,
	}))
}

// GetAccountInfoWithConfig return account's info
func (c *Client) GetMultipleAccountsWithConfig(ctx context.Context, base58Addrs []string, cfg GetMultipleAccountsConfig) ([]AccountInfo, error) {
	return c.processGetMultipleAccounts(c.RpcClient.GetMultipleAccountsWithConfig(ctx, base58Addrs, rpc.GetMultipleAccountsConfig{
		Encoding:   rpc.GetMultipleAccountsConfigEncodingBase64,
		Commitment: cfg.Commitment,
		DataSlice:  cfg.DataSlice,
	}))
}

func (c *Client) processGetMultipleAccounts(res rpc.GetMultipleAccountsResponse, err error) ([]AccountInfo, error) {
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return []AccountInfo{}, err
	}

	return c.rpcMultipleAccountsToClientAccountInfos(res.Result.Value)
}

func (c *Client) rpcMultipleAccountsToClientAccountInfos(values []rpc.GetMultipleAccountsResultValue) ([]AccountInfo, error) {
	res := make([]AccountInfo, len(values))
	for i, v := range values {
		if v == (rpc.GetMultipleAccountsResultValue{}) {
			res[i] = AccountInfo{}
			continue
		}

		data, ok := v.Data.([]interface{})
		if !ok {
			return []AccountInfo{}, fmt.Errorf("failed to cast raw response to []interface{}")
		}
		if data[1] != string(rpc.GetAccountInfoConfigEncodingBase64) {
			return []AccountInfo{}, fmt.Errorf("encoding mistmatch")
		}
		rawData, err := base64.StdEncoding.DecodeString(data[0].(string))
		if err != nil {
			return []AccountInfo{}, fmt.Errorf("failed to base64 decode data")
		}
		res[i] = AccountInfo{
			Lamports:  v.Lamports,
			Owner:     v.Owner,
			Excutable: v.Excutable,
			RentEpoch: v.RentEpoch,
			Data:      rawData,
		}
	}
	return res, nil
}

// DEPRECATED: Please use getFeeForMessage instead This method is expected to be removed in solana-core v2.0
// GetRecentBlockhash return recent blockhash information
func (c *Client) GetRecentBlockhash(ctx context.Context) (rpc.GetRecentBlockHashResultValue, error) {
	res, err := c.RpcClient.GetRecentBlockhash(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return rpc.GetRecentBlockHashResultValue{}, err
	}
	return res.Result.Value, nil
}

type GetLatestBlockhashConfig struct {
	Commitment rpc.Commitment `json:"commitment,omitempty"`
}

// NEW: This method is only available in solana-core v1.9 or newer. Please use getRecentBlockhash for solana-core v1.8
// GetLatestBlockhash returns the latest blockhash
func (c *Client) GetLatestBlockhash(ctx context.Context) (rpc.GetLatestBlockhashValue, error) {
	res, err := c.RpcClient.GetLatestBlockhash(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return rpc.GetLatestBlockhashValue{}, err
	}
	return res.Result.Value, nil
}

// NEW: This method is only available in solana-core v1.9 or newer. Please use getRecentBlockhash for solana-core v1.8
// GetLatestBlockhashWithConfig returns the latest blockhash
func (c *Client) GetLatestBlockhashWithConfig(ctx context.Context, cfg GetLatestBlockhashConfig) (rpc.GetLatestBlockhashValue, error) {
	res, err := c.RpcClient.GetLatestBlockhashWithConfig(ctx, rpc.GetLatestBlockhashConfig{
		Commitment: cfg.Commitment,
	})
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return rpc.GetLatestBlockhashValue{}, err
	}
	return res.Result.Value, nil
}

type IsBlockhashConfig struct {
	Commitment rpc.Commitment `json:"commitment,omitempty"`
}

// NEW: This method is only available in solana-core v1.9 or newer. Please use getFees for solana-core v1.8
// IsBlockhashValid get the fee the network will charge for a particular Message
func (c *Client) IsBlockhashValid(ctx context.Context, blockhash string) (bool, error) {
	res, err := c.RpcClient.IsBlockhashValid(ctx, blockhash)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return false, err
	}
	return res.Result.Value, nil
}

// NEW: This method is only available in solana-core v1.9 or newer. Please use getFees for solana-core v1.8
// IsBlockhashValidWithConfig get the fee the network will charge for a particular Message
func (c *Client) IsBlockhashValidWithConfig(ctx context.Context, blockhash string, cfg IsBlockhashConfig) (bool, error) {
	res, err := c.RpcClient.IsBlockhashValidWithConfig(ctx, blockhash, rpc.IsBlockhashValidConfig{
		Commitment: cfg.Commitment,
	})
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return false, err
	}
	return res.Result.Value, nil
}

type GetFeeForMessageConfig struct {
	Commitment rpc.Commitment `json:"commitment,omitempty"`
}

// NEW: This method is only available in solana-core v1.9 or newer. Please use getFees for solana-core v1.8
// GetFeeForMessage get the fee the network will charge for a particular Message
func (c *Client) GetFeeForMessage(ctx context.Context, message types.Message) (*uint64, error) {
	rawMessage, err := message.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message, err: %v", err)
	}

	res, err := c.RpcClient.GetFeeForMessage(ctx, base64.StdEncoding.EncodeToString(rawMessage))
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result.Value, nil
}

// NEW: This method is only available in solana-core v1.9 or newer. Please use getFees for solana-core v1.8
// GetFeeForMessageWithConfig get the fee the network will charge for a particular Message
func (c *Client) GetFeeForMessageWithConfig(ctx context.Context, message types.Message, cfg GetFeeForMessageConfig) (*uint64, error) {
	rawMessage, err := message.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize message, err: %v", err)
	}

	res, err := c.RpcClient.GetFeeForMessageWithConfig(
		ctx,
		base64.StdEncoding.EncodeToString(rawMessage),
		rpc.GetFeeForMessageConfig{
			Commitment: cfg.Commitment,
		},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result.Value, nil
}

type QuickSendTransactionParam struct {
	Instructions []types.Instruction
	Signers      []types.Account
	FeePayer     common.PublicKey
}

// QuickSendTransaction is a quick way to send tx
func (c *Client) QuickSendTransaction(ctx context.Context, param QuickSendTransactionParam) (string, error) {
	recentBlockhashRes, err := c.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash, err: %v", err)
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			Instructions:    param.Instructions,
			FeePayer:        param.FeePayer,
			RecentBlockhash: recentBlockhashRes.Blockhash,
		}),
		Signers: param.Signers,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new tx, err: %v", err)
	}
	rawTx, err := tx.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize tx, err: %v", err)
	}
	res, err := c.RpcClient.SendTransactionWithConfig(
		ctx,
		base64.StdEncoding.EncodeToString(rawTx),
		rpc.SendTransactionConfig{Encoding: rpc.SendTransactionConfigEncodingBase64},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

// SendTransaction send transaction struct directly
func (c *Client) SendTransaction(ctx context.Context, tx types.Transaction) (string, error) {
	rawTx, err := tx.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize tx, err: %v", err)
	}
	res, err := c.RpcClient.SendTransactionWithConfig(
		ctx,
		base64.StdEncoding.EncodeToString(rawTx),
		rpc.SendTransactionConfig{
			Encoding: rpc.SendTransactionConfigEncodingBase64,
		},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

type SendTransactionConfig struct {
	SkipPreflight       bool
	PreflightCommitment rpc.Commitment
	MaxRetries          uint64
}

// SendTransaction send transaction struct directly
func (c *Client) SendTransactionWithConfig(ctx context.Context, tx types.Transaction, config SendTransactionConfig) (string, error) {
	rawTx, err := tx.Serialize()
	if err != nil {
		return "", fmt.Errorf("failed to serialize tx, err: %v", err)
	}
	res, err := c.RpcClient.SendTransactionWithConfig(
		ctx,
		base64.StdEncoding.EncodeToString(rawTx),
		rpc.SendTransactionConfig{
			Encoding:            rpc.SendTransactionConfigEncodingBase64,
			SkipPreflight:       config.SkipPreflight,
			PreflightCommitment: config.PreflightCommitment,
			MaxRetries:          config.MaxRetries,
		},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

// GetSlot get current slot (finalized)
func (c *Client) GetSlot(ctx context.Context) (uint64, error) {
	res, err := c.RpcClient.GetSlot(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetSlotWithConfig get slot by commitment
func (c *Client) GetSlotWithConfig(ctx context.Context, cfg rpc.GetSlotConfig) (uint64, error) {
	res, err := c.RpcClient.GetSlotWithConfig(ctx, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

type GetTransactionResponse struct {
	Slot        uint64
	Meta        *TransactionMeta
	Transaction types.Transaction
	BlockTime   *int64
}

type TransactionMeta struct {
	Err               interface{}
	Fee               uint64
	PreBalances       []int64
	PostBalances      []int64
	PreTokenBalances  []rpc.TransactionMetaTokenBalance
	PostTokenBalances []rpc.TransactionMetaTokenBalance
	LogMessages       []string
	InnerInstructions []TransactionMetaInnerInstruction
}

type TransactionMetaInnerInstruction struct {
	Index        uint64
	Instructions []types.CompiledInstruction
}

// GetTransaction returns transaction details for a confirmed transaction
func (c *Client) GetTransaction(ctx context.Context, txhash string) (*GetTransactionResponse, error) {
	res, err := c.RpcClient.GetTransactionWithConfig(
		ctx,
		txhash,
		rpc.GetTransactionConfig{
			Encoding: rpc.GetTransactionConfigEncodingBase64,
		},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	if res.Result == nil {
		return nil, nil
	}
	tx, err := getTransaction(res)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetTransactionWithConfig returns transaction details for a confirmed transaction
// will ignore encoding
func (c *Client) GetTransactionWithConfig(ctx context.Context, txhash string, cfg rpc.GetTransactionConfig) (*GetTransactionResponse, error) {
	res, err := c.RpcClient.GetTransactionWithConfig(
		ctx,
		txhash,
		rpc.GetTransactionConfig{
			Encoding:   rpc.GetTransactionConfigEncodingBase64,
			Commitment: cfg.Commitment,
		},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	if res.Result == nil {
		return nil, nil
	}
	tx, err := getTransaction(res)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func getTransaction(res rpc.GetTransactionResponse) (GetTransactionResponse, error) {
	data, ok := res.Result.Transaction.([]interface{})
	if !ok {
		return GetTransactionResponse{}, fmt.Errorf("failed to cast raw response to []interface{}")
	}
	if data[1] != string(rpc.GetTransactionConfigEncodingBase64) {
		return GetTransactionResponse{}, fmt.Errorf("encoding mistmatch")
	}
	rawTx, err := base64.StdEncoding.DecodeString(data[0].(string))
	if err != nil {
		return GetTransactionResponse{}, fmt.Errorf("failed to base64 decode data, err: %v", err)
	}
	tx, err := types.TransactionDeserialize(rawTx)
	if err != nil {
		return GetTransactionResponse{}, fmt.Errorf("failed to deserialize transaction, err: %v", err)
	}

	var transactionMeta *TransactionMeta
	if res.Result.Meta != nil {
		innerInstructions := make([]TransactionMetaInnerInstruction, 0, len(res.Result.Meta.InnerInstructions))
		for _, metaInnerInstruction := range res.Result.Meta.InnerInstructions {
			compiledInstructions := make([]types.CompiledInstruction, 0, len(metaInnerInstruction.Instructions))
			for _, innerInstruction := range metaInnerInstruction.Instructions {
				var data []byte
				if len(innerInstruction.Data) > 0 {
					data, err = base58.Decode(innerInstruction.Data)
					if err != nil {
						return GetTransactionResponse{}, fmt.Errorf("failed to base58 decode data, data: %v, err: %v", innerInstruction.Data, err)
					}
				}
				compiledInstructions = append(compiledInstructions, types.CompiledInstruction{
					ProgramIDIndex: innerInstruction.ProgramIDIndex,
					Accounts:       innerInstruction.Accounts,
					Data:           data,
				})
			}
			innerInstructions = append(innerInstructions, TransactionMetaInnerInstruction{
				Index:        metaInnerInstruction.Index,
				Instructions: compiledInstructions,
			})
		}
		transactionMeta = &TransactionMeta{
			Err:               res.Result.Meta.Err,
			Fee:               res.Result.Meta.Fee,
			PreBalances:       res.Result.Meta.PreBalances,
			PostBalances:      res.Result.Meta.PostBalances,
			PreTokenBalances:  res.Result.Meta.PreTokenBalances,
			PostTokenBalances: res.Result.Meta.PostTokenBalances,
			LogMessages:       res.Result.Meta.LogMessages,
			InnerInstructions: innerInstructions,
		}
	}

	return GetTransactionResponse{
		Slot:        res.Result.Slot,
		BlockTime:   res.Result.BlockTime,
		Transaction: tx,
		Meta:        transactionMeta,
	}, nil
}

type GetBlockResponse struct {
	Blockhash         string
	BlockTime         *int64
	BlockHeight       *int64
	PreviousBlockhash string
	ParentSLot        uint64
	Transactions      []GetBlockTransaction
	Rewards           []rpc.GetBlockReward
}

type GetBlockTransaction struct {
	Meta        *TransactionMeta
	Transaction types.Transaction
}

// GetBlock returns identity and transaction information about a confirmed block in the ledger
func (c *Client) GetBlock(ctx context.Context, slot uint64) (GetBlockResponse, error) {
	res, err := c.RpcClient.GetBlockWithConfig(
		ctx,
		slot,
		rpc.GetBlockConfig{
			Encoding: rpc.GetBlockConfigEncodingBase64,
		},
	)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return GetBlockResponse{}, err
	}
	return getBlock(res)
}

// add test and get block
func getBlock(res rpc.GetBlockResponse) (GetBlockResponse, error) {
	txs := make([]GetBlockTransaction, 0, len(res.Result.Transactions))
	for _, rTx := range res.Result.Transactions {
		data, ok := rTx.Transaction.([]interface{})
		if !ok {
			return GetBlockResponse{}, fmt.Errorf("failed to cast raw response to []interface{}")
		}
		if data[1] != string(rpc.GetTransactionConfigEncodingBase64) {
			return GetBlockResponse{}, fmt.Errorf("encoding mistmatch")
		}
		rawTx, err := base64.StdEncoding.DecodeString(data[0].(string))
		if err != nil {
			return GetBlockResponse{}, fmt.Errorf("failed to base64 decode data, err: %v", err)
		}
		tx, err := types.TransactionDeserialize(rawTx)
		if err != nil {
			return GetBlockResponse{}, fmt.Errorf("failed to deserialize transaction, err: %v", err)
		}

		var transactionMeta *TransactionMeta
		if rTx.Meta != nil {
			innerInstructions := make([]TransactionMetaInnerInstruction, 0, len(rTx.Meta.InnerInstructions))
			for _, metaInnerInstruction := range rTx.Meta.InnerInstructions {
				compiledInstructions := make([]types.CompiledInstruction, 0, len(metaInnerInstruction.Instructions))
				for _, innerInstruction := range metaInnerInstruction.Instructions {
					var data []byte
					if len(innerInstruction.Data) > 0 {
						data, err = base58.Decode(innerInstruction.Data)
						if err != nil {
							return GetBlockResponse{}, fmt.Errorf("failed to base58 decode data, data: %v, err: %v", innerInstruction.Data, err)
						}
					}
					compiledInstructions = append(compiledInstructions, types.CompiledInstruction{
						ProgramIDIndex: innerInstruction.ProgramIDIndex,
						Accounts:       innerInstruction.Accounts,
						Data:           data,
					})
				}
				innerInstructions = append(innerInstructions, TransactionMetaInnerInstruction{
					Index:        metaInnerInstruction.Index,
					Instructions: compiledInstructions,
				})
			}
			transactionMeta = &TransactionMeta{
				Err:               rTx.Meta.Err,
				Fee:               rTx.Meta.Fee,
				PreBalances:       rTx.Meta.PreBalances,
				PostBalances:      rTx.Meta.PostBalances,
				PreTokenBalances:  rTx.Meta.PreTokenBalances,
				PostTokenBalances: rTx.Meta.PostTokenBalances,
				LogMessages:       rTx.Meta.LogMessages,
				InnerInstructions: innerInstructions,
			}
		}

		txs = append(txs,
			GetBlockTransaction{
				Meta:        transactionMeta,
				Transaction: tx,
			},
		)
	}
	return GetBlockResponse{
		Blockhash:         res.Result.Blockhash,
		BlockTime:         res.Result.BlockTime,
		BlockHeight:       res.Result.BlockHeight,
		PreviousBlockhash: res.Result.PreviousBlockhash,
		ParentSLot:        res.Result.ParentSLot,
		Rewards:           res.Result.Rewards,
		Transactions:      txs,
	}, nil
}

// GetMinimumBalanceForRentExemption returns minimum balance required to make account rent exempt
func (c *Client) GetMinimumBalanceForRentExemption(ctx context.Context, dataLen uint64) (uint64, error) {
	res, err := c.RpcClient.GetMinimumBalanceForRentExemption(ctx, dataLen)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetBlockTime returns the estimated production time of a block.
func (c *Client) GetBlockTime(ctx context.Context, slot uint64) (int64, error) {
	res, err := c.RpcClient.GetBlockTime(ctx, slot)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetIdentity returns the identity pubkey for the current node
func (c *Client) GetIdentity(ctx context.Context) (string, error) {
	res, err := c.RpcClient.GetIdentity(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return "", err
	}
	return res.Result.Identity, nil
}

// GetGenesisHash returns the genesis hash
func (c *Client) GetGenesisHash(ctx context.Context) (string, error) {
	res, err := c.RpcClient.GetGenesisHash(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

// GetFirstAvailableBlock returns the slot of the lowest confirmed block that has not been purged from the ledger
func (c *Client) GetFirstAvailableBlock(ctx context.Context) (uint64, error) {
	res, err := c.RpcClient.GetFirstAvailableBlock(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetVersion returns the current solana versions running on the node
func (c *Client) GetVersion(ctx context.Context) (rpc.GetVersionResult, error) {
	res, err := c.RpcClient.GetVersion(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return rpc.GetVersionResult{}, err
	}
	return res.Result, nil
}

// RequestAirdrop requests an airdrop of lamports to a Pubkey
func (c *Client) RequestAirdrop(ctx context.Context, base58Addr string, lamports uint64) (string, error) {
	res, err := c.RpcClient.RequestAirdrop(ctx, base58Addr, lamports)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

// MinimumLedgerSlot returns the lowest slot that the node has information about in its ledger.
// This value may increase over time if the node is configured to purge older ledger data
func (c *Client) MinimumLedgerSlot(ctx context.Context) (uint64, error) {
	res, err := c.RpcClient.MinimumLedgerSlot(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetTransactionCount returns the current Transaction count from the ledger
func (c *Client) GetTransactionCount(ctx context.Context) (uint64, error) {
	res, err := c.RpcClient.GetTransactionCount(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

// GetTransactionCountWithConfig returns the current Transaction count from the ledger
func (c *Client) GetTransactionCountWithConfig(ctx context.Context, cfg rpc.GetTransactionCountConfig) (uint64, error) {
	res, err := c.RpcClient.GetTransactionCountWithConfig(ctx, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return 0, err
	}
	return res.Result, nil
}

type ClusterNode struct {
	Pubkey       common.PublicKey
	Gossip       *string
	Tpu          *string
	Rpc          *string
	Version      *string
	FeatureSet   *uint32
	ShredVersion *uint16
}

// GetClusterNodes returns information about all the nodes participating in the cluster
func (c *Client) GetClusterNodes(ctx context.Context) ([]ClusterNode, error) {
	res, err := c.RpcClient.GetClusterNodes(ctx)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	output := make([]ClusterNode, 0, len(res.Result))
	for _, info := range res.Result {
		output = append(output, ClusterNode{
			Pubkey:       common.PublicKeyFromString(info.Pubkey),
			Gossip:       info.Gossip,
			Tpu:          info.Tpu,
			Rpc:          info.Rpc,
			Version:      info.Version,
			FeatureSet:   info.FeatureSet,
			ShredVersion: info.ShredVersion,
		})
	}
	return output, nil
}

func (c *Client) GetSignatureStatus(ctx context.Context, signature string) (*rpc.GetSignatureStatusesResultValue, error) {
	res, err := c.RpcClient.GetSignatureStatuses(ctx, []string{signature})
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result.Value[0], nil
}

func (c *Client) GetSignatureStatusWithConfig(ctx context.Context, signature string, cfg rpc.GetSignatureStatusesConfig) (*rpc.GetSignatureStatusesResultValue, error) {
	res, err := c.RpcClient.GetSignatureStatusesWithConfig(ctx, []string{signature}, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result.Value[0], nil
}

func (c *Client) GetSignatureStatuses(ctx context.Context, signatures []string) ([]*rpc.GetSignatureStatusesResultValue, error) {
	res, err := c.RpcClient.GetSignatureStatuses(ctx, signatures)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result.Value, nil
}

func (c *Client) GetSignatureStatusesWithConfig(ctx context.Context, signatures []string, cfg rpc.GetSignatureStatusesConfig) ([]*rpc.GetSignatureStatusesResultValue, error) {
	res, err := c.RpcClient.GetSignatureStatusesWithConfig(ctx, signatures, cfg)
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result.Value, nil
}

type SimulateTransaction struct {
	Err      interface{}
	Logs     []string
	Accounts []*AccountInfo
}

type SimulateTransactionConfig struct {
	SigVerify              bool
	Commitment             rpc.Commitment
	ReplaceRecentBlockhash bool
	Addresses              []string
}

func (c *Client) SimulateTransaction(ctx context.Context, tx types.Transaction) (SimulateTransaction, error) {
	rawTx, err := tx.Serialize()
	if err != nil {
		return SimulateTransaction{}, fmt.Errorf("failed to serialize tx, err: %v", err)
	}
	return c.processSimulateTransaction(
		c.RpcClient.SimulateTransactionWithConfig(
			ctx,
			base64.StdEncoding.EncodeToString(rawTx),
			rpc.SimulateTransactionConfig{
				Encoding: rpc.SimulateTransactionConfigEncodingBase64,
			},
		),
	)
}

func (c *Client) SimulateTransactionWithConfig(ctx context.Context, tx types.Transaction, cfg SimulateTransactionConfig) (SimulateTransaction, error) {
	rawTx, err := tx.Serialize()
	if err != nil {
		return SimulateTransaction{}, fmt.Errorf("failed to serialize tx, err: %v", err)
	}

	var accountCfg *rpc.SimulateTransactionConfigAccounts
	if len(cfg.Addresses) > 0 {
		accountCfg = &rpc.SimulateTransactionConfigAccounts{
			Encoding:  rpc.GetAccountInfoConfigEncodingBase64,
			Addresses: cfg.Addresses,
		}
	}

	return c.processSimulateTransaction(
		c.RpcClient.SimulateTransactionWithConfig(
			ctx,
			base64.StdEncoding.EncodeToString(rawTx),
			rpc.SimulateTransactionConfig{
				Encoding:               rpc.SimulateTransactionConfigEncodingBase64,
				SigVerify:              cfg.SigVerify,
				Commitment:             cfg.Commitment,
				ReplaceRecentBlockhash: cfg.ReplaceRecentBlockhash,
				Accounts:               accountCfg,
			},
		),
	)
}

func (c *Client) processSimulateTransaction(res rpc.SimulateTransactionResponse, err error) (SimulateTransaction, error) {
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return SimulateTransaction{}, err
	}

	var accountInfos []*AccountInfo
	if res.Result.Value.Accounts != nil {
		accountInfos = make([]*AccountInfo, 0, len(res.Result.Value.Accounts))
		for _, r := range res.Result.Value.Accounts {
			if r == nil {
				accountInfos = append(accountInfos, nil)
				continue
			}
			accountInfo, err := c.rpcAccountInfoToClientAccountInfo(*r)
			if err != nil {
				return SimulateTransaction{}, err
			}
			accountInfos = append(accountInfos, &accountInfo)
		}
	}

	return SimulateTransaction{
		Err:      res.Result.Value.Err,
		Logs:     res.Result.Value.Logs,
		Accounts: accountInfos,
	}, nil
}

func (c *Client) GetSignaturesForAddress(ctx context.Context, base58Addr string) ([]rpc.GetSignaturesForAddressResult, error) {
	return c.processGetSignaturesForAddress(c.RpcClient.GetSignaturesForAddress(ctx, base58Addr))
}

func (c *Client) GetSignaturesForAddressWithConfig(ctx context.Context, base58Addr string, cfg rpc.GetSignaturesForAddressConfig) ([]rpc.GetSignaturesForAddressResult, error) {
	return c.processGetSignaturesForAddress(c.RpcClient.GetSignaturesForAddressWithConfig(ctx, base58Addr, cfg))
}

func (c *Client) processGetSignaturesForAddress(res rpc.GetSignaturesForAddressResponse, err error) ([]rpc.GetSignaturesForAddressResult, error) {
	err = checkRpcResult(res.GeneralResponse, err)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func checkRpcResult(res rpc.GeneralResponse, err error) error {
	if err != nil {
		return err
	}
	if res.Error != nil {
		errRes, err := json.Marshal(res.Error)
		if err != nil {
			return fmt.Errorf("rpc response error: %v", res.Error)
		}
		return fmt.Errorf("rpc response error: %v", string(errRes))
	}
	return nil
}
