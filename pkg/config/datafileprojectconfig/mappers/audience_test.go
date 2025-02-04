/****************************************************************************
 * Copyright 2019,2021 Optimizely, Inc. and contributors                    *
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

// Package mappers //
package mappers

import (
	"testing"

	datafileEntities "github.com/WolffunService/experiment/pkg/config/datafileprojectconfig/entities"
	"github.com/WolffunService/experiment/pkg/entities"

	"github.com/stretchr/testify/assert"
)

func TestMapAudiencesEmptyList(t *testing.T) {

	audienceMap := MapAudiences(nil)
	expectedAudienceMap := map[string]entities.Audience{}

	assert.Equal(t, expectedAudienceMap, audienceMap)
}

func TestMapAudiences(t *testing.T) {

	expectedConditions := "[\"and\", [\"or\", [\"or\", {\"name\": \"s_foo\", \"type\": \"custom_attribute\", \"value\": \"foo\"}]]]"
	audienceList := []datafileEntities.Audience{{ID: "1", Name: "one", Conditions: expectedConditions}, {ID: "2", Name: "two"},
		{ID: "3", Name: "three"}, {ID: "2", Name: "four"}, {ID: "1", Name: "one"}}
	audienceMap := MapAudiences(audienceList)

	expectedConditionTree := &entities.TreeNode{
		Operator: "and",
		Nodes: []*entities.TreeNode{
			{
				Operator: "or",
				Nodes: []*entities.TreeNode{
					{
						Operator: "or",
						Nodes: []*entities.TreeNode{
							{
								Item: entities.Condition{
									Name:                 "s_foo",
									Type:                 "custom_attribute",
									Value:                "foo",
									StringRepresentation: `{"name":"s_foo","type":"custom_attribute","value":"foo"}`,
								},
							},
						},
					},
				},
			},
		},
	}
	expectedAudienceMap := map[string]entities.Audience{"1": {ID: "1", Name: "one", ConditionTree: expectedConditionTree, Conditions: expectedConditions}, "2": {ID: "2", Name: "two"},
		"3": {ID: "3", Name: "three"}}

	assert.Equal(t, expectedAudienceMap, audienceMap)
}
