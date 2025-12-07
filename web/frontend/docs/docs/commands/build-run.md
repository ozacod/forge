# Build & Run Commands

Build, run, test, and sanity-check with one CLI.

## cpx build
Uses CMake presets when available.
```bash
cpx build                     # debug
cpx build --release           # release (-O2)
cpx build -O3                 # max opt
cpx build -j 8                # parallel
cpx build --clean             # clean then build
cpx build --target <name>     # specific target
cpx build --asan|--tsan|--msan|--ubsan   # sanitizer builds
```

## cpx run
Build then run. Arguments after `--` go to your program.
```bash
cpx run
cpx run --release
cpx run -- --flag value
```

## cpx test
Build and run tests (CTest).
```bash
cpx test
cpx test -v
cpx test --filter <pattern>
```

## cpx bench
Build and run benchmarks.
```bash
cpx bench            # build + run benchmarks
cpx bench --verbose  # show verbose build output
```
Supports Google Benchmark, Catch2 Benchmark, and nanobench frameworks.

## cpx check
Quick sanity builds with sanitizers.
```bash
cpx check --asan   # memory errors
cpx check --tsan   # data races
cpx check --ubsan  # undefined behavior
cpx check --msan   # uninitialized
```
Use one sanitizer at a time; expect slower runs.

## cpx clean
```bash
cpx clean          # remove build artifacts
cpx clean --all    # also remove generated files
```
