# incompatible bop schemas

Files in this directory should reasonably be supported by rainway bebop but are not.

We support these files, however. Documentation is included to describe what each file demonstrates.

## causes

"uint8" is a supported type, but "int8" is not.

Numbers that exceed 64 bits are truncated by rainway (in theory, this needs to be confirmed), and rejected as 
out of range errors by us, using Go's standard library number parsing

Rainway allows overwriting primitive types with message, struct, or enum definitions. We do not.

We do not allow recursive struct definitions, as they will never terminate when encoded or parsed. Rainway does not have a check for this.
