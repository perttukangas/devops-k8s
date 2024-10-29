/*
Copyright 2024.

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

package controller

import (
	"context"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stabledwkv1 "stable.dwk/api/v1"
)

// DummySiteReconciler reconciles a DummySite object
type DummySiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DummySite object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *DummySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the DummySite instance
	var dummySite stabledwkv1.DummySite
	if err := r.Get(ctx, req.NamespacedName, &dummySite); err != nil {
		logger.Error(err, "unable to fetch DummySite")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Read the website URL from the DummySite spec
	websiteURL := dummySite.Spec.WebsiteUrl

	// Fetch the content from the website URL
	resp, err := http.Get(websiteURL)
	if err != nil {
		logger.Error(err, "unable to fetch website content", "url", websiteURL)
		return ctrl.Result{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "unable to read website content", "url", websiteURL)
		return ctrl.Result{}, err
	}

	// Convert the body to a string
	htmlContent := string(body)

	// Log the HTML content (or you can process it further as needed)
	logger.Info("Fetched website content", "url", websiteURL, "content", htmlContent)

	// ??
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummySiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stabledwkv1.DummySite{}).
		Named("dummysite").
		Complete(r)
}
