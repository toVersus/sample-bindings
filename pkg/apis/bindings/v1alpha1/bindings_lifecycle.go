/*
Copyright 2019 The Knative Authors.

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

package v1alpha1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/tracker"

	"github.com/toVersus/sample-bindings/pkg/database"
)

const (
	// DatabaseBindingConditionReady is set when the binding has been applied to the subjects.
	DatabaseBindingConditionReady = apis.ConditionReady
)

var dbCondSet = apis.NewLivingConditionSet()

// GetGroupVersionKind implements kmeta.OwnerRefable
func (db *DatabaseBinding) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("DatabaseBinding")
}

func (db *DatabaseBinding) GetSubject() tracker.Reference {
	return db.Spec.Subject
}

func (dbs *DatabaseBindingStatus) InitializeConditions() {
	dbCondSet.Manage(dbs).InitializeConditions()
}

func (dbs *DatabaseBindingStatus) MarkBindingUnavailable(reason, message string) {
	dbCondSet.Manage(dbs).MarkFalse(
		DatabaseBindingConditionReady, reason, message)
}

func (dbs *DatabaseBindingStatus) MarkServiceAvailable() {
	dbCondSet.Manage(dbs).MarkTrue(DatabaseBindingConditionReady)
}

func (db *DatabaseBinding) Do(ctx context.Context, ps *duckv1.WithPod) {
	// First undo so that we can just unconditionally append below.
	db.Undo(ctx, ps)

	// Make sure the PodSpec has a EnvFrom like this:
	envFrom := corev1.EnvFromSource{
		SecretRef: &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: db.Spec.Secret.Name,
			},
		},
	}
	spec := ps.Spec.Template.Spec
	for i := range spec.Containers {
		spec.Containers[i].EnvFrom = append(spec.Containers[i].EnvFrom, envFrom)
	}
}

func (db *DatabaseBinding) Undo(ctx context.Context, ps *duckv1.WithPod) {
	spec := ps.Spec.Template.Spec

	// Make sure that none of the containers have the database secret ref
	for i, c := range spec.Containers {
		for j, ef := range c.EnvFrom {
			if ef.SecretRef.Name == database.DBPasswordEnvKey {
				spec.Containers[i].EnvFrom = append(spec.Containers[i].EnvFrom[:j], spec.Containers[i].EnvFrom[j+1:]...)
				break
			}
		}
	}
}
