# API Reference

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
| `action_status`  | Status of the operation               | `"success"`    |
| `error_message`  | Error message if any                  | `""`           |
