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

package mappers

import (
	"testing"

	datafileConfig "github.com/WolffunService/experiment/pkg/config/datafileprojectconfig/entities"
	"github.com/WolffunService/experiment/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func TestBuildAudienceConditionTreeEmpty(t *testing.T) {
	conditionString := ""
	var conditions interface{}
	json.Unmarshal([]byte(conditionString), &conditions)
	conditionTree, err := buildAudienceConditionTree(conditions)

	expectedTree := &entities.TreeNode{Operator: "or"}
	assert.NoError(t, err)
	assert.Equal(t, expectedTree, conditionTree)
}

func TestBuildAudienceConditionTreeSimpleAudienceCondition(t *testing.T) {
	conditionString := "[ \"and\", [ \"or\", [ \"or\",  \"12\", \"123\", \"1234\"] ] ]"
	var conditions interface{}
	json.Unmarshal([]byte(conditionString), &conditions)
	conditionTree, err := buildAudienceConditionTree(conditions)
	if err != nil {
		assert.Fail(t, err.Error())
	}

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
								Item: "12",
							},
							{
								Item: "123",
							},
							{
								Item: "1234",
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedConditionTree, conditionTree)
}

func TestBuildAudienceConditionTreeNoOperators(t *testing.T) {
	conditions := []string{"123"}
	expectedConditionTree := &entities.TreeNode{
		Operator: "or",
		Nodes: []*entities.TreeNode{
			{
				Item: "123",
			},
		},
	}

	conditionTree, err := buildAudienceConditionTree(conditions)
	assert.Equal(t, expectedConditionTree, conditionTree)
	assert.NoError(t, err)
}

func TestBuildConditionTreeUsingDatafileAudienceConditions(t *testing.T) {

	audience := datafileConfig.Audience{
		ID:         "12567320080",
		Name:       "message",
		Conditions: "[\"and\", [\"or\", [\"or\", {\"name\": \"s_foo\", \"type\": \"custom_attribute\", \"value\": \"foo\"}]]]",
	}

	conditionTree, err := buildConditionTree(audience.Conditions)
	if err != nil {
		assert.Fail(t, err.Error())
	}

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
	assert.Equal(t, expectedConditionTree, conditionTree)
}

func TestBuildConditionTreeSimpleAudienceCondition(t *testing.T) {
	conditionString := "[ \"and\", [ \"or\", [ \"or\", { \"type\": \"custom_attribute\", \"name\": \"s_foo\", \"match\": \"exact\", \"value\": \"foo\" } ] ] ]"
	var conditions interface{}
	json.Unmarshal([]byte(conditionString), &conditions)
	conditionTree, err := buildConditionTree(conditions)
	if err != nil {
		assert.Fail(t, err.Error())
	}

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
									Match:                "exact",
									Type:                 "custom_attribute",
									Value:                "foo",
									StringRepresentation: `{"match":"exact","name":"s_foo","type":"custom_attribute","value":"foo"}`,
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, expectedConditionTree, conditionTree)
}

func TestBuildConditionTreeWithLeafNode(t *testing.T) {
	conditionString := "{ \"type\": \"custom_attribute\", \"name\": \"s_foo\", \"match\": \"exact\", \"value\": \"foo\" }"
	var conditions interface{}
	json.Unmarshal([]byte(conditionString), &conditions)
	conditionTree, err := buildConditionTree(conditions)
	assert.NoError(t, err)

	expectedConditionTree := &entities.TreeNode{
		Operator: "or",
		Nodes: []*entities.TreeNode{
			{
				Item: entities.Condition{
					Name:                 "s_foo",
					Match:                "exact",
					Type:                 "custom_attribute",
					Value:                "foo",
					StringRepresentation: `{"match":"exact","name":"s_foo","type":"custom_attribute","value":"foo"}`,
				},
			},
		},
	}
	assert.Equal(t, expectedConditionTree, conditionTree)
}
