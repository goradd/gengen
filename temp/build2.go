package temp

// This go file uses the .got files in this directory to build a variety of versions of the maps
// in this directory. You can use this as an example of how to use GoT to build your own custom
// versions of the maps here.

//go:generate gengen -c string_string.json < mapi.tpl > ../map/strstrmapi.go
//go:generate gengen -c string_string.json < stdmap.tpl > ../map/strstrmap.go
//go:generate gengen -c string_string.json < safemap.tpl > ../map/strstrsafemap.go
//go:generate gengen -c string_string.json < syncmap.tpl > ../map/strstrsyncmap.go

//go:generate gengen -c string_string.json < slicemapi.tpl > ../map/strstrslicemapi.go
//go:generate gengen -c string_string.json < slicemap.tpl > ../map/strstrslicemap.go
//go:generate gengen -c string_string.json < safeslicemap.tpl > ../map/strstrsafeslicemap.go
//go:generate gengen -c string_string.json < syncslicemap.tpl > ../map/strstrsyncslicemap.go
