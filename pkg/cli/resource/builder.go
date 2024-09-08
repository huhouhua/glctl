package resource

// Builder provides convenience functions for taking arguments and parameters
// from the command line and converting them to a list of resources to iterate
// over using the Visitor interface.
type Builder struct {
	errs  []error
	paths []Visitor
	dir   bool
}

type FilenameOptions struct {
	Filenames []string
	Recursive bool
}

func (b *Builder) AddError(err error) *Builder {
	if err == nil {
		return b
	}
	b.errs = append(b.errs, err)
	return b
}
