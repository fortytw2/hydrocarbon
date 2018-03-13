package hydrocarbon

import (
	"net/http"
	"testing"
)

func TestGetRemoteIP(t *testing.T) {
	t.Parallel()

	{
		r, _ := http.NewRequest(http.MethodGet, "https://www.hydrocarbon.io/", nil)
		r.Header.Set("X-Forwarded-For", "193.167.12.23, 204.121.12.21")

		ip := getRemoteIP(r)
		if ip != "193.167.12.23" {
			t.Error("x-forwarded-for is broken")
		}
	}

	{
		r, _ := http.NewRequest(http.MethodGet, "https://www.hydrocarbon.io/", nil)
		r.RemoteAddr = "123.34.121.121"

		ip := getRemoteIP(r)
		if ip != "123.34.121.121" {
			t.Error("remote addr fallback is broken")
		}
	}
}
