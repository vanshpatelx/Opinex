package common

import (
	"fmt"
	"holding/pkg/cache"
	"log"
	"math/big"
)

// Cases: User Add more money => Just add to balance
// User remove money => decrese from balance
// if user place order => decrese from balance and add to lb
// if order placed sucessfully if buy order placed => remove from lb
// if order placed sucessfully if sell order placed => add money to balance


func UpdateBalance(action string, userID string, amount float64) error {
	var evalScript string

	switch action {
	case "addBalance":
		// Add amount to balance in JSON format
		evalScript = `
		local data = redis.call("GET", KEYS[1])
		local balance, lb

		if not data then
			balance = 0
			lb = 0
		else
			local decoded = cjson.decode(data)
			balance = decoded.balance
			lb = decoded.lb
		end

		balance = balance + tonumber(ARGV[1])
		local updatedData = cjson.encode({ balance = balance, lb = lb })
		redis.call("SET", KEYS[1], updatedData)
		return updatedData
		`
	case "decreaseBalance":
		// Decrease balance if enough funds exist
		evalScript = `
		local data = redis.call("GET", KEYS[1])
		if not data then return "Insufficient balance" end

		local decoded = cjson.decode(data)
		local balance, lb = decoded.balance, decoded.lb

		if balance < tonumber(ARGV[1]) then return "Insufficient balance" end

		balance = balance - tonumber(ARGV[1])
		local updatedData = cjson.encode({ balance = balance, lb = lb })
		redis.call("SET", KEYS[1], updatedData)
		return updatedData
		`
	case "placeOrder":
		// Deduct from balance and add to locked balance (`lb`)
		evalScript = `
		local data = redis.call("GET", KEYS[1])
		if not data then return "Insufficient balance" end

		local decoded = cjson.decode(data)
		local balance, lb = decoded.balance, decoded.lb

		if balance < tonumber(ARGV[1]) then return "Insufficient balance" end

		balance = balance - tonumber(ARGV[1])
		lb = lb + tonumber(ARGV[1])

		local updatedData = cjson.encode({ balance = balance, lb = lb })
		redis.call("SET", KEYS[1], updatedData)
		return updatedData
		`
	case "buyOrderPlaced":
		// Remove from locked balance (`lb`)
		evalScript = `
		local data = redis.call("GET", KEYS[1])
		if not data then return "No locked balance" end

		local decoded = cjson.decode(data)
		local balance, lb = decoded.balance, decoded.lb

		lb = lb - tonumber(ARGV[1])
		if lb < 0 then lb = 0 end

		local updatedData = cjson.encode({ balance = balance, lb = lb })
		redis.call("SET", KEYS[1], updatedData)
		return updatedData
		`
	case "sellOrderPlaced":
		// Add money back to balance after a successful sell
		evalScript = `
		local data = redis.call("GET", KEYS[1])
		local balance, lb

		if not data then
			balance = 0
			lb = 0
		else
			local decoded = cjson.decode(data)
			balance = decoded.balance
			lb = decoded.lb
		end

		balance = balance + tonumber(ARGV[1])
		local updatedData = cjson.encode({ balance = balance, lb = lb })
		redis.call("SET", KEYS[1], updatedData)
		return updatedData
		`
	default:
		return fmt.Errorf("Invalid action: %s", action)
	}

	// Execute the script on Redis
	key := "Balance:UserID:" + userID
	result, err := cache.EvalUserScript(evalScript, []string{key}, amount)
	if err != nil {
		return fmt.Errorf("Update failed: %v", err)
	}

	log.Println("Balance updated successfully:", result)
	return nil
}


func BuyOrder(Task string, UserID *big.Int, Amount int){
	UpdateBalance(Task, UserID.String(), float64(Amount))
}


func SellOrder(Task string, UserID *big.Int, Amount int){
	UpdateBalance(Task, UserID.String(), float64(Amount))
}

