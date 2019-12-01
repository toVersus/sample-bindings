/*
Copyright 2019 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package fake

import (
	"context"

	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	fake "knative.dev/sample-controller/pkg/client/injection/informers/factory/fake"
	addressableservice "knative.dev/sample-controller/pkg/client/injection/informers/samples/v1alpha1/addressableservice"
)

var Get = addressableservice.Get

func init() {
	injection.Fake.RegisterInformer(withInformer)
}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := fake.Get(ctx)
	inf := f.Samples().V1alpha1().AddressableServices()
	return context.WithValue(ctx, addressableservice.Key{}, inf), inf.Informer()
}
