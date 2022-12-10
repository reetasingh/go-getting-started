package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func getPodList() *v1.PodList {
	p := v1.PodList{}
	pod1 := v1.Pod{}
	pod1.Name = "a"
	pod1.Namespace = "reetasingh"
	var csPod1 []v1.ContainerStatus
	csPod1 = append(csPod1, v1.ContainerStatus{RestartCount: 0})
	pod1.Status.ContainerStatuses = csPod1

	pod2 := v1.Pod{}
	pod2.Name = "b"
	pod2.Namespace = "reetasingh"
	var csPod2 []v1.ContainerStatus
	csPod2 = append(csPod2, v1.ContainerStatus{RestartCount: 5})
	pod2.Status.ContainerStatuses = csPod2

	pod3 := v1.Pod{}
	pod3.Name = "c"
	pod3.Namespace = "reetasingh"
	var csPod3 []v1.ContainerStatus
	csPod3 = append(csPod3, v1.ContainerStatus{RestartCount: 3})
	pod3.Status.ContainerStatuses = csPod3

	// pod2 created first
	pod2.Status.StartTime = &metav1.Time{time.Now()}
	// pod 3 created after pod2
	pod3.Status.StartTime = &metav1.Time{time.Now().Add(time.Second * 100)}
	// pod1 created in the end
	pod1.Status.StartTime = &metav1.Time{time.Now().Add(time.Second * 200)}
	p.Items = []v1.Pod{pod1, pod2, pod3}
	return &p
}

func Test_cluster_getPodList(t *testing.T) {
	podList := getPodList()
	clienset := testclient.NewSimpleClientset(podList)
	s := Cluster{clienset}
	p, err := s.getPodList()
	assert.Equal(t, len(p.Items), 3)
	assert.NoError(t, err)

}

func Test_cluster_podsCount(t *testing.T) {
	podList := getPodList()
	clienset := testclient.NewSimpleClientset(podList)
	s := Cluster{clienset}
	result, err := s.getPodsCount()
	assert.Equal(t, result, 3)
	assert.NoError(t, err)
}

func Test_cluster_getSortedPods_By(t *testing.T) {
	podList := getPodList()
	clienset := testclient.NewSimpleClientset(podList)
	s := Cluster{clienset}

	// by name
	result, err := s.getSortedPods("name")
	assert.Equal(t, len(result), 3)
	assert.Equal(t, result[0].Name, "a")
	assert.Equal(t, result[1].Name, "b")
	assert.Equal(t, result[2].Name, "c")
	assert.NoError(t, err)

	// by age
	result, err = s.getSortedPods("age")
	assert.Equal(t, len(result), 3)
	assert.Equal(t, result[0].Name, "b")
	assert.Equal(t, result[1].Name, "c")
	assert.Equal(t, result[2].Name, "a")
	assert.NoError(t, err)

	// by restart count
	result, err = s.getSortedPods("restarts")
	assert.Equal(t, len(result), 3)
	assert.Equal(t, result[0].Name, "a")
	assert.Equal(t, result[1].Name, "c")
	assert.Equal(t, result[2].Name, "b")
	assert.NoError(t, err)
}
