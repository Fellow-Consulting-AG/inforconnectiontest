package main

import "log"

// debugPrint only prints logs if debugMode is enabled
func debugPrint(format string, v ...interface{}) {
	if debugMode {
		log.Printf(format, v...)
	}
}
