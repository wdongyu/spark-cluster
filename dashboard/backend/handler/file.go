package handler

import (
	"net/http"

	sparkclustercontroller "spark-cluster/pkg/controller/sparkcluster"

	hdfs "github.com/colinmarc/hdfs"
	log "github.com/sirupsen/logrus"
)

func (handler *APIHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	option := hdfs.ClientOptions{
		Addresses: []string{sparkclustercontroller.ShareServer + ":32083"},
		User:      "root"}

	client, err := hdfs.NewClient(option)
	if err != nil {
		log.Info(err)
	}

	//file, err := client.Open("/user/root/input/file1.txt")
	if err != nil {
		log.Info(err)
	}

	buf := make([]byte, 1048576)
	buf, err = client.ReadFile("/user/root/input/file1.txt")
	// length, err := file.Read(buf)
	// log.Info(string(length))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(buf, w, http.StatusOK)
	}

	// body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	// if err != nil {
	// 	responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	// }
	// defer r.Body.Close()
}
