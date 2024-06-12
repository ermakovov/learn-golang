package webserver

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func StartMathWebserver() {
	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, ".log")

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go to /sum"))
	})

	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		xParam := r.URL.Query().Get("x")
		xArg, err := strconv.Atoi(xParam)
		if err != nil {
			logger.WithField("x", xParam).Error("query param parsing")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		yParam := r.URL.Query().Get("y")
		yArg, err := strconv.Atoi(yParam)
		if err != nil {
			logger.WithField("y", yParam).Error("query param parsing")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sum := xArg + yArg
		if abs(xArg) > math.MaxInt64-abs(yArg) {
			logger.WithFields(logrus.Fields{
				"x": xArg,
				"y": yArg,
			}).Warning("Sum overflows int")

			sum = -1
		}

		w.Write([]byte(strconv.Itoa(sum)))
	})

	port := "80"
	logWithPort := logger.WithFields(logrus.Fields{
		"port": port,
	})

	logWithPort.Info("Starting webserver on port")
	logWithPort.Fatal(http.ListenAndServe(":"+port, nil))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
