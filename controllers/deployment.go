package controllers

import (
	"context"

	cachev1alpha1 "github.com/amirvejahat/memcached-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func labels(v *cachev1alpha1.Memcached) map[string]string {
	return map[string]string{
		"app": "memcached",
	}
}

func (r *MemcachedReconciler) ensureDeployment(request reconcile.Request,
	instance *cachev1alpha1.Memcached,
	dep *appsv1.Deployment) (*reconcile.Result, error) {
	// see if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(context.TODO(), dep)

		if err != nil {
			// deployment failed
			return &reconcile.Result{}, err
		} else {
			// deployment was successful
			return nil, nil
		}

	} else if err != nil {
		// deployment might be exist
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *MemcachedReconciler) backendDeployment(v *cachev1alpha1.Memcached) *appsv1.Deployment {

	labels := labels(v)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name,
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &v.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            v.Name,
						Image:           "memcached",
						ImagePullPolicy: corev1.PullAlways,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "memcachedport",
						}},
					}},
				},
			},
		},
	}
	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}
