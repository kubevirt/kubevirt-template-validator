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
	"math"
)

func toInt64(obj interface{}) (int64, bool) {
	if intVal, ok := obj.(int); ok {
		return int64(intVal), true
	}
	if intVal, ok := obj.(int32); ok {
		return int64(intVal), true
	}
	if intVal, ok := obj.(int64); ok {
		return int64(intVal), true
	}
	if intVal, ok := obj.(uint); ok {
		return int64(intVal), true
	}
	if intVal, ok := obj.(uint32); ok {
		return int64(intVal), true
	}
	if intVal, ok := obj.(uint64); ok {
		return int64(intVal), true
	}
	if floatVal, ok := obj.(float32); ok {
		return int64(math.Round(float64(floatVal))), true
	}
	if floatVal, ok := obj.(float64); ok {
		return int64(math.Round(floatVal)), true
	}
	return 0, false
}
