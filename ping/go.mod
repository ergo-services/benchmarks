module ping

go 1.20

replace ergo.services/ergo => ../../ergo

replace ergo.services/logger/colored => /home/taras/devel/ergo.services/logger/colored

require (
	ergo.services/ergo v0.0.0-00010101000000-000000000000
	ergo.services/logger/colored v0.0.0-00010101000000-000000000000
	github.com/klauspost/cpuid/v2 v2.2.6
)

require (
	github.com/fatih/color v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
