# PDF64

PDF64 是一個將 PDF 文件轉換為 Base64 編碼圖像的工具，方便在網頁應用中顯示和處理 PDF 內容。

## 功能特點

- 將 PDF 文件轉換為 Base64 編碼的圖像
- 支持調整圖像密度和質量
- RESTful API 接口
- Docker 容器支持

## 系統需求

- Go 1.23+
- ImageMagick 7
- Ghostscript

## 安裝

### 使用 Docker

```bash
docker pull ghcr.io/elct9620/pdf64:latest
docker run -p 8080:8080 ghcr.io/elct9620/pdf64:latest
```

### 從源碼構建

```bash
git clone https://github.com/elct9620/pdf64.git
cd pdf64
go build -o pdf64 ./cmd
```

## 使用方法

### API 使用

```bash
curl -X POST \
  -F "data=@example.pdf" \
  -F "density=300" \
  -F "quality=90" \
  http://localhost:8080/v1/convert
```

### 響應格式

```json
{
  "id": "unique-file-id",
  "data": [
    "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcG...",
    "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcG..."
  ]
}
```

## 開發

```bash
# 運行測試
go test -v -cover ./...

# 運行 linter
golangci-lint run
```

## 項目結構

- `cmd/`: 主應用程序入口
- `pkg/`: 公共包
  - `apis/`: API 接口定義
- `internal/`: 私有包
  - `app/`: 應用程序配置
  - `controller/`: API 控制器
  - `builder/`: 工廠模式實現
  - `entity/`: 領域實體
  - `service/`: 服務實現
  - `usecase/`: 業務邏輯

## 授權

本項目採用 MIT 授權。
