# .PHONY: all clean

all: curr linuxStatic linuxArm64 winAmd64

# setup:
# 	sed -i 's/const debug = !false/const debug = false/' main.go
#
# clean:
# 	sed -i 's/const debug = false/const debug = !false/' main.go

curr:
	go build -ldflags "-s -w" -o build/

install:
	@echo "Installing in to the system"
	@go install -ldflags "-s -w"
linuxStatic:
	@echo "[+] Building the static Linux version"
	@env GOOS=linux CGO_ENABLED=0 go build -ldflags "-s -w" -o build/clock.static

linuxArm64:
	@echo "[+] Building the Linux ARM64 version"
	@env GOARCH=arm64 GOOS=linux go build -ldflags "-s -w" -o build/clock.arm

winAmd64:
	@echo "[+] Building the Windows version"
	@env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/clock.exe

# This target is always executed, even if there are errors in previous targets.
.PHONY: always
always: clean
