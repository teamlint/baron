# How to develop baron

## Dependencies

1. Everything required to install `baron`
2. go-bindata for compiling templates into binary
```
$ go get github.com/jteeuwen/go-bindata/...
```

## Building

Whenever templates are modified, the templates must be recompiled to binary,
this is done with:

```
$ go generate github.com/teamlint/baron/...
```

Then to build baron and its protoc plugin to your $GOPATH/bin directory:

```
$ go install github.com/teamlint/baron/...
```

Both can be done from the Makefile in the root directory:

```
$ cd $GOPATH/github.com/teamlint/baron
$ make
```

## Testing

Before submitting a pull request always run tests that cover modified code.
Also build baron and run baron's integration test. This can be done by

```
$ cd $GOPATH/src/github.com/teamlint/baron
$ make
$ make test
# If the tests failed and you want to remove generated code
$ make testclean
```

## Structure

baron works as follows:

1. Read in a group of `.proto` files
2. Execute `protoc` with our `protoc-gen-protocast` protoc plugin, which
   outputs the protoc AST representation of the .proto files
3. Parse protoc's AST output and  the `.proto` file with the
   `grpc Service` definition for http annotations using `go-baron/deftree`
4. Use `protoc` and `protoc-gen-go` to generate `.pb.go` files containing
   protobuf structs and transport for golang
5. Use the constructed `deftree` with `gengokit` to template out basic gokit service with grpc
   and http/json transport and empty handlers
6. Generate documentation from comments with `gendocs`

If there was already generated code in the filesystem then baron will not
overwrite user code in the /NAME-service/handlers directory

Additional internal packages of note used by these programs are:

- `deftree`, located in `deftree/`, which makes sense of the protobuf file
  passed to it by `protoc`, and is used by `gengokit` and
  `gendoc`
