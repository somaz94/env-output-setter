# Usage Guide

<br/>

## Basic Usage

The **GitHub Environment/Output Setter** allows you to set multiple key-value pairs in both `$GITHUB_ENV` and `$GITHUB_OUTPUT`. This action is useful for workflows that need to dynamically define environment variables or output values that other steps can reference.

<br/>

## Example Workflow

Below is an example of how to use the **GitHub Environment/Output Setter** action in a GitHub Actions workflow with all available options:

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
          echo "Status: ${{ steps.set_variables.outputs.action_status }}"
          echo "Error (if any): ${{ steps.set_variables.outputs.error_message }}"

      # Error handling
      - name: Check for Errors
        if: steps.set_variables.outputs.action_status == 'failure'
        run: |
          echo "Error occurred: ${{ steps.set_variables.outputs.error_message }}"
          exit 1
```

<br/>

## Common Use Cases

### 1. Multi-Region Deployment
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'GCP_REGION,AWS_REGION,AZURE_REGION'
    env_value: 'asia-northeast1,us-east-1,eastasia'
```

### 2. Environment Configuration
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'ENV,STAGE,VERSION'
    env_value: 'production,prod,v1.0.0'
```

### 3. Status Tracking
```yaml
- uses: somaz94/env-output-setter@v1
  with:
    output_key: 'DEPLOY_STATUS,TEST_STATUS'
    output_value: 'success,passed'
```

<br/>

## Error Handling Examples

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
