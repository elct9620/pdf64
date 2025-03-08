# PDF64

PDF64 is a tool for converting PDF documents to Base64 encoded images, making it convenient to display and process PDF content in web applications.

## Features

- Convert PDF files to Base64 encoded images
- Support for adjusting image density and quality
- Support for password-protected PDF files
- RESTful API interface
- Docker container support

## Requirements

- Go 1.23+
- ImageMagick 7
- Ghostscript
- QPDF

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
# Basic usage
curl -X POST \
  -F "data=@example.pdf" \
  -F "density=300" \
  -F "quality=90" \
  http://localhost:8080/v1/convert

# For password-protected PDF files
curl -X POST \
  -F "data=@encrypted.pdf" \
  -F "password=your_password" \
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

## License

This project is licensed under the [Apache License 2.0](LICENSE).
