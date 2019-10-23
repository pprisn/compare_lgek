TARGET=compare_lgek2dbf.exe

all: clean build

clean:
	rm -rf $(TARGET)

build:
	go build -o $(TARGET) compare_lgek2dbf.go
