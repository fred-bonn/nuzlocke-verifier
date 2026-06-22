package main

import "log"

func vlog(v ...any) {
	if !*verbose {
		return
	}
	log.Print(v...)
}

func vlogln(v ...any) {
	if !*verbose {
		return
	}
	log.Println(v...)
}

func vlogf(format string, v ...any) {
	if !*verbose {
		return
	}
	log.Printf(format, v...)
}
