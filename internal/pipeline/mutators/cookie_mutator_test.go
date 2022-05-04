package mutators

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dadrus/heimdall/internal/heimdall"
	"github.com/dadrus/heimdall/internal/pipeline/subject"
	"github.com/dadrus/heimdall/internal/pipeline/template"
	"github.com/dadrus/heimdall/internal/pipeline/testsupport"
)

func TestCreateCookieMutator(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		uc     string
		config []byte
		assert func(t *testing.T, err error, mut *cookieMutator)
	}{
		{
			uc: "without configuration",
			assert: func(t *testing.T, err error, mut *cookieMutator) {
				t.Helper()

				require.Error(t, err)
				assert.ErrorIs(t, err, heimdall.ErrConfiguration)
				assert.Contains(t, err.Error(), "no cookie")
			},
		},
		{
			uc:     "without cookie configuration",
			config: []byte(``),
			assert: func(t *testing.T, err error, mut *cookieMutator) {
				t.Helper()

				require.Error(t, err)
				assert.ErrorIs(t, err, heimdall.ErrConfiguration)
				assert.Contains(t, err.Error(), "no cookie")
			},
		},
		{
			uc: "with unsupported attributes",
			config: []byte(`
cookies:
  foo: bar
foo: bar
`),
			assert: func(t *testing.T, err error, mut *cookieMutator) {
				t.Helper()

				require.Error(t, err)
				assert.ErrorIs(t, err, heimdall.ErrConfiguration)
				assert.Contains(t, err.Error(), "failed to unmarshal")
			},
		},
		{
			uc: "with valid config",
			config: []byte(`
cookies:
  foo: bar
  bar: "{{ .ID }}"`),
			assert: func(t *testing.T, err error, mut *cookieMutator) {
				t.Helper()

				require.NoError(t, err)
				assert.Len(t, mut.cookies, 2)
				assert.Equal(t, template.Template("bar"), mut.cookies["foo"])
				assert.Equal(t, template.Template("{{ .ID }}"), mut.cookies["bar"])
			},
		},
	} {
		t.Run("case="+tc.uc, func(t *testing.T) {
			conf, err := testsupport.DecodeTestConfig(tc.config)
			require.NoError(t, err)

			// WHEN
			mutator, err := newCookieMutator(conf)

			// THEN
			tc.assert(t, err, mutator)
		})
	}
}

func TestCreateCookieMutatorFromPrototype(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		uc              string
		prototypeConfig []byte
		config          []byte
		assert          func(t *testing.T, err error, prototype *cookieMutator, configured *cookieMutator)
	}{
		{
			uc: "no new configuration provided",
			prototypeConfig: []byte(`
cookies:
  foo: bar
`),
			assert: func(t *testing.T, err error, prototype *cookieMutator, configured *cookieMutator) {
				t.Helper()

				require.NoError(t, err)
				assert.Equal(t, prototype, configured)
			},
		},
		{
			uc: "configuration without cookies provided",
			prototypeConfig: []byte(`
cookies:
  foo: bar
`),
			config: []byte(``),
			assert: func(t *testing.T, err error, prototype *cookieMutator, configured *cookieMutator) {
				t.Helper()

				require.NoError(t, err)
				assert.Equal(t, prototype, configured)
			},
		},
		{
			uc: "new cookies provided",
			prototypeConfig: []byte(`
cookies:
  foo: bar
`),
			config: []byte(`
cookies:
  bar: foo
`),
			assert: func(t *testing.T, err error, prototype *cookieMutator, configured *cookieMutator) {
				t.Helper()

				require.NoError(t, err)
				assert.NotEqual(t, prototype, configured)
				require.NotNil(t, configured)
				assert.NotEmpty(t, configured.cookies)
				assert.Equal(t, template.Template("foo"), configured.cookies["bar"])
			},
		},
	} {
		t.Run("case="+tc.uc, func(t *testing.T) {
			pc, err := testsupport.DecodeTestConfig(tc.prototypeConfig)
			require.NoError(t, err)

			conf, err := testsupport.DecodeTestConfig(tc.config)
			require.NoError(t, err)

			prototype, err := newCookieMutator(pc)
			require.NoError(t, err)

			// WHEN
			mut, err := prototype.WithConfig(conf)

			// THEN
			cookieMut, ok := mut.(*cookieMutator)
			require.True(t, ok)

			tc.assert(t, err, prototype, cookieMut)
		})
	}
}

func TestCookieMutatorExecute(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		uc               string
		config           []byte
		configureContext func(t *testing.T, ctx *testsupport.MockContext)
		createSubject    func(t *testing.T) *subject.Subject
		assert           func(t *testing.T, err error)
	}{
		{
			uc: "with nil subject",
			config: []byte(`
cookies:
  foo: bar
  bar: "{{ .ID }}"
`),
			configureContext: func(t *testing.T, ctx *testsupport.MockContext) {
				// nothing is required here
				t.Helper()
			},
			createSubject: func(t *testing.T) *subject.Subject {
				t.Helper()

				return nil
			},
			assert: func(t *testing.T, err error) {
				t.Helper()

				require.Error(t, err)
				assert.ErrorIs(t, err, heimdall.ErrInternal)
				assert.Contains(t, err.Error(), "'nil' subject")
			},
		},
		{
			uc: "with bad template",
			config: []byte(`
cookies:
  bar: "{{ .ID | foobar }}"
`),
			configureContext: func(t *testing.T, ctx *testsupport.MockContext) {
				// nothing is required here
				t.Helper()
			},
			createSubject: func(t *testing.T) *subject.Subject {
				t.Helper()

				return &subject.Subject{ID: "FooBar"}
			},
			assert: func(t *testing.T, err error) {
				t.Helper()

				require.Error(t, err)
				assert.ErrorIs(t, err, heimdall.ErrInternal)
				assert.Contains(t, err.Error(), "failed to render")
			},
		},
		{
			uc: "with all preconditions satisfied",
			config: []byte(`
cookies:
  foo: "{{ .Attributes.bar }}"
  bar: "{{ .ID }}"
  baz: bar
`),
			configureContext: func(t *testing.T, ctx *testsupport.MockContext) {
				t.Helper()

				ctx.On("AddResponseCookie", "foo", "baz")
				ctx.On("AddResponseCookie", "bar", "FooBar")
				ctx.On("AddResponseCookie", "baz", "bar")
			},
			createSubject: func(t *testing.T) *subject.Subject {
				t.Helper()

				return &subject.Subject{ID: "FooBar", Attributes: map[string]any{"bar": "baz"}}
			},
			assert: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
	} {
		t.Run("case="+tc.uc, func(t *testing.T) {
			// GIVEN
			conf, err := testsupport.DecodeTestConfig(tc.config)
			require.NoError(t, err)

			mctx := &testsupport.MockContext{}
			mctx.On("AppContext").Return(context.Background())

			sub := tc.createSubject(t)

			tc.configureContext(t, mctx)

			mutator, err := newCookieMutator(conf)
			require.NoError(t, err)

			// WHEN
			err = mutator.Execute(mctx, sub)

			// THEN
			tc.assert(t, err)

			mctx.AssertExpectations(t)
		})
	}
}
