package middleware

import (
	"log"
	"net/http"
	"time"
	"webserver/mydatabase"
)

type MiddleWare struct {
	http_handler   http.Handler
	logging        bool
	authentication bool
}

func NewMiddleWare(http http.Handler) *MiddleWare {
	return &MiddleWare{
		http_handler:   http,
		logging:        true,
		authentication: true,
	}
}

type HandlerFunction func(http.ResponseWriter, *http.Request)

func log_duration(md *NewMiddleWare, f HandlerFunction) {
	f := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		f(w, r)
		log.Println("time it took for http handler - %s, url path - %s, %v", r.Method, r.URL.Path, time.Since(t1))

	}
	return f
}

func authentication(md *NewMiddleWare, f HandlerFunction) {

	f := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		sessionId := cookie.Value
		if sessionId == "" {
			log.Println("cookie not found unauthorized user")
			http.Error(w, "Please login using your username and password", http.StatusUnauthorized)
			return
		} else {
			redisclient := mydatabase.GetRedisClient()
			val, err1 := redisclient.GetValue(r.Context(), r.FormValue("username")).Result()
			if err1 != nil {
				log.Println("error getting session from redis my be expired")
				http.Error(w, "Please login using your username and password", http.StatusUnauthorized)
				return
			} else {
				if val == sessionId {
					//multiplex the call to the old mux
					m.http_handle.ServeHTTP(w, r)
					//any other processin needed for the request
					if f != nil {
						f(w, r)
					}

				} else {
					log.Println("session id and cookie not matching")
					http.Error(w, "Please login using your email id and password", http.StatusUnauthorized)
					return
				}
			}
		}
	}
	return f
}

func (m *MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f := authentication(m, nil)
	f = log_duration(m, f)
}
