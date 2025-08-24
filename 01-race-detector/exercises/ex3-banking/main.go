package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TODO: Fix ALL race conditions in this banking system
type Account struct {
	id      int
	balance int // RACE: Concurrent access to balance
}

type Bank struct {
	accounts map[int]*Account // RACE: Concurrent map access
	nextID   int             // RACE: Concurrent access to nextID
}

func (b *Bank) CreateAccount(initialBalance int) int {
	// RACE: Multiple goroutines accessing nextID and accounts map
	id := b.nextID
	b.nextID++
	b.accounts[id] = &Account{id: id, balance: initialBalance}
	return id
}

func (b *Bank) GetBalance(id int) int {
	// RACE: Reading from accounts map and balance
	if account, ok := b.accounts[id]; ok {
		return account.balance
	}
	return 0
}

func (b *Bank) Transfer(fromID, toID, amount int) bool {
	// RACE: Multiple operations on shared data
	fromAccount, ok1 := b.accounts[fromID]
	toAccount, ok2 := b.accounts[toID]
	
	if !ok1 || !ok2 {
		return false
	}
	
	if fromAccount.balance >= amount {
		fromAccount.balance -= amount // RACE: Concurrent balance modification
		toAccount.balance += amount   // RACE: Concurrent balance modification
		return true
	}
	return false
}

func (b *Bank) TotalBalance() int {
	// RACE: Reading from accounts map and balances
	total := 0
	for _, account := range b.accounts {
		total += account.balance
	}
	return total
}

func simulateTransactions(bank *Bank, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Create accounts
	accounts := make([]int, 5)
	for i := 0; i < 5; i++ {
		accounts[i] = bank.CreateAccount(1000)
	}
	
	// Perform random transfers
	for i := 0; i < 100; i++ {
		from := accounts[rand.Intn(5)]
		to := accounts[rand.Intn(5)]
		amount := rand.Intn(100) + 1
		
		bank.Transfer(from, to, amount)
		
		if i%10 == 0 {
			balance := bank.GetBalance(from)
			fmt.Printf("Account %d balance: %d\n", from, balance)
		}
	}
}

func main() {
	bank := &Bank{
		accounts: make(map[int]*Account),
		nextID:   1,
	}
	
	var wg sync.WaitGroup
	
	// Run multiple simulations concurrently
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go simulateTransactions(bank, &wg)
	}
	
	// Periodically check total balance
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				total := bank.TotalBalance()
				fmt.Printf("Total bank balance: %d\n", total)
			case <-done:
				return
			}
		}
	}()
	
	wg.Wait()
	close(done)
	
	fmt.Printf("Final total balance: %d\n", bank.TotalBalance())
}
