package parser

import (
	"strings"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/goliteral"
)

// StripType takes the parser and throws away the return value.
func StripType[T any](p parse.Parser[T]) parse.Parser[interface{}] {
	return parse.Func(func(in *parse.Input) (out interface{}, ok bool, err error) {
		return p.Parse(in)
	})
}

func Must[T any](p parse.Parser[T], msg string) parse.Parser[T] {
	return parse.Func(func(in *parse.Input) (out T, ok bool, err error) {
		out, ok, err = p.Parse(in)
		if err != nil {
			return
		}
		if !ok {
			err = parse.Error(msg, in.Position())
		}
		return out, ok, err
	})
}

func ExpressionOf(p parse.Parser[string]) parse.Parser[Expression] {
	return parse.Func(func(in *parse.Input) (out Expression, ok bool, err error) {
		from := in.Position()

		var exp string
		if exp, ok, err = p.Parse(in); err != nil || !ok {
			return
		}

		return NewExpression(exp, from, in.Position()), true, nil
	})
}

var lt = parse.Rune('<')
var gt = parse.Rune('>')
var openBrace = parse.String("{")
var optionalSpaces = parse.StringFrom(parse.Optional(
	parse.AtLeast(1, parse.Rune(' '))))
var openBraceWithPadding = parse.StringFrom(optionalSpaces,
	openBrace,
	optionalSpaces)
var openBraceWithOptionalPadding = parse.Any(openBraceWithPadding, openBrace)

var closeBrace = parse.String("}")
var closeBraceWithPadding = parse.String(" }")
var closeBraceWithOptionalPadding = parse.Any(closeBraceWithPadding, closeBrace)

var openBracket = parse.String("(")
var closeBracket = parse.String(")")
var closeBracketWithOptionalPadding = parse.StringFrom(optionalSpaces, closeBracket)

var exp = expressionParser{
	startBraceCount: 1,
}

type expressionParser struct {
	startBraceCount int
}

func (p expressionParser) Parse(pi *parse.Input) (s Expression, ok bool, err error) {
	from := pi.Position()

	braceCount := p.startBraceCount

	var sb strings.Builder
loop:
	for {
		var result string

		// Try to read a string literal first.
		if result, ok, err = goliteral.String.Parse(pi); err != nil {
			return
		}
		if ok {
			sb.WriteString(result)
			continue
		}
		// Also try for a rune literal.
		if result, ok, err = goliteral.Rune.Parse(pi); err != nil {
			return
		}
		if ok {
			sb.WriteString(result)
			continue
		}
		// Try opener.
		if result, ok, err = openBrace.Parse(pi); err != nil {
			return
		}
		if ok {
			braceCount++
			sb.WriteString(result)
			continue
		}
		// Try closer.
		startOfCloseBrace := pi.Index()
		if result, ok, err = closeBraceWithOptionalPadding.Parse(pi); err != nil {
			return
		}
		if ok {
			braceCount--
			if braceCount < 0 {
				err = parse.Error("expression: too many closing braces", pi.Position())
				return
			}
			if braceCount == 0 {
				pi.Seek(startOfCloseBrace)
				break loop
			}
			sb.WriteString(result)
			continue
		}

		// Read anything else.
		var c string
		c, ok = pi.Take(1)
		if !ok {
			break loop
		}
		if rune(c[0]) == 65533 { // Invalid Unicode.
			break loop
		}
		sb.WriteString(c)
	}
	if braceCount != 0 {
		err = parse.Error("expression: unexpected brace count", pi.Position())
		return
	}

	return NewExpression(sb.String(), from, pi.Position()), true, nil
}

type functionArgsParser struct {
	startBracketCount int
}

func (p functionArgsParser) Parse(pi *parse.Input) (s Expression, ok bool, err error) {
	from := pi.Position()

	bracketCount := p.startBracketCount

	var sb strings.Builder
loop:
	for {
		var result string

		// Try to read a string literal first.
		if result, ok, err = goliteral.String.Parse(pi); err != nil {
			return
		}
		if ok {
			sb.WriteString(result)
			continue
		}
		// Also try for a rune literal.
		if result, ok, err = goliteral.Rune.Parse(pi); err != nil {
			return
		}
		if ok {
			sb.WriteString(result)
			continue
		}
		// Try opener.
		if result, ok, err = openBracket.Parse(pi); err != nil {
			return
		}
		if ok {
			bracketCount++
			sb.WriteString(result)
			continue
		}
		// Try closer.
		startOfCloseBracket := pi.Index()
		if result, ok, err = closeBracketWithOptionalPadding.Parse(pi); err != nil {
			return
		}
		if ok {
			bracketCount--
			if bracketCount < 0 {
				err = parse.Error("expression: too many closing brackets", pi.Position())
				return
			}
			if bracketCount == 0 {
				pi.Seek(startOfCloseBracket)
				break loop
			}
			sb.WriteString(result)
			continue
		}

		// Read anything else.
		var c string
		c, ok = pi.Take(1)
		if !ok {
			break loop
		}
		if rune(c[0]) == 65533 { // Invalid Unicode.
			break loop
		}
		sb.WriteString(c)
	}
	if bracketCount != 0 {
		err = parse.Error("expression: unexpected bracket count", pi.Position())
		return
	}

	return NewExpression(sb.String(), from, pi.Position()), true, nil
}
