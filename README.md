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
| Variable name | Default value | Description |
|--|--|--|
|`CURIER_STORAGE_PATH`|`/var/lib/curier/uploads/` (Linux/Docker), `?` (Windows)|Absolute path where the file uploads will be stored on disk|
|`CURIER_BASE_URL`|`http://localhost`|Default prefix for the `/download/{id}` URL|
|`CURIER_HOST`|`127.0.0.1`|Network address to bind to|
|`CURIER_PORT`|`8080`|Port to use|