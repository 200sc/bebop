package protos

//go:generate protoc -I=. --go_opt=paths=source_relative --go_out=. ./basic_types.proto ./basic_arrays.proto ./array_of_strings.proto ./jazz.proto ./map_types.proto
