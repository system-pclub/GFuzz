
## Benchmark

### Prerequistes
- Docker
- cloc (setup instruction can be found at http://cloc.sourceforge.net)
- Git
  
### 1. Clone Repositories

```bash
$ clone-repos.sh ./repos
```

### 2. Calculate LOC

```bash
$ loc.sh ./repos
```

### 3. Calculate Number of Tests
```bash
// First, run builder, it will
// 1. create builder docker
// 2. compile all test binary for given repositories at tmp/builder
$ ./build.sh

// Second, run script to count specific app. Same way for the others
$ ./benchmark.sh count-tests --dir /builder/etcd/inst
```


### 4. Calculate Performance
```bash

# If you run before, skip it.
$ ./build.sh

# /builder is the mapped directory of host directory 'tmp/builder', which is output of ./build.sh
$ ./benchmark.sh benchmark --dir /builder/grpc-go/native --mode native --out /builder/out/grpc-go-native.out
$ ./benchmark.sh benchmark --dir /builder/grpc-go/inst --mode inst --out /builder/out/grpc-go-inst.out

# After you have both results, compare common parts of them
$ ./filter.py ../builder/out/grpc-go-native.out ../builder/out/grpc-go-inst.out
```