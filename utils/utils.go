package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/dariubs/percent"
	"github.com/fatih/color"
)

const (
	// WelcomeArt the welcome art for the user
	WelcomeArt = `
=========================================================================
      _         _                  _          __  __  _                  
     | |       | |                (_)        |  \/  |(_)                 
     | |  ___  | |__    ___  ___   _  _ __   | \  / | _ __  __ ___  _ __ 
 _   | | / _ \ | '_ \  / __|/ _ \ | || '_ \  | |\/| || |\ \/ // _ \| '__|
| |__| || (_) || |_) || (__| (_) || || | | | | |  | || | >  <|  __/| |   
 \____/  \___/ |_.__/  \___|\___/ |_||_| |_| |_|  |_||_|/_/\_\\___||_|   																		
=========================================================================
`
)

var (
	// Red color for errors
	Red = color.New(color.FgRed, color.Bold).SprintFunc()
	// Blue color for messages
	Blue = color.New(color.FgBlue, color.Bold).SprintFunc()
	// Green color for success
	Green = color.New(color.FgGreen, color.Bold).SprintFunc()
)

const (
	// Min is the minimum value to wait in milliseconds Min * milliseconds
	Min = 100
	// Max is the maximum value to wait in milliseconds Max * milliseconds
	Max = 900
	// Deposit is the deposit or house address which will receive the unmixed JobCoins
	Deposit = "Deposit"
	// SmallAmount is the very small amount to create a new address
	SmallAmount = 0.000001
	// Fee is the percentage fee that the mixer charges
	Fee = 1
	// BaseURL is the base url for the jobcoin api
	BaseURL = "http://jobcoin.gemini.com/certainly-thursday/api"
	// TransactionsEndpoint is the transactions endpoint
	TransactionsEndpoint = BaseURL + "/transactions"
	// DepositEndpoint is the deposit address endpoint
	DepositEndpoint = BaseURL + "/addresses/" + Deposit
	// HouseAddress is the main house address
	HouseAddress = "House"
)

// CalculateFee gets the fee to charge the user
func CalculateFee(a float64) float64 {
	return percent.PercentFloat(Fee, a)
}

// CalculateTotals gets the amount that should be distributed per address.
// If the random flag is false, we distribute the coins equally across n addresses.
// If the random flag is true, we generate a non uniform distribution of values that add up to the total.
func CalculateTotals(r bool, n int, a float64, f float64) []float64 {
	amounts := make([]float64, n)
	if r == false {
		amount := (a - f) / float64(n)
		for i := range amounts {
			amounts[i] = amount
		}
		return amounts
	}
	var result float64
	amount := (a - f)
	for i := range amounts {
		amounts[i] = rand.Float64()
		result += amounts[i]
	}
	for i := range amounts {
		amounts[i] /= result
		amounts[i] *= amount
	}

	return amounts
}

// StringToFloat return a float64 from a string
func StringToFloat(a string) float64 {
	f, err := strconv.ParseFloat(a, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return f
}
