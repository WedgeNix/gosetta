package gosetta

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

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
	// call        chan bool
	texts2trans chan []text2tran
	bsize       chan int
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

	// call := make(chan bool, 100)
	// for range make([]interface{}, cap(call)) {
	// 	call <- true
	// }

	t2ts := make(chan []text2tran, 1)
	t2ts <- []text2tran{}

	bsize := make(chan int, 1)
	bsize <- 0

	return &Rose{
		ctx:  ctx,
		cl:   cl,
		opts: opts,
		// call:        call,
		texts2trans: t2ts,
		bsize:       bsize,
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
func (r *Rose) MustTranslate(x string) string {
	return <-r.Translate(x)
}

type text2tran struct {
	text string
	out  chan string
}

const (
	charLim = 5000
	textLim = 128
)

// Translate sends buffer loads every five seconds. As long as the buffer isn't empty,
// keep sending. If empty, another call to translate would reignite
// the buffer and it'd start sending again.
func (r *Rose) Translate(x string) <-chan string {

	t2t := text2tran{x, make(chan string, 1)}
	t2ts := <-r.texts2trans

	if r.opts.Source.String() == r.dst.String() {
		t2t.out <- x
		return t2t.out
	}

	if len(t2ts) < 1 {
		r.texts2trans <- []text2tran{t2t}
		tick := time.NewTicker(charLim * time.Millisecond)

		go func() {
			for {
				<-tick.C

				t2ts := <-r.texts2trans

				l := len(t2ts)
				if l < 1 {
					break
				}

				var piv int
				size := 0
				texts := []string{}
				for i, t := range t2ts {
					piv = i

					tlen := len(t.text)
					if size+tlen > charLim || i >= textLim {
						break
					}

					texts = append(texts, t.text)
					size += tlen
				}

				fmt.Println(`<sending `, len(texts), `x items`, texts, `>`)

				trans, err := r.cl.Translate(r.ctx, texts, r.dst, r.opts)
				if err != nil {
					panic(err)
				}
				fmt.Println(`<sent!>`)

				for i, t := range trans {
					t2ts[i].out <- t.Text
				}

				r.texts2trans <- t2ts[piv:]
			}
		}()

	} else {

		r.texts2trans <- append(t2ts, t2t)
	}

	return t2t.out
}

//
//
//

// Translate moves from input from a source o destination language.
// func (r *Rose) Translate(x string) (<-chan string, error) {

// 	t2t := text2tran{x, make(chan string)}
// 	l := len(x)
// 	bs := <-r.bsize
// 	t2ts := <-r.texts2trans

// 	if bs+l < 5000 && len(t2ts) < 128 {
// 		r.texts2trans <- append(t2ts, t2t)
// 		r.bsize <- bs + l
// 		return t2t.out, nil
// 	}

// 	<-r.call
// 	time.Sleep(time.Duration(bs) * time.Millisecond)

// 	textContent := []string{}
// 	for _, t := range t2ts {
// 		textContent = append(textContent, t.text)
// 	}

// 	trans, err := r.cl.Translate(r.ctx, textContent, r.dst, r.opts)
// 	if err != nil {
// 		return t2t.out, err
// 	}

// 	for i, t := range trans {
// 		t2ts[i].out <- t.Text
// 	}

// 	r.texts2trans <- []text2tran{t2t}
// 	r.bsize <- l
// 	r.call <- true

// 	return t2t.out, nil
// }
