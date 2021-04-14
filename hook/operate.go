package hook

import (
	"encoding/json"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
)

//1.# console
//monitoring-platform-access
//2.# cml
//associatedCRP
//4.# warehouse
//istio-injection
//5.# compute
//istio-injection
//6.# implala
//istio-injection
//7.# monitoring
//cdp.cloudera/version

//3.# cml-user
//ds-parent-namespace

//only the # cml-user//ds-parent-namespace//
// cml-engine-resources admission is to add resource limit
func createAddResourcesContextPatch(pod corev1.Pod, availableAnnotations map[string]string, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)

	// update pod container resource
	replaceResourse := patchOperation{
		Op:    "replace",
		Path:  "/spec/containers/0/resources",
		Value: assembleResourceRequirements(pod),
	}
	glog.Infof("update  pod container resource for value: %v", replaceResourse)
	patch = append(patch, replaceResourse)

	return json.Marshal(patch)
}




func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

func assembleResourceRequirements(pod corev1.Pod) corev1.ResourceRequirements {
	res := corev1.ResourceRequirements{}
	// jiexun default requests have value, put the requests value to the limits value
	res.Limits = pod.Spec.Containers[0].Resources.Requests
	res.Requests = pod.Spec.Containers[0].Resources.Requests
	return res
}
