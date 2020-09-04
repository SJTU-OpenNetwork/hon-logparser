.PNOHY:parser
.PHONY:parser-win
parser:
	go build -o parser
parser-win:
	go build -o parser.exe
