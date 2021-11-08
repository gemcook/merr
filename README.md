# merr

[![Analysis](https://github.com/gemcook/merr/actions/workflows/analysis.yml/badge.svg)](https://github.com/gemcook/merr/actions/workflows/analysis.yml) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gemcook/merr)

`merr` is a package that binds multiple errors.

## Features

`merr` has pretty print, which outputs the bound errors in a well-formatted and easy-to-understand format.
It is also possible to change the output destination if necessary, which is very useful for debugging.

## Installation

```sh
go get -u github.com/gemcook/merr
```

## Usage

```go
import "github.com/gemcook/merr"
```

### Appends error

Appends error to list of errors.
If there is no error list, it creates a new one.

```go
multiError := merr.New()

for i := 0; i < 10; i++ {
    if err := something(); err != nil {
        multiError.Append(err)
    }
}
```

### Print list of error

Prints the object of the list of errors, separated by `,\n`.

```go
multiError := merr.New()

for i := 0; i < 10; i++ {
    // something() returns &errors.errorString{s: "something error"}
    if err := something(); err != nil {
        multiError.Append(err)
    }
}

fmt.Println(multiError.Error())
```

```
"something error",
"something error",
  .
  .
  .
"something error"
```

### Pretty print

Prints a list of errors in a well-formatted, easy-to-understand format.

```go
multiError := merr.New()

for i := 0; i < 10; i++ {
    // something() returns &errors.errorString{s: "something error"}
    if err := something(); err != nil {
        multiError.Append(err)
    }
}

multiError.PrettyPrint()
```

```
Errors[
  &errors.errorString{
    s: "something error",
  },
  &errors.errorString{
    s: "something error",
  },
  &errors.errorString{
    s: "something error",
  },
  .
  .
  .
  &errors.errorString{
    s: "something error",
  },
]
```

The default output destination is the **standard output**.
You can also change the output destination.

```go
buf := bytes.NewBuffer(nil)
merr.SetOutput(buf)
```

## License

MIT License
