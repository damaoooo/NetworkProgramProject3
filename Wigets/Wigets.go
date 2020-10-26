package Wigets

import "log"

func ErrHandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}
