package controller

import (
	"context"

	demov1alpha1 "github.com/jamaan01/managed-database-operator/api/v1alpha1"
	"github.com/jamaan01/managed-database-operator/internal/provisioner"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const managedDatabaseFinalizer = "demo.example.com/finalizer"

type ManagedDatabaseReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Provisioner *provisioner.Client
}

// +kubebuilder:rbac:groups=demo.example.com,resources=manageddatabases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=demo.example.com,resources=manageddatabases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=demo.example.com,resources=manageddatabases/finalizers,verbs=update

func (r *ManagedDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	var managedDatabase demov1alpha1.ManagedDatabase

	if err := r.Get(ctx, req.NamespacedName, &managedDatabase); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	if managedDatabase.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&managedDatabase, managedDatabaseFinalizer) {
			controllerutil.AddFinalizer(&managedDatabase, managedDatabaseFinalizer)

			if err := r.Update(ctx, &managedDatabase); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
	}

	if !managedDatabase.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(&managedDatabase, managedDatabaseFinalizer) {
			if managedDatabase.Status.ExternalID != "" {
				if err := r.Provisioner.Delete(ctx, managedDatabase.Status.ExternalID); err != nil {
					return ctrl.Result{}, err
				}
			}

			controllerutil.RemoveFinalizer(&managedDatabase, managedDatabaseFinalizer)

			if err := r.Update(ctx, &managedDatabase); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

	}

	return ctrl.Result{}, nil
}

func (r *ManagedDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1alpha1.ManagedDatabase{}).
		Named("manageddatabase").
		Complete(r)
}
