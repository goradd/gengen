package templates

// This go file uses the .got files in this directory to build a variety of versions of the maps
// in this directory. You can use this as an example of how to use GoT to build your own custom
// versions of the maps here.

//go:generate gengen -c string_string.json  -o ../../pkg/maps/strmapi.go mapi.tmpl

//go:generate gengen -c string_string.json  -o ../../pkg/maps/strmap.go standard_map.tmpl
//go:generate gengen -c safe_string_string.json  -o ../../pkg/maps/safestrmap.go standard_map.tmpl
//go:generate gengen -c string_string.json  -o ../../pkg/maps/strslicemap.go slice_map.tmpl
//go:generate gengen -c safe_string_string.json  -o ../../pkg/maps/safestrslicemap.go slice_map.tmpl

//go:generate gengen -c standard_test.json  -o ../../pkg/maps/strmap_test.go string_string_test.tmpl
//go:generate gengen -c safe_test.json  -o ../../pkg/maps/safestrmap_test.go string_string_test.tmpl
//go:generate gengen -c standard_test.json  -o ../../pkg/maps/strslicemap_test.go string_string_slice_test.tmpl
//go:generate gengen -c safe_test.json  -o ../../pkg/maps/safestrslicemap_test.go string_string_slice_test.tmpl

//go:generate gengen -c string_interface.json  -o ../../pkg/maps/mapi.go mapi.tmpl

//go:generate gengen -c string_interface.json  -o ../../pkg/maps/map.go standard_map.tmpl
//go:generate gengen -c safe_string_interface.json  -o ../../pkg/maps/safemap.go standard_map.tmpl
//go:generate gengen -c string_interface.json  -o ../../pkg/maps/slicemap.go slice_map.tmpl
//go:generate gengen -c safe_string_interface.json  -o ../../pkg/maps/safeslicemap.go slice_map.tmpl

//go:generate gengen -c standard_test.json  -o ../../pkg/maps/map_test.go string_interface_test.tmpl
//go:generate gengen -c safe_test.json  -o ../../pkg/maps/safemap_test.go string_interface_test.tmpl
//go:generate gengen -c standard_test.json  -o ../../pkg/maps/slicemap_test.go string_interface_slice_test.tmpl
//go:generate gengen -c safe_test.json  -o ../../pkg/maps/safeslicemap_test.go string_interface_slice_test.tmpl
