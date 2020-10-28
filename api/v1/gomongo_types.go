/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GoMongoSpec defines the desired state of GoMongo
type GoMongoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of GoMongo. Edit GoMongo_types.go to remove/update
	DbSize           int32  `json:"db_size"`
	WebSize          int32  `json:"web_size"`
	MongoDbAdminPass string `json:"mongodb_admin_pass"`
}

// GoMongoStatus defines the observed state of GoMongo
type GoMongoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//	DbPodsNames  []string `json:"db_pods"`
	//	WebPodsNames []string `json:"web_pods"`
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GoMongo is the Schema for the gomongoes API
type GoMongo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoMongoSpec   `json:"spec,omitempty"`
	Status GoMongoStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GoMongoList contains a list of GoMongo
type GoMongoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoMongo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoMongo{}, &GoMongoList{})
}
