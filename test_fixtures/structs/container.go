package structs

import "github.com/shipyard-run/hclconfig/types"

// TypeContainer is the resource string for a Container resource
const TypeContainer types.ResourceType = "container"

// Container defines a structure for creating Docker containers
type Container struct {
	// embedded type holding name, etc
	types.ResourceInfo `hcl:",remain" mapstructure:",squash"`

	Depends []string `hcl:"depends_on,optional" json:"depends,omitempty"`

	Networks []NetworkAttachment `hcl:"network,block" json:"networks,omitempty"` // Attach to the correct network // only when Image is specified

	Build      *Build            `hcl:"build,block" json:"build"`                             // Enables containers to be built on the fly
	Entrypoint []string          `hcl:"entrypoint,optional" json:"entrypoint,omitempty"`      // entrypoint to use when starting the container
	Command    []string          `hcl:"command,optional" json:"command,omitempty"`            // command to use when starting the container
	Env        map[string]string `hcl:"env,optional" json:"env,omitempty" mapstructure:"env"` // environment variables to set when starting the container
	Volumes    []Volume          `hcl:"volume,block" json:"volumes,omitempty"`                // volumes to attach to the container
	DNS        []string          `hcl:"dns,optional" json:"dns,omitempty"`                    // Add custom DNS servers to the container

	Privileged bool `hcl:"privileged,optional" json:"privileged,omitempty"` // run the container in privileged mode?

	// resource constraints
	Resources *Resources `hcl:"resources,block" json:"resources,omitempty"` // resource constraints for the container

	MaxRestartCount int `hcl:"max_restart_count,optional" json:"max_restart_count,omitempty" mapstructure:"max_restart_count"`

	// User block for mapping the user id and group id inside the container
	RunAs *User `hcl:"run_as,block" json:"run_as,omitempty" mapstructure:"run_as"`
}

type User struct {
	// Username or UserID of the user to run the container as
	User string `hcl:"user" json:"user,omitempty" mapstructure:"user"`
	// Groupname GroupID of the user to run the container as
	Group string `hcl:"group" json:"group,omitempty" mapstructure:"group"`
}

type NetworkAttachment struct {
	Name      string   `hcl:"name" json:"name"`
	IPAddress string   `hcl:"ip_address,optional" json:"ip_address,omitempty" mapstructure:"ip_address"`
	Aliases   []string `hcl:"aliases,optional" json:"aliases,omitempty"` // Network aliases for the resource
}

// Resources allows the setting of resource constraints for the Container
type Resources struct {
	CPU    int   `hcl:"cpu,optional" json:"cpu,omitempty"`                                // cpu limit for the container where 1 CPU = 1000
	CPUPin []int `hcl:"cpu_pin,optional" json:"cpu_pin,omitempty" mapstructure:"cpu_pin"` // pin the container to one or more cpu cores
	Memory int   `hcl:"memory,optional" json:"memory,omitempty"`                          // max memory the container can consume in MB
}

// Volume defines a folder, Docker volume, or temp folder to mount to the Container
type Volume struct {
	Source                      string `hcl:"source" json:"source"`                                                                                                                  // source path on the local machine for the volume
	Destination                 string `hcl:"destination" json:"destination"`                                                                                                        // path to mount the volume inside the container
	Type                        string `hcl:"type,optional" json:"type,omitempty"`                                                                                                   // type of the volume to mount [bind, volume, tmpfs]
	ReadOnly                    bool   `hcl:"read_only,optional" json:"read_only,omitempty" mapstructure:"read_only"`                                                                // specify that the volume is mounted read only
	BindPropagation             string `hcl:"bind_propagation,optional" json:"bind_propagation,omitempty" mapstructure:"bind_propagation"`                                           // propagation mode for bind mounts [shared, private, slave, rslave, rprivate]
	BindPropagationNonRecursive bool   `hcl:"bind_propagation_non_recursive,optional" json:"bind_propagation_non_recursive,omitempty" mapstructure:"bind_propagation_non_recursive"` // recursive bind mount, default true
}

// KV is a key/value type
type KV struct {
	Key   string `hcl:"key" json:"key"`
	Value string `hcl:"value" json:"value"`
}

// Build allows you to define the conditions for building a container
// on run from a Dockerfile
type Build struct {
	File    string `hcl:"file,optional" json:"file,omitempty"` // Location of build file inside build context defaults to ./Dockerfile
	Context string `hcl:"context" json:"context"`              // Path to build context
	Tag     string `hcl:"tag,optional" json:"tag,omitempty"`   // Image tag, defaults to latest
}

// New creates a new Nomad job config resource, implements Resource New method
func (c *Container) New(name string) types.Resource {
	return &Container{ResourceInfo: types.ResourceInfo{Name: name, Type: TypeContainer, Status: types.PendingCreation}}
}

// Info returns the resource info implements the Resource Info method
func (c *Container) Info() *types.ResourceInfo {
	return &c.ResourceInfo
}

func (c *Container) Parse(file string) error {
	return nil
}

func (c *Container) Process() error {
	return nil
}
