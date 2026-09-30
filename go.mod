module github.com/oidoid/void-template

go 1.27.1

tool github.com/go-critic/go-critic/cmd/go-critic

tool (
	github.com/oidoid/void/src/cmd/fat
	github.com/oidoid/void/src/cmd/pack
	github.com/oidoid/void/src/cmd/packatlas
	github.com/oidoid/void/src/cmd/packboards
	honnef.co/go/tools/cmd/staticcheck
)

require github.com/oidoid/void v0.1.11-0.20260928040410-c39ea3ca21da

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/cristalhq/acmd v0.12.0 // indirect
	github.com/evanw/esbuild v0.28.2 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-critic/go-critic v0.14.4 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.2.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/pkgload v1.2.2 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/quasilyte/go-ruleguard v0.4.5 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	golang.org/x/exp/typeparams v0.0.0-20260820142414-ca536658362e // indirect
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/tools v0.49.0 // indirect
	honnef.co/go/tools v0.8.1 // indirect
)
