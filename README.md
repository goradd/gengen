# Gengen

Gengen is a command line generics generator for the GO language suitable for use
in go:generate lines. It uses standard go text templates to
generate GO files that use your custom types to generate a variety of collection
objects, and includes some common objects as well.

Gengen is pretty simple, and just combines a template with a json configuration file you can use to
define values in the template. It expects the template to be the stdin or the only argument
to the command line. Output will be directed to the file you specify with the -o option,
or stdout if no file is specified. 
This design means gengen does not have additional dependencies, is easily installed 
using *go get*, is cross-platform, and works well in go:generate lines.

When specifying the config file or input file, you can use the word "GOPATH" to refer
to the current go path in the path to either file. For example, ```GOPATH/src/myfile.tmpl```.

Configuration files can start with comments. Gengen will look for the first open bracket,
and start reading its json input from there.

## Installation

To just use the pre-built library of collections, execute the following `go get` command and import the library into
your project.

```shell
go get github.com/goradd/gengen
```

If you want to build your own files, either using the pre-made generic templates, or by creating your own, you will
need to install the `gengen` executable by doing the following instead:

```shell
go get -u github.com/goradd/gengen/...
```

## Usage

To use the command line tool to build a generic template into one specific to  your types, do the following in the shell:

```shell
gengen -c <config_file>  [-o out_file] [template_file]
```

`gengen` requires the -c command to specify a json configuration file that sets up the "dot" context of the template.
If you do not specify an out_file, output will be sent the StdOut. If you do not specify a template_file, the template
will be read from StdIn.

File paths are module and package aware. In other words, if you do this:

```shell
gengen -c github.com/goradd/gengen/templates/map_src/safe_test.json
```

`gengen` will see the `github.com/goradd/gengen` as a module or package path, and substitute the real path. This works
whether or not you are using modules (a new feature in GO 1.11).

Environment variables can be inserted into the path using this syntax: `$var` or `${var}`. This works on all platforms.

## Examples

See the `templates/build.go` file for an example of how the included library is built.

## Library

Gengen includes a maintained library of templates that it uses to create some standard
configurations of useful collections, and that you can use to create your own versions
of those using your own types. The library includes its own generated unit test code.

## License

Gengen is licensed under the MIT License.
