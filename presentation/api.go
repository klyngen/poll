package presentation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/klyngen/votomatic-3000/packages/backend/models"
	"golang.org/x/net/websocket"
)

type presentationAPI struct {
	config      models.Configuration
	poll        models.Poll
	configBytes []byte
	listeners   []chan models.PollAnswer
}

func NewPresentationAPI(config models.Configuration) *presentationAPI {
	configBytes, _ := makeConfigBytes(config)
	return &presentationAPI{
		config:      config,
		configBytes: configBytes,
		poll:        *models.NewPoll(config),
	}
}

func (g *presentationAPI) addListener(listener chan models.PollAnswer) {
	g.listeners = append(g.listeners, listener)
}

func (g *presentationAPI) updateListeners(index int, alternative int) {
	update := models.PollAnswer{
		Id:          index,
		Alternative: alternative,
	}
	for _, vs := range g.listeners {
		vs <- update
	}
}

func (p *presentationAPI) listenToSocket(port int) error {
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), websocket.Handler(p.handleSocketConnection))
}

func (p *presentationAPI) handleSocketConnection(connection *websocket.Conn) {
	channel := make(chan models.PollAnswer)
	p.addListener(channel)
	data := make([]byte, 128)
	for {
		if len(data) == 0 {
			log.Println("Client closed the connection")
			break
		}
		select {
		case update := <-channel:
			bytes, err := json.Marshal(update)
			connection.Write(bytes)

			if err == nil {
				fmt.Println(err)
			} else {
				log.Println("unable to marshal json")
			}
		}
	}
}

func (p *presentationAPI) ListenAndServe(restPort int) error {
	log.Printf("Starting rest on port: %v\n", restPort)
	router := chi.NewMux()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/questions", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write(p.configBytes)
	})

	router.Get("/pollresult", func(w http.ResponseWriter, r *http.Request) {
		answers := p.poll.GetAnswers()

		jsonBytes, err := json.Marshal(answers)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonBytes)
	})

	router.Post("/pollresult", func(w http.ResponseWriter, r *http.Request) {
		rawId := r.URL.Query().Get("id")
		rawAlternative := r.URL.Query().Get("alternative")

		id, _ := strconv.ParseInt(rawId, 10, 32)
		alternativeId, _ := strconv.ParseInt(rawAlternative, 10, 32)

		p.poll.AddAnswer(models.PollAnswer{
			Id:          int(id),
			Alternative: int(alternativeId),
		})

		go p.updateListeners(int(id), int(alternativeId))
	})

	router.Patch("/pollresult", func(w http.ResponseWriter, r *http.Request) {
		apikey := r.Header.Get("X-API-KEY")

		if apikey == os.Getenv("API_KEY") {
			p.poll = *models.NewPoll(p.config)
			return
		}

		w.WriteHeader(401)
	})

	fileServerSPA(router, "/", "./public")

	printRoutes(router)

	http.ListenAndServe(fmt.Sprintf(":%v", restPort), router)
	return nil
}

// FileServer para SPA
func fileServerSPA(r chi.Router, public string, static string) {

	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", 301).ServeHTTP)
		public += "/"
	}

	r.Get(public+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))
}

func printRoutes(router *chi.Mux) {
	chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s]:\t'%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
}

func makeConfigBytes(config models.Configuration) ([]byte, error) {
	return json.Marshal(config.GetQuestions())
}
