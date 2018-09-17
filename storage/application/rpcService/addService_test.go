package rpcService

import (
	"context"
	"github.com/NBSChain/go-nbs/utils"
	"github.com/NBSChain/go-nbs/utils/cmdKits/pb"
	"os"
	"path/filepath"
	"testing"
)

//TODO:: test more file size and more rpc chunk size.
func TestImportLittleFile(t *testing.T) {
	fileName := "1.jpg"
	testImportFile(t, fileName, "QmPGVr4zCr9yibh3bEoswTkx3STpHPnXV1tLzwZeRwxqA1")
}
func TestImportNormalFile(t *testing.T) {
	fileName := "3.jpg"
	testImportFile(t, fileName, "QmWnD4GNz8gUsf14sp9JCJx8cLXdLdosUKd25PW5ftgd77")
}

func testImportFile(t *testing.T, fileName string, target string) {

	fileInfo, ok := utils.FileExists(fileName)
	if !ok {
		t.Error("File is not available.")
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Error("Failed to open file")
	}
	fullPath, err := filepath.Abs(fileName)
	if err != nil {
		t.Error(err)
	}

	request := &pb.AddRequest{
		FileName:     fileName,
		FullPath:     fullPath,
		FileSize:     fileInfo.Size(),
		FileType:     pb.FileType_FILE,
		SplitterSize: SplitterSize,
	}

	request.FileData = make([]byte, fileInfo.Size())
	file.Read(request.FileData)

	service := &addService{}

	response, err := service.AddFile(context.Background(), request)
	if err != nil {
		t.Error(err)
	}

	if response.Hash != target {
		t.Error("Hash 值计算错误!")
	} else {
		t.Log("测试通过")
	}

}
