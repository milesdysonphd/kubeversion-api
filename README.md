:warning: This API is unstable, and results/formatting/etc. may change unexpectedly until a stable release.

# Kube Version API
This is a simple API for quickly getting version information related to Kubernetes. The original intention was to quickly and easily get the latest patch version for a given minor release for downloading kubectl (in pipelines, iac, etc).

## TODO
- [ ] Expand binary download information (only kubectl is listed)

## Endpoints

| Method | Endpoint                                           | Description                               |
|--------|----------------------------------------------------|-------------------------------------------|
| `GET`  | [`/v1/versions`](#get-v1versions)                  | Get available versions.                   |
| `GET`  | [`/v1/versions/latest`](#get-v1versionslatest)     | Get the latest version available.         |
| `GET`  | [`/v1/versions/:version`](#get-v1versionsversions) | Get the specified version (if available). |

---

### `GET /v1/versions`
Get all available versions.

#### Query Parameters
| Parameter | Description                               |
|-----------|-------------------------------------------|
| `major`   | The major version to filter on (eg: `1`)  |
| `minor`   | The minor version to filter on (eg: `21`) |

#### Examples

**Query**
```
https://api.kubeversion.com/v1/versions?major=1&minor=21
```
**Output**
```json
{
  "data": {
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
}
```

---

### `GET /v1/versions/:version`
Returns the specified `version` if available.

#### Examples

**Query**
```
https://api.kubeversion.com/v1/versions/v1.21.12
```
**Output**
```json
{
	"data": {
		"version": "1.21.12",
		"downloads": {
			"darwin": [
				{
					"binary": "kubectl",
					"platform": "darwin",
					"architecture": "amd64",
					"url": "https://dl.k8s.io/v1.21.12/bin/darwin/amd64/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/darwin/amd64/kubectl.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "darwin",
					"architecture": "arm64",
					"url": "https://dl.k8s.io/v1.21.12/bin/darwin/arm64/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/darwin/arm64/kubectl.sha256"
				}
			],
			"linux": [
				{
					"binary": "kubectl",
					"platform": "linux",
					"architecture": "386",
					"url": "https://dl.k8s.io/v1.21.12/bin/linux/386/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/386/kubectl.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "linux",
					"architecture": "amd64",
					"url": "https://dl.k8s.io/v1.21.12/bin/linux/amd64/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/amd64/kubectl.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "linux",
					"architecture": "arm",
					"url": "https://dl.k8s.io/v1.21.12/bin/linux/arm/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/arm/kubectl.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "linux",
					"architecture": "arm64",
					"url": "https://dl.k8s.io/v1.21.12/bin/linux/arm64/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/arm64/kubectl.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "linux",
					"architecture": "ppc64le",
					"url": "https://dl.k8s.io/v1.21.12/bin/linux/ppc64le/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/ppc64le/kubectl.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "linux",
					"architecture": "s390x",
					"url": "https://dl.k8s.io/v1.21.12/bin/linux/s390x/kubectl",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/s390x/kubectl.sha256"
				}
			],
			"windows": [
				{
					"binary": "kubectl",
					"platform": "windows",
					"architecture": "386",
					"url": "https://dl.k8s.io/v1.21.12/bin/windows/386/kubectl.exe",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/windows/386/kubectl.exe.sha256"
				},
				{
					"binary": "kubectl",
					"platform": "windows",
					"architecture": "amd64",
					"url": "https://dl.k8s.io/v1.21.12/bin/windows/amd64/kubectl.exe",
					"checksumUrl": "https://dl.k8s.io/v1.21.12/bin/windows/amd64/kubectl.exe.sha256"
				}
			]
		}
	}
}
```

---

### `GET /v1/versions/latest`
Returns the latest version, and related kubectl download information.

#### Query Parameters
| Parameter | Description                               |
|-----------|-------------------------------------------|
| `major`   | The major version to filter on (eg: `1`)  |
| `minor`   | The minor version to filter on (eg: `21`) |

#### Examples

**Query**
```
https://api.kubeversion.com/v1/versions/latest
```
**Output**
```json
{
  "data": {
    "version": "1.23.6",
    "downloads": {
      "darwin": [
        {
          "binary": "kubectl",
          "platform": "darwin",
          "architecture": "amd64",
          "url": "https://dl.k8s.io/v1.23.6/bin/darwin/amd64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/darwin/amd64/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "darwin",
          "architecture": "arm64",
          "url": "https://dl.k8s.io/v1.23.6/bin/darwin/arm64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/darwin/arm64/kubectl.sha256"
        }
      ],
      "linux": [
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "386",
          "url": "https://dl.k8s.io/v1.23.6/bin/linux/386/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/linux/386/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "amd64",
          "url": "https://dl.k8s.io/v1.23.6/bin/linux/amd64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/linux/amd64/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "arm",
          "url": "https://dl.k8s.io/v1.23.6/bin/linux/arm/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/linux/arm/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "arm64",
          "url": "https://dl.k8s.io/v1.23.6/bin/linux/arm64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/linux/arm64/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "ppc64le",
          "url": "https://dl.k8s.io/v1.23.6/bin/linux/ppc64le/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/linux/ppc64le/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "s390x",
          "url": "https://dl.k8s.io/v1.23.6/bin/linux/s390x/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/linux/s390x/kubectl.sha256"
        }
      ],
      "windows": [
        {
          "binary": "kubectl",
          "platform": "windows",
          "architecture": "386",
          "url": "https://dl.k8s.io/v1.23.6/bin/windows/386/kubectl.exe",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/windows/386/kubectl.exe.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "windows",
          "architecture": "amd64",
          "url": "https://dl.k8s.io/v1.23.6/bin/windows/amd64/kubectl.exe",
          "checksumUrl": "https://dl.k8s.io/v1.23.6/bin/windows/amd64/kubectl.exe.sha256"
        }
      ]
    }
  }
}
```

**Query**
```
https://api.kubeversion.com/v1/versions/latest?major=1&minor=21
```
**Output**
```json
{
  "data": {
    "version": "1.21.12",
    "downloads": {
      "darwin": [
        {
          "binary": "kubectl",
          "platform": "darwin",
          "architecture": "amd64",
          "url": "https://dl.k8s.io/v1.21.12/bin/darwin/amd64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/darwin/amd64/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "darwin",
          "architecture": "arm64",
          "url": "https://dl.k8s.io/v1.21.12/bin/darwin/arm64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/darwin/arm64/kubectl.sha256"
        }
      ],
      "linux": [
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "386",
          "url": "https://dl.k8s.io/v1.21.12/bin/linux/386/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/386/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "amd64",
          "url": "https://dl.k8s.io/v1.21.12/bin/linux/amd64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/amd64/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "arm",
          "url": "https://dl.k8s.io/v1.21.12/bin/linux/arm/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/arm/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "arm64",
          "url": "https://dl.k8s.io/v1.21.12/bin/linux/arm64/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/arm64/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "ppc64le",
          "url": "https://dl.k8s.io/v1.21.12/bin/linux/ppc64le/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/ppc64le/kubectl.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "linux",
          "architecture": "s390x",
          "url": "https://dl.k8s.io/v1.21.12/bin/linux/s390x/kubectl",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/linux/s390x/kubectl.sha256"
        }
      ],
      "windows": [
        {
          "binary": "kubectl",
          "platform": "windows",
          "architecture": "386",
          "url": "https://dl.k8s.io/v1.21.12/bin/windows/386/kubectl.exe",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/windows/386/kubectl.exe.sha256"
        },
        {
          "binary": "kubectl",
          "platform": "windows",
          "architecture": "amd64",
          "url": "https://dl.k8s.io/v1.21.12/bin/windows/amd64/kubectl.exe",
          "checksumUrl": "https://dl.k8s.io/v1.21.12/bin/windows/amd64/kubectl.exe.sha256"
        }
      ]
    }
  }
}
```