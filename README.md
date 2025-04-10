# GitHub Environment/Output Setter

[![License](https://img.shields.io/github/license/somaz94/env-output-setter)](https://github.com/somaz94/container-action)
![Latest Tag](https://img.shields.io/github/v/tag/somaz94/env-output-setter)
![Top Language](https://img.shields.io/github/languages/top/somaz94/env-output-setter?color=green&logo=go&logoColor=b)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Environment/Output%20Setter-blue?logo=github)](https://github.com/marketplace/actions/environment-output-setter)

## Overview

The **GitHub Environment/Output Setter** is a GitHub Action that allows you to
set multiple key-value pairs in both `$GITHUB_ENV` and `$GITHUB_OUTPUT`. This
action is useful for workflows that need to dynamically define environment
variables or output values that other steps can reference.

<br/>

## Inputs

| Input               | Required | Description                                         | Default | Example                       |
| ------------------ | -------- | --------------------------------------------------- | ------- | ----------------------------- |
| `env_key`          | Yes      | Comma-separated list of environment variable keys   | -       | `"GCP_REGION,AWS_REGION"`     |
| `env_value`        | Yes      | Comma-separated list of environment variable values | -       | `"asia-northeast1,us-east-1"` |
| `output_key`       | Yes      | Comma-separated list of output keys                 | -       | `"GCP_OUTPUT,AWS_OUTPUT"`     |
| `output_value`     | Yes      | Comma-separated list of output values               | -       | `"gcp_success,aws_success"`   |
| `delimiter`        | No       | Delimiter for separating keys and values            | `,`     | `","`                         |
| `fail_on_empty`    | No       | Fail if any key or value is empty                  | `true`  | `"true"`                      |
| `trim_whitespace`  | No       | Trim whitespace from keys and values               | `true`  | `"true"`                      |
| `case_sensitive`   | No       | Treat keys as case sensitive                       | `true`  | `"true"`                      |
| `error_on_duplicate`| No      | Error if duplicate keys are found                  | `true`  | `"true"`                      |
| `mask_secrets`     | No       | Mask sensitive values in logs                      | `false` | `"true"`                      |
| `mask_pattern`     | No       | Custom pattern for masking (regex)                 | `""`    | `"(password\|secret).*"`      |
| `to_upper`         | No       | Convert values to uppercase                        | `false` | `"true"`                      |
| `to_lower`         | No       | Convert values to lowercase                        | `false` | `"true"`                      |
| `encode_url`       | No       | URL encode values                                  | `false` | `"true"`                      |
| `escape_newlines`  | No       | Escape newlines in values                         | `true`  | `"true"`                      |
| `max_length`       | No       | Maximum allowed length for values (0 for unlimited) | `0`     | `"10"`                        |
| `allow_empty`      | No       | Allow empty values even when fail_on_empty is true  | `false` | `"true"`                      |
| `debug_mode`       | No       | Enable debug logging for troubleshooting           | `false` | `"true"`                      |
| `group_prefix`     | No       | Prefix to group related environment variables      | `""`    | `"CONFIG"`                    |
| `json_support`     | No       | Enable JSON parsing for complex values             | `false` | `"true"`                      |
| `export_as_env`    | No       | Export output variables as environment variables   | `false` | `"true"`                      |

<br/>

## Outputs

| Output            | Description                           | Example        |
| ---------------- | ------------------------------------- | -------------- |
| `set_env_count`  | Number of environment variables set   | `3`            |
| `set_output_count`| Number of outputs set                | `3`            |
| `status`         | Status of the operation               | `"success"`    |
| `error_message`  | Error message if any                  | `""`           |

<br/>

### Example Workflow

Below is an example of how to use the **GitHub Environment/Output Setter**
action in a GitHub Actions workflow with all available options:

```yaml
name: Example Workflow
on: [push]

jobs:
  set-env-output:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout infrastructure repository
        uses: actions/checkout@v4

      - name: Set Environment and Output Variables
        id: set_variables
        uses: somaz94/env-output-setter@v1
        with:
          # Required inputs
          env_key: 'GCP_REGION,AWS_REGION'
          env_value: 'asia-northeast1,us-east-1'
          output_key: 'GCP_OUTPUT,AWS_OUTPUT'
          output_value: 'gcp_success,aws_success'
          
          # Optional inputs with defaults
          delimiter: ',' 
          fail_on_empty: 'true'
          trim_whitespace: 'true'
          case_sensitive: 'true'
          error_on_duplicate: 'true'

      - name: Display Variables and Status
        run: |
          # Environment Variables
          echo "GCP_REGION: ${{ env.GCP_REGION }}"
          echo "AWS_REGION: ${{ env.AWS_REGION }}"
          
          # Outputs
          echo "GCP_OUTPUT: ${{ steps.set_variables.outputs.GCP_OUTPUT }}"
          echo "AWS_OUTPUT: ${{ steps.set_variables.outputs.AWS_OUTPUT }}"
          
          # Action Results
          echo "Variables Set: ${{ steps.set_variables.outputs.set_env_count }}"
          echo "Outputs Set: ${{ steps.set_variables.outputs.set_output_count }}"
          echo "Status: ${{ steps.set_variables.outputs.status }}"
          echo "Error (if any): ${{ steps.set_variables.outputs.error_message }}"

      # Error handling
      - name: Check for Errors
        if: steps.set_variables.outputs.status == 'failure'
        run: |
          echo "Error occurred: ${{ steps.set_variables.outputs.error_message }}"
          exit 1
```

<br/>

## Features

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

## Advanced Usage

### Error Handling Examples

```yaml
# Handle empty values
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'KEY1,KEY2'
    env_value: ',value2'  # KEY1 is empty
    fail_on_empty: 'true' # This will fail the action

# Case sensitivity example
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'key1,Key1'
    env_value: 'value1,value2'
    case_sensitive: 'false' # This will treat key1 and Key1 as the same key

# Custom delimiter
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'key1;key2'
    env_value: 'value1;value2'
    delimiter: ';'
```

<br/>

### Common Use Cases

1. **Multi-Region Deployment**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'GCP_REGION,AWS_REGION,AZURE_REGION'
    env_value: 'asia-northeast1,us-east-1,eastasia'
```

2. **Environment Configuration**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'ENV,STAGE,VERSION'
    env_value: 'production,prod,v1.0.0'
```

3. **Status Tracking**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    output_key: 'DEPLOY_STATUS,TEST_STATUS'
    output_value: 'success,passed'
```

<br/>

### Value Transformation and Masking

This action supports value transformation and masking of sensitive data:

#### Features
- Mask sensitive values in logs
- Convert values to uppercase/lowercase
- URL encode values
- Custom masking patterns using regex
- Escape newlines in values
- Limit value lengths
- Handle empty values

#### Example Usage
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

#### Value Processing Behavior
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

### JSON Support and Complex Data Structures

This action now supports working with JSON values, allowing you to:
- Parse JSON objects and extract individual properties
- Create separate environment variables for nested JSON keys
- Process complex data structures including nested objects and arrays

When using JSON, select a delimiter that doesn't appear in your JSON content:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'CONFIG_JSON|SIMPLE_VALUE'
    env_value: '{"api_url":"https://api.example.com","timeout":30,"nested":{"key1":"value1","key2":123}}|simple_text'
    output_key: 'config_data|simple_output'
    output_value: '{"api_url":"https://api.example.com","timeout":30,"nested":{"key1":"value1","key2":123}}|simple_text'
    
    # Important: Use a different delimiter with JSON
    delimiter: '|'
    json_support: 'true'
```

This will create the following environment variables:
- `CONFIG_JSON`: The full JSON object
- `CONFIG_JSON_api_url`: "https://api.example.com"
- `CONFIG_JSON_timeout`: "30"
- `CONFIG_JSON_nested_key1`: "value1"
- `CONFIG_JSON_nested_key2`: "123"
- `SIMPLE_VALUE`: "simple_text"

<br/>

#### JSON Arrays and Complex Nesting

The action also fully supports JSON arrays and deeply nested structures:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'API_SETTINGS'
    env_value: '{"endpoints":[{"name":"users","path":"/api/users"},{"name":"products","path":"/api/products"}],"version":"2.0"}'
    json_support: 'true'
    debug_mode: 'true'
```

This will create:
- `API_SETTINGS`: The full JSON object
- `API_SETTINGS_endpoints_0_name`: "users"
- `API_SETTINGS_endpoints_0_path`: "/api/users"
- `API_SETTINGS_endpoints_1_name`: "products"
- `API_SETTINGS_endpoints_1_path`: "/api/products"
- `API_SETTINGS_version`: "2.0"

<br/>

#### Updated Group Prefix Behavior

The `group_prefix` option can be used to add a prefix to all environment variables:

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'COMPLEX_CONFIG|API_SETTINGS'
    env_value: '{"server":{"host":"example.com","port":8080}}|{"version":"2.0"}'
    delimiter: '|'
    json_support: 'true'
    group_prefix: 'APP'
```

When using `group_prefix`, variable names will include the prefix, but the JSON property extraction will maintain the original structure:

- `COMPLEX_CONFIG`: The full JSON object
- `COMPLEX_CONFIG_server_host`: "example.com"
- `COMPLEX_CONFIG_server_port`: "8080"

Rather than:
- `APP_COMPLEX_CONFIG_server_host`: "example.com"

This change makes variable naming more intuitive and consistent with other environment variables.

<br/>

#### Best Practices for Working with JSON

1. **Use a unique delimiter** that doesn't appear in your JSON content (pipe `|` is recommended)
2. **Keep JSON structures manageable** - very deep nesting can lead to long environment variable names
3. **Enable `debug_mode`** when first setting up to see exactly how variables are processed
4. **Validate your JSON** before using it in the action to avoid parsing errors
5. **Access nested properties directly** using the flattened naming convention (e.g., `${{ env.CONFIG_server_host }}`)

<br/>

### Group Prefix and Variable Organization

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

### Export Outputs as Environment Variables

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

### Advanced Usage Examples

1. **Handling Multiline Text**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'CONFIG_DATA'
    env_value: 'line1\nline2\nline3'
    escape_newlines: 'true'
```

2. **Length-Limited Values**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'DESCRIPTION'
    env_value: 'This is a very long description text'
    max_length: '20'
```

3. **Optional Values**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'OPTIONAL_VALUE,REQUIRED_VALUE'
    env_value: ',important_data'
    fail_on_empty: 'true'
    allow_empty: 'true'
```

<br/>

### Multiline and Special Characters

When working with multiline values or values containing commas, you need to consider the following:

1. **Multiline Input**
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

2. **Values Containing Commas**
If your values contain commas, you should use a different delimiter to avoid parsing issues:
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "KEY1::KEY2::KEY3"
    env_value: "value1, with comma::value2::value3"
    delimiter: "::"  # Use a different delimiter when values contain commas
```

3. **Working with JSON Data**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "CONFIG_JSON|SETTINGS_JSON"
    env_value: '{"server":"example.com","port":8080}|{"logging":true,"debug":false}'
    delimiter: "|"  # Use a delimiter that's not in your JSON
    json_support: "true"
```

4. **Complex JSON Structures**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "COMPLEX_JSON"
    env_value: '{"server":{"host":"example.com","port":8080},"auth":{"enabled":true,"methods":["oauth","basic"]}}'
    json_support: "true"
    group_prefix: "APP"
    debug_mode: "true"
```

#### Important Notes:
- When values contain the default delimiter (comma), use a different delimiter like `::`
- For JSON values, choose a delimiter that won't appear in your JSON data (e.g., `|`)
- Multiline values are automatically normalized
- Use `escape_newlines: true` to properly handle newline characters
- The same delimiter must be used consistently for both keys and values

<br/>

### Debugging and Troubleshooting

<br/>

#### Debug Mode
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

#### Empty Values
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

## Troubleshooting

Common issues and solutions:

1. **Duplicate Keys**
   - Error message: `duplicate key found: KEY_NAME`
   - Solution: Ensure all keys are unique or set `error_on_duplicate: 'false'`

2. **Empty Values**
   - Error message: `empty value found for key: KEY_NAME`
   - Solution: Provide values for all keys or set `fail_on_empty: 'false'`

3. **File Write Issues**
   - Error message: `failed to write to file`
   - Solution: Action will automatically retry up to 3 times

4. **JSON Parsing Errors**
   - Error message: `Invalid JSON format`
   - Solution: Ensure JSON strings are valid and properly escaped

5. **Delimiter Conflicts in JSON**
   - Error message: `env_key and env_value must have the same number of entries`
   - Solution: Use a delimiter that doesn't appear in your JSON (e.g., `|`)

<br/>

### Debug Output Format

When `debug_mode` is enabled, you'll see detailed information about how your inputs are being processed:

<br/>

## License

This project is licensed under the [MIT License](LICENSE) file for details.

<br/>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.