package main

import (
	"fmt"
	"github.com/Flacy/fext/utils"
)

func main() {
	size := utils.GetDirSize("/usr/lib/python3.8/site-packages/requests")
	fmt.Println(size / 1024)
}
