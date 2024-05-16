package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

type createOptions struct {
	// The path of kubeconfig configuration file.
	count int
	// Use the object as the template of objects to create.
	source string
}

func newCreateOptions() *createOptions {
	return &createOptions{
		count: 1,
	}
}

func newCreateCommand(o *options) *cobra.Command {
	co := newCreateOptions()
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create fake resources using kubetest",
	}
	createNodeCommand := newCreateNodeCommand(o, co)
	cmd.AddCommand(createNodeCommand)

	cmd.PersistentFlags().IntVarP(&co.count, "count", "c", co.count, "Number of objects to create")
	cmd.PersistentFlags().StringVarP(&co.source, "source", "s", co.source, "The source of objects to copy from")
	return cmd
}

func newCreateNodeCommand(o *options, co *createOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node",
		Aliases: []string{"nodes"},
		Short:   "Create fake nodes using kubetest",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(o)
			if err != nil {
				return err
			}
			return createNode(cmd, client, co.count, co.source)
		},
	}
	return cmd
}

func createNode(cmd *cobra.Command, clientset kubernetes.Interface, count int, source string) error {
	var nodeTemplate *corev1.Node
	if source != "" {
		node, err := clientset.CoreV1().Nodes().Get(context.TODO(), source, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get the provided Node: %w", err)
		}
		nodeTemplate = &corev1.Node{Spec: node.Spec, Status: node.Status}
	} else {
		nodeTemplate = &corev1.Node{
			Spec: corev1.NodeSpec{
				PodCIDR: "10.0.0.0/24",
				PodCIDRs: []string{
					"10.0.0.0/24",
				},
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:               corev1.NodeMemoryPressure,
						Status:             corev1.ConditionFalse,
						Reason:             "KubeletHasSufficientMemory",
						Message:            "kubelet has sufficient memory available",
						LastHeartbeatTime:  metav1.Now(),
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               corev1.NodeDiskPressure,
						Status:             corev1.ConditionFalse,
						Reason:             "KubeletHasNoDiskPressure",
						Message:            "kubelet has no disk pressure",
						LastHeartbeatTime:  metav1.Now(),
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               corev1.NodePIDPressure,
						Status:             corev1.ConditionFalse,
						Reason:             "KubeletHasSufficientPID",
						Message:            "kubelet has sufficient PID available",
						LastHeartbeatTime:  metav1.Now(),
						LastTransitionTime: metav1.Now(),
					},
					{
						Type:               corev1.NodeReady,
						Status:             corev1.ConditionTrue,
						Reason:             "KubeletReady",
						Message:            "kubelet is posting ready status. AppArmor enabled",
						LastHeartbeatTime:  metav1.Now(),
						LastTransitionTime: metav1.Now(),
					},
				},
				Images: []corev1.ContainerImage{
					{
						Names: []string{
							"golang@sha256:f43c6f049f04cbbaeb28f0aad3eea15274a7d0a7899a617d0037aec48d7ab010",
							"golang:latest",
						},
						SizeBytes: 822340085,
					},
					{
						Names: []string{
							"registry.k8s.io/etcd@sha256:dd75ec974b0a2a6f6bb47001ba09207976e625db898d1b16735528c009cb171c",
							"registry.k8s.io/etcd:3.5.6-0",
						},
						SizeBytes: 299475478,
					},
					{
						Names: []string{
							"registry.k8s.io/kube-proxy@sha256:a9f441a6b440c634ccfe62530ab1c7ff2ea7ed3f577f91f6a71c7e2f51256410",
							"registry.k8s.io/kube-proxy:v1.26.15",
						},
						SizeBytes: 72051242,
					},
					{
						Names: []string{
							"registry.k8s.io/kube-scheduler@sha256:6447dce5ea569c857b161436235292bc30280b3f83fda5df730b23b0812336dc",
							"registry.k8s.io/kube-scheduler:v1.26.15",
						},
						SizeBytes: 56870145,
					},
					{
						Names: []string{
							"registry.k8s.io/pause@sha256:3d380ca8864549e74af4b29c10f9cb0956236dfb01c40ca076fb6c37253234db",
							"registry.k8s.io/pause:3.6",
						},
						SizeBytes: 682696,
					},
					{
						Names: []string{
							"registry.k8s.io/kube-apiserver@sha256:0dc6d5ba5863218a391de0952d27701d2715254b0fbfb3670cadd3074b057f8f",
							"registry.k8s.io/kube-apiserver:v1.26.15",
						},
						SizeBytes: 138240930,
					},
					{
						Names: []string{
							"registry.k8s.io/kube-controller-manager@sha256:ea4dd4c0905132110aca01e638d87f861dfa9db229c7022c583f4076ade0c23a",
							"registry.k8s.io/kube-controller-manager:v1.26.15",
						},
						SizeBytes: 127210753,
					},
					{
						Names: []string{
							"registry.k8s.io/e2e-test-images/busybox@sha256:2e0f836850e09b8b7cc937681d6194537a09fbd5f6b9e08f4d646a85128e8937",
							"registry.k8s.io/e2e-test-images/busybox:1.29-4",
						},
						SizeBytes: 731990,
					},
					{
						Names: []string{
							"registry.k8s.io/coredns/coredns@sha256:8e352a029d304ca7431c6507b56800636c321cb52289686a581ab70aaa8a2e2a",
							"registry.k8s.io/coredns/coredns:v1.9.3",
						},
						SizeBytes: 14837849,
					},
				},
			},
		}
	}
	nodeTemplate.Labels = map[string]string{
		"app": "kubetest",
	}
	nodeTemplate.Spec.Unschedulable = true
	nodeTemplate.Spec.Taints = []corev1.Taint{
		{
			Key:       corev1.TaintNodeUnreachable,
			Effect:    corev1.TaintEffectNoSchedule,
			TimeAdded: ptr.To(metav1.Now()),
		},
		{
			Key:       corev1.TaintNodeUnreachable,
			Effect:    corev1.TaintEffectNoExecute,
			TimeAdded: ptr.To(metav1.Now()),
		},
	}

	for i := 0; i < count; i++ {
		node := nodeTemplate.DeepCopy()
		node.Name = fmt.Sprintf("%s-%s-%s", fakeObjectPrefix, "node", rand.String(6))

		node, err := clientset.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		klog.V(2).InfoS("Created Node", "Node", klog.KObj(node))
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Created %d Nodes.\n", count)
	return nil
}
