package main

import "github.com/changefouk/finalexam/customer"

func main() {
	r := customer.SetupRouter()
	r.Run(":2019")
}
