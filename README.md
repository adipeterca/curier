# curier 🚚

A small Go server for sharing files across the internet.

## How to setup

### Docker environment

You can download the Dockerfile and build the image yourself, or simply pull it from the repo:
```bash
docker pull ghcr.io/adipeterca/curier:latest
docker run -p 8080:8080 ghcr.io/adipeterca/curier:latest
```

### Linux environment

**I strongly recommend using Docker, as it simplifies the configuration a lot**.

If you want to use a precompiled binary, please refer to the [Release](https://github.com/adipeterca/curier/releases) section.

### Windows environment

Because not many servers run Windows, the support I can provide for this platform is limited.
You can download a precompiled binary from the [Release](https://github.com/adipeterca/curier/releases) section or use a Docker container.

## Configuration

You can configure some aspects of the service via environment variables prefixed with **CURIER_**.  
Some information will be exposed via the `/config/` endpoint for better UX.  
For default values, check [config.go](https://github.com/adipeterca/curier/blob/main/config.go).  

| Variable name | Description |
|--|--|
|`CURIER_STORAGE_PATH`|Absolute path where the file uploads will be stored on disk|
|`CURIER_BASE_URL`|Default prefix for the `/download/{id}` URL|
|`CURIER_HOST`|Network address to bind to|
|`CURIER_PORT`|Port to use (also affects the port used inside the container)|
|`CURIER_MAX_FILE_SIZE`|Maximum allowed size for each file upload|
|`CURIER_ALLOWED_FILE_EXTENSIONS`|A `;` separated list of file extensions, lowercase only (the last entry needs to have a `;` too)|