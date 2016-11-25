package util

import (
	"github.com/donnie4w/go-logger/logger"
	"runtime"
)

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		logger.Error(x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			logger.Error("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			logger.Error("EXRAS#%v DATA:%v\n", k, extras[k])
		}
	}
}
