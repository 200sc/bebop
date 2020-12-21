package bebop

import (
	"io"
)

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
			w.Write(formatStruct(tr, readOnly))
			newlineBeforeNextRecord = true
		case tokenKindMessage:
			if newlineBeforeNextRecord {
				w.Write([]byte{'\n'})
			}
			w.Write(formatMessage(tr, readOnly))
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

func formatStruct(tr *tokenReader, readonly bool) []byte {
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
			cmt := append([]byte("\t"), t.concrete...)
			structBytes = append(structBytes, cmt...)
		case tokenKindBlockComment:
			cmt := append([]byte("\t"), t.concrete...)
			cmt = append(cmt, []byte("\n")...)
			structBytes = append(structBytes, cmt...)

		case tokenKindOpenSquare:
			// deprecated, next tokens are 'deprecated', '(', string lit, ')', ']'
			deprecatedBytes := append([]byte{'\t'}, t.concrete...)
			for j := 0; j < 5; j++ {
				tr.Next()
				deprecatedBytes = append(deprecatedBytes, tr.Token().concrete...)
			}
			deprecatedBytes = append(deprecatedBytes, '\n')
			structBytes = append(structBytes, deprecatedBytes...)
		case tokenKindIdent, tokenKindMap, tokenKindArray:
			// <TYPE> <ID>;
			fieldBytes := formatType(tr)
			fieldBytes = append([]byte{'\t'}, fieldBytes...)
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
			structBytes = append(structBytes, t.concrete...)
			structBytes = append(structBytes, '\n')
			break tokenLoop
		}
	}

	return structBytes
}

func formatMessage(tr *tokenReader, readonly bool) []byte {
	// [readonly] message <ID> {\n

	msgBytes := tr.Token().concrete
	if readonly {
		msgBytes = append([]byte("readonly "), msgBytes...)
	}
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
			msgBytes := append([]byte("\t//"), t.concrete...)
			msgBytes = append(msgBytes, '\n')
		case tokenKindBlockComment:
			msgBytes := append([]byte("\t"), t.concrete...)
			msgBytes = append(msgBytes, []byte("\n")...)
		case tokenKindOpenSquare:
			// deprecated, next tokens are 'deprecated', '(', string lit, ')', ']'
			deprecatedBytes := append([]byte{'\t'}, t.concrete...)
			for j := 0; j < 5; j++ {
				tr.Next()
				deprecatedBytes = append(deprecatedBytes, tr.Token().concrete...)
			}
			deprecatedBytes = append(deprecatedBytes, '\n')
			msgBytes = append(msgBytes, deprecatedBytes...)
		case tokenKindInteger:
			// <NUM> -> <TYPE> <ID>;
			fieldBytes := append([]byte{'\t'}, t.concrete...)
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
			msgBytes = append(msgBytes, t.concrete...)
			msgBytes = append(msgBytes, '\n')
			break tokenLoop
		}
	}

	return msgBytes
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
