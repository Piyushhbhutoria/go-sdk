/****************************************************************************
 * Copyright 2019-2020, Optimizely, Inc. and contributors                   *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

// Package matchers //
package matchers

import (
	"fmt"

	"github.com/optimizely/go-sdk/pkg/decision/evaluator/matchers/utils"
	"github.com/optimizely/go-sdk/pkg/entities"
	"github.com/optimizely/go-sdk/pkg/logging"
)

// GtMatcher matches against the "gt" match type
type GtMatcher struct {
	Condition entities.Condition
	Logger    logging.OptimizelyLogProducer
}

// Match returns true if the user's attribute is greater than the condition's string value
func (m GtMatcher) Match(user entities.UserContext) (bool, error) {
	return matchGtOrLt(user, m.Condition, m.Logger, true)
}

func matchGtOrLt(user entities.UserContext, condition entities.Condition, logger logging.OptimizelyLogProducer, gtMatch bool) (bool, error) {
	if !user.CheckAttributeExists(condition.Name) {
		logger.Debugf(logging.NullUserAttribute.String(), condition.StringRepresentation, condition.Name)
		return false, fmt.Errorf(`no attribute named "%s"`, condition.Name)
	}

	if floatValue, ok := utils.ToFloat(condition.Value); ok {
		attributeValue, err := user.GetFloatAttribute(condition.Name)
		if err != nil {
			val, _ := user.GetAttribute(condition.Name)
			logger.Warningf(logging.InvalidAttributeValueType.String(), condition.StringRepresentation, val, condition.Name)
			return false, err
		}
		if gtMatch {
			return floatValue < attributeValue, nil
		}
		return floatValue > attributeValue, nil
	}

	logger.Warningf(logging.UnsupportedConditionValue.String(), condition.StringRepresentation)
	return false, fmt.Errorf("audience condition %s evaluated to NULL because the condition value type is not supported", condition.Name)
}
