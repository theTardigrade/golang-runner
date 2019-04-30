# runner

Allows commands from the Golang bin directory to be run indefinitely, or for a given number of iterations, while storing a log of any errors in the temporary directory.

## example

```bash
$(echo $GOBIN)/runner -command=test -arguments="--one --two --three" -log -verbose -sleep=10ms -iterations=5
```
