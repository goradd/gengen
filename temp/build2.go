package map_src

// This go file uses the .got files in this directory to build a variety of versions of the maps
// in this directory. You can use this as an example of how to use GoT to build your own custom
// versions of the maps here.

//go:generate gengen -c string_string.json  -o ../maps/strmapi2.go mapi2.tmpl
//go:generate gengen -c string_interface.json  -o ../maps/mapi2.go mapi2.tmpl

//go:generate gengen -c string_string.json  -o ../maps/strmapi.go mapi.tmpl
//go:generate gengen -c string_string.json  -o ../maps/strmap.go standard_map.tmpl
//go:generate gengen -c safe_string_string.json  -o ../maps/safestrmap.go standard_map.tmpl
//go:generate gengen -c string_string.json  -o ../maps/strslicemap.go slice_map.tmpl
//go:generate gengen -c string_string.json  -o ../maps/safestrslicemap.go safe_slice_map.tmpl

//go:generate gengen -c standard_test.json  -o ../maps/strmap_test.go string_string_test.tmpl
//go:generate gengen -c safe_test.json  -o ../maps/safestrmap_test.go string_string_test.tmpl
//go:generate gengen -c standard_test.json  -o ../maps/strslicemap_test.go string_string_slice_test.tmpl
//go:generate gengen -c safe_test.json  -o ../maps/safestrslicemap_test.go string_string_slice_test.tmpl

//go:generate gengen -c string_interface.json  -o ../maps/mapi.go mapi.tmpl
//go:generate gengen -c string_interface.json  -o ../maps/map.go standard_map.tmpl
//go:generate gengen -c safe_string_string.json  -o ../maps/safemap.go standard_map.tmpl
//go:generate gengen -c string_interface.json  -o ../maps/slicemap.go slice_map.tmpl
//go:generate gengen -c string_interface.json  -o ../maps/safeslicemap.go safe_slice_map.tmpl

//go:generate gengen -c standard_test.json  -o ../maps/map_test.go string_interface_test.tmpl
//go:generate gengen -c safe_test.json  -o ../maps/safemap_test.go string_interface_test.tmpl
