source ~/.bash_profile

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o abime_lookup_liunx64
mv abime_lookup_liunx64 ~/Documents/abime_lookup_liunx64