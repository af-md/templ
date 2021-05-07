// Code generated by templ DO NOT EDIT.

package testswitch

import "html"
import "context"
import "io"

func render(ctx context.Context, w io.Writer, input string) (err error) {
	switch input {
	case "a":
		_, err = io.WriteString(w, html.EscapeString("it was 'a'"))
		if err != nil {
			return err
		}
	default:		_, err = io.WriteString(w, html.EscapeString("it was something else"))
		if err != nil {
			return err
		}
	}
	return err
}
