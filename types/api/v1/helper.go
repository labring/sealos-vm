/*
Copyright 2022 The labring Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"fmt"

	"github.com/labring/sealvm/pkg/utils/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// ToAggregate converts the ErrorList into an errors.Aggregate.
func ToAggregate(list []error) utilerrors.Aggregate {
	errs := make([]error, 0, len(list))
	errorMsgs := sets.NewString()
	for _, err := range list {
		msg := fmt.Sprintf("%v", err)
		if errorMsgs.Has(msg) {
			continue
		}
		errorMsgs.Insert(msg)
		errs = append(errs, err)
	}
	return utilerrors.NewAggregate(errs)
}

func IsConditionTrue(conditions []Condition, condition Condition) bool {
	for _, con := range conditions {
		if con.Type == condition.Type && con.Status == condition.Status {
			return true
		}
	}
	return false
}
func IsConditionsTrue(conditions []Condition) bool {
	if len(conditions) == 0 {
		return false
	}
	for _, condition := range conditions {
		if condition.Type == "Ready" {
			continue
		}
		if condition.Status != v1.ConditionTrue {
			return false
		}
	}
	return true
}

// UpdateCondition updates condition in cluster conditions using giving condition
// adds condition if not existed
func UpdateCondition(conditions []Condition, condition Condition) []Condition {
	if conditions == nil {
		conditions = make([]Condition, 0)
	}
	hasCondition := false
	for i, cond := range conditions {
		if cond.Type == condition.Type {
			hasCondition = true
			if cond.Reason != condition.Reason || cond.Status != condition.Status || cond.Message != condition.Message {
				conditions[i] = condition
			}
		}
	}
	if !hasCondition {
		conditions = append(conditions, condition)
	}
	return conditions
}

func DeleteCondition(conditions []Condition, conditionType string) []Condition {
	if conditions == nil {
		conditions = make([]Condition, 0)
	}
	newConditions := make([]Condition, 0)
	for _, cond := range conditions {
		if cond.Type == conditionType {
			continue
		}
		newConditions = append(newConditions, cond)
	}
	conditions = newConditions
	return conditions
}

func SetConditionError(condition *Condition, reason string, err error) {
	condition.LastHeartbeatTime = metav1.Now()
	condition.Status = v1.ConditionFalse
	condition.Reason = reason
	condition.Message = err.Error()
	logger.Error("Exec failed reason: %s,error: %+v", reason, err)
}
