# gnark-fibonacci-with-constans
## Installation
Fallow https://go.dev/doc/install for installation golang.

## Overview
You can use the following commands to run circuit.
```bash
go mod tidy
go mod run
```
This circuit computes nth number of fibonacci numbers "a" and "b" where a and b are secret and "result" is public.

In order to change secrets you can edit following lines
```go
firstSecret := *big.NewInt(0)
secondSecret := *big.NewInt(1)
```

In order to change result you can change the following line
```go
result := "354224848179261915075"
```