package handler

import (
	"net/http"
	"os"

	jobcontroller "spark-cluster/pkg/controller/job"
	sparkclustercontroller "spark-cluster/pkg/controller/sparkcluster"

	"k8s.io/client-go/tools/remotecommand"
)

func (handler *APIHandler) GetDroneLog(w http.ResponseWriter, r *http.Request) {
	// config := new(oauth2.Config)
	// auther := config.Client(
	// 	oauth2.NoContext,
	// 	&oauth2.Token{
	// 		AccessToken: token,
	// 	},
	// )

	// // create the drone client with authenticator
	// client := drone.NewClient(host, auther)

	// // gets the named repository information
	// lines, err := client.Logs("wdongyu", "SparkHive", 38, 1, 2)
	// body := []string{}
	// for _, line := range lines {
	// 	body = append(body, line.Message)
	// }
	// fmt.Println(body, err)
	// if err == nil {
	// 	responseJSON(body, w, http.StatusOK)
	// }

}

func (handler *APIHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	req := handler.kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name("user-master").
		Namespace(handler.resourcesNamespace).
		SubResource("exec").
		Param("container", sparkclustercontroller.Master).
		//Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("command", "/bin/cat").
		Param("command", "/root/start.sh")

	executor, err := remotecommand.NewSPDYExecutor(
		handler.kubeConfig, http.MethodPost, req.URL(),
	)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	os.Stdout.Sync()
	os.Stderr.Sync()

	body := executor.Stream(remotecommand.StreamOptions{
		Stdin:             nil,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               false,
		TerminalSizeQueue: nil,
	})
	responseJSON(body, w, http.StatusOK)
}

func (handler *APIHandler) GetPodLog(w http.ResponseWriter, r *http.Request) {
	controller, err := jobcontroller.New()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}

	err = controller.Execute()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
}
