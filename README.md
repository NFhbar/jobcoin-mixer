# Jobcoin Mixer

The Jobcoin mixer, or *tumbler*, allows for mixing potentially identifiable coins with a common pool of coins, so as to further **anonymize** the origin of the funds.

## Requirements

- [go](https://golang.org/doc/install). Greater than `1.15.6`

## Install

`make install`

## Usage

The `cli` provides one command `mix`:

`$ mixer mix`

There are 3 flags that are supported:

- `--random`
- `--delay`
- `--help`

`--random`

When `true`, the mixer will find a non uniform distribution of amounts that add up to the total to be received minus the fee. When `false` the mixer will simply divide the total to be received over `n`, where `n` is the number of destination addresses.

`--delay`

When `true`, the mixer will add a random delay in between distributions. 

`--help`

To display the help messages.

## UI

A useful dashboard can be found here [https://jobcoin.gemini.com/certainly-thursday](https://jobcoin.gemini.com/certainly-thursday).

Available API is here [https://jobcoin.gemini.com/certainly-thursday/api#addresses](https://jobcoin.gemini.com/certainly-thursday/api#addresses).

## Detecting Deposits 

The main issue with detecting the proper deposits in a single `deposit` address is the possibility of deposits with almost identical data, for example, imagine our `deposit` address already processed a mixing job from `Bob`:

```json
{
    "balance": "0",
    "transactions": 
    [
        {
            "timestamp": "2014-04-22T13:10:01.210Z",
            "fromAddress": "Bob",
            "toAddress": "Deposit",
            "amount": "50"
        },
        {
            "timestamp": "2014-04-23T18:25:43.511Z",
            "fromAddress": "Deposit",
            "toAddress": "House",
            "amount": "50"
        }
    ]
}
```

Now, if `Bob` wanted to mix the same amount of coins again, the resulting transactions would be:

```json
{
    "balance": "50",
    "transactions": 
    [
        {
            "timestamp": "2014-04-22T13:10:01.210Z",
            "fromAddress": "Bob",
            "toAddress": "Deposit",
            "amount": "50"
        },
        {
            "timestamp": "2014-04-23T18:25:43.511Z",
            "fromAddress": "Deposit",
            "toAddress": "House",
            "amount": "50"
        },
        {
            "timestamp": "2014-05-22T13:10:01.210Z",
            "fromAddress": "Bob",
            "toAddress": "Deposit",
            "amount": "50"
        },
    ]
}
```

How would the system detect that a new batch needs to be mixed? Meaning, if we are scanning for transactions coming from `Bob` we would falsely identified the first even transaction as incoming even before `Bob` sent the new one. There could be a solution by collecting the `fromAddress` from the user as well as the `amount` and then working with the `timestamp` to identify the newer deposit and match this to the corresponding job. 

The simpler solution is to provide a brand new `deposit` address per mixing job. That way, we will always be guaranteed that the incoming transaction matches the intended outgoing addresses, since only the user has the randomly generated `deposit` address. Additionally, this way we collect the minimum amount of data from the user possible.

There was no `createAddress` endpoint so I simply repurposed the `transactions` endpoint to create a new address with a negligible amount, which is then swept back to the main house address. This simulates the creation of a new address with zero balance.