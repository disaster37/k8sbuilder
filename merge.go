package k8sbuilder

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// MergeK8s permit to merge kubernetes resources
func MergeK8s(dst any, src, new any) (err error) {
  if dst != nil && reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer of instanciated object")
	}

	if src == nil  ||  (reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil()) {
    return errors.New("src can't be null")
	}

	if new == nil || (reflect.ValueOf(new).Kind() == reflect.Ptr && reflect.ValueOf(new).IsNil()) {
    var value reflect.Value
		if reflect.ValueOf(src).Kind() == reflect.Ptr {
			value = reflect.ValueOf(src).Elem()
		} else {
			value = reflect.ValueOf(src)
		}
		reflect.ValueOf(dst).Elem().Set(value)
	}

	dstByte, err := json.Marshal(dst)
	if err != nil {
		return err
	}
	newByte, err := json.Marshal(new)
	if err != nil {
		return err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(dstByte, newByte, reflect.ValueOf(dst).Elem().Interface())
	if err != nil {
		return err
	}

	expectedByte, err := strategicpatch.StrategicMergePatch(dstByte, patch, reflect.ValueOf(dst).Elem().Interface())
	if err != nil {
		return err
	}

	if err = json.Unmarshal(expectedByte, dst); err != nil {
		return err
	}

	return nil
}


// MergeSliceOrDie permit to merge some slice on dst
// It avoid to set the same item based on key value
func MergeSliceOrDie(dst *[]any, key string,  src ...[]any) {
	if dst == nil {
		panic("dst can't be nil")
	}
	
	for _, src :=  range src {
		loopExpected: for _, expectedItem := range src {
			for _, currentItem := range *dst {
				if funk.Get(currentItem, key) == funk.Get(expectedItem, key) {
					continue loopExpected
				}
			}
			*dst = append(*dst, expectedItem)
		}
	}
}