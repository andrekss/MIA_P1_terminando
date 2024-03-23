package main

import (
	"os"
	"strings"
)

func main() {

	switch { // uso de comandos
	case strings.EqualFold(os.Args[1], "mkdisk"):
		MKDisk()
		os.Exit(1)
	case strings.EqualFold(os.Args[1], "rmdisk"):
		RMDisk()
		os.Exit(1)
	case strings.EqualFold(os.Args[1], "fdisk"):
		FDisk()
		os.Exit(1)
	}

}
