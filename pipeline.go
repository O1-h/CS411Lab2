// Lab 2: Race Detectives -- starter (BUGGY ON PURPOSE).
//
// This file has TWO independent, deliberate race conditions: (1)
// unsynchronized access to accountBalances, and (2) an unsynchronized
// totalProcessed counter. Find both, prove them with `go test -race
// ./...`, and fix them. Do not change the function signatures below
// -- pipeline_test.go calls them directly.
package pipeline

import (
	"sync"
	"sync/atomic"
)

// Transaction is a single account credit/debit to apply.
type Transaction struct {
	AccountID string
	Amount    int
}

// Pipeline holds the shared state every worker goroutine touches.
type Pipeline struct {
	// TODO (Requirement 2): add a sync.Mutex (or sync.RWMutex) here to
	// protect accountBalances. Lock only around the critical section --
	// do not hold it across the whole worker loop body.
	newmut          sync.Mutex
	accountBalances map[string]int

	// TODO (Requirement 3): protect this with either the same mutex
	// above, a separate mutex, or sync/atomic (atomic.Int64). Document
	// your choice with a comment explaining why.
	totalProcessed atomic.Int64
}

// NewPipeline returns a Pipeline ready to process transactions.
func NewPipeline() *Pipeline {
	return &Pipeline{
		accountBalances: make(map[string]int),
	}
}

// applyTransaction is called concurrently by every worker goroutine.
// THIS IS WHERE BOTH RACES LIVE.
func (p *Pipeline) applyTransaction(tx Transaction) {
	// BUG 1: unsynchronized read-modify-write on a shared map.
	// Concurrent map writes in Go don't just "lose an update" the way
	// a plain int counter does -- see Part A, Question 3.
	p.newmut.Lock()
	p.accountBalances[tx.AccountID] += tx.Amount
	p.newmut.Unlock()
	// BUG 2: unsynchronized read-modify-write on a shared int.
	p.totalProcessed.Add(1) //chose atomic since its a single counter

}

// Run spawns numWorkers goroutines that pull Transactions off `txs`
// until the channel is closed, applying each one. Blocks until every
// worker has finished, then returns the final balances and the total
// number of transactions processed.
func (p *Pipeline) Run(txs <-chan Transaction, numWorkers int) (map[string]int, int) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for tx := range txs {
				p.applyTransaction(tx)
			}
		}()
	}

	wg.Wait()
	return p.accountBalances, int(p.totalProcessed.Load())
}
