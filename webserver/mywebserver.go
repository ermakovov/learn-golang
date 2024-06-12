package webserver

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

var courses = map[int64]string{
	1: "First course",
	2: "Second course",
}

func StartWebserver() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/courses/description", CoursesDescHandler)

	port := "80"
	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting server on port")
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go to /courses/description"))
}

func CoursesDescHandler(w http.ResponseWriter, r *http.Request) {
	courseIdParam := r.URL.Query().Get("course_id")

	courseId, err := strconv.ParseInt(courseIdParam, 10, 64)
	if err != nil {
		logrus.WithError(err).Error("courseIdParam parsing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	course, ok := courses[courseId]
	if !ok {
		logrus.Info("Course doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(course))
}
