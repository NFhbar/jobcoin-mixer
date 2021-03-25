package jobcoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"mixer/utils"
	"net/http"
	"os"
	"time"
)

// AddressInfo data structure for response object from jobcoin API
type AddressInfo struct {
	Balance      string `json:"balance"`
	Transactions []struct {
		Timestamp   string `json:"timestamp"`
		FromAddress string `json:"fromAddress"`
		ToAddress   string `json:"toAddress"`
		Amount      string `json:"amount"`
	} `json:"transactions"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GenerateNewDepositAddress a new deposit address.
func GenerateNewDepositAddress() string {
	rand.Seed(time.Now().UnixNano())
	newAddress := utils.Deposit + "-" + randSeq(10)
	postTransaction(utils.HouseAddress, newAddress, utils.SmallAmount)

	return newAddress
}

// QueryDepositAddress GET request to the deposit address
func QueryDepositAddress(target interface{}, address string) error {
	depositEndpoint := utils.BaseURL + "/addresses/" + address
	r, err := http.Get(depositEndpoint)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

// TransferToHomeBase POST request from deposit address to house address
func TransferToHomeBase(depositAddress string, amount float64) {
	postTransaction(depositAddress, utils.HouseAddress, amount)
}

// TransferToDestination POST request from house address to destination addresses
func TransferToDestination(depositAddresses []string, amounts []float64, delay bool) {
	for i, address := range depositAddresses {
		if delay == true {
			random := rand.Intn(utils.Max-utils.Min) + utils.Min
			time.Sleep(time.Duration(random) * time.Millisecond)
		}
		postTransaction(utils.HouseAddress, address, amounts[i])
	}
}

func postTransaction(from string, to string, amountFloat float64) {
	amount := fmt.Sprintf("%f", amountFloat)
	values := map[string]string{"fromAddress": from, "toAddress": to, "amount": amount}
	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(utils.TransactionsEndpoint, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
