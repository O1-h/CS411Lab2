// Lab 2: Race Detectives -- verification suite (given, do not edit).
//
// Run: go test -race -v ./...
// The -v is required: `go test` swallows a passing package's stdout
// (including this token) unless -v is set.
// Prints "SUCCESS TOKEN: RACE-FREE-PIPELINE-CLEARED" only if every
// check below passes AND the race detector finds nothing.
package pipeline

import (
	"fmt"
	"os"
	"testing"
)

// TestMain gates the Success Token on the whole package's exit code, not
// just one test's own assertions -- m.Run() returns non-zero if ANY test
// fails OR the race detector reports a data race anywhere in the run, even
// one that doesn't trip a t.Errorf check directly.
func TestMain(m *testing.M) {
	code := m.Run()
	if code == 0 {
		fmt.Println("SUCCESS TOKEN: RACE-FREE-PIPELINE-CLEARED")
	}
	os.Exit(code)
}

func TestPipelineCorrectness(t *testing.T) {
	p := NewPipeline()

	const numWorkers = 20
	const txPerAccount = 50
	accounts := []string{"alice", "bob", "carla", "dave"}

	txs := make(chan Transaction)
	go func() {
		for _, acct := range accounts {
			for i := 0; i < txPerAccount; i++ {
				txs <- Transaction{AccountID: acct, Amount: 1}
			}
		}
		close(txs)
	}()

	balances, total := p.Run(txs, numWorkers)

	expectedTotal := len(accounts) * txPerAccount
	if total != expectedTotal {
		t.Errorf("totalProcessed = %d, want %d (this is the classic lost-update symptom of an unsynchronized counter)", total, expectedTotal)
	}

	for _, acct := range accounts {
		if balances[acct] != txPerAccount {
			t.Errorf("balances[%q] = %d, want %d", acct, balances[acct], txPerAccount)
		}
	}
}

// TestPipelineHighContention hammers a SINGLE account from many
// workers simultaneously -- the scenario most likely to expose a
// lost update on accountBalances even if TestPipelineCorrectness
// above happens not to catch it on a given run.
func TestPipelineHighContention(t *testing.T) {
	p := NewPipeline()
	const numWorkers = 50
	const txPerWorker = 100

	txs := make(chan Transaction)
	go func() {
		for i := 0; i < numWorkers*txPerWorker; i++ {
			txs <- Transaction{AccountID: "hot-account", Amount: 1}
		}
		close(txs)
	}()

	balances, total := p.Run(txs, numWorkers)

	expected := numWorkers * txPerWorker
	if balances["hot-account"] != expected {
		t.Errorf("balances[hot-account] = %d, want %d", balances["hot-account"], expected)
	}
	if total != expected {
		t.Errorf("totalProcessed = %d, want %d", total, expected)
	}
}
