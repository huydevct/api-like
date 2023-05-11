package utils

import (
	"fmt"
	"log"
	"os"
)

// LogFileService :
type LogFileService struct {
	client *os.File
}

// NewLogFileService : Tạo mới đối tượng LogFileService
func NewLogFileService(name string) *LogFileService {

	return &LogFileService{
		client: initLogService(name),
	}
}

func initLogService(name string) (LogFileService *os.File) {

	now := GetUTCNow()
	//
	dir := fmt.Sprintf("log/%s/%s", name, now.Format("02-01-2006"))
	// Kiểm tra folder đả tồn tại hay chưa ?
	_, errCheckExisted := os.Stat(dir)
	if os.IsNotExist(errCheckExisted) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			err = fmt.Errorf("Create directory: %s fail: %s", dir, err)
			log.Fatal(err)
		}
	}
	// Tạo file log
	documentPath := fmt.Sprintf("%s/%s.txt", dir, now.Format("15:04"))
	LogFileService, err := os.Create(documentPath)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Write ..
func (u *LogFileService) Write(input interface{}) {
	fmt.Println(fmt.Sprintf("%v", input))
	u.client.WriteString(fmt.Sprintf("%v \n", input))
}

// Close ..
func (u *LogFileService) Close() {
	u.client.Close()
}
