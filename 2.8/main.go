package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	time, err := ntp.Time("pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get the time: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(time.Format("2006-01-02 15:04:05 MST"))
}
