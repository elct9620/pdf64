# PDF64

PDF64 is a tool for converting PDF documents to Base64 encoded images, making it convenient to display and process PDF content in web applications.

## Features

- Convert PDF files to Base64 encoded images
- Support for adjusting image density and quality
- RESTful API interface
- Docker container support

## Requirements

- Go 1.23+
- ImageMagick 7
- Ghostscript

## Installation

### Using Docker

```bash
docker pull ghcr.io/elct9620/pdf64:latest
docker run -p 8080:8080 ghcr.io/elct9620/pdf64:latest
```

### Building from Source

```bash
git clone https://github.com/elct9620/pdf64.git
cd pdf64
go build -o pdf64 ./cmd
```

## Usage

### API Usage

```bash
curl -X POST \
  -F "data=@example.pdf" \
  -F "density=300" \
  -F "quality=90" \
  http://localhost:8080/v1/convert
```

### Response Format

```json
{
  "id": "unique-file-id",
  "data": [
    "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcG...",
    "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcG..."
  ]
}
```

## Development

```bash
# Run tests
go test -v -cover ./...

# Run linter
golangci-lint run
```

## Project Structure

- `cmd/`: Main application entry point
- `pkg/`: Public packages
  - `apis/`: API interface definitions
- `internal/`: Private packages
  - `app/`: Application configuration
  - `controller/`: API controllers
  - `builder/`: Factory pattern implementations
  - `entity/`: Domain entities
  - `service/`: Service implementations
  - `usecase/`: Business logic

## License

This project is licensed under the MIT License.
