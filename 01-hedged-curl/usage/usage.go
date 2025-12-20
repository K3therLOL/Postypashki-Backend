package usage

import (
	"os"
	"flag"
	"fmt"
)

var Help bool
var Time int

func Reference() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [-t SECONDS] URL1 [URL2 ...]\n", os.Args[0])    
	flag.Usage()
}

func Arguments() []string {
	return flag.Args()
}

func init() {
	flag.BoolVar(&Help, "h", false, "output help reference")
	flag.BoolVar(&Help, "help", false, "output help reference (long form)")
	flag.IntVar(&Time, "t", 15, "set timeout for all HTTP requests in seconds")
	flag.IntVar(&Time, "time", 15, "set timeout for all HTTP requests in seconds (long form)")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "OPTIONS")
		flag.VisitAll(func(f *flag.Flag) {
			var flagName string
			if len(f.Name) > 1 {
				flagName = "--" + f.Name
			} else {
				flagName = "-" + f.Name
			}

			fmt.Fprintf(flag.CommandLine.Output(), "%6s\t%s\n", flagName, f.Usage)
		})
	}
	
	flag.Parse()
}
