# Test Coverage Report

<br/>

## Test Execution

Run all tests:
```bash
go test ./... -v
```

Check coverage:
```bash
go test ./... -cover
```

Test specific packages:
```bash
go test -v ./internal/config
go test -v ./internal/printer
go test -v ./internal/transformer
go test -v ./internal/writer
```

Run benchmarks:
```bash
go test -bench=. ./internal/transformer
go test -bench=. ./internal/writer
```