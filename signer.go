package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const usageString = `
evm-tx-signer [options] <unsigned_tx_string>

	evm-tx-signer signs an usinged transaction using supplied private key for EVM based blockchains like Ethereum and Polygon 

Options:
	--privatekey		Hex encoded private key (required)
	--chainid			The chain ID to use for signing transaction (required)
    --signer			The signer to use [london|legacy] (london by default)
`

var (
	chainId    int64
	privateKey string
	signer     string
)

func init() {
	flag.Usage = showUsage
	flag.StringVar(&privateKey, "privatekey", "", "Hex encoded private key")
	flag.Int64Var(&chainId, "chainid", 0, "Chain ID to sign transaction")
	flag.StringVar(&signer, "signer", "london", "The signer to use (london by default)")
}

func main() {
	flag.Parse()
	args := flag.Args()
	argsLen := len(args)
	if argsLen == 0 {
		showUsage()
	}

	if privateKey == "" {
		showUsage()
	}
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("failed to read private key: %s", err)
	}

	if chainId == 0 {
		showUsage()
	}
	chainIdBig := big.NewInt(chainId)

	var txSigner types.Signer
	switch signer {
	case "london":
		txSigner = types.NewLondonSigner(chainIdBig)
	case "legacy":
		txSigner = types.NewEIP155Signer(chainIdBig)
	default:
		showUsage()
	}

	unsignedHexTx := args[argsLen-1]
	rawTx, err := hex.DecodeString(unsignedHexTx)
	if err != nil {
		printErrorResult(fmt.Sprintf("failed to hex decode unsigned tx: %s", err))
	}

	unsignedTx := new(types.Transaction)
	if err := unsignedTx.UnmarshalBinary(rawTx); err != nil {
		printErrorResult(fmt.Sprintf("failed to decode unsigned tx: %s", err))
	}

	signedRawTx, err := types.SignTx(unsignedTx, txSigner, privKey)
	if err != nil {
		printErrorResult(fmt.Sprintf("failed to sign transaction: %s", err))
	}

	signedTx, err := signedRawTx.MarshalBinary()
	if err != nil {
		printErrorResult(fmt.Sprintf("failed to encoded signed transaction: %s", err))
	}

	printResult(signedTx)
}

type result struct {
	Error    string `json:"error,omitempty"`
	Signer   string `json:"signer"`
	ChainId  int64  `json:"chain_id"`
	SignedTx string `json:"signed_tx,omitempty"`
}

func printErrorResult(errStr string) {
	data := &result{
		Error:   errStr,
		Signer:  signer,
		ChainId: chainId,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(errStr)
	}

	fmt.Println(string(jsonData))
	os.Exit(1)
}

func printResult(signedTx []byte) {
	data := &result{
		Signer:   signer,
		ChainId:  chainId,
		SignedTx: hex.EncodeToString(signedTx),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(signedTx)
		os.Exit(0)
	}

	fmt.Println(string(jsonData))
	os.Exit(0)
}

func showUsage() {
	fmt.Print(usageString)
	os.Exit(1)
}
