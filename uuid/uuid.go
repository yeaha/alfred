package main

import (
	"fmt"
	"strings"

	"github.com/ruedap/go-alfred"
	"github.com/satori/go.uuid"
)

func main() {
	u := uuid.NewV4()

	r := alfred.NewResponse()

	r.AddItem(&alfred.ResponseItem{
		Valid: true,
		Title: u.String(),
	})

	r.AddItem(&alfred.ResponseItem{
		Valid: true,
		Title: strings.ToUpper(u.String()),
	})

	xml, err := r.ToXML()
	if err != nil {
		title := fmt.Sprintf("Error: %v", err.Error())
		subtitle := "UUID Workflow Error"
		arg := title
		errXML := alfred.ErrorXML(title, subtitle, arg)
		fmt.Println(errXML)
	}

	fmt.Println(xml)
}
