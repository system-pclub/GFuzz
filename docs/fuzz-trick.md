# Some tricks for fuzzing specific targets

Some of the unit tests will check external goroutines leak, some of them may
break testing before oracle checking blocking bugs. 

## General Tricks
- If a test requires tasks like listoning port that are exclusive, add --parallel 1.


## prometheus

util/testutil/testing.go

```go
func TolerantVerifyLeak(m *testing.M) {
	m.Run()
    // comment out leak verification
}
```

## gRPC

internal/leakcheck/leakcheck.go

```go
func check(efer Errorfer, timeout time.Duration) {
    // comment out or add oraclert as exception
}
```