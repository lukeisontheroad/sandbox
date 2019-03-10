package main

import (
	"errors"
	"io"
	"log"

	"github.com/fletaio/common/util"
	"github.com/fletaio/core/amount"
	"github.com/fletaio/core/data"
	"github.com/fletaio/core/transaction"
)

//////////////////////////////////////////////////////////////////////
// Sandbox Area Begin
//////////////////////////////////////////////////////////////////////

// consts
const (
	GameCommandChannelSize = 10
)

// errors
var (
	ErrInvalidCount = errors.New("invalid count")
)

type WebAddCountReq struct {
	UTXO    uint64 `json:"utxo"` // DO NOT CHANGE
	Payload string `json:"payload"`
}

type WebNotify struct {
	Height  int    `json:"height"` // DO NOT CHANGE
	Type    string `json:"type"`   // DO NOT CHANGE
	Payload string `json:"payload"`
	UTXO    int    `json:"utxo"`  // DO NOT CHANGE
	Error   string `json:"error"` // DO NOT CHANGE
}

type WebGameRes struct {
	Height  int    `json:"height"` // DO NOT CHANGE
	Payload string `json:"payload"`
}

// transaction_type transaction types
const (
	// You can define [11 - 59] Transaction type
	// In sandbox, Transaction fee is not calculated
	// Game Transactions
	AddCountTransactionType = transaction.Type(11)
)

func initSandboxComponent(act *data.Accounter, tran *data.Transactor) error {
	TxFeeTable := map[string]*txFee{
		"sandbox.AddCount": &txFee{AddCountTransactionType, amount.COIN.MulC(10)},
		// ADD YOUR OWN TRANSACTION TO HERE
	}
	for name, item := range TxFeeTable {
		if err := tran.RegisterType(name, item.Type, item.Fee); err != nil {
			log.Println(name, item, err)
			return err
		}
	}
	return nil
}

// GameData stores all data of the game
type GameData struct {
	Payload string
}

// NewGameData returns a GameData
func NewGameData() *GameData {
	gd := &GameData{
		Payload: "",
	}
	return gd
}

// WriteTo is a serialization function
func (gd *GameData) WriteTo(w io.Writer) (int64, error) {
	var wrote int64
	if n, err := util.WriteString(w, gd.Payload); err != nil {
		return wrote, err
	} else {
		wrote += n
	}
	return wrote, nil
}

// ReadFrom is a deserialization function
func (gd *GameData) ReadFrom(r io.Reader) (int64, error) {
	var read int64
	if v, n, err := util.ReadString(r); err != nil {
		return read, err
	} else {
		read += n
		gd.Payload = v
	}
	return read, nil
}

//////////////////////////////////////////////////////////////////////
// Sandbox Area End
//////////////////////////////////////////////////////////////////////
