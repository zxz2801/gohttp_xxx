package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// default data
const (
	FileRotateSize  = 1024 * 1024 * 500
	FileRotateCount = 5
	DateFormat      = "2006-01-02" // DateFormat ...
)

// rotate type
const (
	RotateNull = iota
	RotateDate
	RotateSize
)

// FileWriter implement interface LogOutput
type FileWriter struct {
	// config meta
	FilePath    string `json:"filepath"`
	FileName    string `json:"filename"`
	RotateSize  uint64 `json:"rotatesize"`
	RotateCount int16  `json:"rotatecount"`
	// contorl meta
	mux      *sync.Mutex
	file     *os.File  // current file fd
	fileDate time.Time // current file date
	curName  string    // current full name
	curSize  int64     // current file size
}

var (
	gwriter *FileWriter
	once    sync.Once
)

func (w *FileWriter) write(when time.Time, msg []byte) (int, error) {
	w.checkFile(when)
	w.mux.Lock()
	defer w.mux.Unlock()
	n, err := w.file.Write(msg)
	if err != nil {
		return n, err
	}
	w.curSize += int64(n)
	return n, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("Write[%d]", len(p))
	return w.write(time.Now(), p)
}

// Destroy close file
func (w *FileWriter) Destroy() {
	w.file.Close()
}

func (w *FileWriter) checkFile(when time.Time) {
	rotate := w.checkRotate(when)
	switch rotate {
	case RotateDate:
		w.createNewFile()
	case RotateSize:
		w.rotateNewFile()
	default:
		return
	}
}

func (w *FileWriter) checkRotate(when time.Time) int {
	if w.checkRotateTime(when) {
		return RotateDate
	}
	if w.checkRotateSize() {
		return RotateSize
	}
	return RotateNull
}

func (w *FileWriter) checkRotateTime(when time.Time) bool {
	ts, _ := time.Parse(DateFormat, when.Format(DateFormat))
	return ts.After(w.fileDate)
}

func (w *FileWriter) checkRotateSize() bool {
	return uint64(w.curSize) > w.RotateSize
}

func (w *FileWriter) rotateNewFile() error {

	w.mux.Lock()
	defer w.mux.Unlock()

	// check is conflict with rotateTime
	if !w.checkRotateSize() {
		return fmt.Errorf("Please check rotate confilct")
	}
	if w.file == nil {
		return fmt.Errorf("File is nil when rotate")
	}

	for i := w.RotateCount - 1; i >= 0; i-- {
		oldName := fmt.Sprintf("%s.%d", w.curName, i)
		newName := fmt.Sprintf("%s.%d", w.curName, i+1)
		_, err := os.Stat(oldName)
		if err != nil {
			continue
		}
		err = os.Rename(oldName, newName)
		if err != nil {
			fmt.Println("Rename error.", err.Error())
		}
	}

	w.file.Close()
	err := os.Rename(w.curName, fmt.Sprintf("%s.%d", w.curName, 0))
	if err != nil {
		return fmt.Errorf("Rename error. %s", err.Error())
	}
	file, err := os.OpenFile(w.curName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("OpenFile error. %s", err.Error())
	}
	w.file = file
	w.curSize = 0

	return nil
}

func (w *FileWriter) createNewFile() error {
	w.mux.Lock()
	defer w.mux.Unlock()

	// check is conflict with rotatesize
	if !w.checkRotateTime(time.Now()) {
		return fmt.Errorf("Please check rotate confilct")
	}
	if w.file != nil {
		w.file.Close()
	}

	now, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	today := fmt.Sprintf("%.4d%.2d%.2d", now.Year(), now.Month(), now.Day())
	w.curName = fmt.Sprintf("%s/%s.%s", w.FilePath, w.FileName, today)
	file, err := os.OpenFile(w.curName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	w.file = file
	w.fileDate = now
	w.curSize = GetFileSize(w.curName)
	return nil
}

// GetFileSize for check file size
func GetFileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}
