package main

import (
	"context"
	"fmt"
	"net/http"

	fluxhttp "github.com/fluxcd/flux/pkg/http"
	fluxclient "github.com/fluxcd/flux/pkg/http/client"
	"github.com/fluxcd/flux/pkg/job"
	"github.com/fluxcd/flux/pkg/policy"
	"github.com/fluxcd/flux/pkg/resource"
	"github.com/fluxcd/flux/pkg/update"
	"k8s.io/client-go/kubernetes"
)

type Workload struct {
	Name      string
	Version   string
	Namespace string
	Author    string
}

type Client struct {
	api *fluxclient.Client
	k   *kubernetes.Clientset
}

func (w *Client) Unlock(fc *FluxCommand) error {

	d := Workload{
		Name:      fc.args.Name,
		Namespace: fc.args.Namespace,
	}

	_, err := w.lockUnlock(context.Background(), false, d, fc.args.User, fc.args.Namespace)
	if err != nil {
		return err
	}
	return err
}

func (w *Client) Lock(fc *FluxCommand) error {

	d := Workload{
		Name:      fc.args.Name,
		Namespace: fc.args.Namespace,
	}

	_, err := w.lockUnlock(context.Background(), true, d, fc.args.User, fc.args.Namespace)
	if err != nil {
		return err
	}
	return err
}

func (w *Client) lockUnlock(ctx context.Context, lock bool, d Workload, user string, ns string) (job.ID, error) {
	var msg string
	resourceID, err := resource.ParseIDOptionalNamespace(ns, fmt.Sprintf("deployment/%s", d.Name))
	if err != nil {
		return job.ID("no-id"), fmt.Errorf("cannot lock workload %w", err)
	}

	pol := policy.Set{}
	pol = pol.Add(policy.Locked)
	changes := resource.PolicyUpdate{}

	if lock {
		msg = fmt.Sprintf("Locking deployment/%s", d.Name)
		changes.Add = pol
	} else {
		msg = fmt.Sprintf("Unlocking deployment/%s", d.Name)
		changes.Remove = pol
	}

	return w.api.UpdateManifests(ctx, update.Spec{
		Type: update.Policy,
		Cause: update.Cause{
			Message: msg,
			User:    user,
		},
		Spec: resource.PolicyUpdates{
			resourceID: changes,
		},
	})
}

func New(fluxURL string, k *kubernetes.Clientset) *Client {
	return &Client{
		fluxclient.New(http.DefaultClient, fluxhttp.NewAPIRouter(), fluxURL, fluxclient.Token("")),
		k,
	}
}
