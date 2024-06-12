package webserver

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

var courses = map[int64]string{
	1: "First course",
	2: "Second course",
}

func StartWebserver() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/courses/description", CoursesDescHandler)

	port := "80"
	log.WithFields(log.Fields{
		"port": port,
	}).Info("Starting server on port")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go to /courses/description"))
}

func CoursesDescHandler(w http.ResponseWriter, r *http.Request) {
	courseIdParam := r.URL.Query().Get("course_id")

	courseId, err := strconv.ParseInt(courseIdParam, 10, 64)
	if err != nil {
		log.WithError(err).Error("courseIdParam parsing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	course, ok := courses[courseId]
	if !ok {
		log.Info("Course doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(course))
}
