package Utils

import "log"

func Errhandle(e interface{}) {
	if e != nil {
		log.Fatal(e)
	}
}
