# Test Coverage Report

<br/>

## Test Execution

모든 테스트 실행:
```bash
go test ./... -v
```

커버리지 확인:
```bash
go test ./... -cover
```

특정 패키지 테스트:
```bash
go test -v ./internal/config
go test -v ./internal/printer
go test -v ./internal/transformer
go test -v ./internal/writer
```

벤치마크 실행:
```bash
go test -bench=. ./internal/transformer
go test -bench=. ./internal/writer
```