package installer

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/fake"
	clitest "k8s.io/client-go/testing"
)

func Test_getGVRWithObject(t *testing.T) {
	type args struct {
		object *unstructured.Unstructured
		gvr    schema.GroupVersionResource
		d      discovery.DiscoveryInterface
	}
	tests := []struct {
		name    string
		args    args
		want    *schema.GroupVersionResource
		wantErr bool
	}{
		{
			name: "discovery and object api version is same",
			args: args{
				object: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "tekton.dev/v1",
					},
				},
				gvr: schema.GroupVersionResource{
					Group:    tektonGroup,
					Resource: "Task",
				},
				d: &fake.FakeDiscovery{
					Fake: &clitest.Fake{
						Resources: []*metav1.APIResourceList{
							{
								GroupVersion: "tekton.dev/v1",
								APIResources: []metav1.APIResource{
									{
										Name:    "tasks",
										Kind:    "Task",
										Group:   tektonGroup,
										Version: "v1",
									},
								},
							},
						},
					},
				},
			},
			want: &schema.GroupVersionResource{
				Group:    tektonGroup,
				Resource: "tasks",
				Version:  "v1",
			},
		},
		{
			name: "discovery and object api version is different",
			args: args{
				object: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "tekton.dev/v1beta1",
					},
				},
				gvr: schema.GroupVersionResource{
					Group:    tektonGroup,
					Resource: "Task",
				},
				d: &fake.FakeDiscovery{
					Fake: &clitest.Fake{
						Resources: []*metav1.APIResourceList{
							{
								GroupVersion: "tekton.dev/v1",
								APIResources: []metav1.APIResource{
									{
										Name:    "tasks",
										Kind:    "Task",
										Group:   tektonGroup,
										Version: "v1",
									},
								},
							},
						},
					},
				},
			},
			want: &schema.GroupVersionResource{
				Group:    tektonGroup,
				Resource: "tasks",
				Version:  "v1beta1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGVRWithObject(tt.args.object, tt.args.gvr, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGVRWithObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Group != tt.want.Group {
				t.Errorf("getGVRWithObject() got = %v, want %v", got, tt.want)
			}
			if got.Resource != tt.want.Resource {
				t.Errorf("getGVRWithObject() got = %v, want %v", got, tt.want)
			}
			if got.Version != tt.want.Version {
				t.Errorf("getGVRWithObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}
