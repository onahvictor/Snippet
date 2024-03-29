package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	chain "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	//initialize the router
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	//Leave the static files route unchanged No need for a session
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// "/static" prefix before the request reaches the file server.
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := chain.New(app.sessionManager.LoadAndSave, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// return app.recoverPanic(app.logRequest(secureHeaders(router)))
	standard := chain.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
