package dto

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// StackConfig contiene le informazioni necessarie per deployare uno stack
type StackConfig struct {
	ModuleDir       string
	EnvironmentsMap pulumi.StringMap
	Environment     Environments
	DefaultTags     pulumi.StringMap
	Action          PulumiAction
}

// PulumiAction rappresenta le principali operazioni eseguibili tramite la CLI di Pulumi.
type PulumiAction string

const (
	// Up (Update): Executes the preview, confirms changes, and applies modifications to the stack
	// to bring it to the state defined in the code. This is the most common action.
	Up PulumiAction = "up" // Alias: update

	// Destroy: Destroys all resources managed by the stack and removes the stack from the backend state.
	Destroy PulumiAction = "destroy"

	// Preview: Shows an execution plan of the changes that would be made by a 'pulumi up',
	// without applying anything.
	Preview PulumiAction = "preview"

	// Refresh: Synchronizes the local stack state with the real state of resources in the cloud
	// (e.g., AWS, Azure, etc.), without applying code modifications.
	Refresh PulumiAction = "refresh"

	// Import: Allows existing cloud resources to be taken and imported into a Pulumi project's management.
	Import PulumiAction = "import"

	// Convert: Converts a Pulumi project from one language to another (e.g., TypeScript to Go)
	// or from another IaC tool (e.g., CloudFormation or Terraform) to Pulumi.
	Convert PulumiAction = "convert"
)

type SubnetVisibility string

const (
	Public   SubnetVisibility = "Public"
	Private  SubnetVisibility = "Private"
	Isolated SubnetVisibility = "Isolated"
)
