package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"

	sparkclusterv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"

	"github.com/docker/docker/pkg/term"
	log "github.com/sirupsen/logrus"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/util/interrupt"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Terminal struct {
	conn     sockjs.Session
	sizeChan chan *remotecommand.TerminalSize

	Context   string
	Namespace string
	Pod       v1.Pod
	Container string

	kubeClient kubernetes.Interface
	kubeConfig *rest.Config
}

func NewTerminal(kubeClient kubernetes.Interface, kubeConfig *rest.Config) *Terminal {
	return &Terminal{
		sizeChan:   make(chan *remotecommand.TerminalSize),
		kubeClient: kubeClient,
		kubeConfig: kubeConfig,
	}
}

// Read handles pty->process messages (stdin, resize)
// Called in a loop from remotecommand as long as the process is running
func (t *Terminal) Read(p []byte) (int, error) {
	var reply string
	var msg map[string]uint16
	reply, err := t.conn.Recv()
	if err != nil {
		return 0, err
	}
	if err := json.Unmarshal([]byte(reply), &msg); err != nil {
		return copy(p, reply), nil
	} else {
		t.sizeChan <- &remotecommand.TerminalSize{
			msg["cols"],
			msg["rows"],
		}
		return 0, nil
	}
}

// Write handles process->pty stdout
// Called from remotecommand whenever there is any output
func (t *Terminal) Write(p []byte) (int, error) {
	err := t.conn.Send(string(p))
	return len(p), err
}

// Next handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t *Terminal) Next() *remotecommand.TerminalSize {
	size := <-t.sizeChan
	log.Debugf("terminal resize to width: %d, height: %d", size.Width, size.Height)
	return size
}

func (handler *APIHandler) CreateTerminal(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, ok := query["sparkcluster"]
	if !ok || len(name) == 0 {
		responseJSON(Message{"Missing required field \"sparkcluster\""}, w, http.StatusInternalServerError)
		return
	}
	namespace := Namespace

	// Get sparkcluster according to path parameter.
	sparkcluster := &sparkclusterv1alpha1.SparkCluster{}
	err := handler.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name[0]}, sparkcluster)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("sparkcluster \"%s\" not found", name)
		responseJSON(Message{err.Error()}, w, http.StatusNotFound)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	// Get pods which belong to this workspace
	// Use the first pod as terminal pod.
	var podList v1.PodList
	// labels := workspacecontroller.DefaultLabels(workspace)
	opts := &client.ListOptions{}
	opts.SetLabelSelector(fmt.Sprintf("app=%s", "hadoop-spark-master"))
	opts.InNamespace(Namespace)
	err = handler.client.List(context.TODO(), opts, &podList)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	} else if len(podList.Items) == 0 {
		err := fmt.Errorf("sparkcluster \"%s\" is not running", name)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	t := &Terminal{
		Namespace: Namespace,
		Pod:       podList.Items[0],
		Container: "hadoop-spark-master",
	}

	tmpl, err := template.ParseFiles(path.Join(handler.frontDir, "terminal.html"))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	tmpl.Execute(w, t)
}

func (t *Terminal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	podName, ok1 := query["pod"]
	container, ok2 := query["container"]
	if !ok1 || !ok2 || len(podName) == 0 || len(container) == 0 {
		responseJSON(Message{"missing required fields."}, w, http.StatusInternalServerError)
		return
	}

	if pod, err := t.kubeClient.CoreV1().Pods(Namespace).Get(podName[0], metav1.GetOptions{}); err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	} else {
		t.Pod = *pod
		t.Container = container[0]
	}

	sockjsHandler := func(session sockjs.Session) {
		t.conn = session
		if err := t.exec("/bin/bash"); err != nil {
			fmt.Println(err)
			responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
			return
		}
	}

	sockjs.NewHandler("/terminal/ws", sockjs.DefaultOptions, sockjsHandler).ServeHTTP(w, r)
}

func (t *Terminal) exec(command string) error {
	fn := func() error {
		req := t.kubeClient.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(t.Pod.Name).
			Namespace(t.Pod.Namespace).
			SubResource("exec").
			Param("container", t.Container).
			Param("stdin", "true").
			Param("stdout", "true").
			Param("stderr", "true").
			Param("command", command).Param("tty", "true")
		req.VersionedParams(
			&v1.PodExecOptions{
				Container: t.Container,
				Command:   []string{},
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				TTY:       true,
			},
			scheme.ParameterCodec,
		)

		executor, err := remotecommand.NewSPDYExecutor(
			t.kubeConfig, http.MethodPost, req.URL(),
		)
		if err != nil {
			return err
		}

		// Connect this process' terminal to the remote shell process.
		return executor.Stream(remotecommand.StreamOptions{
			Stdin:             t,
			Stdout:            t,
			Stderr:            t,
			Tty:               true,
			TerminalSizeQueue: t,
		})
	}

	inFd, _ := term.GetFdInfo(t.conn)
	state, err := term.SaveState(inFd)
	if err != nil {
		return err
	}
	return interrupt.Chain(nil, func() {
		term.RestoreTerminal(inFd, state)
	}).Run(fn)
}
