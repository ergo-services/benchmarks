module memusage

replace ergo.services/ergo => ../../ergo3

replace ergo.services/logger/colored => ../../logger/colored

go 1.20

require (
	ergo.services/ergo v0.0.0-00010101000000-000000000000
	ergo.services/logger/colored v0.0.0-00010101000000-000000000000
)

require (
	github.com/fatih/color v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	golang.org/x/sys v0.10.0 // indirect
)