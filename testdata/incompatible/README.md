# incompatible bop schemas

Files in this directory should reasonably be supported by rainway bebop but are not, and we support, or vice versa.

## Sample causes

"uint8" is a supported type, but "int8" is not.

Rainway allows overwriting primitive types with message, struct, or enum definitions. We do not.

We do not allow recursive struct definitions, as they will never terminate when encoded or parsed. Rainway does not have a check for this.

See `compatibility_test.go` for which files specifically fail for which compilers.

## 2 to 3 compatibility

bebopc-go supports both bebop 2 and bebop 3. The changes from 2 to 3 are not backwards compatible, but the grammar changes can be supported at the same time without adding ambiguity. 
TODO: a flag can be provided to require bebop 2 or bebop 3 grammars are required. 