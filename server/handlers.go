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
	NextRun  string
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if s.currentClip == nil {
		http.Error(w, "No clip available", http.StatusNotFound)
		return
	}
	t, err := s.job.NextRun()
	if err != nil {
		s.logger.Errorw("error getting next run time", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	data := HomePageData{
		ClipName: s.currentClip.Name,
		NextRun:  t.Format("Mon Jan 2 15:04:05 EST"),
	}
	err = homeTmpl.Execute(w, data)
	if err != nil {
		s.logger.Errorw("Error rendering template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	s.logger.Infow("homepage served", "clip", s.currentClip.Name)
}

func (s *Server) serveClip(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.currentClip.Path)
}
