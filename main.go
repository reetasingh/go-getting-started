package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/pat"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Cluster struct {
	kubeClient kubernetes.Interface
}

type server struct {
	cluster Cluster
}

func newServer(clientset kubernetes.Interface) *server {
	s := new(server)
	s.cluster = Cluster{clientset}
	return s
}

func main() {
	fmt.Println("Starting hello-world server...")
	clientset, _ := newKubeClient("")
	s := newServer(clientset)
	mux := pat.New()
	mux.Get("/pods/count", s.podsCount)
	mux.Get("/pods/list/{sort}", s.podsList)
	mux.Get("/", s.defaultHandler)
	http.Handle("/", mux)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func (s *server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world from Okteto!")
}

func (s *server) podsList(w http.ResponseWriter, r *http.Request) {
	sortParam := r.URL.Query().Get(":sort")
	var sortParamValue string
	if sortParam == "" {
		sortParamValue = "name"
	} else {
		sortParamValue = strings.ToLower(sortParam)
	}
	if sortParamValue != "name" && sortParamValue != "age" && sortParamValue != "restarts" {
		w.WriteHeader(500)
		fmt.Fprint(w, "invalid value for sortBy param, should be name, age or restarts. default is name")
		return
	}
	result, err := s.cluster.getSortedPods(sortParamValue)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, string(b))
}

func (s *server) podsCount(w http.ResponseWriter, r *http.Request) {
	result, err := s.cluster.getPodsCount()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, result)
}

type Pods struct {
	pList  []v1.Pod
	sortBy string
}

func (p *Pods) Len() int {
	return len(p.pList)
}

func (p *Pods) Swap(i, j int) {
	p.pList[i], p.pList[j] = p.pList[j], p.pList[i]
}

func getPodRestartCount(pod v1.Pod) int {
	restartCount := 0
	containerStatus := pod.Status.ContainerStatuses
	for _, c := range containerStatus {
		restartCount = restartCount + int(c.RestartCount)
	}

	return restartCount
}

func (p *Pods) Less(i, j int) bool {
	p.pList[i], p.pList[j] = p.pList[j], p.pList[i]
	if p.sortBy == "name" {
		return p.pList[i].Name < p.pList[j].Name
	} else if p.sortBy == "age" {
		t1 := p.pList[i].Status.StartTime
		t2 := p.pList[j].Status.StartTime
		return t1.Before(t2)
	} else if p.sortBy == "restarts" {
		return getPodRestartCount(p.pList[i]) < getPodRestartCount(p.pList[j])
	}
	return false
}

func (s *Cluster) getPodsCount() (int, error) {
	pods, err := s.getPodList()
	if err != nil {
		return 0, err
	}
	return len(pods.Items), nil
}

func (s *Cluster) getPodList() (*v1.PodList, error) {
	// TODO: find a way to get namespac
	pods, err := s.kubeClient.CoreV1().Pods("reetasingh").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (s *Cluster) getSortedPods(sortBy string) ([]v1.Pod, error) {
	// TODO: find a way to get namespac
	pods, err := s.getPodList()
	if err != nil {
		return nil, err
	}
	podList := pods.Items
	p := &Pods{podList, sortBy}
	sort.Sort(p)
	return p.pList, nil
}

func newKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientsets
	clientset, err := kubernetes.NewForConfig(config)

	return clientset, err
}
