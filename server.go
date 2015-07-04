package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"code.google.com/p/go.net/context"
	"github.com/go-gorp/gorp"
	gtx "github.com/goji/context"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func dbmContext(dbm *gorp.DbMap) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, "dbm", dbm)
}

func dbmMiddleware(ctx context.Context) web.MiddlewareType {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			gtx.Set(c, ctx)
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Server(bind, src string) error {
	dbm, err := InitDb(src)
	if err != nil {
		return err
	}
	defer dbm.Db.Close()

	flag.Set("bind", bind)

	m := goji.DefaultMux
	m.Use(dbmMiddleware(dbmContext(dbm)))
	m.Use(middleware.Recoverer)
	m.Use(middleware.NoCache)

	m.Get("/", indexHandler)
	m.Get("/api/equipments", equipmentsHandler)

	goji.Serve()

	return nil
}

func serverError(w http.ResponseWriter, err error) {
	log.Println(err.Error())

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
}

func indexHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	a, _ := Asset("html/index.html")
	fmt.Fprintln(w, string(a))
}

func equipmentsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	dbm := gtx.FromC(c).Value("dbm").(*gorp.DbMap)
	results, err := SelectAllEquipmentInfo(dbm)
	if err != nil {
		serverError(w, err)
		return
	}

	res, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		serverError(w, err)
		return
	}

	fmt.Fprintln(w, string(res))
}
