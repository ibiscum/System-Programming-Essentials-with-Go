package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	// ...

	f, err := os.Create("memprofile.out")
	if err != nil {
		// Handle error
		log.Fatal(err)
	}
	defer f.Close()
	runtime.GC()
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	// ... (Rest of your code)
}
