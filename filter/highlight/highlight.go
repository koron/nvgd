package highlight

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/resource"
)

//go:embed highlight.html
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "filter/highlight", "")

type doc struct {
	lexer     chroma.Lexer
	style     *chroma.Style
	tokens    chroma.Iterator
	formatter *html.Formatter
}

func (d *doc) CSS() (template.CSS, error) {
	bb := &bytes.Buffer{}
	err := d.formatter.WriteCSS(bb, d.style)
	if err != nil {
		return "", err
	}
	return template.CSS(bb.String()), nil
}

func (d *doc) Body() (template.HTML, error) {
	bb := &bytes.Buffer{}
	err := d.formatter.Format(bb, d.style, d.tokens)
	if err != nil {
		return "", err
	}
	return template.HTML(bb.String()), nil
}

func (d *doc) Lexer() string {
	return d.lexer.Config().Name
}

func (d *doc) Style() string {
	return d.style.Name
}

func highlight(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	defer r.Close()
	paramLexer := p.String("lexer", "")
	paramStyle := p.String("style", "github")

	// Load all
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Determine the lexer
	var lexer chroma.Lexer
	if paramLexer != "" {
		lexer = lexers.Get(paramLexer)
	} else {
		lexer = lexers.Analyse(string(b))
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Determine the style
	style := styles.Get(paramStyle)
	if style == nil {
		style = styles.Fallback
	}

	tokens, err := lexer.Tokenise(nil, string(b))
	if err != nil {
		return nil, err
	}

	formatter := html.New(
		html.WithCustomCSS(nil),
		html.WithClasses(true),
		html.WithLineNumbers(true),
		html.WithLinkableLineNumbers(true, "L"),
		html.LineNumbersInTable(false),
	)

	d := doc{
		lexer:     lexer,
		style:     style,
		tokens:    tokens,
		formatter: formatter,
	}

	// Write HTML
	tmpl, err := template.ParseFS(assetFS, "highlight.html")
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, &d); err != nil {
		return nil, err
	}

	return r.Wrap(io.NopCloser(buf)).
		PutContentType("text/html").
		Put(resource.SkipFilters, true), nil
}

func init() {
	filter.MustRegister("highlight", highlight)
}
