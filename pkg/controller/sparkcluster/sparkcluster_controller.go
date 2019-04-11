/*
Copyright 2019 wdongyu.

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

package sparkcluster

import (
	"context"
	"fmt"
	"reflect"

	sparkv1alpha1 "spark-cluster/pkg/apis/spark/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SparkCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSparkCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sparkcluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to SparkCluster
	err = c.Watch(&source.Kind{Type: &sparkv1alpha1.SparkCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by SparkCluster - change this for objects you create

	// Watch for changes to Pod
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sparkv1alpha1.SparkCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sparkv1alpha1.SparkCluster{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to Pvc
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sparkv1alpha1.SparkCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSparkCluster{}

// ReconcileSparkCluster reconciles a SparkCluster object
type ReconcileSparkCluster struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SparkCluster object and makes changes based on the state read
// and what is in the SparkCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/exec,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=spark.k8s.io,resources=sparkclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=spark.k8s.io,resources=sparkclusters/status,verbs=get;update;patch
func (r *ReconcileSparkCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the SparkCluster instance
	instance := &sparkv1alpha1.SparkCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.ObjectMeta.DeletionTimestamp != nil {
		return reconcile.Result{}, nil
	}

	opts := &client.ListOptions{}
	opts.SetLabelSelector(fmt.Sprintf("app=%s", instance.Spec.ClusterPrefix+"-slave"))
	opts.InNamespace(request.NamespacedName.Namespace)
	podList := &corev1.PodList{}
	err = r.List(context.TODO(), opts, podList)
	if err != nil {
		return reconcile.Result{}, err
	}
	foundNum := len(podList.Items)
	expectNum := instance.Spec.SlaveNum

	// TODO(user): Change this to be the object type created by your controller

	// Master
	if instance.Spec.PvcEnable {
		pvcDeploy := r.getMasterPvc(instance)
		if err := r.checkPvc(instance, pvcDeploy); err != nil {
			return reconcile.Result{}, err
		}
	}
	podDeploy := r.getMasterPod(instance)
	if err := r.checkPod(instance, podDeploy); err != nil {
		return reconcile.Result{}, err
	}
	serviceDeploy := r.getMasterService(instance)
	if err := r.checkService(instance, serviceDeploy); err != nil {
		return reconcile.Result{}, err
	}
	serviceDeploy = r.getUIService(instance)
	if err := r.checkService(instance, serviceDeploy); err != nil {
		return reconcile.Result{}, err
	}

	// Slaves
	if foundNum > expectNum {
		for i := foundNum; i > expectNum; i-- {
			pod := &corev1.Pod{}
			err := r.Get(context.TODO(), types.NamespacedName{Name: slaveName(instance, i), Namespace: instance.Namespace}, pod)
			if err != nil && !errors.IsNotFound(err) {
				return reconcile.Result{}, err
			}
			log.Info("Deleting Pod", "namespace", instance.Namespace, "name", slaveName(instance, i))
			r.Delete(context.TODO(), pod)

			svc := &corev1.Service{}
			err = r.Get(context.TODO(), types.NamespacedName{Name: slaveName(instance, i), Namespace: instance.Namespace}, svc)
			if err != nil && !errors.IsNotFound(err) {
				return reconcile.Result{}, err
			}
			log.Info("Deleting Service", "namespace", instance.Namespace, "name", slaveName(instance, i))
			r.Delete(context.TODO(), svc)

			if instance.Spec.PvcEnable {
				pvc := &corev1.PersistentVolumeClaim{}
				err := r.Get(context.TODO(), types.NamespacedName{Name: SlavePvc + "-" + fmt.Sprintf("%d", i), Namespace: instance.Namespace}, pvc)
				if err != nil && !errors.IsNotFound(err) {
					return reconcile.Result{}, err
				}
				log.Info("Deleting Pvc", "namespace", instance.Namespace, "name", SlavePvc+"-"+fmt.Sprintf("%d", i))
				r.Delete(context.TODO(), pvc)
			}
		}
	} else {
		for i := 1; i <= expectNum; i++ {
			if instance.Spec.PvcEnable {
				pvcDeploy := r.getSlavePvc(instance, i)
				if err := r.checkPvc(instance, pvcDeploy); err != nil {
					return reconcile.Result{}, err
				}
			}
			podDeploy := r.getSlavePod(instance, i)
			if err := r.checkPod(instance, podDeploy); err != nil {
				return reconcile.Result{}, err
			}
			serviceDeploy := r.getSlaveService(instance, i)
			if err := r.checkService(instance, serviceDeploy); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, r.updateStatus(instance)
}

func (r *ReconcileSparkCluster) checkPvc(instance *sparkv1alpha1.SparkCluster, deploy *corev1.PersistentVolumeClaim) error {
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return err
	}

	found := &corev1.PersistentVolumeClaim{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Pvc", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Create(context.TODO(), deploy)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if !reflect.DeepEqual(deploy.Spec.Resources, found.Spec.Resources) {
			log.Info("Spec : " + deploy.Spec.String() + found.Spec.String())
			found.Spec = deploy.Spec
			err = r.Update(context.TODO(), found)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ReconcileSparkCluster) checkService(instance *sparkv1alpha1.SparkCluster, deploy *corev1.Service) error {
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Service already exists
	found := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Service", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Create(context.TODO(), deploy)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {

		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes

		// UI Service
		if deploy.Name == "hadoop-ui-service" {
			var deployPorts, foundPorts []int32
			for _, Ports := range deploy.Spec.Ports {
				deployPorts = append(deployPorts, Ports.Port)
			}
			for _, Ports := range found.Spec.Ports {
				foundPorts = append(foundPorts, Ports.Port)
			}
			if !reflect.DeepEqual(deployPorts, foundPorts) {
				// log.Info("Spec : " + deploy.Spec.String() + found.Spec.String())
				deploy.Spec.ClusterIP = found.Spec.ClusterIP
				found.Spec = deploy.Spec
				log.Info("Updating Service", "namespace", deploy.Namespace, "name", deploy.Name)
				err = r.Update(context.TODO(), found)
				if err != nil {
					return err
				}
			}
		}
		// else {
		// 	if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		// 		// log.Info("Spec : " + deploy.Spec.String() + found.Spec.String())
		// 		found.Spec = deploy.Spec
		// 		log.Info("Updating Service", "namespace", deploy.Namespace, "name", deploy.Name)
		// 		err = r.Update(context.TODO(), found)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }
	}
	return nil
}

func (r *ReconcileSparkCluster) checkPod(instance *sparkv1alpha1.SparkCluster, deploy *corev1.Pod) error {
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return err
	}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	found := &corev1.Pod{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Pod", "namespace", deploy.Namespace, "name", deploy.Name)
		err = r.Create(context.TODO(), deploy)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {

		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes
		// if !reflect.DeepEqual(deploy.Spec.Containers[0].Command, found.Spec.Containers[0].Command) ||
		sameResource := r.sameResources(deploy.Spec.Containers[0].Resources, found.Spec.Containers[0].Resources)
		sameVolumes := r.sameVolumesMounts(deploy.Spec.Containers[0].VolumeMounts, found.Spec.Containers[0].VolumeMounts)
		if !sameResource || (!sameVolumes) {
			// found.Spec = deploy.Spec

			log.Info("Updating Pod", "namespace", deploy.Namespace, "name", deploy.Name)
			// err = r.Update(context.TODO(), found)
			err1 := r.Delete(context.TODO(), found)
			//r.Create(context.TODO(), deploy)
			if err1 != nil {
				return err1
			}
		}
	}
	return nil
}

func (r *ReconcileSparkCluster) sameVolumesMounts(deploy []corev1.VolumeMount, found []corev1.VolumeMount) bool {
	l1 := len(deploy)
	l2 := len(found)
	if l2-l1 == 1 {
		for _, d := range deploy {
			in := false
			for _, f := range found {
				if d.Name == f.Name {
					in = true
				}
			}

			if !in {
				return false
			}
		}
		return true
	}
	return false

}

func (r *ReconcileSparkCluster) sameResources(deploy corev1.ResourceRequirements, found corev1.ResourceRequirements) bool {
	q1 := deploy.Requests[corev1.ResourceMemory]
	q2 := deploy.Requests[corev1.ResourceCPU]
	q3 := deploy.Limits[corev1.ResourceMemory]
	q4 := deploy.Limits[corev1.ResourceCPU]
	if q1.Cmp(found.Requests[corev1.ResourceMemory]) != 0 ||
		q2.Cmp(found.Requests[corev1.ResourceCPU]) != 0 ||
		q3.Cmp(found.Limits[corev1.ResourceMemory]) != 0 ||
		q4.Cmp(found.Limits[corev1.ResourceCPU]) != 0 {

		return false
	}
	return true
}
