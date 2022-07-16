package http

import "net/http"

func (s *Server) CheckProfiling() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if s.config.Server.Profiling {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "{\"error\": \"pprof is off\"}", http.StatusServiceUnavailable)
			}
		}

		return http.HandlerFunc(fn)
	}
}
