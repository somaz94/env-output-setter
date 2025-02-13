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

## Outputs

| Output            | Description                           | Example        |
| ---------------- | ------------------------------------- | -------------- |
| `set_env_count`  | Number of environment variables set   | `3`            |
| `set_output_count`| Number of outputs set                | `3`            |
| `status`         | Status of the operation               | `"success"`    |
| `error_message`  | Error message if any                  | `""`           |

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

## Features

- Set multiple environment variables and outputs in one step
- Configurable delimiter for key-value pairs
- Whitespace trimming option
- Case sensitivity control for keys
- Duplicate key detection
- Empty value validation
- Detailed operation status and error reporting
- Retry mechanism for file operations

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

### Value Transformation and Masking

This action supports value transformation and masking of sensitive data:

#### Features
- Mask sensitive values in logs
- Convert values to uppercase/lowercase
- URL encode values
- Custom masking patterns using regex

#### Example Usage
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'API_KEY,USERNAME,PASSWORD'
    env_value: 'secret123,admin,pass123'
    mask_secrets: 'true' # Masks sensitive values
    mask_pattern: '(password|secret).' # Custom masking pattern
    to_upper: 'true' # Convert to uppercase
    to_lower: 'false' # Convert to lowercase
    encode_url: 'false' # URL encode values
```

#### Masking Behavior
- When `mask_secrets` is enabled, sensitive values are masked in logs
- Custom `mask_pattern` allows regex-based masking
- Default masking shows first 2 characters followed by asterisks
- Short values (<4 characters) are completely masked

#### Value Transformations
- `to_upper`: Converts values to uppercase
- `to_lower`: Converts values to lowercase  
- `encode_url`: Applies URL encoding to values
- Transformations are applied in order: case conversion -> URL encoding

Note: Masking only affects log output, not the actual values set in environment variables or outputs.

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.