package matcher

import (
	"github.com/dadrus/heimdall/internal/heimdall"
	"github.com/dadrus/heimdall/internal/x"
)

type ErrorConditionMatcher struct {
	Error  *ErrorTypeMatcher `mapstructure:"error"`
	CIDR   *CIDRMatcher      `mapstructure:"request_cidr"`
	Header *HeaderMatcher    `mapstructure:"request_header"`
}

func (ecm ErrorConditionMatcher) Match(ctx heimdall.Context, err error) bool {
	errorMatched := x.IfThenElseExec(ecm.Error != nil,
		func() bool { return ecm.Error.Match(err) },
		func() bool { return true })

	ipMatched := x.IfThenElseExec(ecm.CIDR != nil,
		func() bool { return ecm.CIDR.Match(ctx.RequestClientIPs()...) },
		func() bool { return true })

	headerMatched := x.IfThenElseExec(ecm.Header != nil,
		func() bool { return ecm.Header.Match(ctx.RequestHeaders()) },
		func() bool { return true })

	return errorMatched && ipMatched && headerMatched
}
