package driver

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func InterceptRequests(page *rod.Page) *rod.HijackRouter {
	router := page.HijackRequests()

	cancelReq := func(ctx *rod.Hijack) {
		if ctx.Request.Type() == proto.NetworkResourceTypeImage {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	}

	router.MustAdd("*.png", cancelReq)
	router.MustAdd("*.svg", cancelReq)
	router.MustAdd("*.gif", cancelReq)
	router.MustAdd("*.css", cancelReq)

	return router
}

func Sel(el *rod.Element, value string) error {
	regex := fmt.Sprintf("^%s$", value)
	return el.Select([]string{regex}, true, rod.SelectorTypeRegex)
}
