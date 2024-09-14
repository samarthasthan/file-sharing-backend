package main

import (
	"github.com/samarthasthan/21BRS1248_Backend/api"
)

func main() {
	f := api.NewFiberHandler()
	f.Handle()
	err := f.Start("1248")
	if err != nil {
		panic(err)
	}
}
