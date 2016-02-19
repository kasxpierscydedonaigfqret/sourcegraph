// +build ignore

package main

import (
	"text/template"

	"src.sourcegraph.com/sourcegraph/gen"
)

func main() {
	svcs := []string{
		"../go-sourcegraph/sourcegraph/sourcegraph.pb.go",
		"../vendor/sourcegraph.com/sourcegraph/srclib/store/pb/srcstore.pb.go",
		"../gitserver/gitpb/git_transport.pb.go",
	}
	gen.Generate("context.go", tmpl, svcs, nil, "")
}

var tmpl = template.Must(template.New("").Delims("<<<", ">>>").Parse(`// GENERATED CODE - DO NOT EDIT!
//
// Generated by:
//
//   go run gen_context.go
//
// Called via:
//
//   go generate
//

package svc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"src.sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"src.sourcegraph.com/sourcegraph/gitserver/gitpb"
	"sourcegraph.com/sourcegraph/srclib/store/pb"
)

type contextKey int

const (
	<<<range .>>>_<<<.Name>>>Key contextKey = iota
	<<<end>>>
)

// Services contains fields for all existing services.
type Services struct {
	<<<range .>>><<<.Name>>> <<<.TypeName>>>
	<<<end>>>
}

// RegisterAll calls all of the the RegisterXxxServer funcs.
func RegisterAll(s *grpc.Server, svcs Services) {
	<<<range .>>>
		if svcs.<<<.Name>>> != nil {
			<<<.PkgName>>>.Register<<<.Name>>>Server(s, svcs.<<<.Name>>>)
		}
	<<<end>>>
}

// WithServices returns a copy of parent with the given services. If a service's field value is nil, its previous value is inherited from parent in the new context.
func WithServices(ctx context.Context, s Services) context.Context {
	<<<range .>>>
		if s.<<<.Name>>> != nil {
			ctx = With<<<.Name>>>(ctx, s.<<<.Name>>>)
		}
	<<<end>>>
	return ctx
}

<<<range .>>>
	// With<<<.Name>>> returns a copy of parent that uses the given <<<.Name>>> service.
	func With<<<.Name>>>(ctx context.Context, s <<<.TypeName>>>) context.Context {
		return context.WithValue(ctx, _<<<.Name>>>Key, s)
	}

	// <<<.Name>>> gets the context's <<<.Name>>> service. If the service is not present, it panics.
	func <<<.Name>>>(ctx context.Context) <<<.TypeName>>> {
		s, ok := ctx.Value(_<<<.Name>>>Key).(<<<.TypeName>>>)
		if !ok || s == nil {
			panic("no <<<.Name>>> set in context")
		}
		return s
	}

	// <<<.Name>>>OrNil returns the context's <<<.Name>>> service if present, or else nil.
	func <<<.Name>>>OrNil(ctx context.Context) <<<.TypeName>>> {
		s, ok := ctx.Value(_<<<.Name>>>Key).(<<<.TypeName>>>)
		if ok {
			return s
		}
		return nil
	}
<<<end>>>
`))
