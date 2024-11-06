# GitHub Environment/Output Setter

[![License](https://img.shields.io/github/license/somaz94/env-output-setter)](https://github.com/somaz94/container-action)
![Latest Tag](https://img.shields.io/github/v/tag/somaz94/env-output-setter)
![Top Language](https://img.shields.io/github/languages/top/somaz94/env-output-setter?color=green&logo=go&logoColor=b)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Environment/Output%20Setter-blue?logo=github)](https://github.com/marketplace/actions/environment-output-setter)


## Overview

The **GitHub Environment/Output Setter** is a GitHub Action that allows you to set multiple key-value pairs in both `$GITHUB_ENV` and `$GITHUB_OUTPUT`. This action is useful for workflows that need to dynamically define environment variables or output values that other steps can reference.


## Inputs

| Input          | Required | Description                                                   | Example                     |
|----------------|----------|---------------------------------------------------------------|-----------------------------|
| `env_key`      | Yes      | Comma-separated list of environment variable keys             | `"GCP_REGION,AWS_REGION"`   |
| `env_value`    | Yes      | Comma-separated list of environment variable values           | `"asia-northeast1,us-east-1"` |
| `output_key`   | Yes      | Comma-separated list of output keys                           | `"GCP_OUTPUT,AWS_OUTPUT"`   |
| `output_value` | Yes      | Comma-separated list of output values                         | `"gcp_success,aws_success"` |

### Example Workflow

Below is an example of how to use the **GitHub Environment/Output Setter** action in a GitHub Actions workflow. This example sets environment variables and output variables and then prints the success message.

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
        uses: somaz94/github-env-output-setter@v1
        with:
          env_key: "GCP_REGION,AWS_REGION"
          env_value: "asia-northeast1,us-east-1"
          output_key: "GCP_OUTPUT,AWS_OUTPUT"
          output_value: "gcp_success,aws_success"

      - name: Display Env and Output Variables
        run: |
          echo "GCP_REGION: ${{ env.GCP_REGION }}" # asis-northeast1
          echo "AWS_REGION: ${{ env.AWS_REGION }}" # us-east-1
          echo "GCP_OUTPUT: ${{ steps.set_variables.outputs.GCP_OUTPUT }} # gcp_success
          echo "GCP_OUTPUT: ${{ steps.set_variables.outputs.AWS_OUTPUT }} # aws_success
```
