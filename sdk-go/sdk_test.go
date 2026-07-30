package sdk

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/peterrk/simple-abtest/engine/core"
)

type staticTransport struct {
	body string
}

func (t staticTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(t.body)),
	}, nil
}

func TestUpdateNormalizesExperimentFilters(t *testing.T) {
	client := &Client{
		srcUrl: "http://unit.test/app/1",
		client: &http.Client{Transport: staticTransport{body: `[
			{"filter":[{"op":4,"dtype":1,"key":"bucket","ss":["A","B"]}]},
			{"filter":[{"op":5,"dtype":1,"key":"bucket","ss":["A","B"]}]}
		]`}},
		ctx: context.Background(),
	}
	if err := client.update(); err != nil {
		t.Fatal(err)
	}

	exps := client.data.Load().exps
	if !core.EvalExpr(exps[0].Filter, map[string]string{"bucket": "A"}) {
		t.Fatal("IN did not match a set member after snapshot update")
	}
	if core.EvalExpr(exps[1].Filter, map[string]string{"bucket": "A"}) {
		t.Fatal("NOT IN matched a set member after snapshot update")
	}
}
