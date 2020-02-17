package phaser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultExecutePanics(t *testing.T) {
	p := Phase{}

	// The default Phase should panic
	assert.Panics(t, func() { _, _ = p.execute(struct{}{}) })
}

func TestDefaultRunPanics(t *testing.T) {
	p := Phase{}

	// The default Phase should panic
	assert.Panics(t, func() { _, _ = p.run(struct{}{}) })
}

func TestAddPreHooks(t *testing.T) {
	p := Phase{}

	// Test hooks
	hooks := []PhaseHook{
		func(value interface{}) (interface{}, error) { return 0, nil },
		func(value interface{}) (interface{}, error) { return 1, nil },
		func(value interface{}) (interface{}, error) { return 2, nil },
	}

	require.Equal(t, len(p.preHooks), 0)

	// Add hooks
	p.appendPreHook(hooks[1])
	p.prependPreHook(hooks[0])
	p.appendPreHook(hooks[2])

	require.Equal(t, len(p.preHooks), len(hooks))

	// Check hook order
	for i, hook := range p.preHooks {
		val, _ := hook(nil)
		require.Equal(t, val.(int), i)
	}
}

func TestAddPostHooks(t *testing.T) {
	p := Phase{}

	// Test hooks
	hooks := []PhaseHook{
		func(value interface{}) (interface{}, error) { return 0, nil },
		func(value interface{}) (interface{}, error) { return 1, nil },
		func(value interface{}) (interface{}, error) { return 2, nil },
	}

	require.Equal(t, len(p.postHooks), 0)

	// Add hooks
	p.appendPostHook(hooks[1])
	p.prependPostHook(hooks[0])
	p.appendPostHook(hooks[2])

	require.Equal(t, len(p.postHooks), len(hooks))

	// Check hook order
	for i, hook := range p.postHooks {
		val, _ := hook(nil)
		require.Equal(t, val.(int), i)
	}
}

func TestProcessHooksPass(t *testing.T) {
	p := Phase{}

	// Test hooks
	p.preHooks = []PhaseHook{
		func(value interface{}) (interface{}, error) { return value.(int) + 1, nil },
		func(value interface{}) (interface{}, error) { return value.(int) + 2, nil },
		func(value interface{}) (interface{}, error) { return value.(int) + 3, nil },
	}

	value, err := p.processHooks(0, &p.preHooks)
	assert.NoError(t, err)
	assert.Equal(t, value, 6)
}

func TestProcessHooksFail(t *testing.T) {
	p := Phase{}

	// Test hooks
	p.preHooks = []PhaseHook{
		func(value interface{}) (interface{}, error) { return value.(int) + 1, nil },
		func(value interface{}) (interface{}, error) { return value.(int) + 2, assert.AnError },
		func(value interface{}) (interface{}, error) { return value.(int) + 3, nil },
	}

	value, err := p.processHooks(0, &p.preHooks)
	assert.Nil(t, value)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestTestPhaseExecute(t *testing.T) {
	p := Phase{
		execute: func (value interface{}) (interface{}, error) {
			return value, nil
		},
	}
	val := 1
	value, err := p.execute(val)

	// The test Phase should not panic and return the same value it receives
	assert.NoError(t, err)
	assert.Equal(t, value.(int), val)
}

func TestTestPhaseRun(t *testing.T) {
	val := 1
	p := Phase{
		preHooks: []PhaseHook{
			// Multiply the input value by two
			func(value interface{}) (interface{}, error) {
				return value.(int) * 2, nil
			},
			// Check that the input value is now two times val
			func(value interface{}) (interface{}, error) {
				var err error
				if value.(int) != val * 2 {
					err = assert.AnError
				}
				return value, err
			},
		},
		// Check that the input value is two times val
		execute: func(value interface{}) (interface{}, error) {
			var err error
			if value.(int) != val * 2 {
				err = assert.AnError
			}
			return value, err
		},
		// Divide output value by two
		postHooks: []PhaseHook{
			// Check that the output value is two times val
			func(value interface{}) (interface{}, error) {
				var err error
				if value.(int) != val * 2 {
					err = assert.AnError
				}
				return value, err
			},
			// Divide the output value by two
			func(value interface{}) (interface{}, error) {
				return value.(int) / 2, nil
			},
			// Check that the output value is now val
			func(value interface{}) (interface{}, error) {
				var err error
				if value.(int) != val {
					err = assert.AnError
				}
				return value, err
			},
		},
	}

	value, err := p.run(val)

	// The test Phase should not return an error and should return the same
	// value it receives
	assert.NoError(t, err)
	assert.Equal(t, value.(int), val)
}

