package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	logFile *os.File
	 
	sync.Mutex
}

var Default *Logger

func InitDefault(path string) error {
	var err error
	Default,err = NewLogger(path)
	return err
}

func NewLogger(path_to_logfile string) (*Logger,error) {
	// Opening or creating the logfile
	f,err := os.Create(path_to_logfile)
	if err != nil {
		return nil,err
	}

	ret := new(Logger)
	ret.logFile = f

	return ret,nil
}

func (l *Logger) Log(level string,format string, args ...any) {
	l.Lock()
	formated := fmt.Sprintf(format,args...)
	str := fmt.Sprintf("[%s]\t%s\t: "+formated+"\n", level, time.Now().Format(time.DateTime))
	fmt.Fprint(l.logFile,str)
	l.Unlock()
	fmt.Print(str)
}

func (l *Logger) LogPanic(format string, args ...any,) {
	l.Lock()
	formated := fmt.Sprintf(format,args...)
	str := fmt.Sprintf("[FATAL]\t%s\t: "+formated+"\n", time.Now().Format(time.DateTime))
	fmt.Fprint(l.logFile,str)
	l.Unlock()
	panic(str)
}