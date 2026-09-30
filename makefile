include config.make

out := dist/index.wasm
tileset_manifest := dist/tileset-manifest.json
void_module := github.com/oidoid/void
bundle_version = $(shell git describe --always --dirty)
bundle_published = $(shell TZ=UTC git log -1 --format=%cd --date=format-local:%Y%m%d)
bundle_id = $(bundle_version)+$(bundle_published)
tinygo_nodebug := --no-debug
go_tags := $(if $(value DEBUG),--tags=debug,)
# pick fastest CPU with `lscpu --extended`.
bench_test := \
	trap 'trap - exit int term; sudo cpupower set --turbo-boost=1' exit int term; \
	sudo cpupower set --turbo-boost=0; \
	GOMAXPROCS=1 powerprofilesctl launch --profile performance -- \
	taskset --cpu-list 3 \
	go test --bench=. --count=15 --cpu=1 --p=1 --parallel=1 --run='^$$' $(go_tags) ./src/...
go_test_filter = \
	grep --color=always --extended --line-buffered '^--- FAIL: [^ ]+|$$'| \
	sed --regexp-extended --unbuffered $(if $(value V),'','/^ok |\[no test files\]$$|PASS$$|^goos: |^goarch: |^pkg: |^cpu: /d')
# precise GC lets the collector skip scanning pointer-free heap objects instead
# of conservatively treating every word as a possible pointer.
tinygo_flags += $(go_tags) --ldflags="-X github.com/oidoid/void-template/src/game.Version=$(bundle_id)" --scheduler=none --gc=precise $(if $(value DEBUG),,$(tinygo_nodebug) --panic=trap) $(if $(value V),--print-allocs=.,)
# $(1) flags
pack = go tool $(void_module)/src/cmd/pack --out=dist/ --tsconfig=tsconfig.json $(1) src/web/assets/index.html
# $(1) flags
packatlas = go tool packatlas --name=atlas --img-out=dist/ --atlas-out=src/assets/atlas_bin.go --tags-out=src/tags/tags.go --tileset-manifest-out=$(tileset_manifest) $(1) src/assets/atlas/
# $(1) flags
packboards = go tool packboards --tileset-manifest=$(tileset_manifest) --out=src/boards $(1) src/assets/boards/
fat = go tool fat
fat_files := dist/atlas.webp dist/index.css dist/index.html dist/index.js dist/index.wasm
favicon = \
	mkdir -p dist/favicon; \
	for scale in 1 2 3 4 12 32; do \
		favicon=dist/favicon/favicon$$((scale * 16)); \
		aseprite src/web/assets/favicon.aseprite --batch --color-mode=indexed --scale=$$scale --save-as=$$favicon.png; \
		cwebp -exact -lossless -mt -quiet -z 9 $$favicon.png -o $$favicon.webp; \
	done

.PHONY: bench build build-atlas build-boards build-favicon build-go build-web check clean dependencies fat-analyze fat-check fat-save fmt fmt-go fmt-mod fmt-web install link-void lint lint-critic lint-static lint-vet lint-web test test-fmt-go test-fmt-mod test-go test-web typecheck-web watch watch-atlas watch-boards watch-go watch-web

watch: export DEBUG := 1
watch: dependencies build-atlas .WAIT build-boards build-favicon .WAIT watch-go watch-atlas watch-boards watch-web
watch-go:; watchexec --exts=go --quiet --watch=src/ -- $(MAKE) build-go
watch-atlas:; $(call packatlas,--watch)
watch-boards:; $(call packboards,--watch)
watch-web: link-void; $(call pack,--watch)

build: build-atlas .WAIT build-boards build-favicon .WAIT build-go build-web
build-go:
	# no concurrency.
	GOOS=wasip1 GOARCH=wasm tinygo build $(tinygo_flags) -o $(out) ./src/web/
	$(if $(value DEBUG),,wasm-opt -o $(out) -Oz --strip-debug --strip-producers $(out))
build-atlas:; $(call packatlas,)
build-boards: build-atlas; $(call packboards,)
build-favicon:; $(favicon)
build-web: build-go link-void; $(call pack,--minify --one-file)

clean:; rm --force --recursive dist/ src/assets/atlas_bin.go src/boards/ src/tags/tags.go

dependencies:
	for exe in aseprite cwebp go mono node shader_minifier.exe tinygo wasm-opt watchexec; do
		command -v $$exe > /dev/null || { echo "no $$exe" >&2; false; }
	done

fat-analyze: tinygo_nodebug :=
fat-analyze: tinygo_flags += --size full
fat-analyze: build
fat-check:; $(fat) check
fat-save:; $(fat) save $(fat_files)

fmt: fmt-mod fmt-go fmt-web
fmt-mod:; go mod tidy
fmt-go:; go fmt ./src/...
fmt-web: link-void; npx lint --fix > /dev/null

install:; go mod download; npm install; $(MAKE) link-void
link-void:
	if [ -d ../void ]; then
		printf 'go 1.27.1\n\nuse (\n\t.\n\t../void\n)\n' > go.work
	fi
	void_dir=$$(go list -m -f '{{.Dir}}' $(void_module))
	npm install \
		--no-audit \
		--no-fund \
		--no-save \
		--package-lock=false \
		--install-links=false \
		"$$void_dir/src/void/vweb"

lint: lint-critic lint-static lint-vet lint-web
lint-critic:; go tool go-critic check --enableAll --disable=unnamedResult ./src/...
lint-static:; go tool staticcheck ./src/...
lint-vet:; go vet ./src/...
lint-web: link-void; npx lint > /dev/null

check: test-fmt-go test-fmt-mod lint test-go test-web typecheck-web
test: dependencies .WAIT build .WAIT check .WAIT fat-check
test-fmt-go:
	out=$$(go fmt ./src/...)
	[ -z "$$out" ] || { printf >&2 "unformatted files:\n%s\n" "$$out"; false; }
test-fmt-mod:; go mod tidy -diff
test-go:; go test $(go_tags) ./src/... | $(go_test_filter)
test-web: link-void
	FORCE_COLOR=3 npm run test:unit|
	sed --unbuffered $(if $(value V),'','1,/✖ failing tests:/ {/[✔ℹ▶✖] /d}')
typecheck-web: link-void; npm run typecheck

bench:; $(bench_test) | $(go_test_filter)
