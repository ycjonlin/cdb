test: test/test.ps main
	./main test test/test.ps
	go run test/main/main.go

main: main.go pkg/*
	go build main.go

bench:
	go test -bench=.

clean:
	rm test/*.go