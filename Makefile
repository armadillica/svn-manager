all: mocks

mocks: go-get svnman/svnman_mock.go

go-get:
	go get github.com/golang/mock/mockgen

clean:
	rm -v */*_mock.go

%_mock.go: %.go go-get
	@# notdir-realpath-dir takes the last directory component.
	mockgen -package $(notdir $(realpath $(dir $<))) -source $< -destination $@

.PHONY: mocks clean
