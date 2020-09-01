mkdir -p releases

GOOS=darwin GOARCH=amd64 packr build && mv ./gore ./releases/darwin-gore \
  && GOOS=linux GOARCH=amd64 packr build && mv ./gore ./releases/linux-gore \
  && GOOS=windows GOARCH=386 packr build && mv ./gore.exe ./releases/gore.exe \
  && packr clean