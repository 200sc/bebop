package bebop

import (
	"io"
)

// Format reads a .bop file from r and writes out a formatted version of that file to out.
func Format(r io.Reader, out io.Writer) {
	format(newTokenReader(r), out)
}

func format(tr *tokenReader, w io.Writer) {
	var readOnly bool
	var newlineBeforeNextRecord bool
	for tr.Next() {
		t := tr.Token()
		switch t.kind {
		case tokenKindOpenSquare:
			if newlineBeforeNextRecord {
				w.Write([]byte{'\n'})
			}
			// opcode, next tokens are 'opcode', '(', hex or string lit, ')', ']'
			opCodeBytes := t.concrete
			for j := 0; j < 5; j++ {
				tr.Next()
				opCodeBytes = append(opCodeBytes, tr.Token().concrete...)
			}
			// inject newline after opcodes
			opCodeBytes = append(opCodeBytes, '\n')
			w.Write(opCodeBytes)
			newlineBeforeNextRecord = false
		case tokenKindLineComment:
			cmtBytes := t.concrete
			w.Write(cmtBytes)
			newlineBeforeNextRecord = false
		case tokenKindBlockComment:
			cmtBytes := t.concrete
			cmtBytes = append(cmtBytes, []byte("\n")...)
			w.Write(cmtBytes)
			newlineBeforeNextRecord = false
		case tokenKindReadOnly:
			readOnly = true
			continue
		case tokenKindEnum:
			if newlineBeforeNextRecord {
				w.Write([]byte{'\n'})
			}
			w.Write(formatEnum(tr))
			newlineBeforeNextRecord = true
		case tokenKindStruct:
			if newlineBeforeNextRecord {
				w.Write([]byte{'\n'})
			}
			w.Write(formatStruct(tr, readOnly, "\t"))
			newlineBeforeNextRecord = true
		case tokenKindMessage:
			if newlineBeforeNextRecord {
				w.Write([]byte{'\n'})
			}
			w.Write(formatMessage(tr, "\t"))
			newlineBeforeNextRecord = true
		case tokenKindUnion:
			if newlineBeforeNextRecord {
				w.Write([]byte{'\n'})
			}
			w.Write(formatUnion(tr, "\t"))
			newlineBeforeNextRecord = true
		}
		readOnly = false
	}
}

func formatEnum(tr *tokenReader) []byte {
	// enum <ID> {\n
	enumBytes := tr.Token().concrete
	for j := 0; j < 2; j++ {
		enumBytes = append(enumBytes, ' ')
		tr.Next()
		enumBytes = append(enumBytes, tr.Token().concrete...)
	}
	enumBytes = append(enumBytes, '\n')

tokenLoop:
	for tr.Next() {
		t := tr.Token()
		switch t.kind {
		case tokenKindLineComment:
			cmtBytes := append([]byte("\t"), t.concrete...)
			enumBytes = append(enumBytes, cmtBytes...)
		case tokenKindBlockComment:
			cmtBytes := append([]byte("\t"), t.concrete...)
			enumBytes = append(enumBytes, cmtBytes...)
			enumBytes = append(enumBytes, []byte("\n")...)
		case tokenKindOpenSquare:
			// deprecated, next tokens are 'deprecated', '(', string lit, ')', ']'
			deprecatedBytes := append([]byte{'\t'}, t.concrete...)
			for j := 0; j < 5; j++ {
				tr.Next()
				deprecatedBytes = append(deprecatedBytes, tr.Token().concrete...)
			}
			deprecatedBytes = append(deprecatedBytes, '\n')
			enumBytes = append(enumBytes, deprecatedBytes...)
		case tokenKindIdent:
			// <ID> = <NUM>;
			optBytes := append([]byte{'\t'}, t.concrete...)
			for j := 0; j < 2; j++ {
				optBytes = append(optBytes, ' ')
				tr.Next()
				optBytes = append(optBytes, tr.Token().concrete...)
			}
			tr.Next()
			optBytes = append(optBytes, []byte(";\n")...)
			enumBytes = append(enumBytes, optBytes...)
		case tokenKindCloseCurly:
			enumBytes = append(enumBytes, t.concrete...)
			enumBytes = append(enumBytes, '\n')
			break tokenLoop
		}
	}

	return enumBytes
}

func formatStruct(tr *tokenReader, readonly bool, prefix string) []byte {
	// [readonly] struct <ID> {\n

	structBytes := tr.Token().concrete
	if readonly {
		structBytes = append([]byte("readonly "), structBytes...)
	}
	for j := 0; j < 2; j++ {
		structBytes = append(structBytes, ' ')
		tr.Next()
		structBytes = append(structBytes, tr.Token().concrete...)
	}
	structBytes = append(structBytes, '\n')

tokenLoop:
	for tr.Next() {
		t := tr.Token()
		switch t.kind {
		case tokenKindLineComment:
			cmt := append([]byte(prefix), t.concrete...)
			structBytes = append(structBytes, cmt...)
		case tokenKindBlockComment:
			cmt := append([]byte(prefix), t.concrete...)
			cmt = append(cmt, []byte("\n")...)
			structBytes = append(structBytes, cmt...)

		case tokenKindOpenSquare:
			// deprecated, next tokens are 'deprecated', '(', string lit, ')', ']'
			deprecatedBytes := append([]byte(prefix), t.concrete...)
			for j := 0; j < 5; j++ {
				tr.Next()
				deprecatedBytes = append(deprecatedBytes, tr.Token().concrete...)
			}
			deprecatedBytes = append(deprecatedBytes, '\n')
			structBytes = append(structBytes, deprecatedBytes...)
		case tokenKindIdent, tokenKindMap, tokenKindArray:
			// <TYPE> <ID>;
			fieldBytes := formatType(tr)
			fieldBytes = append([]byte(prefix), fieldBytes...)
			tr.Next()
			fieldBytes = append(fieldBytes, ' ')
			fieldBytes = append(fieldBytes, tr.Token().concrete...)
			tr.Next()
			fieldBytes = append(fieldBytes, []byte(";")...)
			structBytes = append(structBytes, fieldBytes...)
			tr.Next()
			t := tr.Token()
			if t.kind == tokenKindLineComment {
				structBytes = append(structBytes, ' ')
				structBytes = append(structBytes, t.concrete...)
			} else {
				tr.UnNext()
				structBytes = append(structBytes, '\n')
			}
		case tokenKindCloseCurly:
			structBytes = append(structBytes, append([]byte(prefix)[:len(prefix)-1], t.concrete...)...)
			structBytes = append(structBytes, '\n')
			break tokenLoop
		}
	}

	return structBytes
}

func formatMessage(tr *tokenReader, prefix string) []byte {
	msgBytes := tr.Token().concrete
	for j := 0; j < 2; j++ {
		msgBytes = append(msgBytes, ' ')
		tr.Next()
		msgBytes = append(msgBytes, tr.Token().concrete...)
	}
	msgBytes = append(msgBytes, '\n')

tokenLoop:
	for tr.Next() {
		t := tr.Token()
		switch t.kind {
		case tokenKindLineComment:
			msgBytes := append([]byte(prefix), t.concrete...)
			msgBytes = append(msgBytes, '\n')
		case tokenKindBlockComment:
			msgBytes := append([]byte(prefix), t.concrete...)
			msgBytes = append(msgBytes, []byte("\n")...)
		case tokenKindOpenSquare:
			// deprecated, next tokens are 'deprecated', '(', string lit, ')', ']'
			deprecatedBytes := append([]byte(prefix), t.concrete...)
			for j := 0; j < 5; j++ {
				tr.Next()
				deprecatedBytes = append(deprecatedBytes, tr.Token().concrete...)
			}
			deprecatedBytes = append(deprecatedBytes, '\n')
			msgBytes = append(msgBytes, deprecatedBytes...)
		case tokenKindInteger:
			// <NUM> -> <TYPE> <ID>;
			fieldBytes := append([]byte(prefix), t.concrete...)
			fieldBytes = append(fieldBytes, ' ')
			tr.Next()
			fieldBytes = append(fieldBytes, tr.Token().concrete...)
			tr.Next()
			fieldBytes = append(fieldBytes, ' ')
			typeBytes := formatType(tr)
			fieldBytes = append(fieldBytes, typeBytes...)
			tr.Next()
			fieldBytes = append(fieldBytes, ' ')
			fieldBytes = append(fieldBytes, tr.Token().concrete...)
			tr.Next()
			fieldBytes = append(fieldBytes, []byte(";\n")...)
			msgBytes = append(msgBytes, fieldBytes...)
		case tokenKindCloseCurly:
			msgBytes = append(msgBytes, append([]byte(prefix)[:len(prefix)-1], t.concrete...)...)
			msgBytes = append(msgBytes, '\n')
			break tokenLoop
		}
	}

	return msgBytes
}

func formatUnion(tr *tokenReader, prefix string) []byte {
	unionBytes := tr.Token().concrete
	for j := 0; j < 2; j++ {
		unionBytes = append(unionBytes, ' ')
		tr.Next()
		unionBytes = append(unionBytes, tr.Token().concrete...)
	}
	unionBytes = append(unionBytes, '\n')

tokenLoop:
	for tr.Next() {
		t := tr.Token()
		switch t.kind {
		case tokenKindLineComment:
			unionBytes := append([]byte(prefix), t.concrete...)
			unionBytes = append(unionBytes, '\n')
		case tokenKindBlockComment:
			unionBytes := append([]byte(prefix), t.concrete...)
			unionBytes = append(unionBytes, []byte("\n")...)
		case tokenKindOpenSquare:
			// deprecated, next tokens are 'deprecated', '(', string lit, ')', ']'
			deprecatedBytes := append([]byte(prefix), t.concrete...)
			for j := 0; j < 5; j++ {
				tr.Next()
				deprecatedBytes = append(deprecatedBytes, tr.Token().concrete...)
			}
			deprecatedBytes = append(deprecatedBytes, '\n')
			unionBytes = append(unionBytes, deprecatedBytes...)
		case tokenKindInteger:
			// <NUM> -> ???;
			fieldBytes := append([]byte(prefix), t.concrete...)
			fieldBytes = append(fieldBytes, ' ')
			tr.Next()
			fieldBytes = append(fieldBytes, tr.Token().concrete...)
			tr.Next()
			fieldBytes = append(fieldBytes, ' ')
			tk := tr.Token()
			switch tk.kind {
			case tokenKindUnion:
				fieldBytes = append(fieldBytes, formatUnion(tr, prefix+"\t")...)
			case tokenKindMessage:
				fieldBytes = append(fieldBytes, formatMessage(tr, prefix+"\t")...)
			case tokenKindStruct:
				fieldBytes = append(fieldBytes, formatStruct(tr, false, prefix+"\t")...)
			}
			unionBytes = append(unionBytes, fieldBytes...)
		case tokenKindCloseCurly:
			unionBytes = append(unionBytes, append([]byte(prefix)[:len(prefix)-1], t.concrete...)...)
			unionBytes = append(unionBytes, '\n')
			break tokenLoop
		}
	}

	return unionBytes
}

func formatType(tr *tokenReader) []byte {
	typeBytes := []byte{}
	t := tr.Token()
	switch t.kind {
	case tokenKindIdent:
		// simple!
		typeBytes = append(typeBytes, t.concrete...)
	case tokenKindMap:
		// map[<ID>, <TYPE>]
		typeBytes = append(typeBytes, t.concrete...)
		for j := 0; j < 3; j++ {
			tr.Next()
			typeBytes = append(typeBytes, tr.Token().concrete...)
		}
		typeBytes = append(typeBytes, ' ')
		tr.Next()
		valBytes := formatType(tr)
		typeBytes = append(typeBytes, valBytes...)
		tr.Next()
		typeBytes = append(typeBytes, tr.Token().concrete...)
	case tokenKindArray:
		// array[<TYPE>]
		typeBytes = append(typeBytes, t.concrete...)
		tr.Next()
		typeBytes = append(typeBytes, tr.Token().concrete...)
		tr.Next()
		valBytes := formatType(tr)
		typeBytes = append(typeBytes, valBytes...)
		tr.Next()
		typeBytes = append(typeBytes, tr.Token().concrete...)
	}

	// ...[]?
	tr.Next()
	if tr.Token().kind == tokenKindOpenSquare {
		tr.Next()
		typeBytes = append(typeBytes, []byte("[]")...)
	} else {
		tr.UnNext()
	}

	return typeBytes
}
