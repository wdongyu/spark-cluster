package handler

import (
	"net/http"
)

func (handler *APIHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	// req := handler.kubeClient.CoreV1().RESTClient().Post().
	// 	Resource("pods").
	// 	Name(sparkclustercontroller.MasterName).
	// 	Namespace(Namespace).
	// 	SubResource("exec").
	// 	Param("container", sparkclustercontroller.MasterName).
	// 	Param("stdin", "true").
	// 	Param("stdout", "true").
	// 	Param("stderr", "true").
	// 	Param("command", "/bin/bash")
	// //Param("tty", "true")
	// req.VersionedParams(
	// 	&v1.PodExecOptions{
	// 		Container: sparkclustercontroller.MasterName,
	// 		Command:   []string{"cat", "~/start.sh"},
	// 		Stdin:     true,
	// 		Stdout:    true,
	// 		Stderr:    true,
	// 		// TTY:       true,
	// 	},
	// 	scheme.ParameterCodec,
	// )

	// executor, err := remotecommand.NewSPDYExecutor(
	// 	handler.kubeConfig, http.MethodPost, req.URL(),
	// )
	// if err != nil {
	// 	responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	// }
	// body := executor.Stream(remotecommand.StreamOptions{
	// 	Stdin:  os.Stdin,
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stderr,
	// 	Tty:    true,
	// })
	// responseJSON(body, w, http.StatusOK)

	// body, err := req.DoRaw()
	// if err != nil {
	// 	log.Warningf("failed to cat log file: %v", err)
	// 	responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	// } else {
	// 	responseJSON(body, w, http.StatusOK)
	// }
	// fn := func() error {

	// 	executor, err := remotecommand.NewSPDYExecutor(
	// 		handler.kubeConfig, http.MethodPost, req.URL(),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Connect this process' terminal to the remote shell process.
	// 	return executor.Stream(remotecommand.StreamOptions{})
	// }

	// err := interrupt.Chain(nil, func() {
	// 	return
	// }).Run(fn)
}
