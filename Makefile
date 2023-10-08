AppName = cj
AppVer = v0.1.3

Builder = go
Flags = --ldflags='-s -w'
Dist = dist



run:
	$(Builder) run .
	$(Builder) run . 額意
	$(Builder) run . 的 我 人

build: preClean build4windows build4mac build4linux afterClean

build4windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(Builder) build $(Flags) -o "$(Dist)/$(AppName)_$(AppVer)_windows_amd64/cj.exe"
	cd "$(Dist)" && zip -r "$(AppName)_$(AppVer)_windows_amd64.zip" "$(AppName)_$(AppVer)_windows_amd64"

build4mac:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(Builder) build $(Flags) -o "$(Dist)/$(AppName)_$(AppVer)_darwin_amd64/cj"
	cd "$(Dist)" && tar zcvf "$(AppName)_$(AppVer)_darwin_amd64.tar.gz" "$(AppName)_$(AppVer)_darwin_amd64"

build4linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(Builder) build $(Flags) -o "$(Dist)/$(AppName)_$(AppVer)_linux_amd64/cj"
	cd "$(Dist)" && tar zcvf "$(AppName)_$(AppVer)_linux_amd64.tar.gz" "$(AppName)_$(AppVer)_linux_amd64"

preClean:
	rm -rf $(Dist)

afterClean:
	find $(Dist) -mindepth 1 -type d -exec rm -rf {} \+


