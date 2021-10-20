package server

import (
	"encoding/gob"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

func dbLocation() string {
	env := os.Getenv("DB_LOCATION")
	if env == "" {
		res := "store.db"
		log.Printf("using %s as db location", res)

		return res
	} else {
		return env
	}
}

func StartServer() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	gob.Register(SessionUser{})

	store, err := NewStore(dbLocation())
	if err != nil {
		return err
	}
	defer store.Close()

	lm, err := NewLoginManager(store)
	if err != nil {
		return err
	}

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string
			Password string
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		su, err := lm.LogIn(User{
			Name:     body.Username,
			Password: []byte(body.Password),
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tu, err := NewTokenUser(su)

		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(&tu)
		if err != nil {
			log.Printf("%v", err)
		}
	})

	r.Put("/password", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Password string
			Token string
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := lm.DecodeToken(body.Token)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(body.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("password too short (less than 8)"))
			return
		}

		err = lm.ChangePassword(user, body.Password)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	r.Put("/admin", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string
			Admin bool
			Token string
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := lm.DecodeToken(body.Token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !user.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if body.Name == user.Name {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("can't change your own admin status"))
			return
		}

		err = lm.SetAdmin(body.Name, body.Admin)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	})

	r.Delete("/user", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string
			Token string
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := lm.DecodeToken(body.Token)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !user.Admin && user.Name != body.Name {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userToDelete, err := store.GetUser(body.Name)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		num, err := store.CountAdminUsers()
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}


		if num == 1 && userToDelete.Admin {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("can't delete last (admin) user"))
			return
		}

		err = store.RmUser(body.Name)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	})

	r.Post("/user/get-all", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Token string
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := lm.DecodeToken(body.Token)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}


		if !user.Admin {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		users, err := store.GetUsers()
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i := 0; i < len(users); i++ {
			users[i].Password = nil
		}

		err = json.NewEncoder(w).Encode(&users)
		if err != nil {
			log.Printf("%v", err)
		}
	})

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string
			Password string
			Admin bool
			Token string
		}

		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := lm.DecodeToken(body.Token)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !user.Admin {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if body.Username == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("username cannot be empty"))
			return
		}

		newUser := User{
			Name:     body.Username,
			Password: []byte(body.Password),
			Admin:    body.Admin,
			Tabs:     nil,
		}
		exists, err := lm.CreateUser(newUser)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if exists {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("user with this name exists"))
			return
		}


		newUser.Password = nil
		err = json.NewEncoder(w).Encode(&newUser)
		if err != nil {
			log.Printf("%v", err)
		}
	})

	r.Route("/tab", func(r chi.Router) {
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Id string
				Token string
			}

			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			user, err := lm.DecodeToken(body.Token)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusUnauthorized)
			}


			tabId, err := uuid.Parse(body.Id)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			alias, err := store.GetTab(tabId)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if alias.Owner != user.Name && !user.Admin {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			err = store.RmTab(alias)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		})

		r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Token string
			}

			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			user, err := lm.DecodeToken(body.Token)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tab := Tab{
				Id:       uuid.New(),
				Owner:    user.Name,
				Public:   false,
				Contents: "",
			}

			err = store.CreateTab(tab)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(tab)
			if err != nil {
				log.Printf("%v", err)
			}
			return
		})

		r.Post("/all-for-user", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Token string
			}

			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			user, err := lm.DecodeToken(body.Token)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			res, err := store.GetUserTabs(&user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				log.Printf("%v", err)
			}
		})

		r.Post("/all-public", func(w http.ResponseWriter, r *http.Request) {
			res, err := store.GetPublicTabs()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				log.Printf("%v", err)
			}
		})

		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Token string
				Data string
				Id string
			}

			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			user, err := lm.DecodeToken(body.Token)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			id, err := uuid.Parse(body.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			tab, err := store.GetTab(id)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if tab.Owner != user.Name {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tab.Contents = body.Data

			err = store.SetTab(tab.Id, tab)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})

		r.Put("/public", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Token string
				Id string
				Public bool
			}

			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			user, err := lm.DecodeToken(body.Token)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			id, err := uuid.Parse(body.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			tab, err := store.GetTab(id)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if tab.Owner != user.Name {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tab.Public = body.Public

			err = store.SetTab(tab.Id, tab)
			if err != nil {
				log.Printf("%v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})


		r.Post("/get", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Id string
				Token string
			}

			err = json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			id, err := uuid.Parse(body.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			res, err := store.GetTab(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !res.Public {
				user, err := lm.DecodeToken(body.Token)
				if err != nil {
					log.Printf("%v", err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if user.Name != res.Owner {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				log.Printf("%v", err)
			}
		})
	})


	url := "0.0.0.0:3000"
	log.Printf("listening on %s", url)
	return http.ListenAndServe(url, r)
}
