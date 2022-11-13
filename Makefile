commsbar:
	go build bin/comms/commsbar.go

primarybar:
	./build.sh bin/primary/primarybar.go

all: primarybar commsbar

