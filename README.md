# GitHub Environment/Output Setter

[![License](https://img.shields.io/github/license/somaz94/env-output-setter)](https://github.com/somaz94/container-action)
![Latest Tag](https://img.shields.io/github/v/tag/somaz94/env-output-setter)
![Top Language](https://img.shields.io/github/languages/top/somaz94/env-output-setter?color=green&logo=go&logoColor=b)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Environment/Output%20Setter-blue?logo=github)](https://github.com/marketplace/actions/environment-output-setter)

## Overview

The **GitHub Environment/Output Setter** is a GitHub Action that allows you to set multiple key-value pairs in both `$GITHUB_ENV` and `$GITHUB_OUTPUT`. This action is useful for workflows that need to dynamically define environment variables or output values that other steps can reference.

<br/>

## Quick Start

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: 'GCP_REGION,AWS_REGION'
    env_value: 'asia-northeast1,us-east-1'
    output_key: 'GCP_OUTPUT,AWS_OUTPUT'
    output_value: 'gcp_success,aws_success'
```

<br/>

## Documentation

- [Usage Guide](docs/USAGE.md) - Comprehensive examples and common use cases
- [Features](docs/FEATURES.md) - Feature overview and advanced usage
- [JSON Support](docs/JSON_SUPPORT.md) - Working with JSON data structures
- [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues and solutions
- [API Reference](docs/API.md) - Complete inputs and outputs reference

<br/>

## Key Features

- Set multiple environment variables and outputs in one step
- Value transformation (uppercase, lowercase, URL encoding)
- Mask sensitive values in logs
- JSON support for complex data structures
- Group related variables with prefixes
- Retry mechanism for file operations
- Debug mode for troubleshooting

<br/>

## Basic Example

```yaml
- uses: somaz94/env-output-setter@v1
  id: set_vars
  with:
    env_key: 'GCP_REGION,AWS_REGION'
    env_value: 'asia-northeast1,us-east-1'
    output_key: 'DEPLOY_STATUS'
    output_value: 'success'

- name: Use variables
  run: |
    echo "GCP: ${{ env.GCP_REGION }}"
    echo "AWS: ${{ env.AWS_REGION }}"
    echo "Status: ${{ steps.set_vars.outputs.DEPLOY_STATUS }}"
```

For more examples, see the [Usage Guide](docs/USAGE.md).

<br/>

<br/>

## Development

```bash
make build        # Build the binary
make test         # Run unit tests with coverage
make test-local   # Run local integration tests
make test-all     # Run all tests
make cover        # Generate coverage report
make cover-html   # Open coverage in browser
make bench        # Run benchmarks
make lint         # Run go vet
make fmt          # Format code
make clean        # Remove build artifacts
```

<br/>

## License

This project is licensed under the [Apache License 2.0](LICENSE).

<br/>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

For detailed information, please refer to the documentation in the [docs](docs/) folder.