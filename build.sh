mkdir -p releases

GOOS=darwin GOARCH=amd64 packr build
tar -czvf ./releases/darwin-gore.tar.gz ./gore
rm ./gore

GOOS=linux GOARCH=amd64 packr build
tar -czvf ./releases/linux-gore.tar.gz ./gore
rm ./gore

GOOS=windows GOARCH=386 packr build
tar -czvf ./releases/gore-exe.tar.gz ./gore.exe
rm ./gore.exe

packr clean