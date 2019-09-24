// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrSignatureInvalid is returned when the webhook
	// signature is invalid or cannot be calculated.
	ErrSignatureInvalid = errors.New("Invalid webhook signature")
)

type (
	// Webhook defines a webhook for repository events.
	Webhook interface {
		Repository() Repository
	}

	// WebhookUnmarshaler wraps Webhook and assigns a type for unmarshalling.
	// Use this if you need to deserialize Webhooks of uknown concrete type.
	WebhookUnmarshaler struct {
		Type    string
		Webhook Webhook
	}

	// Label on a PR
	Label struct {
		URL         string
		Name        string
		Description string
		Color       string
	}

	// PushCommit represents general info about a commit.
	PushCommit struct {
		ID       string
		Message  string
		Added    []string
		Removed  []string
		Modified []string
	}

	// PushHook represents a push hook, eg push events.
	PushHook struct {
		Ref     string
		BaseRef string
		Repo    Repository
		Before  string
		After   string
		Created bool
		Deleted bool
		Forced  bool
		Compare string
		Commits []PushCommit
		Commit  Commit
		Sender  User
		GUID    string
	}

	// BranchHook represents a branch or tag event,
	// eg create and delete github event types.
	BranchHook struct {
		Ref    Reference
		Repo   Repository
		Action Action
		Sender User
	}

	// TagHook represents a tag event, eg create and delete
	// github event types.
	TagHook struct {
		Ref    Reference
		Repo   Repository
		Action Action
		Sender User
	}

	// IssueHook represents an issue event, eg issues.
	IssueHook struct {
		Action Action
		Repo   Repository
		Issue  Issue
		Sender User
	}

	// IssueCommentHook represents an issue comment event,
	// eg issue_comment.
	IssueCommentHook struct {
		Action  Action
		Repo    Repository
		Issue   Issue
		Comment Comment
		Sender  User
	}

	PullRequestHookBranchFrom struct {
		From string
	}

	PullRequestHookBranch struct {
		Ref  PullRequestHookBranchFrom
		Sha  PullRequestHookBranchFrom
		Repo Repository
	}

	PullRequestHookChanges struct {
		Base PullRequestHookBranch
	}

	// PullRequestHook represents an pull request event,
	// eg pull_request.
	PullRequestHook struct {
		Action      Action
		Repo        Repository
		Label       Label
		PullRequest PullRequest
		Sender      User
		Changes     PullRequestHookChanges
		GUID        string
	}

	// PullRequestCommentHook represents an pull request
	// comment event, eg pull_request_comment.
	PullRequestCommentHook struct {
		Action      Action
		Repo        Repository
		PullRequest PullRequest
		Comment     Comment
		Sender      User
	}

	// ReviewCommentHook represents a pull request review
	// comment, eg pull_request_review_comment.
	ReviewCommentHook struct {
		Action      Action
		Repo        Repository
		PullRequest PullRequest
		Review      Review
	}

	// DeployHook represents a deployment event. This is
	// currently a GitHub-specific event type.
	DeployHook struct {
		Data      interface{}
		Desc      string
		Ref       Reference
		Repo      Repository
		Sender    User
		Target    string
		TargetURL string
		Task      string
	}

	// SecretFunc provides the Webhook parser with the
	// secret key used to validate webhook authenticity.
	SecretFunc func(webhook Webhook) (string, error)

	// WebhookService provides abstract functions for
	// parsing and validating webhooks requests.
	WebhookService interface {
		// Parse returns the parsed the repository webhook payload.
		Parse(req *http.Request, fn SecretFunc) (Webhook, error)
	}
)

// Repository() defines the repository webhook and provides
// a convenient way to get the associated repository without
// having to cast the type.

func (h *PushHook) Repository() Repository               { return h.Repo }
func (h *BranchHook) Repository() Repository             { return h.Repo }
func (h *DeployHook) Repository() Repository             { return h.Repo }
func (h *TagHook) Repository() Repository                { return h.Repo }
func (h *IssueHook) Repository() Repository              { return h.Repo }
func (h *IssueCommentHook) Repository() Repository       { return h.Repo }
func (h *PullRequestHook) Repository() Repository        { return h.Repo }
func (h *PullRequestCommentHook) Repository() Repository { return h.Repo }
func (h *ReviewCommentHook) Repository() Repository      { return h.Repo }

// MarshalJSON adds a type field to a serialized Webhook so that it can be deserialized into the same concrete type.
func (wu *WebhookUnmarshaler) MarshalJSON() ([]byte, error) {

	marshaledHook := make(map[string]interface{})

	var genericWebhook interface{}
	genericWebhook = wu.Webhook

	if _, ok := genericWebhook.(PushHook); ok {
		marshaledHook["type"] = "pushHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(BranchHook); ok {
		marshaledHook["type"] = "branchHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(DeployHook); ok {
		marshaledHook["type"] = "deployHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(TagHook); ok {
		marshaledHook["type"] = "tagHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(IssueHook); ok {
		marshaledHook["type"] = "issueHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(IssueCommentHook); ok {
		marshaledHook["type"] = "issueCommentHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(PullRequestHook); ok {
		marshaledHook["type"] = "pullRequestHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(PullRequestCommentHook); ok {
		marshaledHook["type"] = "pullRequestCommentHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	if _, ok := genericWebhook.(ReviewCommentHook); ok {
		marshaledHook["type"] = "reviewCommentHook"
		marshaledHook["webhook"] = wu.Webhook
		return json.Marshal(marshaledHook)
	}

	return nil, fmt.Errorf("WebhookUnmarshaler.Webhook does not implement a concrete Webhook type")
}

// UnmarshalJSON supports deserialization of GitEventSpec.ParsedWebhook into a concrete implementation of scm.Webhook
func (wu *WebhookUnmarshaler) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMessage *json.RawMessage
	var webhookMap map[string]string
	err = json.Unmarshal(*rawMessage, &webhookMap)
	if err != nil {
		return err
	}

	if webhookMap["type"] == "pushHook" {

		var h *PushHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "branchHook" {

		var h *BranchHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "deployHook" {

		var h *DeployHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "tagHook" {

		var h *TagHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "issueHook" {

		var h *IssueHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "issueCommentHook" {

		var h *IssueHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "pullRequestHook" {

		var h *IssueHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "pullRequestCommentHook" {

		var h *IssueHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	} else if webhookMap["type"] == "reviewCommentHook" {

		var h *IssueHook
		err = json.Unmarshal(*rawMessage, h)
		if err != nil {
			return err
		}
		wu.Webhook = h

	}

	return nil
}
