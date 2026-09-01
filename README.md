# Lab 2: Race Detectives

See the assignment for the full Part A/Part B requirements.

## Run
```bash
go test -race ./... 2>&1 | tee race_report_before.txt   # run BEFORE fixing anything
# ...fix pipeline.go...
go test -race -v ./...   # -v is required -- `go test` swallows a passing package's
                          # stdout, including the Success Token, unless -v is set
```
Fix the two deliberate races in `pipeline.go` (unsynchronized map,
unsynchronized counter) marked with `// BUG` comments. Do not modify
`pipeline_test.go` or the function signatures in `pipeline.go`.
Success Token: `RACE-FREE-PIPELINE-CLEARED`.

## Submit
1. `Lab2_Theory.pdf` (or `.md`)
2. `pipeline.go` (your fixed version)
3. `race_report_before.txt`
4. The Success Token from `go test -race -v ./...`
