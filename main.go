package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/joho/godotenv"
)

func etherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(params.Ether))
}

func weiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, big.NewInt(params.Ether))
}

func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Load environment variables from settings.env
	err := godotenv.Load("settings.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	nodeEndpoints := os.Getenv("NODE_ENDPOINTS")
	targetPrivateKey := os.Getenv("TARGET_PRIVATE_KEYS")
	hqAddress := os.Getenv("HQ_ADDRESS")

	if nodeEndpoints == "" || targetPrivateKey == "" || hqAddress == "" {
		log.Fatal("One or more environment variables are not set")
	}

	privateKey, err := crypto.HexToECDSA(targetPrivateKey)
	if err != nil {
		log.Fatalf("Error parsing private key: %v", err)
	}

	// Split node endpoints if multiple are provided
	nodeEndpointsList := strings.Split(nodeEndpoints, ",")

	for _, nodeEndpoint := range nodeEndpointsList {
		nodeEndpoint = strings.TrimSpace(nodeEndpoint)
		go func(endpoint string) {
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				log.Fatalf("Error connecting to the Ethereum client: %v", err)
			}

			checkrate := time.Second // 1 second

			for {
				time.Sleep(checkrate)

				publicKey := privateKey.Public()
				publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
				if !ok {
					log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
				}

				fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
				nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
				checkerr(err)

				balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
				checkerr(err)

				fmt.Println("Nonce: ", nonce)
				fmt.Println("Account balance:", balance)

				gasLimit := uint64(21000) // in units
				gasPrice, err := client.SuggestGasPrice(context.Background())
				checkerr(err)

				gasLimitBigInt := new(big.Int).SetUint64(gasLimit)
				gasExpense := new(big.Int).Mul(gasLimitBigInt, gasPrice)
				fmt.Println("Gas Expense: ", gasExpense)

				value := new(big.Int).Sub(balance, gasExpense)
				fmt.Println("Value: ", value)

				if value.Cmp(big.NewInt(0)) > 0 {
					fmt.Println("Valid balance: Initialize transaction")

					toAddress := common.HexToAddress(hqAddress)
					var data []byte
					tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

					chainID, err := client.NetworkID(context.Background())
					checkerr(err)

					signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
					checkerr(err)

					err = client.SendTransaction(context.Background(), signedTx)
					checkerr(err)

					fmt.Printf("TX sent: %s\n", signedTx.Hash().Hex())
				} else {
					fmt.Println("Balance too low. Waiting...")
				}
			}
		}(nodeEndpoint)
	}

	// Keep the main function running
	select {}
}
