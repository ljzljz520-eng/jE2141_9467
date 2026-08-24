# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	volunteerhours/cmd/volunteerhours	[no test files]
?   	volunteerhours/report	[no test files]
--- FAIL: TestInvalidVolunteerHoursAreRejected (0.01s)
    runner_test.go:61: expected invalid hours error
FAIL
FAIL	volunteerhours/cli	0.042s
ok  	volunteerhours/domain	0.002s
ok  	volunteerhours/query	0.002s
ok  	volunteerhours/service	0.040s
ok  	volunteerhours/storage	0.027s
ok  	volunteerhours/workflow	0.073s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/volunteerhours): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/volunteerhours): exit `0`
