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

// Package evaluator //
package evaluator

import (
	"fmt"

	"github.com/optimizely/go-sdk/pkg/decision/evaluator/matchers"
	"github.com/optimizely/go-sdk/pkg/entities"
	"github.com/optimizely/go-sdk/pkg/logging"
)

const (
	exactMatchType     = "exact"
	existsMatchType    = "exists"
	ltMatchType        = "lt"
	gtMatchType        = "gt"
	substringMatchType = "substring"
)

// ItemEvaluator evaluates a condition against the given user's attributes
type ItemEvaluator interface {
	Evaluate(interface{}, *entities.TreeParameters) (bool, error)
}

// CustomAttributeConditionEvaluator evaluates conditions with custom attributes
type CustomAttributeConditionEvaluator struct {
	logger logging.OptimizelyLogProducer
}

// NewCustomAttributeConditionEvaluator creates a custom attribute condition
func NewCustomAttributeConditionEvaluator(logger logging.OptimizelyLogProducer) *CustomAttributeConditionEvaluator {
	return &CustomAttributeConditionEvaluator{logger: logger}
}

// Evaluate returns true if the given user's attributes match the condition
func (c CustomAttributeConditionEvaluator) Evaluate(condition entities.Condition, condTreeParams *entities.TreeParameters) (bool, error) {
	// We should only be evaluating custom attributes

	if condition.Type != customAttributeType {

		c.logger.Warningf(logging.UnknownConditionType.String(), condition.StringRepresentation)
		return false, fmt.Errorf(`unable to evaluate condition of type "%s"`, condition.Type)
	}

	var matcher matchers.Matcher
	matchType := condition.Match
	if matchType == "" {
		matchType = exactMatchType
	}
	switch matchType {
	case exactMatchType:
		matcher = matchers.ExactMatcher{
			Condition: condition,
			Logger:    c.logger,
		}
	case existsMatchType:
		matcher = matchers.ExistsMatcher{
			Condition: condition,
		}
	case ltMatchType:
		matcher = matchers.LtMatcher{
			Condition: condition,
			Logger:    c.logger,
		}
	case gtMatchType:
		matcher = matchers.GtMatcher{
			Condition: condition,
			Logger:    c.logger,
		}
	case substringMatchType:
		matcher = matchers.SubstringMatcher{
			Condition: condition,
			Logger:    c.logger,
		}
	default:
		c.logger.Warningf(logging.UnknownMatchType.String(), condition.StringRepresentation)
		return false, fmt.Errorf(`invalid Condition matcher "%s"`, condition.Match)
	}

	user := *condTreeParams.User
	result, err := matcher.Match(user)
	return result, err
}

// AudienceConditionEvaluator evaluates conditions with audience condition
type AudienceConditionEvaluator struct {
	logger logging.OptimizelyLogProducer
}

// NewAudienceConditionEvaluator creates a audience condition evaluator
func NewAudienceConditionEvaluator(logger logging.OptimizelyLogProducer) *AudienceConditionEvaluator {
	return &AudienceConditionEvaluator{logger: logger}
}

// Evaluate returns true if the given user's attributes match the condition
func (c AudienceConditionEvaluator) Evaluate(audienceID string, condTreeParams *entities.TreeParameters) (bool, error) {

	if audience, ok := condTreeParams.AudienceMap[audienceID]; ok {
		c.logger.Debugf(logging.AudienceEvaluationStarted.String(), audienceID)
		condTree := audience.ConditionTree
		conditionTreeEvaluator := NewMixedTreeEvaluator(c.logger)
		retValue, isValid := conditionTreeEvaluator.Evaluate(condTree, condTreeParams)
		if !isValid {
			return false, fmt.Errorf(`an error occurred while evaluating nested tree for audience ID "%s"`, audienceID)
		}
		c.logger.Debugf(logging.AudienceEvaluatedTo.String(), audienceID, retValue)
		return retValue, nil
	}

	return false, fmt.Errorf(`unable to evaluate nested tree for audience ID "%s"`, audienceID)
}
