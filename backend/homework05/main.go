package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	blockNum := flag.Uint64("blockNum", 0, "Block number to fetch")
	flag.Parse()

	rpcUrl := os.Getenv("ETH_RPC_URL")

	if rpcUrl == "" {
		log.Fatalf("ETH_RPC_URL is not set")
	}

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	log.Printf("Connected to chain ID: %v", chainID)

	var block *types.Block
	if *blockNum == 0 {
		block, err = client.BlockByNumber(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to get the latest block: %v", err)
		}
	} else {
		block, err = client.BlockByNumber(ctx, big.NewInt(int64(*blockNum)))
		if err != nil {
			log.Fatalf("Failed to get block %d: %v", *blockNum, err)
		}
	}
	log.Printf("Block hash: %s", block.Hash().Hex())
	log.Printf("Block number: %d", block.Number().Uint64())
	log.Printf("Block timestamp: %s", time.Unix(int64(block.Time()), 0).Format(time.RFC3339))
	log.Printf("Block transactions: %d", len(block.Transactions()))

	privateKeyHex := os.Getenv("ETH_PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatalf("ETH_PRIVATE_KEY is not set")
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	log.Printf("From address: %s", fromAddress.Hex())
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}
	log.Printf("Pending nonce: %d", nonce)

	toAddress := "0x7360a27F3Be8eF7f85b65E35461f0c96d9fD4E7F"
	log.Printf("To address: %s", toAddress)
	amount := big.NewInt(1e15)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}
	log.Printf("Suggested gas price: %s", gasPrice.String())
	gasLimit := uint64(21000) // standard gas limit for ETH transfer
	log.Printf("Gas limit: %d", gasLimit)
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, gasLimit, gasPrice, nil)
	log.Printf("Created transaction: %v", tx)

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}
	log.Printf("Signed transaction: %v", signedTx)

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}
	log.Printf("Transaction sent: %s", signedTx.Hash().Hex())
}
