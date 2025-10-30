# JSON Support

<br/>


## Overview

This action now supports working with JSON values, allowing you to:
- Parse JSON objects and extract individual properties
- Create separate environment variables for nested JSON keys
- Process complex data structures including nested objects and arrays

<br/>

## Basic JSON Usage

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

## JSON Arrays and Complex Nesting

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

## Working with JSON and Group Prefix

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

## Best Practices for Working with JSON

1. **Use a unique delimiter** that doesn't appear in your JSON content (pipe `|` is recommended)
2. **Keep JSON structures manageable** - very deep nesting can lead to long environment variable names
3. **Enable `debug_mode`** when first setting up to see exactly how variables are processed
4. **Validate your JSON** before using it in the action to avoid parsing errors
5. **Access nested properties directly** using the flattened naming convention (e.g., `${{ env.CONFIG_server_host }}`)

<br/>

## Complex JSON Example

```yaml
- uses: somaz94/env-output-setter@v1
  with:
    env_key: "COMPLEX_JSON"
    env_value: '{"server":{"host":"example.com","port":8080},"auth":{"enabled":true,"methods":["oauth","basic"]}}'
    json_support: "true"
    group_prefix: "APP"
    debug_mode: "true"
```

This creates:
- `COMPLEX_JSON`: Full JSON
- `COMPLEX_JSON_server_host`: "example.com"
- `COMPLEX_JSON_server_port`: "8080"
- `COMPLEX_JSON_auth_enabled`: "true"
- `COMPLEX_JSON_auth_methods_0`: "oauth"
- `COMPLEX_JSON_auth_methods_1`: "basic"
