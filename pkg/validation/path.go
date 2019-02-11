/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2019 Red Hat, Inc.
 */

package validation

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/util/jsonpath"

	k6tv1 "kubevirt.io/kubevirt/pkg/api/v1"
)

var (
	ErrMismatchingTypes error = fmt.Errorf("Mismatching type(s)")
	ErrWrongTypes       error = fmt.Errorf("Wrong type(s)")
)

type Path struct {
	jp      *jsonpath.JSONPath
	results [][]reflect.Value
}

func (p *Path) Len() int {
	return len(p.results)
}

func (p *Path) AsString() ([]string, error) {
	var ret []string
	for i := range p.results {
		res := p.results[i]
		for j := range res {
			obj := res[j].Interface()
			strObj, ok := obj.(string)
			if ok {
				ret = append(ret, strObj)
				continue
			}
			return nil, fmt.Errorf("mismatching type: %v, not string", res[j].Type().Name())
		}
	}
	return ret, nil
}

func (p *Path) AsInt64() ([]int64, error) {
	var ret []int64
	for i := range p.results {
		res := p.results[i]
		for j := range res {
			obj := res[j].Interface()
			if intObj, ok := obj.(int64); ok {
				ret = append(ret, intObj)
				continue
			}
			if quantityObj, ok := obj.(resource.Quantity); ok {
				v, ok := quantityObj.AsInt64()
				if ok {
					ret = append(ret, v)
					continue
				}
			}
			return nil, fmt.Errorf("mismatching type: %v, not int or resource.Quantity", res[j].Type().Name())
		}
	}
	return ret, nil
}

func Find(vm *k6tv1.VirtualMachine, expr string) (*Path, error) {
	jp := jsonpath.New(expr) // unique name
	err := jp.Parse(expr)
	if err != nil {
		return nil, err
	}

	results, err := jp.FindResults(vm)
	if err != nil {
		return nil, err
	}
	return &Path{
		jp:      jp,
		results: results,
	}, nil
}
