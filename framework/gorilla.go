package framework

import "github.com/gorilla/mux"

func SetupMux() {
	_ = mux.NewRouter()
}
