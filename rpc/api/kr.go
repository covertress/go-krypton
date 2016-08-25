// Copyright 2015 The go-krypton Authors
// This file is part of the go-krypton library.
//
// The go-krypton library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-krypton library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-krypton library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"bytes"
	"encoding/json"
	"math/big"

	"fmt"

	"github.com/krypton/go-krypton/common"
	"github.com/krypton/go-krypton/common/natspec"
	"github.com/krypton/go-krypton/kr"
	"github.com/krypton/go-krypton/rlp"
	"github.com/krypton/go-krypton/rpc/codec"
	"github.com/krypton/go-krypton/rpc/shared"
	"github.com/krypton/go-krypton/xkr"
	"gopkg.in/fatih/set.v0"
)

const (
	KrApiVersion = "1.0"
)

// kr api provider
// See https://github.com/krypton/wiki/wiki/JSON-RPC
type krApi struct {
	xkr     *xkr.XKr
	krypton *kr.Krypton
	methods  map[string]krhandler
	codec    codec.ApiCoder
}

// kr callback handler
type krhandler func(*krApi, *shared.Request) (interface{}, error)

var (
	krMapping = map[string]krhandler{
		"eth_accounts":                            (*krApi).Accounts,
		"eth_blockNumber":                         (*krApi).BlockNumber,
		"eth_getBalance":                          (*krApi).GetBalance,
		"eth_protocolVersion":                     (*krApi).ProtocolVersion,
		"eth_coinbase":                            (*krApi).Coinbase,
		"eth_mining":                              (*krApi).IsMining,
		"eth_syncing":                             (*krApi).IsSyncing,
		"eth_gasPrice":                            (*krApi).GasPrice,
		"eth_getStorage":                          (*krApi).GetStorage,
		"eth_storageAt":                           (*krApi).GetStorage,
		"eth_getStorageAt":                        (*krApi).GetStorageAt,
		"eth_getTransactionCount":                 (*krApi).GetTransactionCount,
		"eth_getBlockTransactionCountByHash":      (*krApi).GetBlockTransactionCountByHash,
		"eth_getBlockTransactionCountByNumber":    (*krApi).GetBlockTransactionCountByNumber,
		"eth_getUncleCountByBlockHash":            (*krApi).GetUncleCountByBlockHash,
		"eth_getUncleCountByBlockNumber":          (*krApi).GetUncleCountByBlockNumber,
		"eth_getData":                             (*krApi).GetData,
		"eth_getCode":                             (*krApi).GetData,
		"eth_getNatSpec":                          (*krApi).GetNatSpec,
		"eth_sign":                                (*krApi).Sign,
		"eth_sendRawTransaction":                  (*krApi).SubmitTransaction,
		"eth_submitTransaction":                   (*krApi).SubmitTransaction,
		"eth_sendTransaction":                     (*krApi).SendTransaction,
		"eth_signTransaction":                     (*krApi).SignTransaction,
		"eth_transact":                            (*krApi).SendTransaction,
		"eth_estimateGas":                         (*krApi).EstimateGas,
		"eth_call":                                (*krApi).Call,
		"eth_flush":                               (*krApi).Flush,
		"eth_getBlockByHash":                      (*krApi).GetBlockByHash,
		"eth_getBlockByNumber":                    (*krApi).GetBlockByNumber,
		"eth_getTransactionByHash":                (*krApi).GetTransactionByHash,
		"eth_getTransactionByBlockNumberAndIndex": (*krApi).GetTransactionByBlockNumberAndIndex,
		"eth_getTransactionByBlockHashAndIndex":   (*krApi).GetTransactionByBlockHashAndIndex,
		"eth_getUncleByBlockHashAndIndex":         (*krApi).GetUncleByBlockHashAndIndex,
		"eth_getUncleByBlockNumberAndIndex":       (*krApi).GetUncleByBlockNumberAndIndex,
		"eth_getCompilers":                        (*krApi).GetCompilers,
		"eth_compileSolidity":                     (*krApi).CompileSolidity,
		"eth_newFilter":                           (*krApi).NewFilter,
		"eth_newBlockFilter":                      (*krApi).NewBlockFilter,
		"eth_newPendingTransactionFilter":         (*krApi).NewPendingTransactionFilter,
		"eth_uninstallFilter":                     (*krApi).UninstallFilter,
		"eth_getFilterChanges":                    (*krApi).GetFilterChanges,
		"eth_getFilterLogs":                       (*krApi).GetFilterLogs,
		"eth_getLogs":                             (*krApi).GetLogs,
		"eth_hashrate":                            (*krApi).Hashrate,
		"eth_getWork":                             (*krApi).GetWork,
		"eth_submitWork":                          (*krApi).SubmitWork,
		"eth_submitHashrate":                      (*krApi).SubmitHashrate,
		"eth_resend":                              (*krApi).Resend,
		"eth_pendingTransactions":                 (*krApi).PendingTransactions,
		"eth_getTransactionReceipt":               (*krApi).GetTransactionReceipt,
	}
)

// create new krApi instance
func NewKrApi(xkr *xkr.XKr, kr *kr.Krypton, codec codec.Codec) *krApi {
	return &krApi{xkr, kr, krMapping, codec.New(nil)}
}

// collection with supported methods
func (self *krApi) Methods() []string {
	methods := make([]string, len(self.methods))
	i := 0
	for k := range self.methods {
		methods[i] = k
		i++
	}
	return methods
}

// Execute given request
func (self *krApi) Execute(req *shared.Request) (interface{}, error) {
	if callback, ok := self.methods[req.Method]; ok {
		return callback(self, req)
	}

	return nil, shared.NewNotImplementedError(req.Method)
}

func (self *krApi) Name() string {
	return shared.KrApiName
}

func (self *krApi) ApiVersion() string {
	return KrApiVersion
}

func (self *krApi) Accounts(req *shared.Request) (interface{}, error) {
	return self.xkr.Accounts(), nil
}

func (self *krApi) Hashrate(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xkr.HashRate()), nil
}

func (self *krApi) BlockNumber(req *shared.Request) (interface{}, error) {
	num := self.xkr.CurrentBlock().Number()
	return newHexNum(num.Bytes()), nil
}

func (self *krApi) GetBalance(req *shared.Request) (interface{}, error) {
	args := new(GetBalanceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xkr.AtStateNum(args.BlockNumber).BalanceAt(args.Address), nil
}

func (self *krApi) ProtocolVersion(req *shared.Request) (interface{}, error) {
	return self.xkr.KrVersion(), nil
}

func (self *krApi) Coinbase(req *shared.Request) (interface{}, error) {
	return newHexData(self.xkr.Coinbase()), nil
}

func (self *krApi) IsMining(req *shared.Request) (interface{}, error) {
	return self.xkr.IsMining(), nil
}

func (self *krApi) IsSyncing(req *shared.Request) (interface{}, error) {
	origin, current, height := self.krypton.Downloader().Progress()
	if current < height {
		return map[string]interface{}{
			"startingBlock": newHexNum(big.NewInt(int64(origin)).Bytes()),
			"currentBlock":  newHexNum(big.NewInt(int64(current)).Bytes()),
			"highestBlock":  newHexNum(big.NewInt(int64(height)).Bytes()),
		}, nil
	}
	return false, nil
}

func (self *krApi) GasPrice(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xkr.DefaultGasPrice().Bytes()), nil
}

func (self *krApi) GetStorage(req *shared.Request) (interface{}, error) {
	args := new(GetStorageArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xkr.AtStateNum(args.BlockNumber).State().SafeGet(args.Address).Storage(), nil
}

func (self *krApi) GetStorageAt(req *shared.Request) (interface{}, error) {
	args := new(GetStorageAtArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return self.xkr.AtStateNum(args.BlockNumber).StorageAt(args.Address, args.Key), nil
}

func (self *krApi) GetTransactionCount(req *shared.Request) (interface{}, error) {
	args := new(GetTxCountArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	count := self.xkr.AtStateNum(args.BlockNumber).TxCountAt(args.Address)
	return fmt.Sprintf("%#x", count), nil
}

func (self *krApi) GetBlockTransactionCountByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xkr.KrBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *krApi) GetBlockTransactionCountByNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xkr.KrBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Transactions())), nil
}

func (self *krApi) GetUncleCountByBlockHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xkr.KrBlockByHash(args.Hash)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *krApi) GetUncleCountByBlockNumber(req *shared.Request) (interface{}, error) {
	args := new(BlockNumArg)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xkr.KrBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return fmt.Sprintf("%#x", len(block.Uncles())), nil
}

func (self *krApi) GetData(req *shared.Request) (interface{}, error) {
	args := new(GetDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v := self.xkr.AtStateNum(args.BlockNumber).CodeAtBytes(args.Address)
	return newHexData(v), nil
}

func (self *krApi) Sign(req *shared.Request) (interface{}, error) {
	args := new(NewSigArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	v, err := self.xkr.Sign(args.From, args.Data, false)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *krApi) SubmitTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewDataArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	v, err := self.xkr.PushTx(args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// JsonTransaction is returned as response by the JSON RPC. It contains the
// signed RLP encoded transaction as Raw and the signed transaction object as Tx.
type JsonTransaction struct {
	Raw string `json:"raw"`
	Tx  *tx    `json:"tx"`
}

func (self *krApi) SignTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	tx, err := self.xkr.SignTransaction(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}

	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return JsonTransaction{"0x" + common.Bytes2Hex(data), newTx(tx)}, nil
}

func (self *krApi) SendTransaction(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	// nonce may be nil ("guess" mode)
	var nonce string
	if args.Nonce != nil {
		nonce = args.Nonce.String()
	}

	var gas, price string
	if args.Gas != nil {
		gas = args.Gas.String()
	}
	if args.GasPrice != nil {
		price = args.GasPrice.String()
	}
	v, err := self.xkr.Transact(args.From, args.To, nonce, args.Value.String(), gas, price, args.Data)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (self *krApi) GetNatSpec(req *shared.Request) (interface{}, error) {
	args := new(NewTxArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	var jsontx = fmt.Sprintf(`{"params":[{"to":"%s","data": "%s"}]}`, args.To, args.Data)
	notice := natspec.GetNotice(self.xkr, jsontx, self.krypton.HTTPClient())

	return notice, nil
}

func (self *krApi) EstimateGas(req *shared.Request) (interface{}, error) {
	_, gas, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if len(gas) == 0 {
		return newHexNum(0), nil
	} else {
		return newHexNum(common.String2Big(gas)), err
	}
}

func (self *krApi) Call(req *shared.Request) (interface{}, error) {
	v, _, err := self.doCall(req.Params)
	if err != nil {
		return nil, err
	}

	// TODO unwrap the parent method's ToHex call
	if v == "0x0" {
		return newHexData([]byte{}), nil
	} else {
		return newHexData(common.FromHex(v)), nil
	}
}

func (self *krApi) Flush(req *shared.Request) (interface{}, error) {
	return nil, shared.NewNotImplementedError(req.Method)
}

func (self *krApi) doCall(params json.RawMessage) (string, string, error) {
	args := new(CallArgs)
	if err := self.codec.Decode(params, &args); err != nil {
		return "", "", err
	}

	return self.xkr.AtStateNum(args.BlockNumber).Call(args.From, args.To, args.Value.String(), args.Gas.String(), args.GasPrice.String(), args.Data)
}

func (self *krApi) GetBlockByHash(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByHashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	block := self.xkr.KrBlockByHash(args.BlockHash)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xkr.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *krApi) GetBlockByNumber(req *shared.Request) (interface{}, error) {
	args := new(GetBlockByNumberArgs)
	if err := json.Unmarshal(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	block := self.xkr.KrBlockByNumber(args.BlockNumber)
	if block == nil {
		return nil, nil
	}
	return NewBlockRes(block, self.xkr.Td(block.Hash()), args.IncludeTxs), nil
}

func (self *krApi) GetTransactionByHash(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	tx, bhash, bnum, txi := self.xkr.KrTransactionByHash(args.Hash)
	if tx != nil {
		v := NewTransactionRes(tx)
		// if the blockhash is 0, assume this is a pending transaction
		if bytes.Compare(bhash.Bytes(), bytes.Repeat([]byte{0}, 32)) != 0 {
			v.BlockHash = newHexData(bhash)
			v.BlockNumber = newHexNum(bnum)
			v.TxIndex = newHexNum(txi)
		}
		return v, nil
	}
	return nil, nil
}

func (self *krApi) GetTransactionByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xkr.KrBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xkr.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Transactions[args.Index], nil
	}
}

func (self *krApi) GetTransactionByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xkr.KrBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xkr.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Transactions)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Transactions[args.Index], nil
}

func (self *krApi) GetUncleByBlockHashAndIndex(req *shared.Request) (interface{}, error) {
	args := new(HashIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xkr.KrBlockByHash(args.Hash)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xkr.Td(raw.Hash()), false)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		// return NewValidationError("Index", "does not exist")
		return nil, nil
	}
	return block.Uncles[args.Index], nil
}

func (self *krApi) GetUncleByBlockNumberAndIndex(req *shared.Request) (interface{}, error) {
	args := new(BlockNumIndexArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	raw := self.xkr.KrBlockByNumber(args.BlockNumber)
	if raw == nil {
		return nil, nil
	}
	block := NewBlockRes(raw, self.xkr.Td(raw.Hash()), true)
	if args.Index >= int64(len(block.Uncles)) || args.Index < 0 {
		return nil, nil
	} else {
		return block.Uncles[args.Index], nil
	}
}

func (self *krApi) GetCompilers(req *shared.Request) (interface{}, error) {
	var lang string
	if solc, _ := self.xkr.Solc(); solc != nil {
		lang = "Solidity"
	}
	c := []string{lang}
	return c, nil
}

func (self *krApi) CompileSolidity(req *shared.Request) (interface{}, error) {
	solc, _ := self.xkr.Solc()
	if solc == nil {
		return nil, shared.NewNotAvailableError(req.Method, "solc (solidity compiler) not found")
	}

	args := new(SourceArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	contracts, err := solc.Compile(args.Source)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (self *krApi) NewFilter(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	id := self.xkr.NewLogFilter(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)
	return newHexNum(big.NewInt(int64(id)).Bytes()), nil
}

func (self *krApi) NewBlockFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xkr.NewBlockFilter()), nil
}

func (self *krApi) NewPendingTransactionFilter(req *shared.Request) (interface{}, error) {
	return newHexNum(self.xkr.NewTransactionFilter()), nil
}

func (self *krApi) UninstallFilter(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xkr.UninstallFilter(args.Id), nil
}

func (self *krApi) GetFilterChanges(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	switch self.xkr.GetFilterType(args.Id) {
	case xkr.BlockFilterTy:
		return NewHashesRes(self.xkr.BlockFilterChanged(args.Id)), nil
	case xkr.TransactionFilterTy:
		return NewHashesRes(self.xkr.TransactionFilterChanged(args.Id)), nil
	case xkr.LogFilterTy:
		return NewLogsRes(self.xkr.LogFilterChanged(args.Id)), nil
	default:
		return []string{}, nil // reply empty string slice
	}
}

func (self *krApi) GetFilterLogs(req *shared.Request) (interface{}, error) {
	args := new(FilterIdArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	return NewLogsRes(self.xkr.Logs(args.Id)), nil
}

func (self *krApi) GetLogs(req *shared.Request) (interface{}, error) {
	args := new(BlockFilterArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return NewLogsRes(self.xkr.AllLogs(args.Earliest, args.Latest, args.Skip, args.Max, args.Address, args.Topics)), nil
}

func (self *krApi) GetWork(req *shared.Request) (interface{}, error) {
	self.xkr.SetMining(true, 0)
	ret, err := self.xkr.RemoteMining().GetWork()
	if err != nil {
		return nil, shared.NewNotReadyError("mining work")
	} else {
		return ret, nil
	}
}

func (self *krApi) SubmitWork(req *shared.Request) (interface{}, error) {
	args := new(SubmitWorkArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}
	return self.xkr.RemoteMining().SubmitWork(args.Nonce, common.HexToHash(args.Digest), common.HexToHash(args.Header)), nil
}

func (self *krApi) SubmitHashrate(req *shared.Request) (interface{}, error) {
	args := new(SubmitHashRateArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return false, shared.NewDecodeParamError(err.Error())
	}
	self.xkr.RemoteMining().SubmitHashrate(common.HexToHash(args.Id), args.Rate)
	return true, nil
}

func (self *krApi) Resend(req *shared.Request) (interface{}, error) {
	args := new(ResendArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	from := common.HexToAddress(args.Tx.From)

	pending := self.krypton.TxPool().GetTransactions()
	for _, p := range pending {
		if pFrom, err := p.FromFrontier(); err == nil && pFrom == from && p.SigHash() == args.Tx.tx.SigHash() {
			self.krypton.TxPool().RemoveTx(common.HexToHash(args.Tx.Hash))
			return self.xkr.Transact(args.Tx.From, args.Tx.To, args.Tx.Nonce, args.Tx.Value, args.GasLimit, args.GasPrice, args.Tx.Data)
		}
	}

	return nil, fmt.Errorf("Transaction %s not found", args.Tx.Hash)
}

func (self *krApi) PendingTransactions(req *shared.Request) (interface{}, error) {
	txs := self.krypton.TxPool().GetTransactions()

	// grab the accounts from the account manager. This will help with determining which
	// transactions should be returned.
	accounts, err := self.krypton.AccountManager().Accounts()
	if err != nil {
		return nil, err
	}

	// Add the accouns to a new set
	accountSet := set.New()
	for _, account := range accounts {
		accountSet.Add(account.Address)
	}

	var ltxs []*tx
	for _, tx := range txs {
		if from, _ := tx.FromFrontier(); accountSet.Has(from) {
			ltxs = append(ltxs, newTx(tx))
		}
	}

	return ltxs, nil
}

func (self *krApi) GetTransactionReceipt(req *shared.Request) (interface{}, error) {
	args := new(HashArgs)
	if err := self.codec.Decode(req.Params, &args); err != nil {
		return nil, shared.NewDecodeParamError(err.Error())
	}

	txhash := common.BytesToHash(common.FromHex(args.Hash))
	tx, bhash, bnum, txi := self.xkr.KrTransactionByHash(args.Hash)
	rec := self.xkr.GetTxReceipt(txhash)
	// We could have an error of "not found". Should disambiguate
	// if err != nil {
	// 	return err, nil
	// }
	if rec != nil && tx != nil {
		v := NewReceiptRes(rec)
		v.BlockHash = newHexData(bhash)
		v.BlockNumber = newHexNum(bnum)
		v.TransactionIndex = newHexNum(txi)
		return v, nil
	}

	return nil, nil
}
