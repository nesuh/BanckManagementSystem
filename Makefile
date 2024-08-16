build:
	go build -o D:\Go\BankManage\bin\BankManage

run: build
	.\bin\BankManage

test:
	go test -v ./...
