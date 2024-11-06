# GitHub Environment/Output Setter

[![License](https://img.shields.io/github/license/somaz94/ienv-output-setter)](https://github.com/somaz94/container-action)
![Latest Tag](https://img.shields.io/github/v/tag/somaz94/env-output-setter)
![Top Language](https://img.shields.io/github/languages/top/somaz94/env-output-setter?color=green&logo=go&logoColor=b)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Environment/Output%20Setter-blue?logo=github)](https://github.com/marketplace/actions/env-output-setter)


## Overview

The **GitHub Environment/Output Setter** is a GitHub Action that allows you to set multiple key-value pairs in both `$GITHUB_ENV` and `$GITHUB_OUTPUT`. This action is useful for workflows that need to dynamically define environment variables or output values that other steps can reference.


## Inputs

| Input          | Required | Description                                                   | Example                     |
|----------------|----------|---------------------------------------------------------------|-----------------------------|
| `env_key`      | Yes      | Comma-separated list of environment variable keys             | `"GCP_REGION,AWS_REGION"`   |
| `env_value`    | Yes      | Comma-separated list of environment variable values           | `"asia-northeast1,us-east-1"` |
| `output_key`   | Yes      | Comma-separated list of output keys                           | `"GCP_OUTPUT,AWS_OUTPUT"`   |
| `output_value` | Yes      | Comma-separated list of output values                         | `"gcp_success,aws_success"` |

## Outputs

| Output            | Description                                                                |
|-------------------|----------------------------------------------------------------------------|
| `success_message` | Confirmation message indicating successful setting of environment and output variables. |


### Example Workflow

Below is an example of how to use the **GitHub Environment/Output Setter** action in a GitHub Actions workflow. This example sets environment variables and output variables and then prints the success message.

```yaml
name: Example Workflow
on: [push]

jobs:
  set-env-output:
    runs-on: ubuntu-latest
    steps:
      - name: Set Environment and Output Variables
        id: set_variables
        uses: somaz94/github-env-output-setter@v1
        with:
          env_key: "GCP_REGION,AWS_REGION"
          env_value: "asia-northeast1,us-east-1"
          output_key: "GCP_OUTPUT,AWS_OUTPUT"
          output_value: "gcp_success,aws_success"

      - name: Display Success Message
        run: |
          echo "Success: ${{ steps.set_variables.outputs.success_message }}"

    outputs:
      GCP_REGION: ${{ steps.run.outputs.GCP_OUTPUT }}
      AWS_OUTPUT: ${{ steps.run.outputs.AWS_OUTPUT }}
```

### Additional Information

- **Icon**: settings
- **Color**: blue
- **Author**: somaz94

This action is packaged in a Docker container, making it portable and easy to run on any compatible GitHub Actions runner. 