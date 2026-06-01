# curier 🚚

A small Go server for sharing files across the internet.

## How to setup

### Linux environment

_I need to add this part, via a dedicated Bash script_

### Windows environment

_I need to add this part, via a dedicated Powershell script_

### Docker environment

_I need to add this part, either as a Dockerfile or an already built container_

## Configuration

You can configure some aspects of the service via environment variables prefixed with **CURIER_**.
Some information will be exposed via the `/config/` endpoint for better UX.
For default values, check [config.go](https://github.com/adipeterca/curier/blob/main/config.go).

| Variable name | Description |
|--|--|
|`CURIER_STORAGE_PATH`|Absolute path where the file uploads will be stored on disk|
|`CURIER_BASE_URL`|Default prefix for the `/download/{id}` URL|
|`CURIER_HOST`|Network address to bind to|
|`CURIER_PORT`|Port to use|
|`CURIER_MAX_FILE_SIZE`|Maximum allowed size for each file upload|
|`CURIER_ALLOWED_FILE_EXTENSIONS`|A `;` separated list of file extensions, lowercase only (the last entry needs to have a `;` too)|