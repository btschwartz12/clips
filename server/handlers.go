package server

import (
	"html/template"
	"net/http"

	"github.com/btschwartz12/clips/server/assets"
)

var (
	homeTmpl = template.Must(template.ParseFS(
		assets.Templates,
		"templates/home.html.tmpl",
	))
)

type HomePageData struct {
	ClipName string
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	data := HomePageData{
		ClipName: s.currentClip.Name,
	}
	err := homeTmpl.Execute(w, data)
	if err != nil {
		s.logger.Errorw("Error rendering template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	s.logger.Infow("homepage served", "clip", s.currentClip.Name)
}

func (s *Server) serveClip(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.currentClip.Path)
}

