# Troubleshooting

<br/>

## Common Issues and Solutions

<br/>

### 1. Duplicate Keys
- **Error message**: `duplicate key found: KEY_NAME`
- **Solution**: Ensure all keys are unique or set `error_on_duplicate: 'false'`

### 2. Empty Values
- **Error message**: `empty value found for key: KEY_NAME`
- **Solution**: Provide values for all keys or set `fail_on_empty: 'false'`

### 3. File Write Issues
- **Error message**: `failed to write to file`
- **Solution**: Action will automatically retry up to 3 times

### 4. JSON Parsing Errors
- **Error message**: `Invalid JSON format`
- **Solution**: Ensure JSON strings are valid and properly escaped

### 5. Delimiter Conflicts in JSON
- **Error message**: `env_key and env_value must have the same number of entries`
- **Solution**: Use a delimiter that doesn't appear in your JSON (e.g., `|`)

<br/>

## Debugging and Troubleshooting

### Debug Mode
You can enable debug mode to see detailed logging of how your inputs are being processed:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "KEY1::KEY2"
    env_value: "value1::value2"
    debug_mode: true
```

Debug output includes:
- Raw input values
- Normalized values after whitespace processing
- Final key-value pairs
- Delimiter being used
- JSON parsing results (if json_support is enabled)

<br/>

### Empty Values
The action provides two ways to handle empty values:

1. Using `fail_on_empty`:
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "KEY1,KEY2"
    env_value: ",value2"
    fail_on_empty: true  # This will fail
```

2. Using `allow_empty`:
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "KEY1,KEY2"
    env_value: ",value2"
    fail_on_empty: true
    allow_empty: true  # This will allow empty values to pass
```

<br/>

## Multiline and Special Characters

When working with multiline values or values containing commas, you need to consider the following:

### 1. Multiline Input
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: |
      MULTI_KEY1,
      MULTI_KEY2,
      MULTI_KEY3
    env_value: |
      first value,
      second value,
      third value
```

### 2. Values Containing Commas
If your values contain commas, you should use a different delimiter to avoid parsing issues:
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "KEY1::KEY2::KEY3"
    env_value: "value1, with comma::value2::value3"
    delimiter: "::"  # Use a different delimiter when values contain commas
```

### 3. Working with JSON Data
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "CONFIG_JSON|SETTINGS_JSON"
    env_value: '{"server":"example.com","port":8080}|{"logging":true,"debug":false}'
    delimiter: "|"  # Use a delimiter that's not in your JSON
    json_support: "true"
```

### 4. Complex JSON Structures
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "COMPLEX_JSON"
    env_value: '{"server":{"host":"example.com","port":8080},"auth":{"enabled":true,"methods":["oauth","basic"]}}'
    json_support: "true"
    group_prefix: "APP"
    debug_mode: "true"
```

### Important Notes:
- When values contain the default delimiter (comma), use a different delimiter like `::`
- For JSON values, choose a delimiter that won't appear in your JSON data (e.g., `|`)
- Multiline values are automatically normalized
- Use `escape_newlines: true` to properly handle newline characters
- The same delimiter must be used consistently for both keys and values

<br/>

## Debug Output Format

When `debug_mode` is enabled, you'll see detailed information about how your inputs are being processed.
