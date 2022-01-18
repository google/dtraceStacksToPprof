# DTraceStacksToPprof

This is a tool used to convert performance profiles collected from
[DTrace](http://dtrace.org/) to [pprof](http://github.com/google/pprof),
enabling graphical representation of the frequency of occuring stacks.

## Getting started

First clone the repo,

```
$ git clone https://github.com/google/dtraceStacksToPprof.git
```

The tool requires Go, which can be downloaded at the [Go homepage](https://golang.org/)

`dtraceStacksToPprof` can be installed to the `GOPATH` using
```
go install github.com/google/dtraceStacksToPprof
```

or run directly in the repo using
```
go run main.go
```

## How to use

The converter supports stack output produced by DTrace, in the following pseudo BNF syntax:

```
  <input> ::= <prologue> <stacks>
  <prologue> = { JUNK_LINE }
  <stacks> ::= { <stack> }
  <stack> ::= STACK_HEADER <frames> COUNT
  <frames> ::= { <frame> }
  <frame> ::= HEX_ADDRESS | FUNCTION
  JUNK_LINE ::= (string not matching STACK_HEADER)
  STACK_HEADER ::= ^\s*[\S\s]+\s*:$  (a sentence ending with a colon, optionally
                                      surrounded by white space)
  COUNT ::= ^\s*[0-9]+\s*$  (a sole decimal digit, optionally surrounded by
                             white space)
  HEX_ADDRESS ::= ^\s*0x[0-9a-fA-F]+\s*$ (a hexadecimal address, optionally
                                          surrounded by white space)
  FUNCTION ::= ^\s*[\S]+`.*\+?\S*$ (non-whitespace, backtick, non-whitespace,
                                    followed by an optional plus sign and
                                    optional non-whitespace after, all
                                    optionally surrounded by white space)
```

It's designed to take the output of dtrace scripts like the following:
```
#!/bin/dtrace -s
syscall:::entry {
  @aggr["Label:", ustack()] = count();
}
tick-1s {
  printa(@aggr);
}
```

The `dtraceStacksToPprof` tool reads text from stdin and produces a pprof
protobuf. To produce a pprof directly from dtrace, simply pipe dtrace's output
into `dtraceStacksToPprof` like
```
$ sudo dtrace -s <script name> -p <target pid> | dtraceStacksToPprof
```
or store to an intermediate text file and then convert.
```
$ dtraceStacksToPprof < <intermediate file>
```
The tool writes to `profile.pb.gz` by default, but the output file name can be
customized via the `--output` command line flag:
```
$ dtraceStacksToPprof --output InterestingBehavior.pb.gz < <intermediate file>
```

# Disclaimer
This is not an officially supported Google product.
