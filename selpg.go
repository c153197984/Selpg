package main

import (
    "flag"
    "fmt"
    "os"
	"bufio"
	"os/exec"
	"strings"
)

type selpgArgs struct {
	start int
	end   int
	pagelen   int
	page_seperator bool
	dest string
	inputFile string
}

var progname string

func getSelpg(arg *selpgArgs) {
	flag.IntVar(&(arg.start), "s", -1, "Start page..")
	flag.IntVar(&(arg.end), "e", -1, "End page.")
	flag.IntVar(&(arg.pagelen), "l", 72, "Line number per page.")
	flag.BoolVar(&(arg.page_seperator), "f", false, "use [-f=true] to seperate pages")
	flag.StringVar(&(arg.dest), "d", "", "Destionation of output.")
	flag.Parse()

	if arg.start == -1 || arg.end == -1 {
		printErr("not enough arguments!")
	}
	if len(flag.Args()) > 1 {
		printErr("can only read one file!")
	}
	if arg.start < 1 {
		printErr("start page can't less than 1!")
	}
	if arg.end < arg.start {
		printErr("start page should <= end page!")
	}
	if arg.pagelen < 1 {
		printErr("page length can't less than 1!")
	}

	if len(flag.Args()) == 0 {
		arg.inputFile = ""
	} else {
		arg.inputFile = flag.Args()[0]
	}
}

func printErr(err string) {
	fmt.Fprintf(os.Stderr, err+"\n"+
		"\nUSAGE: %s -s start_page -e end_page [-f=true|false | -l pageline] [-d dest] [inputFileName]\n", progname)
	os.Exit(1)
}

func run() {
	var args selpgArgs
	progname := os.Args[0]
	getSelpg(&args)

	var err error
	var cmd *exec.Cmd
	fin := os.Stdin
	fout := os.Stdout
	if args.inputFile != "" {
		fin, err = os.Open(args.inputFile)
		if err != nil {
			printErr("can't open file\"" + args.inputFile + "\"!")
		}
	}
	if args.dest != "" {
		t := fmt.Sprintf("%s", args.dest)
		cmd = exec.Command("sh", "-c", t)
		if err != nil {
			printErr("can't open file to \"" + t + "\"!")
		}
	}

	inReader := bufio.NewReader(fin)
	outfile_text := ""
	pageNum := 1
	lineNum := 0
	if args.page_seperator == false {
		var line string
		for true {
			line, err = inReader.ReadString('\n')
			if err != nil {
				break
			}
			lineNum++
			if lineNum > args.pagelen {
				pageNum++
				lineNum = 1
			}
			if pageNum >= args.start && pageNum <= args.end {
				if args.dest == "" {
					fmt.Fprintf(fout, "%s", line)
				} else {
					outfile_text += line
				}
			}
		}
	} else {
		for true {
			c, _, erro := inReader.ReadRune()
			if erro != nil {
				break
			}
			if c == '\f' {
				pageNum++
			}
			if pageNum >= args.start && pageNum <= args.end {
				if args.dest == "" {
					fmt.Fprintf(fout, "%c", c)
				} else {
					outfile_text += string(c)
				}
			}
		}
	}

	if args.dest != "" {
		cmd.Stdin = strings.NewReader(outfile_text)
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			printErr("print error!")
		}
	}

	if pageNum < args.start {
		fmt.Fprintf(os.Stderr, "%s: start_page greater than total pages, no output\n", progname)
	}
	if pageNum >= args.start && pageNum < args.end {
		fmt.Fprintf(os.Stderr, "%s: end_page greater than total pages, less output\n", progname)
	}

	fmt.Fprintf(os.Stderr, "%s: Commond Done\n", progname)
	fin.Close()
	fout.Close()
}

func main() {
	run()
}
