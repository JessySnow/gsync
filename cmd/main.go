package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	config2 "gync/internal/config"
	"io"
	"os"
	"path/filepath"
)

var config *config2.Config

func init() {
	path, err := os.Executable()
	if err != nil {
		fmt.Println("can't get executable path")
		os.Exit(1)
	}

	dir := filepath.Dir(path)
	fileName := "gsync.json"
	filePath := filepath.Join(dir, fileName)
	configFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("can't open config file: %v\n", err)
		os.Exit(1)
	}

	defer configFile.Close()

	bytes, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Printf("read config file failed: %v\n", err)
		os.Exit(1)
	}

	config, err = config2.Parse(bytes)
	if err != nil {
		fmt.Printf("parse config file failed: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	gsync := figure.NewFigure("GSYNC", "", true)
	gsync.Print()
}
