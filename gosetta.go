package gosetta

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// Rose holds translation information.
type Rose struct {
	apiKey string
	ctx    context.Context
	cl     *translate.Client
	opts   *translate.Options
}

// New creates a rose for translating.
func New(src language.Tag) (*Rose, error) {
	key := os.Getenv("TRANSLATE_API_KEY")
	if len(key) < 1 {
		return nil, errors.New("missing TRANSLATE_API_KEY")
	}
	ctx := context.Background()
	cl, err := translate.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	opts := &translate.Options{Source: src}
	return &Rose{key, ctx, cl, opts}, nil
}

// Source sets the new source language for the rose.
func (r *Rose) Source(src language.Tag) {
	r.opts.Source = src
}

// Translate moves from input from a source to destination language.
func (r Rose) Translate(in []string, dst language.Tag) ([]string, error) {
	// defer wg.Done()
	trans, err := r.cl.Translate(r.ctx, in, dst, r.opts)
	if err != nil {
		return nil, err
	}
	lang := []string{}
	for _, tran := range trans {
		lang = append(lang, tran.Text)
	}
	return lang, nil

}
