ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    CP=copy
    RM=del
else
    SET=export
    CP=cp
    RM=rm
endif

NAME:=$(notdir $(CURDIR))
EXE:=$(shell go env GOEXE)
VERSION:=$(shell git describe --tags || echo v0.0.0)

TARGET:=$(NAME)$(EXE)
RUNLISP?="./$(TARGET)"
GENERATES=sort-world newtypes.go stringer.go

all:
	go fmt ./...
	go build
	go fmt cmd/gmnlisp/main.go
	go build -C cmd/gmnlisp -o "$(CURDIR)/$(TARGET)" -ldflags "-s -w -X main.version=$(VERSION)"

generate: $(GENERATES)

### test ###

test:
	$(RUNLISP) test.lsp
	go fmt
	go test

### Updating documents

README.md: tools/prepro.lsp _readme.md
	$(RUNLISP) $< < _README.md > $@

### Cleaning ###

clean:
	$(RM) $(TARGET)

clean-gen:
	$(RM) $(GENERATES)

### Formating sources ###

sort-world:
	$(RUNLISP) tools/gosort.lsp < world.go > world.go_
	-cmp world.go world.go_ || $(CP) world.go_ world.go
	$(RM) world.go_

### Generating sources ###

newtypes.go : tools/newtypes.lsp Makefile
	$(RUNLISP) $< $(NAME) "*StringBuilder" "*inputStream" "*_OutputFileStream" "*_Macro" "_ReaderNode" "_WriterNode" > $@

stringer.go : tools/stringer.lsp Makefile
	$(RUNLISP) $< $(NAME) ErrorNode Float Integer _WriterNode _ReaderNode _Macro _OutputFileStream inputStream _JoinedForm LispString SpecialF _Lambda _TrueType Cons Keyword Rune _NullType Array Function _Hash > $@

### Packaging ###

_package:
	$(SET) "CGO_ENABLED=0" && $(MAKE) clean && $(MAKE) all
	zip -9 $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(TARGET)

package:
	$(SET) "GOOS=linux" && $(SET) "GOARCH=386"   && $(MAKE) _package
	$(SET) "GOOS=linux" && $(SET) "GOARCH=amd64" && $(MAKE) _package
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _package
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _package

manifest:
	go run mkmanifest.go *-windows-*.zip > gmnlisp.json

release:
	gh release create -d --notes "" -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

.PHONY: generate test clean clean-gen sort-world _package package manifest
