// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hxjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/coolsnady/hxd/hxjson"
)

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("addnode", "127.0.0.1", hxjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return hxjson.NewAddNodeCmd("127.0.0.1", hxjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &hxjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: hxjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []hxjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hxjson.NewCreateRawTransactionCmd(txInputs, amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123}],"id":1}`,
			unmarshalled: &hxjson.CreateRawTransactionCmd{
				Inputs:  []hxjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1,"tree":0}]`,
					`{"456":0.0123}`, int64(12312333333), int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []hxjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return hxjson.NewCreateRawTransactionCmd(txInputs, amounts, hxjson.Int64(12312333333), hxjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1,"tree":0}],{"456":0.0123},12312333333,12312333333],"id":1}`,
			unmarshalled: &hxjson.CreateRawTransactionCmd{
				Inputs:   []hxjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: hxjson.Int64(12312333333),
				Expiry:   hxjson.Int64(12312333333),
			},
		},
		{
			name: "debuglevel",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("debuglevel", "trace")
			},
			staticCmd: func() interface{} {
				return hxjson.NewDebugLevelCmd("trace")
			},
			marshalled: `{"jsonrpc":"1.0","method":"debuglevel","params":["trace"],"id":1}`,
			unmarshalled: &hxjson.DebugLevelCmd{
				LevelSpec: "trace",
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hxjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hxjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return hxjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &hxjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "estimatesmartfee",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("estimatesmartfee", 6, hxjson.EstimateSmartFeeConservative)
			},
			staticCmd: func() interface{} {
				return hxjson.NewEstimateSmartFeeCmd(6, hxjson.EstimateSmartFeeConservative)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"estimatesmartfee","params":[6,"conservative"],"id":1}`,
			unmarshalled: &hxjson.EstimateSmartFeeCmd{Confirmations: 6, Mode: hxjson.EstimateSmartFeeConservative},
		},
		{
			name: "generate",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("generate", 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGenerateCmd(1)
			},
			marshalled: `{"jsonrpc":"1.0","method":"generate","params":[1],"id":1}`,
			unmarshalled: &hxjson.GenerateCmd{
				NumBlocks: 1,
			},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &hxjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetAddedNodeInfoCmd(true, hxjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &hxjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: hxjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblock",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getbestblock")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBestBlockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblock","params":[],"id":1}`,
			unmarshalled: &hxjson.GetBestBlockCmd{},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &hxjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &hxjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hxjson.Bool(true),
				VerboseTx: hxjson.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				// Intentionally use a source param that is
				// more pointers than the destination to
				// exercise that path.
				verbosePtr := hxjson.Bool(true)
				return hxjson.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockCmd("123", hxjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &hxjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hxjson.Bool(true),
				VerboseTx: hxjson.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockCmd("123", hxjson.Bool(true), hxjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &hxjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   hxjson.Bool(true),
				VerboseTx: hxjson.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &hxjson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &hxjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &hxjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: hxjson.Bool(true),
			},
		},
		{
			name: "getblocksubsidy",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblocksubsidy", 123, 256)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockSubsidyCmd(123, 256)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocksubsidy","params":[123,256],"id":1}`,
			unmarshalled: &hxjson.GetBlockSubsidyCmd{
				Height: 123,
				Voters: 256,
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &hxjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := hxjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return hxjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &hxjson.GetBlockTemplateCmd{
				Request: &hxjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hxjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hxjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hxjson.GetBlockTemplateCmd{
				Request: &hxjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := hxjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return hxjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &hxjson.GetBlockTemplateCmd{
				Request: &hxjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getcfilter", "123", "extended")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetCFilterCmd("123", "extended")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123","extended"],"id":1}`,
			unmarshalled: &hxjson.GetCFilterCmd{
				Hash:       "123",
				FilterType: "extended",
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getcfilterheader", "123", "extended")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetCFilterHeaderCmd("123", "extended")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123","extended"],"id":1}`,
			unmarshalled: &hxjson.GetCFilterHeaderCmd{
				Hash:       "123",
				FilterType: "extended",
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &hxjson.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &hxjson.GetConnectionCountCmd{},
		},
		{
			name: "getcurrentnet",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getcurrentnet")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetCurrentNetCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getcurrentnet","params":[],"id":1}`,
			unmarshalled: &hxjson.GetCurrentNetCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &hxjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &hxjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &hxjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetInfoCmd{},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &hxjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &hxjson.GetNetworkHashPSCmd{
				Blocks: hxjson.Int(120),
				Height: hxjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNetworkHashPSCmd(hxjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &hxjson.GetNetworkHashPSCmd{
				Blocks: hxjson.Int(200),
				Height: hxjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNetworkHashPSCmd(hxjson.Int(200), hxjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &hxjson.GetNetworkHashPSCmd{
				Blocks: hxjson.Int(200),
				Height: hxjson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawMempoolCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &hxjson.GetRawMempoolCmd{
				Verbose: hxjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawMempoolCmd(hxjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &hxjson.GetRawMempoolCmd{
				Verbose: hxjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional 2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawmempool", false, "all")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawMempoolCmd(hxjson.Bool(false), hxjson.String("all"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false,"all"],"id":1}`,
			unmarshalled: &hxjson.GetRawMempoolCmd{
				Verbose: hxjson.Bool(false),
				TxType:  hxjson.String("all"),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &hxjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hxjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawTransactionCmd("123", hxjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &hxjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: hxjson.Int(1),
			},
		},
		{
			name: "getstakeversions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getstakeversions", "deadbeef", 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetStakeVersionsCmd("deadbeef", 1)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getstakeversions","params":["deadbeef",1],"id":1}`,
			unmarshalled: &hxjson.GetStakeVersionsCmd{
				Hash:  "deadbeef",
				Count: 1,
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &hxjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hxjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetTxOutCmd("123", 1, hxjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &hxjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: hxjson.Bool(true),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &hxjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getvoteinfo",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getvoteinfo", 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetVoteInfoCmd(1)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getvoteinfo","params":[1],"id":1}`,
			unmarshalled: &hxjson.GetVoteInfoCmd{
				Version: 1,
			},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &hxjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetWorkCmd(hxjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &hxjson.GetWorkCmd{
				Data: hxjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return hxjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &hxjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return hxjson.NewHelpCmd(hxjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &hxjson.HelpCmd{
				Command: hxjson.String("getblock"),
			},
		},
		{
			name: "node option remove",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("node", hxjson.NRemove, "1.1.1.1")
			},
			staticCmd: func() interface{} {
				return hxjson.NewNodeCmd("remove", "1.1.1.1", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["remove","1.1.1.1"],"id":1}`,
			unmarshalled: &hxjson.NodeCmd{
				SubCmd: hxjson.NRemove,
				Target: "1.1.1.1",
			},
		},
		{
			name: "node option disconnect",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("node", hxjson.NDisconnect, "1.1.1.1")
			},
			staticCmd: func() interface{} {
				return hxjson.NewNodeCmd("disconnect", "1.1.1.1", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["disconnect","1.1.1.1"],"id":1}`,
			unmarshalled: &hxjson.NodeCmd{
				SubCmd: hxjson.NDisconnect,
				Target: "1.1.1.1",
			},
		},
		{
			name: "node option connect",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("node", hxjson.NConnect, "1.1.1.1", "perm")
			},
			staticCmd: func() interface{} {
				return hxjson.NewNodeCmd("connect", "1.1.1.1", hxjson.String("perm"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"node","params":["connect","1.1.1.1","perm"],"id":1}`,
			unmarshalled: &hxjson.NodeCmd{
				SubCmd:        hxjson.NConnect,
				Target:        "1.1.1.1",
				ConnectSubCmd: hxjson.String("perm"),
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return hxjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &hxjson.PingCmd{},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(1),
				Skip:        hxjson.Int(0),
				Count:       hxjson.Int(100),
				VinExtra:    hxjson.Int(0),
				Reverse:     hxjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address",
					hxjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(0),
				Skip:        hxjson.Int(0),
				Count:       hxjson.Int(100),
				VinExtra:    hxjson.Int(0),
				Reverse:     hxjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address",
					hxjson.Int(0), hxjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(0),
				Skip:        hxjson.Int(5),
				Count:       hxjson.Int(100),
				VinExtra:    hxjson.Int(0),
				Reverse:     hxjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address",
					hxjson.Int(0), hxjson.Int(5), hxjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(0),
				Skip:        hxjson.Int(5),
				Count:       hxjson.Int(10),
				VinExtra:    hxjson.Int(0),
				Reverse:     hxjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address",
					hxjson.Int(0), hxjson.Int(5), hxjson.Int(10), hxjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(0),
				Skip:        hxjson.Int(5),
				Count:       hxjson.Int(10),
				VinExtra:    hxjson.Int(1),
				Reverse:     hxjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address",
					hxjson.Int(0), hxjson.Int(5), hxjson.Int(10),
					hxjson.Int(1), hxjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(0),
				Skip:        hxjson.Int(5),
				Count:       hxjson.Int(10),
				VinExtra:    hxjson.Int(1),
				Reverse:     hxjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return hxjson.NewSearchRawTransactionsCmd("1Address",
					hxjson.Int(0), hxjson.Int(5), hxjson.Int(10),
					hxjson.Int(1), hxjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &hxjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     hxjson.Int(0),
				Skip:        hxjson.Int(5),
				Count:       hxjson.Int(10),
				VinExtra:    hxjson.Int(1),
				Reverse:     hxjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &hxjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hxjson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendRawTransactionCmd("1122", hxjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &hxjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: hxjson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &hxjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hxjson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSetGenerateCmd(true, hxjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &hxjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: hxjson.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return hxjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &hxjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &hxjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := hxjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return hxjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &hxjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &hxjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hxjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hxjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return hxjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &hxjson.VerifyChainCmd{
				CheckLevel: hxjson.Int64(3),
				CheckDepth: hxjson.Int64(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return hxjson.NewVerifyChainCmd(hxjson.Int64(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &hxjson.VerifyChainCmd{
				CheckLevel: hxjson.Int64(2),
				CheckDepth: hxjson.Int64(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return hxjson.NewVerifyChainCmd(hxjson.Int64(2), hxjson.Int64(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &hxjson.VerifyChainCmd{
				CheckLevel: hxjson.Int64(2),
				CheckDepth: hxjson.Int64(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return hxjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &hxjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := hxjson.MarshalCmd("1.0", testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = hxjson.MarshalCmd("1.0", testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request hxjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = hxjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &hxjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &hxjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        hxjson.Error{Code: hxjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &hxjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        hxjson.Error{Code: hxjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error type - got `%T` (%v), got `%T`",
				i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(hxjson.Error); ok {
			gotErrorCode := err.(hxjson.Error).Code
			if gotErrorCode != terr.Code {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.Code)
				continue
			}
		}
	}
}
