package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// fmt.Println(Jsonify([]string{"sd"}, "q"))
	// dd := func1()
	// fmt.Println(dd)
	fmt.Println(os.Hostname())
	fmt.Println(filepath.Base(os.Args[0]))

	SetDebugPro(true)
	// Perf("wwwwwwwwwwwww", func1)
	DbgPro("prefix", "[EXEC]%s", "ss")

}
