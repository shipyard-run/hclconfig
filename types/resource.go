package types

// Processable defines an optional interface that allows a resource to define a callback
// that is executed when the resources is processed by the DAG.
type Processable interface {
	// Process is called by the parser when the DAG is resolved
	Process() error
}

// Resource is an interface that all
type Resource interface {
	// return the resource Metadata
	Metadata() *ResourceMetadata
}

// ResourceMetadata is the embedded type for any config resources
// it deinfes common meta data that all resources share
type ResourceMetadata struct {
	// Name is the name of the resource
	Name string `json:"name"`

	// Type is the type of resource, this is the text representation of the golang type
	Type string `json:"type"`

	// Module is the name of the module if a resource has been loaded from a module
	Module string `json:"module,omitempty"`

	// Linked resources which must be set before this config can be processed
	ResourceLinks []string `json:"resource_links,omitempty"`

	// DependsOn is a user configurable list of dependencies for this resource
	DependsOn []string `hcl:"depends_on,optional" json:"depends_on,omitempty"`

	// Enabled determines if a resource is enabled and should be processed
	Disabled bool `hcl:"disabled,optional" json:"disabled,omitempty"`
}

func (r *ResourceMetadata) Metadata() *ResourceMetadata {
	return r
}
