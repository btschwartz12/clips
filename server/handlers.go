package server

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/btschwartz12/clips/server/assets"
	"github.com/go-chi/chi/v5"
)

type HomePageData struct {
	Tabs []HomePageTab
}

type HomePageTab struct {
	Name      string
	NextRun   string
	ClipName  string
	Guessable bool
}

var (
	homeTmpl = template.Must(template.ParseFS(
		assets.Templates,
		"templates/home.html.tmpl",
	))
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	var tabs []HomePageTab

	for i := range s.config.Medias {
		m := &s.config.Medias[i]
		if m.CurrentClip.IsAbsent() {
			continue
		}
		clip := m.CurrentClip.MustGet()

		nextRun, err := m.GetNextRun()
		nextRunStr := "N/A"
		if err == nil {
			nextRunStr = nextRun.Format("Mon Jan 2 15:04:05 MST")
		}

		tabs = append(tabs, HomePageTab{
			Name:      m.Name,
			NextRun:   nextRunStr,
			ClipName:  clip.Name,
			Guessable: m.Guessable,
		})
	}

	data := HomePageData{Tabs: tabs}
	if err := homeTmpl.Execute(w, data); err != nil {
		s.logger.Errorw("Error rendering template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) serveClip(w http.ResponseWriter, r *http.Request) {
	mediaNameEscaped := chi.URLParam(r, "name")
	mediaName, err := url.PathUnescape(mediaNameEscaped)
	if err != nil {
		http.Error(w, "Invalid media name", http.StatusBadRequest)
		return
	}

	for i := range s.config.Medias {
		m := &s.config.Medias[i]
		if m.Name == mediaName {
			if m.CurrentClip.IsAbsent() {
				http.Error(w, "No clip available", http.StatusNotFound)
				return
			}
			http.ServeFile(w, r, m.CurrentClip.MustGet().Path)
			return
		}
	}
	http.Error(w, "No matching media found", http.StatusNotFound)
}
