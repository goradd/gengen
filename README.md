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

## Library

Gengen includes a maintained library of templates that it uses to create some standard
configurations of useful collections, and that you can use to create your own versions
of those using your own types. 

## License

Gengen is licensed under the MIT License.
