package entity

import (
	"bytes"
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

var (
	pods = `
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
`
)

func TestWorkloadsTestWorkloads(t *testing.T) {
	var p v1.Pod
	_ = yaml.Unmarshal(bytes.NewBufferString(pods).Bytes(), &p)
	fmt.Printf("******** %+v\n", p)

}
