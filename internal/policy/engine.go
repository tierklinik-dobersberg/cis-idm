package policy

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
)

// A list of common package names.
const (
	// PackageForwardAuth is the package name for all policies related to
	// forward-authentication using a supported reverse proxy.
	PackageForwardAuth = "cisidm.forward_auth"
)

var (
	ErrNoResults         = errors.New("query returned no results")
	ErrToManyResults     = errors.New("query returned to many results")
	ErrInvalidResultType = errors.New("query returned an invalid result type")
)

type Engine struct {
	compiler *ast.Compiler

	option
}

type option struct {
	extraModules map[string]string
	debug        bool
}

type EngineOption func(*option)

func WithDebug() EngineOption {
	return func(o *option) {
		o.debug = true
	}
}

func WithRawPolicy(path, content string) EngineOption {
	return func(o *option) {
		if o.extraModules == nil {
			o.extraModules = make(map[string]string)
		}

		o.extraModules[path] = content
	}
}

func NewEngine(ctx context.Context, paths []string, opts ...EngineOption) (*Engine, error) {
	var options option

	for _, fn := range opts {
		fn(&options)
	}

	modules, err := loader.AllRegos(paths)
	if err != nil {
		return nil, fmt.Errorf("failed to load rego files: %w", err)
	}

	moduleMap := modules.ParsedModules()

	for name, content := range options.extraModules {
		parsed, err := ast.ParseModule(name, content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse extra module %s: %w", name, err)
		}

		moduleMap[name] = parsed
	}

	compiler := ast.NewCompiler()

	compiler.Compile(moduleMap)
	if compiler.Failed() {
		return nil, fmt.Errorf("failed to compile rego policies: %w", compiler.Errors)
	}

	log.L(ctx).
		WithField("modules", len(compiler.Modules)).
		Infof("policy engine prepared")

	e := &Engine{
		compiler: compiler,
		option:   options,
	}

	return e, nil
}

func (engine *Engine) Query(
	ctx context.Context,
	query string,
	input any,
) (rego.ResultSet, error) {
	options := []func(*rego.Rego){
		rego.Imports([]string{"rego.v1"}),
		rego.Query(query),
		rego.Input(input),
		rego.Compiler(engine.compiler),

		// we always add a print hook so users can debug their policies
		// without enabling the whole tracing and dumping thing...
		rego.PrintHook(topdown.NewPrintHook(os.Stderr)),
	}

	if engine.option.debug {
		log.L(ctx).Infof("rego-tracer: preparing to query %q in debug mode", query)

		tracer := new(topdown.BufferTracer)

		options = append(options, rego.Trace(true))
		options = append(options, rego.QueryTracer(tracer))
		options = append(options, rego.Dump(os.Stderr))

		defer func() {
			for _, evt := range *tracer {
				log.L(ctx).Infof("rego-tracer: %s", evt.String())
			}
		}()
	}

	eval := rego.New(options...)

	result, err := eval.Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate rego query: %w", err)
	}

	if engine.option.debug {
		log.L(ctx).WithField("result", result).Infof("rego-tracer: policy result")
	}

	return result, nil
}

func (engine *Engine) QueryOne(
	ctx context.Context,
	query string,
	input any,
	target any,
) error {
	res, err := engine.Query(ctx, query, input)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return ErrNoResults
	}

	if len(res) > 1 {
		return ErrToManyResults
	}

	if len(res[0].Expressions) == 0 {
		// TODO(ppacher): return a different error here?
		return ErrNoResults
	}

	if len(res[0].Expressions) > 1 {
		// TODO(ppacher): return a different error here?
		return ErrToManyResults
	}

	expr := res[0].Expressions[0].Value
	if _, ok := expr.(map[string]any); !ok {
		return fmt.Errorf("%w: expected a map[string]any, got %T", ErrInvalidResultType, expr)
	}

	if err := mapstructure.WeakDecode(expr, target); err != nil {
		return fmt.Errorf("failed to convert result to target: %w", err)
	}

	return nil
}
