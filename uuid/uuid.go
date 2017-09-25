package main

import (
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
)

func main() {
	u := uuid.NewV4()

	fmt.Println(u.String())
	fmt.Println(strings.ToUpper(u.String()))
}
