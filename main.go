package main

import (
	"log"

	"github.com/topxeq/xsend/v9/src/cli"
)

func main() {
	// "github.com/pkg/profile"
	// go func() {
	// 	for {
	// 		f, err := os.Create("croc.pprof")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		runtime.GC() // get up-to-date statistics
	// 		if err := pprof.WriteHeapProfile(f); err != nil {
	// 			panic(err)
	// 		}
	// 		f.Close()
	// 		time.Sleep(3 * time.Second)
	// 		fmt.Println("wrote profile")
	// 	}
	// }()
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
