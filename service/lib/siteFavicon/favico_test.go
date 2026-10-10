package siteFavicon

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFallbackAndResolution(t *testing.T) {
	// 页面没有 <link rel=icon>：应回退到 /favicon.ico，并取回标题
	mux := http.NewServeMux()
	mux.HandleFunc("/no-link", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title> NoLink </title></head><body></body></html>`))
	})
	// 页面只有 apple-touch-icon 和一个 icon：应优先 icon
	mux.HandleFunc("/has-link", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>Has</title>
<link rel="apple-touch-icon" href="/apple.png">
<link rel="icon" href="assets/icon.ico?v=2"></head></html>`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	u, title, err := GetOneFaviconURLAndTitle(srv.URL + "/no-link")
	if err != nil {
		t.Fatalf("no-link err: %v", err)
	}
	if title != "NoLink" {
		t.Fatalf("title=%q", title)
	}
	if u != srv.URL+"/favicon.ico" {
		t.Fatalf("fallback url=%q", u)
	}

	u2, _, err := GetOneFaviconURLAndTitle(srv.URL + "/has-link")
	if err != nil {
		t.Fatalf("has-link err: %v", err)
	}
	if u2 != srv.URL+"/assets/icon.ico?v=2" {
		t.Fatalf("resolved url=%q", u2)
	}
}
