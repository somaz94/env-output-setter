# Test Suite for Environment/Output Setter

Comprehensive test suite for validating the Environment/Output Setter functionality.

<br/>

## Running Tests

<br/>

### Unit Tests

```bash
# Run all unit tests with coverage
go test ./internal/... -v -cover

# Generate coverage profile
go test ./internal/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

<br/>

### Integration Tests (Local)

```bash
# Run from project root
go run tests/test_local.go

# Or from tests directory
cd tests
go run test_local.go
```

<br/>

### Features

- **Basic Tests**: Environment variables, custom delimiters
- **Whitespace Tests**: Trimming, normalization
- **Transformation Tests**: Uppercase, lowercase, URL encoding, length limits, newline escaping
- **Validation Tests**: Duplicate keys, empty values
- **JSON Tests**: Simple JSON, nested objects, arrays
- **Combined Tests**: Multiple transformations at once
- Color-coded output
- Detailed error messages
- Category organization

<br/>

## Test Categories

<br/>

### 1. Basic Tests (2 tests)
- Basic environment variables
- Custom delimiter (::)

<br/>

### 2. Whitespace Tests (1 test)
- Trim whitespace

<br/>

### 3. Transformation Tests (5 tests)
- Convert to uppercase
- Convert to lowercase
- URL encoding
- Max length limitation
- Escape newlines

<br/>

### 4. Validation Tests (2 tests)
- Duplicate keys (should fail)
- Empty values with allow_empty

<br/>

### 5. JSON Tests (3 tests)
- Simple JSON parsing
- Nested JSON parsing
- JSON array parsing

<br/>

### 6. Combined Tests (1 test)
- Multiple transformations (uppercase + URL encode + max length)

<br/>

## Test Structure

Each test case includes:

```go
TestCase{
    Name:         "Test name",
    Category:     "Category name",
    EnvKeys:      "KEY1,KEY2",
    EnvValues:    "value1,value2",
    OutputKeys:   "OUT1,OUT2",
    OutputValues: "out1,out2",
    Delimiter:    ",",
    ExpectedEnvs: map[string]string{...},
    ExpectedOuts: map[string]string{...},
    ShouldFail:   false,
    Options:      TestOptions{...},
}
```

<br/>

## Test Options

```go
TestOptions{
    FailOnEmpty:      bool
    TrimWhitespace:   bool
    CaseSensitive:    bool
    ErrorOnDuplicate: bool
    MaskSecrets:      bool
    ToUpper:          bool
    ToLower:          bool
    EncodeURL:        bool
    EscapeNewlines:   bool
    MaxLength:        int
    AllowEmpty:       bool
    JsonSupport:      bool
    ExportAsEnv:      bool
    GroupPrefix:      string
}
```

<br/>

## Adding New Tests

To add a new test, append to the `createTestSuite()` function:

```go
tests = append(tests, TestCase{
    Name:         "Your test name",
    Category:     "Your category",
    EnvKeys:      "KEY",
    EnvValues:    "value",
    // ... other fields
})
```

<br/>

## Debugging

Enable debug mode in test config:

```go
cfg := &config.Config{
    // ... other settings
    DebugMode: true,
}
```

<br/>

## CI/CD Integration

Tests are automatically run in GitHub Actions:

```yaml
- name: Run Go unit tests with coverage
  run: go test ./internal/... -v -cover -coverprofile=coverage.out

- name: Run Local Test Suite
  run: go run tests/test_local.go
```

<br/>

## Related Documentation

- [Main README](../README.md)
- [Action Configuration](../action.yml)
- [CI Workflow](../.github/workflows/ci.yml)
