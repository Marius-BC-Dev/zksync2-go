package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
)

var (
	h bool

	ZkSyncProvider   string
	EthereumProvider string
	PrivateKey       string
	Amount           int64
	// ZkSyncProvider   = "https://sepolia.era.zksync.dev"   // zkSync Era testnet
	// EthereumProvider = "https://rpc.ankr.com/eth_sepolia" // Sepolia testnet
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.StringVar(&ZkSyncProvider, "rpc_l2", "http://127.0.0.1:3050", "l2 host")
	flag.StringVar(&EthereumProvider, "rpc_l1", "http://127.0.0.1:8545", "l1 host")
	flag.StringVar(&PrivateKey, "pv", "0x0", "private key")
	flag.Int64Var(&Amount, "a", 0, "withdraw amount")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	// Connect to zkSync network
	client, err := clients.Dial(ZkSyncProvider)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	// Connect to Ethereum network
	ethClient, err := ethclient.Dial(EthereumProvider)
	if err != nil {
		log.Panic(err)
	}
	defer ethClient.Close()

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Block number: ", blockNumber)

	w, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	if err != nil {
		log.Panic(err)
	}

	balanceL1, err := w.BalanceL1(nil, utils.EthAddress) // balance on goerli network
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Balance L1: ", balanceL1)

	tx, err := w.Withdraw(nil, accounts.WithdrawalTransaction{
		To:     w.Address(),
		Amount: big.NewInt(Amount),
		Token:  utils.EthAddress,
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Withdraw transaction: ", tx.Hash())
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: withdraw [-h [-rpc_l1 host] [-rpc_l2 host]

Options:
`)
	flag.PrintDefaults()
}
