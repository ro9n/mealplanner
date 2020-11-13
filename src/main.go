package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"menuplanner/pkg"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Assumption 1. All three arguments are always passed
	if len(os.Args[1:]) < 3 {
		log.Fatal("Invalid argument length, expected 3")
	}

	query, err := pkg.NewQuery(os.Args[1], os.Args[2], os.Args[3])

	if err != nil {
		log.Fatal(err)
	}

	// "io/ioutil"
	// NOTE comment this line to enable logging
	log.SetOutput(ioutil.Discard)

	// Assumption 2. Data directory and it's relative ordering does not change
	users := query.Apply(filepath.Dir("../data"))

	fmt.Println(strings.Join(users[:], ", "))
	os.Exit(0)
}
