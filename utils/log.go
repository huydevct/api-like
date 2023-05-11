package utils

import (
	"fmt"
	"log"
	"os"
)

// LogFile :
type LogFile struct {
	client *os.File
}

// NewLogFile : Tạo mới đối tượng logFile
func NewLogFile() *LogFile {

	return &LogFile{
		client: initLog(),
	}
}

func initLog() (logFile *os.File) {

	// now := GetUTCNow()
	//
	dir := fmt.Sprintf("log")
	// Kiểm tra folder đả tồn tại hay chưa ?
	_, errCheckExisted := os.Stat(dir)
	if os.IsNotExist(errCheckExisted) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			err = fmt.Errorf("Create directory: %s fail: %s", dir, err)
			log.Fatal(err)
		}
	}
	// Tạo file log
	// documentPath := fmt.Sprintf("%s/%s.txt", dir, now.Format("02-01-2006"))
	documentPath := fmt.Sprintf("%s/api.log", dir)
	logFile, err := os.Create(documentPath)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Write ..
func (u *LogFile) Write(input interface{}) {
	u.client.WriteString(fmt.Sprintf("%v \n", input))
}

// Close ..
func (u *LogFile) Close() {
	u.client.Close()
}
