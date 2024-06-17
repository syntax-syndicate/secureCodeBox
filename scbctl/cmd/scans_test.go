// SPDX-FileCopyrightText: the secureCodeBox authors
//
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"errors"
	"testing"

	v1 "github.com/secureCodeBox/secureCodeBox/operator/apis/execution/v1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type MockClientProvider struct {
	client.Client
	namespace string
	err       error
}

func (m *MockClientProvider) GetClient(_ *genericclioptions.ConfigFlags) (client.Client, string, error) {
	return m.Client, m.namespace, m.err
}

func TestScanCommand(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(v1.AddToScheme(scheme))

	testcases := []struct {
		name          string
		args          []string
		expectedError error
	}{
		{
			name:          "Valid entry ",
			args:          []string{"nmap", "--", "scanme.nmap.org"},
			expectedError: nil,
		},
		{
			name:          "Valid entry with namespace",
			args:          []string{"nmap", "--", "scanme.nmap.org"},
			expectedError: nil,
		},
		{
			name:          "No scan parameters provided",
			args:          []string{"nmap"},
			expectedError: errors.New("You must use '--' to separate scan parameters"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			client := fake.NewClientBuilder().WithScheme(scheme).Build()
			clientProvider = &MockClientProvider{
				Client:    client,
				namespace: "default",
				err:       nil,
			}

			cmd := &cobra.Command{
				Use:     ScanCmd.Use,
				Short:   ScanCmd.Short,
				Long:    ScanCmd.Long,
				Example: ScanCmd.Example,
				RunE:    ScanCmd.RunE,
			}

			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, but got: %v", err)
			}

		})
	}
}
