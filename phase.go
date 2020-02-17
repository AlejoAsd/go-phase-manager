package phaser

import "fmt"

// PhaseHook is the hook type used by Phaser implementations.
type PhaseHook func(value interface{}) (interface{}, error)

// Phaser is an interface for phases. You should rarely need to implement Phaser
// from scratch. Instead, include the Phase struct in your own struct and
// override the necessary methods.
type Phaser interface {
	// run runs the phase. It calls the phase's pre-hooks, followed by its
	// execute method, and finally its post-hooks.
	run(value interface{}) (interface{}, error)
	// handleError handles any errors returned during any point in the phase. It
	// should cleanup and tear down resources if necessary.
	handleError(err error) error
	// prependHook prepends a PhaseHook function to the target PhaseHook slice
	prependHook(hooks *[]PhaseHook, newHook PhaseHook)
	// prependPreHook prepends a PhaseHook function to the PreHook slice
	prependPreHook(hook PhaseHook)
	// prependPreHook prepends a PhaseHook function to the PostHook slice
	prependPostHook(hook PhaseHook)
	// appendHook appends a PhaseHook function to the target PhaseHook slice
	appendHook(hooks *[]PhaseHook, newHook PhaseHook)
	// appendPreHook appends a PhaseHook function to the PreHook slice
	appendPreHook(hook PhaseHook)
	// appendPreHook appends a PhaseHook function to the PostHook slice
	appendPostHook(hook PhaseHook)
}

type Phase struct {
	// Name contains the name of the phase. This value should be unique as it
	// will be the phase identifier
	Name string
	// preHooks contains the hooks ran before the execution phase. Used to
	// validate/preprocess phase input data
	preHooks []PhaseHook
	// execute performs the phase's action.
	execute func (value interface{}) (interface{}, error)
	// postHooks contains the hooks ran after the execution phase. Used to
	// validate/postprocess phase output data
	postHooks []PhaseHook
}

func (p *Phase) run(value interface{}) (interface{}, error) {
	var err error

	// Process pre-hooks
	if value, err = p.processHooks(value, &p.preHooks); err != nil {
		return value, err
	}
	// Execute phase
	if p.execute == nil {
		panic(fmt.Sprintf("phase %s not implemented", p.Name))
	}
	if value, err = p.execute(value); err != nil {
		return p.handleError(err)
	}
	// Process post-hooks
	if value, err = p.processHooks(value, &p.postHooks); err != nil {
		return value, err
	}

	return value, nil
}

// handleError handles any errors that may come up. If not overriden, it will
// simply return the raised error.
func (p *Phase) handleError(err error) (interface{}, error) {
	return nil, err
}

// processHooks receives an input value and processes it using a list of hook
// functions
func (p *Phase) processHooks(value interface{}, hooks *[]PhaseHook) (interface{}, error) {
	var err error

	for _, hook := range *hooks {
		if value, err = hook(value); err != nil {
			return p.handleError(err)
		}
	}

	return value, nil
}


func (p *Phase) prependHook(hooks *[]PhaseHook, newHook PhaseHook) {
	*hooks = append([]PhaseHook{newHook}, *hooks...)
}

func (p *Phase) prependPreHook(hook PhaseHook) {
	p.prependHook(&p.preHooks, hook)
}

func (p *Phase) appendPreHook(hook PhaseHook) {
	p.appendHook(&p.preHooks, hook)
}

func (p *Phase) appendHook(hooks *[]PhaseHook, newHook PhaseHook) {
	*hooks = append(*hooks, newHook)
}

func (p *Phase) appendPostHook(hook PhaseHook) {
	p.appendHook(&p.postHooks, hook)
}

func (p *Phase) prependPostHook(hook PhaseHook) {
	p.prependHook(&p.postHooks, hook)
}


