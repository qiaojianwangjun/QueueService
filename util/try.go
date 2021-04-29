package util

import "fmt"

type Exception struct {
	ErrorMsg interface{}
}

func (e *Exception) Error() string {
	return fmt.Sprint(e.ErrorMsg)
}

// Try 尝试执行，类似try/catch
func Try(action func(), catch func(err error)) {
	defer func() {
		if err := recover(); err != nil {
			if catch != nil {
				err2, ok := err.(error)
				if !ok {
					err2 = &Exception{ErrorMsg: err}
				}
				catch(err2)
			}
		}
	}()
	action()
}
