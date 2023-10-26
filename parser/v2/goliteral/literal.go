package goliteral

import "github.com/a-h/parse"

// Letters and digits

var octal_digit = parse.RuneIn("01234567")
var hex_digit = parse.RuneIn("0123456789ABCDEFabcdef")

// https://go.dev/ref/spec#Rune_literals

var Rune = parse.StringFrom(
	parse.Rune('\''),
	parse.StringFrom(parse.Until(
		parse.Any(unicode_value_rune, byte_value),
		parse.Rune('\''),
	)),
	parse.Rune('\''),
)
var unicode_value_rune = parse.Any(little_u_value, big_u_value, escaped_char, parse.RuneNotIn("'"))

// byte_value       = octal_byte_value | hex_byte_value .
var byte_value = parse.Any(octal_byte_value, hex_byte_value)

// octal_byte_value = `\` octal_digit octal_digit octal_digit .
var octal_byte_value = parse.StringFrom(
	parse.String(`\`),
	octal_digit, octal_digit, octal_digit,
)

// hex_byte_value   = `\` "x" hex_digit hex_digit .
var hex_byte_value = parse.StringFrom(
	parse.String(`\x`),
	hex_digit, hex_digit,
)

// little_u_value   = `\` "u" hex_digit hex_digit hex_digit hex_digit .
var little_u_value = parse.StringFrom(
	parse.String(`\u`),
	hex_digit, hex_digit,
	hex_digit, hex_digit,
)

// big_u_value      = `\` "U" hex_digit hex_digit hex_digit hex_digit
var big_u_value = parse.StringFrom(
	parse.String(`\U`),
	hex_digit, hex_digit, hex_digit, hex_digit,
	hex_digit, hex_digit, hex_digit, hex_digit,
)

// escaped_char     = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .
var escaped_char = parse.StringFrom(
	parse.Rune('\\'),
	parse.Any(
		parse.Rune('a'),
		parse.Rune('b'),
		parse.Rune('f'),
		parse.Rune('n'),
		parse.Rune('r'),
		parse.Rune('t'),
		parse.Rune('v'),
		parse.Rune('\\'),
		parse.Rune('\''),
		parse.Rune('"'),
	),
)

// https://go.dev/ref/spec#String_literals

var String = parse.Any(interpreted_string_lit, raw_string_lit)

var interpreted_string_lit = parse.StringFrom(
	parse.Rune('"'),
	parse.StringFrom(parse.Until(
		parse.Any(unicode_value_interpreted, byte_value),
		parse.Rune('"'),
	)),
	parse.Rune('"'),
)
var unicode_value_interpreted = parse.Any(little_u_value, big_u_value, escaped_char, parse.RuneNotIn("\n\""))

var raw_string_lit = parse.StringFrom(
	parse.Rune('`'),
	parse.StringFrom(parse.Until(
		unicode_value_raw,
		parse.Rune('`'),
	)),
	parse.Rune('`'),
)
var unicode_value_raw = parse.Any(little_u_value, big_u_value, escaped_char, parse.RuneNotIn("`"))
