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
	ctx  context.Context
	cl   *translate.Client
	opts *translate.Options
	dst  language.Tag
}

// New creates a rose for translating.
func New(src language.Tag) (*Rose, error) {
	key, found := os.LookupEnv("TRANSLATE_API_KEY")
	if !found {
		return nil, errors.New("TRANSLATE_API_KEY not found")
	}
	ctx := context.Background()
	cl, err := translate.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	opts := &translate.Options{Source: src}
	return &Rose{
		ctx:  ctx,
		cl:   cl,
		opts: opts,
	}, nil
}

// Source sets the new source language for the rose.
func (r *Rose) Source(src language.Tag) {
	r.opts.Source = src
}

// Destination sets the new destination language for the rose.
func (r *Rose) Destination(dst language.Tag) {
	r.dst = dst
}

// MustTranslate translates from source to destination or else it panics.
func (r Rose) MustTranslate(x string) string {
	y, err := r.Translate(x)
	if err != nil {
		panic(err)
	}
	return y
}

// Translate moves from input from a source to destination language.
func (r Rose) Translate(x string) (string, error) {
	trans, err := r.cl.Translate(r.ctx, []string{x}, r.dst, r.opts)
	if err != nil {
		return "", err
	}
	if len(trans) != 1 {
		return "", errors.New("not just one translation")
	}
	return trans[0].Text, nil
}
