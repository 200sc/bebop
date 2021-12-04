# incompatible bop schemas

Files in this directory should reasonably be supported by rainway bebop but are not, and we support, or vice versa.

## Sample causes

"uint8" is a supported type, but "int8" is not.

Rainway allows overwriting primitive types with message, struct, or enum definitions. We do not.

We do not allow recursive struct definitions, as they will never terminate when encoded or parsed. Rainway does not have a check for this.

See `compatability_test.go` for which files specifically fail for which compilers.