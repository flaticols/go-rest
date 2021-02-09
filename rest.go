package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"net/http"

	"github.com/pkg/errors"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

// RenderJSON sends data as JSON
func RenderJSON(w http.ResponseWriter, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	enc.SetEscapeHTML(true)

	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(buf.Bytes())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//RenderJSONFromBytes sends binary data as JSON
func RenderJSONFromBytes(w http.ResponseWriter, r *http.Request, data []byte) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(data); err != nil {
		return errors.Wrapf(err, "failed to send response to %s", r.RemoteAddr)
	}
	return nil
}

//SPAHandler handle all path to render SPA. Use it before API routes!
func SPAHandler(w http.ResponseWriter, r *http.Request, public string, static string) {
	p := strings.Replace(r.URL.Path, public, "", 1)
	p = filepath.Join(static, filepath.Clean(p))
	log.Printf("File: %s", p)

	if info, err := os.Stat(p); err != nil {
		log.Printf("[ERROR] get path error: %s", err.Error())
		index := filepath.Join(static, "index.html")
		http.ServeFile(w, r, index)
		return
	} else if info.IsDir() {
		index := filepath.Join(p, "index.html")
		http.ServeFile(w, r, index)
	} else {
		fp := filepath.Join(static, info.Name())
		log.Printf("[DBG] serve file: %s", fp)
		http.ServeFile(w, r, fp)
	}
}