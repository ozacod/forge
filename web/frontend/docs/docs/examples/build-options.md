# Build Options

Common build variants:
```bash
cpx build                 # debug
cpx build --release       # release (-O2)
cpx build -O3             # max opt
cpx build -j 8            # parallel
cpx build --clean         # clean rebuild
```

Sanitizers for debugging tricky bugs:
```bash
cpx build --asan
cpx build --tsan
cpx build --msan
cpx build --ubsan
```
Use one sanitizer at a time.
