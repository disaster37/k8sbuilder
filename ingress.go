package k8sbuilder

import (
	"reflect"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	networkingv1 "k8s.io/api/networking/v1"
)

// IngressBuilder is the ingress builder interface
type IngressBuilder interface{
	WithIngressSpec(is *networkingv1.IngressSpec, opts ...WithOption) IngressBuilder
	WithLabels(labels map[string]string, opts ...WithOption) IngressBuilder
	WithAnnotations(annotations map[string]string, opts ...WithOption) IngressBuilder
	WithName(name string, opts ...WithOption) IngressBuilder
	WithNamespace(namespace string, opts ...WithOption) IngressBuilder
	Build() (i *networkingv1.Ingress, err error)
}

// IngressBuilderDefault is the default implementation for ingress builder
type IngressBuilderDefault struct {
	i *networkingv1.Ingress
	operations []Operation
}

// NewIngressBuilder permit to get the default ingress builder
func NewIngressBuilder() IngressBuilder {
	return &IngressBuilderDefault{
		i: &networkingv1.Ingress{},
		operations: make([]Operation, 0),
	}
}

// Build permit to build the expected object
// It will execute all pending operation in the same order
// At the end, it will clean all pending operations
func (h *IngressBuilderDefault) Build() (i *networkingv1.Ingress, err error) {

	rv := reflect.ValueOf(h)

	for _, o :=  range h.operations {
		if o.Name != "" {
			m := rv.MethodByName(o.Name)
			if m.IsZero() {
				return nil, errors.Errorf("Method %s not found", o.Name)
			}
			args := make([]reflect.Value, 0, len(o.Args))
			for _, argv := range o.Args {
				args = append(args, reflect.ValueOf(argv))
			} 
			res := m.Call(args)

			for _, r := range res {
				if _, ok := r.Interface().(*error); ok {
					if !r.IsNil() {
						return nil, r.Interface().(error)
					}
				}
			}
		}
	}

	h.operations = make([]Operation, 0)

	return h.i, nil
}

// WithIngressSpec permit to initialize ingress from ingress Spec
func (h *IngressBuilderDefault) WithIngressSpec(is *networkingv1.IngressSpec, opts ...WithOption) IngressBuilder {
	
	o := Operation{
		Name: "withIngressSpec",
		Args: append([]any{is}, opts),
	}
	h.operations = append(h.operations, o)
	
	return h
}

// WithLabels permit to set labels
func (h *IngressBuilderDefault) WithLabels(labels map[string]string, opts ...WithOption) IngressBuilder {
	
	o := Operation{
		Name: "withLabels",
		Args: append([]any{labels}, opts),
	}
	h.operations = append(h.operations, o)
	
	return h
}

// WithAnnotations permit to set annotation
func (h *IngressBuilderDefault) WithAnnotations(annotations map[string]string, opts ...WithOption) IngressBuilder {
	
	o := Operation{
		Name: "withAnnotations",
		Args: append([]any{annotations}, opts),
	}
	h.operations = append(h.operations, o)
	
	return h
}

// WithName permit to set name
func (h *IngressBuilderDefault) WithName(name string, opts ...WithOption) IngressBuilder {

	o := Operation{
		Name: "withName",
		Args: append([]any{name}, opts),
	}
	h.operations = append(h.operations, o)

	return h
}

// WithNamespace permit to set namespace
func (h *IngressBuilderDefault) WithNamespace(namespace string, opts ...WithOption) IngressBuilder {

	o := Operation{
		Name: "withNamespace",
		Args: append([]any{namespace}, opts),
	}
	h.operations = append(h.operations, o)

	return h
}

func (h *IngressBuilderDefault) withName(name string, opts ...WithOption) (err error) {

	// Overwrite
	if IsOverwrite(opts) || IsMerge(opts) || h.i.Name == "" {
		h.i.Name = name
	}

	return nil
}

func (h *IngressBuilderDefault) withNamespace(namespace string, opts ...WithOption) (err error) {
	
	// Overwrite
	if IsOverwrite(opts) || IsMerge(opts) || h.i.Namespace == "" {
		h.i.Namespace = namespace
	}
	
	return nil
}

func (h *IngressBuilderDefault) withLabels(labels map[string]string, opts ...WithOption) (err error) {
	
	// Overwrite
	if IsOverwrite(opts) || h.i.Labels == nil {
		h.i.Labels = labels
		return nil
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.i.Labels).Elem().IsZero() {
		h.i.Labels = labels
		return nil
	}

	// Merge
	if IsMerge(opts) && labels != nil {
		if err := mergo.Merge(h.i.Labels, labels); err != nil {
			return errors.Wrap(err, "Error when merge labels")
		}
	}
	
	return nil
}

func (h *IngressBuilderDefault) withAnnotations(annotations map[string]string, opts ...WithOption) (err error) {
	
	// Overwrite
	if IsOverwrite(opts) || h.i.Annotations == nil {
		h.i.Annotations = annotations
		return nil
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.i.Labels).Elem().IsZero() {
		h.i.Annotations = annotations
		return nil
	}

	// Merge
	if IsMerge(opts) && annotations != nil {
		if err := mergo.Merge(h.i.Annotations, annotations); err != nil {
			return errors.Wrap(err, "Error when merge annotations")
		}
	}
	
	return nil
}

func (h *IngressBuilderDefault) withIngressSpec(is *networkingv1.IngressSpec, opts ...WithOption) (err error) {
	
	if is == nil {
		return nil
	}

	// Overwrite
	if IsOverwrite(opts) {
		h.i.Spec = *is
		return nil
	}

	// Overwrite only if not default
	if IsOverwriteIfDefaultValue(opts) && reflect.ValueOf(h.i.Spec).Elem().IsZero() {
		 h.i.Spec = *is
			return nil
		}

	// Merge
	if IsMerge(opts) {
		//orgIngressSpec := h.i.Spec.DeepCopy()

		if err := MergeK8s(&h.i.Spec, h.i.Spec, is); err != nil {
			return errors.Wrap(err, "Error when merge ingress spec")
		}
	}
	
	
	return nil
}