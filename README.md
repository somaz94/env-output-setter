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

### Common Use Cases

1. **Configuration File Processing**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'CONFIG_CONTENT'
    env_value: ${{ steps.read_config.outputs.content }}
    escape_newlines: 'true'
    max_length: '1000'
```

2. **API Response Handling**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'API_RESPONSE'
    env_value: ${{ steps.api_call.outputs.response }}
    max_length: '500'
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

3. **Multiline Values with Special Characters**
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "MULTILINE_TEXT::CONFIG_JSON"
    env_value: |
      Hello\nWorld::{"key": "value", "array": [1,2,3]}
    delimiter: "::"
    escape_newlines: true
```

4. **Multiline Values with Special Characters (\n)**
```yaml
      - name: Run Performance Analysis
        id: analysis
        uses: somaz94/github-action-analyzer@v1
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          workflow_file: ci.yml
          repository: ${{ github.repository }}
          analysis_depth: '20'
          timeout: '15'
          ignore_patterns: 'checkout,setup'

      - name: Set Analysis Results
        uses: somaz94/env-output-setter@v1
        with:
          env_key: |
            METRICS_SUMMARY
            PERFORMANCE_SUMMARY
            CACHE_RECOMMENDATIONS
            DOCKER_OPTIMIZATIONS
            STATUS
          env_value: |
            ${{ steps.analysis.outputs.metrics_summary }}
            ${{ steps.analysis.outputs.performance_summary }}
            ${{ steps.analysis.outputs.cache_recommendations }}
            ${{ steps.analysis.outputs.docker_optimizations }}
            ${{ steps.analysis.outputs.status }}
          delimiter: "\n"
          trim_whitespace: true
          debug_mode: true
```

#### Important Notes:
- When values contain the default delimiter (comma), use a different delimiter like `::`
- Multiline values are automatically normalized
- Use `escape_newlines: true` to properly handle newline characters
- The same delimiter must be used consistently for both keys and values

<br/>

### Debugging and Troubleshooting

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

<br/>

### Debug Output Format

When `debug_mode` is enabled, you'll see detailed information about how your inputs are being processed:

```
🔍 Debug Information (Env)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📥 Input Values:
  • Keys:      "MULTILINE_TEXT::MAX_LENGTH_TEST::EMPTY_VALUE"
  • Values:    "Hello\nWorld::ThisIsAVeryLongTextThatShouldBeTruncated::   "
  • Delimiter: "::"

📋 Processed Values:
  • Keys:   [MULTILINE_TEXT MAX_LENGTH_TEST EMPTY_VALUE]
  • Values: [Hello World ThisIsAVeryLongTextThatShouldBeTruncated ]

✍️  Writing Values:
  • env: MULTILINE_TEXT = Hello Worl
  • env: MAX_LENGTH_TEST = ThisIsAVer
  • env: EMPTY_VALUE = 
```

Even without debug mode, you'll still see the basic operation output:

```
==================================================
🚀 Setting Env Variables
  • env: MULTILINE_TEXT = Hello Worl
  • env: MAX_LENGTH_TEST = ThisIsAVer
  • env: EMPTY_VALUE = 

==================================================
✅ Execution Complete
Mode: GitHub Actions
```

This helps you understand:
- How your inputs are being processed
- What transformations are being applied
- The final values being set
- Any issues that might arise during processing

### Output Colors

The action uses colors in the console output to help distinguish different types of information:
- 🔵 Blue: Information and section headers
- 🟢 Green: Successful operations
- 🔴 Red: Errors and warnings

Note: Colors may not be visible in all CI environments or when output is redirected to a file.

<br/>

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<br/>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

