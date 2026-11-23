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

// The provider serves its OWN word on EVERY builder leg, not only resolve: the collect-context and
// reverse replies equal kit's own for the literal word "aur", so a dispatch routed to another
// builder fails here, and an op no builder serves is still refused.
func TestInvokeServesThisBuilderWord(t *testing.T) {
	ctx := context.Background()

	collect := spec.BuilderCollectInput{Candy: "c", Builder: "aur"}
	params, _ := json.Marshal(collect)
	reply, err := provider{}.Invoke(ctx, &pb.InvokeRequest{Op: "collect-context", ParamsJson: params})
	if err != nil {
		t.Fatalf("Invoke(collect-context) for %q: %v", builderWord, err)
	}
	wantCtx, err := kit.BuilderCollectContext("aur", collect)
	if err != nil {
		t.Fatalf("kit.BuilderCollectContext(%q): %v", "aur", err)
	}
	wantCollect, _ := json.Marshal(spec.BuilderCollectReply{Context: wantCtx, Package: kit.BuilderPackageOutput("aur", "c")})
	if !bytes.Equal(reply.GetResultJson(), wantCollect) {
		t.Errorf("collect-context served another builder's reply:\n got %s\nwant %s", reply.GetResultJson(), wantCollect)
	}

	// Every key a word's reverse reads, so the leg yields ops whichever builder it is routed to; a
	// dispatch to a different word then answers different ops, which the equality below catches.
	rev := spec.BuilderReverseInput{Candy: "c", Builder: "aur", Context: map[string]any{
		"packages": []any{"x"}, "env_name": "e", "tools": []any{"t"}, "manifest": "m.toml",
		"lib_dir": "/l", "ld_conf": "/e.conf",
	}}
	params, _ = json.Marshal(rev)
	reply, err = provider{}.Invoke(ctx, &pb.InvokeRequest{Op: "reverse", ParamsJson: params})
	if err != nil {
		t.Fatalf("Invoke(reverse): %v", err)
	}
	ops := kit.BuilderReverse("aur", rev)
	if len(ops) == 0 {
		t.Fatalf("kit.BuilderReverse(%q) yielded no teardown ops for a full context; the pin would be vacuous", "aur")
	}
	wantRev, _ := json.Marshal(spec.BuilderReverseReply{ReverseOps: ops})
	if !bytes.Equal(reply.GetResultJson(), wantRev) {
		t.Errorf("reverse served another builder's teardown:\n got %s\nwant %s", reply.GetResultJson(), wantRev)
	}

	if _, err := (provider{}).Invoke(ctx, &pb.InvokeRequest{Op: "bogus"}); err == nil {
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

// The reply this plugin serves carries the REQUIRED output discriminator (spec#424), asserted
// here rather than delegated: against the pre-change sdk this file does not even compile (no
// `Output` field, no kit.BuilderOutput* vocabulary), and a plugin that stopped routing through
// sdk.InvokeBuilder would answer with it empty — which the host now refuses at its first message
// (opencharly/sdk#615).
func TestInvokeResolveReplyCarriesTheOutputDiscriminator(t *testing.T) {
	in := spec.BuilderResolveInput{Candy: "c", Builder: "aur", StageName: "c-aur-build", UID: 1000, GID: 1000, Home: "/home/u", User: "u"}
	params, _ := json.Marshal(in)
	reply, err := provider{}.Invoke(context.Background(), &pb.InvokeRequest{Op: "resolve", ParamsJson: params})
	if err != nil {
		t.Fatalf("Invoke(resolve): %v", err)
	}
	var out spec.BuilderResolveReply
	if err := json.Unmarshal(reply.GetResultJson(), &out); err != nil {
		t.Fatalf("decode reply: %v", err)
	}
	if out.Output != kit.BuilderOutputPackage {
		t.Errorf("the resolve reply carries output=%q, want %q: the host refuses a reply without it", out.Output, "package")
	}
}
