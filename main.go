package main

import (
	"fmt"
	"os"
)

func main() {
	to, domain := os.Getenv("DEMO_EMAIL_TO"), os.Getenv("SENDING_DOMAIN")
	if to == "" || domain == "" {
		panic("DEMO_EMAIL_TO and SENDING_DOMAIN are required")
	}
	c, err := newClient()
	if err != nil {
		panic(err)
	}
	id, err := sendSubscriberUpdate(c, domain, digitalAsset{Title: "Issue 12: Field Notes", Body: "Three new edits from this week's studio session."}, subscriber{Email: to})
	if err != nil {
		panic(err)
	}
	fmt.Println("subscriber update sent:", id)
}
