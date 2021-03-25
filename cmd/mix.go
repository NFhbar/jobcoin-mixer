package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"mixer/jobcoin"
	"mixer/utils"
)

type mixer struct {
	Amount    float64
	Fee       float64
	Addresses []string
	Total     []float64
	Completed bool
}

func addressValidate(address string) error {
	s := strings.Split(address, ",")
	fmt.Printf("testing %s", s)
	if len(s) > 5 {
		return errors.New("Maximum 5 addresses allowed")
	}
	for _, r := range s {
		if len(r) < 3 || len(r) > 10 {
			return errors.New("Address must be between 3 and 10 characters")
		}
		if (r < "a" || r > "z") && (r < "A" || r > "Z") {
			return errors.New("Invalid address. Only alpha characters")
		}
	}
	return nil
}

func generateMixer(random bool, amount string, outgoingArray []string) mixer {
	amountString := utils.StringToFloat(amount)
	fee := utils.CalculateFee(amountString)
	return mixer{
		Amount:    amountString,
		Fee:       fee,
		Addresses: outgoingArray,
		Total:     utils.CalculateTotals(random, len(outgoingArray), amountString, fee),
		Completed: false,
	}
}

func outGoingPrompt() []string {
	o := outgoingInput()
	outgoing, err := o.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outgoingArray := strings.Split(outgoing, ",")

	fmt.Printf("Your mixed Jobcoins will be sent the following address: %s\n", utils.Green(outgoingArray))
	return outgoingArray
}

func outgoingInput() promptui.Prompt {
	prompt := promptui.Prompt{
		Label:    utils.Blue("Enter up to 5 addresses -comma separated- where receive your mixed Jobcoins"),
		Validate: addressValidate,
	}
	return prompt
}

func executeJob(random bool, delay bool, outgoingArray []string) {
	found := false
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Start()
	depositAddress := jobcoin.GenerateNewDepositAddress()
	fmt.Printf("Please send your Jobcoins to this addresses: %s\n", utils.Green(depositAddress))
	fmt.Println("Scanning for incoming Jobcoins")
	var amount string
	for found == false {
		info := jobcoin.AddressInfo{}
		jobcoin.QueryDepositAddress(&info, depositAddress)
		for _, t := range info.Transactions {
			if utils.StringToFloat(t.Amount) > utils.SmallAmount {
				amount = t.Amount
				found = true
				fmt.Printf("%s\n", utils.Green("Found your deposit"))
				break
			}
		}

	}
	s.Stop()
	m := generateMixer(random, amount, outgoingArray)
	jobcoin.TransferToHomeBase(depositAddress, m.Amount+utils.SmallAmount)
	jobcoin.TransferToDestination(m.Addresses, m.Total, delay)
}

// mixCmd represents the mix command
var mixCmd = &cobra.Command{
	Use:   "mix",
	Short: "The mix commands allows for mixing jobcoins into given addresses.",
	Long: `The mix command will prompt for the following:
1. Up to 5 addresses to send the mixed coins

After this, the program will give you a unique address to deposit your coins. After your 
coins are received, the program will send your coins to the specified addresses minus a fee
for mixing.`,
	Run: func(cmd *cobra.Command, args []string) {

		// welcome the user
		fmt.Println(utils.WelcomeArt)
		// get the flags from the user
		random, _ := cmd.Flags().GetBool("random")
		delay, _ := cmd.Flags().GetBool("delay")
		fmt.Printf("Randomize distribution(s): %s\n", utils.Green(random))
		fmt.Printf("Include random delay in distribution(s): %s\n", utils.Green(delay))
		fmt.Printf("Mixer fee : %s%%\n", utils.Green(utils.Fee))
		// start the prompt
		outgoingArray := outGoingPrompt()
		executeJob(random, delay, outgoingArray)
		fmt.Printf("%s\n", utils.Green("Your coins have been mixed!"))
	},
}

func init() {
	rootCmd.AddCommand(mixCmd)
	mixCmd.Flags().BoolP("random", "r", false, "Randomize the outgoing amounts")
	mixCmd.Flags().BoolP("delay", "d", false, "Add a random delay to the distributions")
}
