# selpg
[selpg](https://github.com/zhengbaic/selpg) 是由GO执行的简单CLI



### a struct of selpg
```
type selpgArgs struct {
	start int
	end   int
	pagelen   int
	page_seperator bool
	dest string
	inputFile string
}
```
### init selpg
```
flag.IntVar(&(arg.start), "s", -1, "Start page..")
flag.IntVar(&(arg.end), "e", -1, "End page.")
flag.IntVar(&(arg.pagelen), "l", 72, "Line number per page.")
flag.BoolVar(&(arg.page_seperator), "f", false, "use [-f=true] to seperate pages")
flag.StringVar(&(arg.dest), "d", "", "Destionation of output.")
flag.Parse()
```
###对selpg命令输入进行简单的合理判断
```
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
```

###main process
通过pageNum和lineNumN记录当前页数和行数，将inputfile的中的内容按照命令行的要求进行输出。
```
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
}
```
同时可以通过pageNum来进行error判断，如，selpg -s 2 -e 10  inputfile 2>error_file 指令中。
代码如下：
```
if pageNum < args.start {
	fmt.Fprintf(os.Stderr, "%s: start_page greater than total pages, no output\n", progname)
}
if pageNum >= args.start && pageNum < args.end {
	fmt.Fprintf(os.Stderr, "%s: end_page greater than total pages, less output\n", progname)
}
```


#Bye 
