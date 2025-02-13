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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.