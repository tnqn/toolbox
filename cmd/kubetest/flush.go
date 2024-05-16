package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func newFlushCommand(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flush",
		Short: "Flush resources created by kubetest",
	}
	flushNodeCommand := newFlushNodeCommand(o)
	cmd.AddCommand(flushNodeCommand)
	return cmd
}

func newFlushNodeCommand(o *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node",
		Aliases: []string{"nodes"},
		Short:   "Delete all nodes created by kubetest",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(o)
			if err != nil {
				return err
			}
			return flushNode(cmd, client)
		},
	}
	return cmd
}

func flushNode(cmd *cobra.Command, clientset kubernetes.Interface) error {
	err := clientset.CoreV1().Nodes().DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "app=kubetest",
	})
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Deleted all Nodes created by kubetest.")
	return nil
}
