package main

import (
	"fmt"
	"log"
)

func vlogln(v ...any) {
	if !*verbose {
		return
	}
	fmt.Println(v...)
}

func vlogf(format string, v ...any) {
	if !*verbose {
		return
	}
	fmt.Printf(format+"\n", v...)
}

func vlogWithPrefix(prefix, format string, v ...any) {
	if !*verbose {
		return
	}
	fmt.Printf("%s %s\n", prefix, fmt.Sprintf(format, v...))
}

func vlogMove(prio, speed int, format string, v ...any) {
	vlogWithPrefix(fmt.Sprintf("[MOVE, %d/%d]", prio, speed), format, v...)
}

func vlogSwitch(format string, v ...any) {
	vlogWithPrefix("[SWITCH]", format, v...)
}

func vlogReplace(format string, v ...any) {
	vlogWithPrefix("[REPLACE]", format, v...)
}

func vlogItem(format string, v ...any) {
	vlogWithPrefix("[ITEM]", format, v...)
}

func elogln(v ...any) {
	log.Println(v...)
}

func elogf(format string, v ...any) {
	log.Printf(format, v...)
}

func elogFatalf(format string, v ...any) {
	log.Fatalf(format, v...)
}
