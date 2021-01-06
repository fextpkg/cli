package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("C:\\Users\\hz\\AppData\\Local\\Programs\\Python\\Python39\\Lib\\site-packages\\")
	defer file.Close()
	names, _ := file.Readdirnames(0)
	fmt.Println(names)
}
