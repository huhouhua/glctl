package resource

import (
	"github.com/huhouhua/glctl/pkg/runtime"
	"github.com/xanzy/go-gitlab"
)

// Info contains temporary info to execute a REST call, or show the results
// of an already completed REST call.
type Info struct {
	// Client will only be present if this builder was not local
	Client *gitlab.Client
	// Namespace will be set if the object is namespaced and has a specified value.
	Group string
	Name  string

	// Optional, Source is the filename or URL to template file yaml,
	// or stdin to use to handle the resource
	Source string
	// Optional, this is the most recent value returned by the server if available. It will
	// typically be in unstructured or internal forms, depending on how the Builder was
	// defined. If retrieved from the server, the Builder expects the mapping client to
	// decide the final form. Use the AsVersioned, AsUnstructured, and AsInternal helpers
	// to alter the object versions.

	Object runtime.Object
	// Optional, this is the most recent resource version the server knows about for
	// this type of resource. It may not match the resource version of the object,
	// but if set it should be equal to or newer than the resource version of the
	// object (however the server defines resource version).
	ResourceVersion string
}
