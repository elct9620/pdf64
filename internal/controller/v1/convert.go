package v1

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/elct9620/pdf64/internal/usecase"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
	"github.com/google/uuid"
)

func (s *Service) Convert(ctx context.Context, req *v1.ConvertRequest) (*v1.ConvertResponse, error) {
	convert := usecase.NewConvertUsecase()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	// 創建臨時檔案
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, id.String()+".pdf")
	
	// 建立臨時檔案
	tmpFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()
	
	// 將上傳的檔案內容複製到臨時檔案
	_, err = io.Copy(tmpFile, req.File)
	if err != nil {
		return nil, err
	}
	
	// 確保檔案被寫入並關閉
	tmpFile.Close()
	
	// 在函數結束時刪除臨時檔案
	defer os.Remove(filePath)
	
	// 執行轉換用例
	out, err := convert.Execute(ctx, &usecase.ConvertInput{
		FileId:   id.String(),
		FilePath: filePath,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConvertResponse{
		Id:   out.FileId,
		Data: out.EncodedImages,
	}, nil
}
