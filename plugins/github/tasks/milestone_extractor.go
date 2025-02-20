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

var ExtractMilestonesMeta = core.SubTaskMeta{
	Name:             "extractMilestones",
	EntryPoint:       ExtractMilestones,
	EnabledByDefault: true,
	Description:      "Extract raw milestone data into tool layer table github_milestones",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

type MilestonesResponse struct {
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	LabelsUrl   string `json:"labels_url"`
	Id          int    `json:"id"`
	NodeId      string `json:"node_id"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Creator     struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"creator"`
	OpenIssues   int                 `json:"open_issues"`
	ClosedIssues int                 `json:"closed_issues"`
	State        string              `json:"state"`
	CreatedAt    helper.Iso8601Time  `json:"created_at"`
	UpdatedAt    helper.Iso8601Time  `json:"updated_at"`
	DueOn        *helper.Iso8601Time `json:"due_on"`
	ClosedAt     *helper.Iso8601Time `json:"closed_at"`
}

func ExtractMilestones(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_MILESTONE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			response := &MilestonesResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, response))
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			results = append(results, convertGithubMilestone(response, data.Options.ConnectionId, data.Options.GithubId))
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertGithubMilestone(response *MilestonesResponse, connectionId uint64, repositoryId int) *models.GithubMilestone {
	milestone := &models.GithubMilestone{
		ConnectionId: connectionId,
		MilestoneId:  response.Id,
		RepoId:       repositoryId,
		Number:       response.Number,
		URL:          response.Url,
		Title:        response.Title,
		OpenIssues:   response.OpenIssues,
		ClosedIssues: response.ClosedIssues,
		State:        response.State,
		ClosedAt:     helper.Iso8601TimeToTime(response.ClosedAt),
		CreatedAt:    response.CreatedAt.ToTime(),
		UpdatedAt:    response.UpdatedAt.ToTime(),
	}
	return milestone
}
