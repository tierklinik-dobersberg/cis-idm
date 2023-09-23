package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/pprof"
	"time"
)

// handleCPUProfile returns the CPU profile.
func CPUProfileHandler(w http.ResponseWriter, req *http.Request) {
	// Parse duration.
	duration := 10 * time.Second
	if durationOption := req.URL.Query().Get("duration"); durationOption != "" {
		parsedDuration, err := time.ParseDuration(durationOption)
		if err != nil {
			http.Error(w, "invalid duration", http.StatusBadRequest)

			return
		}

		duration = parsedDuration
	}

	// Indicate download and filename.
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="cpuprofile.pprof"`),
	)

	// Start CPU profiling.
	buf := new(bytes.Buffer)
	if err := pprof.StartCPUProfile(buf); err != nil {
		http.Error(w, "failed to start CPU profile", http.StatusInternalServerError)

		return
	}

	// Wait for the specified duration.
	select {
	case <-time.After(duration):
	case <-req.Context().Done():
	}

	// Stop CPU profiling and return data.
	pprof.StopCPUProfile()

	w.WriteHeader(http.StatusOK)

	w.Write(buf.Bytes())
}
