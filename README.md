:warning: This API is unstable, and results/formatting/etc. may change unexpectedly until a stable release.

# Kube Version API
This is a simple API for quickly getting version information related to Kubernetes.

## Endpoints

### `GET /versions`
Get all available versions.

#### Query Parameters
| Parameter | Description                               |
|-----------|-------------------------------------------|
| `major`   | The major version to filter on (eg: `1`)  |
| `minor`   | The minor version to filter on (eg: `21`) |

#### Examples

**Query**
```
https://api.kubeversion.com/versions?major=1&minor=21
```
**Output**
```json
{
  "versions": [
    "1.21.0",
    "1.21.1",
    "1.21.2",
    "1.21.3",
    "1.21.4",
    "1.21.5",
    "1.21.6",
    "1.21.7",
    "1.21.8",
    "1.21.9",
    "1.21.10",
    "1.21.11",
    "1.21.12"
  ]
}
```


### `GET /versions/latest`
Returns the latest version.

#### Query Parameters
| Parameter | Description                               |
|-----------|-------------------------------------------|
| `major`   | The major version to filter on (eg: `1`)  |
| `minor`   | The minor version to filter on (eg: `21`) |

#### Examples

**Query**
```
https://api.kubeversion.com/versions/latest?major=1&minor=21
```
**Output**
```json
{
  "versions": "1.21.12"
}
```

**Query**
```
https://api.kubeversion.com/versions/latest?major=1
```
**Output**
```json
{
  "versions": "1.23.6"
}
```