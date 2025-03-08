package builder

import (
	"github.com/elct9620/pdf64/internal/entity"
	"github.com/google/uuid"
)

// FileBuilder 實現 usecase.FileBuilder 接口
type FileBuilder struct{}

// NewFileBuilder 創建一個新的 FileBuilder 實例
func NewFileBuilder() *FileBuilder {
	return &FileBuilder{}
}

// BuildFromPath 從文件路徑創建 File 實體
func (b *FileBuilder) BuildFromPath(path string) (*entity.File, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	
	return entity.NewFile(id.String(), path), nil
}
