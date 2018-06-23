// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
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

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
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
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hxjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &hxjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hxjson.NewAddMultisigAddressCmd(2, keys, hxjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &hxjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   hxjson.String("test"),
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return hxjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &hxjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "createnewaccount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("createnewaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewCreateNewAccountCmd("acct",0)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createnewaccount","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.CreateNewAccountCmd{
				Account: "acct",
				AccountType: 0,
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return hxjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &hxjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &hxjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &hxjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &hxjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &hxjson.GetBalanceCmd{
				Account: nil,
				MinConf: hxjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBalanceCmd(hxjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.GetBalanceCmd{
				Account: hxjson.String("acct"),
				MinConf: hxjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetBalanceCmd(hxjson.String("acct"), hxjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &hxjson.GetBalanceCmd{
				Account: hxjson.String("acct"),
				MinConf: hxjson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNewAddressCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &hxjson.GetNewAddressCmd{
				Account:   nil,
				GapPolicy: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getnewaddress", "acct", "ignore")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetNewAddressCmd(hxjson.String("acct"), hxjson.String("ignore"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct","ignore"],"id":1}`,
			unmarshalled: &hxjson.GetNewAddressCmd{
				Account:   hxjson.String("acct"),
				GapPolicy: hxjson.String("ignore"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &hxjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetRawChangeAddressCmd(hxjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.GetRawChangeAddressCmd{
				Account: hxjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: hxjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetReceivedByAccountCmd("acct", hxjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &hxjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: hxjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hxjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: hxjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetReceivedByAddressCmd("1Address", hxjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &hxjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: hxjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &hxjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewGetTransactionCmd("123", hxjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &hxjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: hxjson.Bool(true),
			},
		},
		{
			name: "importaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importaddress","params":["1Address"],"id":1}`,
			unmarshalled: &hxjson.ImportAddressCmd{
				Address: "1Address",
				Rescan:  hxjson.Bool(true),
			},
		},
		{
			name: "importaddress optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importaddress", "1Address", false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportAddressCmd("1Address", hxjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importaddress","params":["1Address",false],"id":1}`,
			unmarshalled: &hxjson.ImportAddressCmd{
				Address: "1Address",
				Rescan:  hxjson.Bool(false),
			},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportPrivKeyCmd("abc", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &hxjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  hxjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportPrivKeyCmd("abc", hxjson.String("label"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &hxjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   hxjson.String("label"),
				Rescan:  hxjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportPrivKeyCmd("abc", hxjson.String("label"), hxjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &hxjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   hxjson.String("label"),
				Rescan:  hxjson.Bool(false),
			},
		},
		{
			name: "importprivkey optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importprivkey", "abc", "label", false, 12345)
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportPrivKeyCmd("abc", hxjson.String("label"), hxjson.Bool(false), hxjson.Int(12345))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false,12345],"id":1}`,
			unmarshalled: &hxjson.ImportPrivKeyCmd{
				PrivKey:  "abc",
				Label:    hxjson.String("label"),
				Rescan:   hxjson.Bool(false),
				ScanFrom: hxjson.Int(12345),
			},
		},
		{
			name: "importpubkey",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importpubkey", "031234")
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportPubKeyCmd("031234", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importpubkey","params":["031234"],"id":1}`,
			unmarshalled: &hxjson.ImportPubKeyCmd{
				PubKey: "031234",
				Rescan: hxjson.Bool(true),
			},
		},
		{
			name: "importpubkey optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("importpubkey", "031234", false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewImportPubKeyCmd("031234", hxjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importpubkey","params":["031234",false],"id":1}`,
			unmarshalled: &hxjson.ImportPubKeyCmd{
				PubKey: "031234",
				Rescan: hxjson.Bool(false),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return hxjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &hxjson.KeyPoolRefillCmd{
				NewSize: hxjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return hxjson.NewKeyPoolRefillCmd(hxjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &hxjson.KeyPoolRefillCmd{
				NewSize: hxjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &hxjson.ListAccountsCmd{
				MinConf: hxjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListAccountsCmd(hxjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &hxjson.ListAccountsCmd{
				MinConf: hxjson.Int(6),
			},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &hxjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAccountCmd{
				MinConf:          hxjson.Int(1),
				IncludeEmpty:     hxjson.Bool(false),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAccountCmd(hxjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAccountCmd{
				MinConf:          hxjson.Int(6),
				IncludeEmpty:     hxjson.Bool(false),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAccountCmd(hxjson.Int(6), hxjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAccountCmd{
				MinConf:          hxjson.Int(6),
				IncludeEmpty:     hxjson.Bool(true),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAccountCmd(hxjson.Int(6), hxjson.Bool(true), hxjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAccountCmd{
				MinConf:          hxjson.Int(6),
				IncludeEmpty:     hxjson.Bool(true),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAddressCmd{
				MinConf:          hxjson.Int(1),
				IncludeEmpty:     hxjson.Bool(false),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAddressCmd(hxjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAddressCmd{
				MinConf:          hxjson.Int(6),
				IncludeEmpty:     hxjson.Bool(false),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAddressCmd(hxjson.Int(6), hxjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAddressCmd{
				MinConf:          hxjson.Int(6),
				IncludeEmpty:     hxjson.Bool(true),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListReceivedByAddressCmd(hxjson.Int(6), hxjson.Bool(true), hxjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &hxjson.ListReceivedByAddressCmd{
				MinConf:          hxjson.Int(6),
				IncludeEmpty:     hxjson.Bool(true),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &hxjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: hxjson.Int(1),
				IncludeWatchOnly:    hxjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListSinceBlockCmd(hxjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &hxjson.ListSinceBlockCmd{
				BlockHash:           hxjson.String("123"),
				TargetConfirmations: hxjson.Int(1),
				IncludeWatchOnly:    hxjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListSinceBlockCmd(hxjson.String("123"), hxjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &hxjson.ListSinceBlockCmd{
				BlockHash:           hxjson.String("123"),
				TargetConfirmations: hxjson.Int(6),
				IncludeWatchOnly:    hxjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListSinceBlockCmd(hxjson.String("123"), hxjson.Int(6), hxjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &hxjson.ListSinceBlockCmd{
				BlockHash:           hxjson.String("123"),
				TargetConfirmations: hxjson.Int(6),
				IncludeWatchOnly:    hxjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &hxjson.ListTransactionsCmd{
				Account:          nil,
				Count:            hxjson.Int(10),
				From:             hxjson.Int(0),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListTransactionsCmd(hxjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &hxjson.ListTransactionsCmd{
				Account:          hxjson.String("acct"),
				Count:            hxjson.Int(10),
				From:             hxjson.Int(0),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListTransactionsCmd(hxjson.String("acct"), hxjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &hxjson.ListTransactionsCmd{
				Account:          hxjson.String("acct"),
				Count:            hxjson.Int(20),
				From:             hxjson.Int(0),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListTransactionsCmd(hxjson.String("acct"), hxjson.Int(20),
					hxjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &hxjson.ListTransactionsCmd{
				Account:          hxjson.String("acct"),
				Count:            hxjson.Int(20),
				From:             hxjson.Int(1),
				IncludeWatchOnly: hxjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListTransactionsCmd(hxjson.String("acct"), hxjson.Int(20),
					hxjson.Int(1), hxjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &hxjson.ListTransactionsCmd{
				Account:          hxjson.String("acct"),
				Count:            hxjson.Int(20),
				From:             hxjson.Int(1),
				IncludeWatchOnly: hxjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return hxjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &hxjson.ListUnspentCmd{
				MinConf:   hxjson.Int(1),
				MaxConf:   hxjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListUnspentCmd(hxjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &hxjson.ListUnspentCmd{
				MinConf:   hxjson.Int(6),
				MaxConf:   hxjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return hxjson.NewListUnspentCmd(hxjson.Int(6), hxjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &hxjson.ListUnspentCmd{
				MinConf:   hxjson.Int(6),
				MaxConf:   hxjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return hxjson.NewListUnspentCmd(hxjson.Int(6), hxjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &hxjson.ListUnspentCmd{
				MinConf:   hxjson.Int(6),
				MaxConf:   hxjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hxjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return hxjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1,"tree":0}]],"id":1}`,
			unmarshalled: &hxjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []hxjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "renameaccount",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("renameaccount", "oldacct", "newacct")
			},
			staticCmd: func() interface{} {
				return hxjson.NewRenameAccountCmd("oldacct", "newacct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"renameaccount","params":["oldacct","newacct"],"id":1}`,
			unmarshalled: &hxjson.RenameAccountCmd{
				OldAccount: "oldacct",
				NewAccount: "newacct",
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &hxjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hxjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendFromCmd("from", "1Address", 0.5, hxjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &hxjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hxjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendFromCmd("from", "1Address", 0.5, hxjson.Int(6),
					hxjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &hxjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hxjson.Int(6),
				Comment:     hxjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendFromCmd("from", "1Address", 0.5, hxjson.Int(6),
					hxjson.String("comment"), hxjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &hxjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     hxjson.Int(6),
				Comment:     hxjson.String("comment"),
				CommentTo:   hxjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hxjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &hxjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hxjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hxjson.NewSendManyCmd("from", amounts, hxjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &hxjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hxjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return hxjson.NewSendManyCmd("from", amounts, hxjson.Int(6), hxjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &hxjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     hxjson.Int(6),
				Comment:     hxjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &hxjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSendToAddressCmd("1Address", 0.5, hxjson.String("comment"),
					hxjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &hxjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   hxjson.String("comment"),
				CommentTo: hxjson.String("commentto"),
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &hxjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &hxjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &hxjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    hxjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hxjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return hxjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &hxjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]hxjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    hxjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []hxjson.RawTxInput{}
				privKeys := []string{"abc"}
				return hxjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &hxjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]hxjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    hxjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []hxjson.RawTxInput{}
				privKeys := []string{}
				return hxjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					hxjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &hxjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]hxjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    hxjson.String("ALL"),
			},
		},
		{
			name: "sweepaccount - optionals provided",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sweepaccount", "default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu", 6, 0.05)
			},
			staticCmd: func() interface{} {
				return hxjson.NewSweepAccountCmd("default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
					func(i uint32) *uint32 { return &i }(6),
					func(i float64) *float64 { return &i }(0.05))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sweepaccount","params":["default","DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",6,0.05],"id":1}`,
			unmarshalled: &hxjson.SweepAccountCmd{
				SourceAccount:         "default",
				DestinationAddress:    "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
				RequiredConfirmations: func(i uint32) *uint32 { return &i }(6),
				FeePerKb:              func(i float64) *float64 { return &i }(0.05),
			},
		},
		{
			name: "sweepaccount - optionals omitted",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("sweepaccount", "default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu")
			},
			staticCmd: func() interface{} {
				return hxjson.NewSweepAccountCmd("default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sweepaccount","params":["default","DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu"],"id":1}`,
			unmarshalled: &hxjson.SweepAccountCmd{
				SourceAccount:      "default",
				DestinationAddress: "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
			},
		},
		{
			name: "verifyseed",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("verifyseed", "abc")
			},
			staticCmd: func() interface{} {
				return hxjson.NewVerifySeedCmd("abc", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc"],"id":1}`,
			unmarshalled: &hxjson.VerifySeedCmd{
				Seed:    "abc",
				Account: nil,
			},
		},
		{
			name: "verifyseed optional",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("verifyseed", "abc", 5)
			},
			staticCmd: func() interface{} {
				account := hxjson.Uint32(5)
				return hxjson.NewVerifySeedCmd("abc", account)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc",5],"id":1}`,
			unmarshalled: &hxjson.VerifySeedCmd{
				Seed:    "abc",
				Account: hxjson.Uint32(5),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return hxjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &hxjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return hxjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &hxjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return hxjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return hxjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &hxjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
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
