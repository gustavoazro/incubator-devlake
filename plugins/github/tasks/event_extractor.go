/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiEventsMeta = core.SubTaskMeta{
	Name:             "extractApiEvents",
	EntryPoint:       ExtractApiEvents,
	EnabledByDefault: true,
	Description:      "Extract raw Events data into tool layer table github_issue_events",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

type IssueEvent struct {
	GithubId int `json:"id"`
	Event    string
	Actor    *GithubAccountResponse
	Issue    struct {
		Id int
	}
	GithubCreatedAt helper.Iso8601Time `json:"created_at"`
}

func ExtractApiEvents(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_EVENTS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			body := &IssueEvent{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			if body.GithubId == 0 || body.Actor == nil {
				return nil, nil
			}
			githubIssueEvent := &models.GithubIssueEvent{
				ConnectionId:    data.Options.ConnectionId,
				GithubId:        body.GithubId,
				IssueId:         body.Issue.Id,
				Type:            body.Event,
				GithubCreatedAt: body.GithubCreatedAt.ToTime(),
			}

			if body.Actor != nil {
				githubIssueEvent.AuthorUsername = body.Actor.Login

				githubAccount, err := convertAccount(body.Actor, data.Options.GithubId, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubAccount)
			}

			if err != nil {
				return nil, err
			}
			results = append(results, githubIssueEvent)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
