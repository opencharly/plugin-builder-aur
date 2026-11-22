package builderaur

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/opencharly/sdk/kit"
	pb "github.com/opencharly/spec/proto"
	"github.com/opencharly/spec/spec"
)

// The provider serves its OWN word through the shared builder dispatch (sdk.InvokeBuilder): a
// collect-context call answers for "aur", and an op no builder serves is refused, never ignored.
func TestInvokeServesThisBuilderWord(t *testing.T) {
	in, _ := json.Marshal(spec.BuilderCollectInput{Candy: "c", Builder: builderWord})
	reply, err := provider{}.Invoke(context.Background(), &pb.InvokeRequest{Op: "collect-context", ParamsJson: in})
	if err != nil {
		t.Fatalf("Invoke(collect-context) for %q: %v", builderWord, err)
	}
	var got spec.BuilderCollectReply
	if err := json.Unmarshal(reply.GetResultJson(), &got); err != nil {
		t.Fatalf("reply is not a BuilderCollectReply: %v (%s)", err, reply.GetResultJson())
	}
	if _, err := (provider{}).Invoke(context.Background(), &pb.InvokeRequest{Op: "bogus"}); err == nil {
		t.Error("an unsupported op was not refused")
	}
}

// The provider resolves its OWN word's stage: the reply its Invoke serves equals kit.BuilderResolve
// for the literal word "aur", so a dispatch that routed this plugin to any other builder fails here.
func TestInvokeResolvesThisBuilderWordsOwnStage(t *testing.T) {
	in := spec.BuilderResolveInput{Candy: "c", Builder: "aur", BuilderRef: "builder:v1", StageName: "c-aur-build", CopySrc: "candy/c", UID: 1000, GID: 1000, Home: "/home/u", User: "u", Manifest: "manifest"}
	params, _ := json.Marshal(in)
	reply, err := provider{}.Invoke(context.Background(), &pb.InvokeRequest{Op: "resolve", ParamsJson: params})
	if err != nil {
		t.Fatalf("Invoke(resolve): %v", err)
	}
	want, err := kit.BuilderResolve("aur", in)
	if err != nil {
		t.Fatalf("kit.BuilderResolve(%q): %v", "aur", err)
	}
	wantJSON, _ := json.Marshal(want)
	if !bytes.Equal(reply.GetResultJson(), wantJSON) {
		t.Errorf("Invoke(resolve) served another builder's stage:\n got %s\nwant %s", reply.GetResultJson(), wantJSON)
	}
}
