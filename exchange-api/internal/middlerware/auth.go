package middlerware

import (
	"common"
	"common/tools"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/net/context"
	"net/http"
)

func Auth(secret string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			result := common.NewResult()
			result.Fail(4000, "no login")
			token := r.Header.Get("x-auth-token")
			if token == "" {
				httpx.WriteJson(w, 200, result)
				return
			}
			userId, err := tools.ParseToken(token, secret)
			if err != nil {
				httpx.WriteJson(w, 200, result)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", userId)
			r = r.WithContext(ctx)
			next(w, r)
		}
	}
}
