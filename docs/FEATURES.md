# Features

## Core Features

- Set multiple environment variables and outputs in one step
- Configurable delimiter for key-value pairs
- Whitespace trimming option
- Case sensitivity control for keys
- Duplicate key detection
- Empty value validation
- Detailed operation status and error reporting
- Retry mechanism for file operations
- JSON support for complex data structures
- Group related variables with common prefixes
- Export output variables as environment variables

<br/>

## Value Transformation and Masking

This action supports value transformation and masking of sensitive data:

### Features
- Mask sensitive values in logs
- Convert values to uppercase/lowercase
- URL encode values
- Custom masking patterns using regex
- Escape newlines in values
- Limit value lengths
- Handle empty values

### Example Usage
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'API_KEY,MULTILINE_TEXT,LONG_TEXT'
    env_value: 'secret123,Hello\nWorld,ThisIsAVeryLongText'
    mask_secrets: 'true'
    mask_pattern: '(password|secret).*'
    escape_newlines: 'true'
    max_length: '10'
    allow_empty: 'true'
```

### Value Processing Behavior
- When `escape_newlines` is enabled, newlines are converted to `\n`
- `max_length` truncates values to specified length (0 for unlimited)
- `allow_empty` permits empty values even when `fail_on_empty` is true
- Transformations are applied in order:
  1. Case conversion (upper/lower)
  2. URL encoding
  3. Newline escaping
  4. Length limiting

- Note: Masking only affects log output, not the actual values set in environment variables or outputs.

<br/>

## Group Prefix and Variable Organization

The `group_prefix` option helps organize related variables:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'DATABASE,API,CACHE'
    env_value: 'postgres,graphql,redis'
    group_prefix: 'SYS'
```

When combined with JSON support, it intelligently groups JSON properties under common prefixes:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'CONFIG_DATA'
    env_value: '{"server":{"host":"example.com","port":8080}}'
    json_support: 'true'
    group_prefix: 'APP'
```

This creates variables with consistent naming:
- `CONFIG_DATA`
- `CONFIG_DATA_server_host`
- `CONFIG_DATA_server_port`

<br/>

## Export Outputs as Environment Variables

With `export_as_env: true`, output variables are also set as environment variables:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'ENV_ONLY_VAR'
    env_value: 'env_value'
    output_key: 'OUTPUT_VAR1,OUTPUT_VAR2'
    output_value: 'output_value1,output_value2'
    export_as_env: 'true'
```

This creates:
- Environment Variables: `ENV_ONLY_VAR`, `OUTPUT_VAR1`, `OUTPUT_VAR2`
- Outputs: `OUTPUT_VAR1`, `OUTPUT_VAR2`

This feature provides flexibility in how you access variables in subsequent steps.

<br/>

## Advanced Usage Examples

### 1. Handling Multiline Text
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'CONFIG_DATA'
    env_value: 'line1\nline2\nline3'
    escape_newlines: 'true'
```

### 2. Length-Limited Values
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'DESCRIPTION'
    env_value: 'This is a very long description text'
    max_length: '20'
```

### 3. Optional Values
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'OPTIONAL_VALUE,REQUIRED_VALUE'
    env_value: ',important_data'
    fail_on_empty: 'true'
    allow_empty: 'true'
```
