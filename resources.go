package gobilla

import (
	"fmt"
	"strconv"
)

// Canonical URI constants.
const (
	websitesURI = "/live/websites"
	buttonURI   = websitesURI + "/button"
	campaignURI = websitesURI + "/campaign"

	appsURI = "/live/apps"

	emailURI       = "/live/email"
	emailButtonURI = emailURI + "/button"
)

var (
	feedbackURI        = "%s/%s/feedback"
	campaignResultsURI = campaignURI + "/%s/results"
	campaignStatsURI   = campaignURI + "/%s/stats"
)

type resource struct {
	auth auth
}

// Buttons represents the button resource of Usabilla API.
type Buttons struct {
	resource
}

// Get function of Buttons resource returns all the buttons
// taking into account the specified query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (b *Buttons) Get(params map[string]string) (*ButtonResponse, error) {
	request := Request{
		method: "GET",
		auth:   b.auth,
		uri:    buttonURI,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewButtonResponse(data)
}

// Feedback encapsulates the feedback item resource.
func (b *Buttons) Feedback() *FeedbackItems {
	return &FeedbackItems{
		resource: resource{
			auth: b.auth,
		},
	}
}

// FeedbackItems represents the feedback item subresource of Usabilla API.
type FeedbackItems struct {
	resource
}

// Get function of FeedbackItem resource returns all the feedback items
// for a specific button, taking into account the provided query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (f *FeedbackItems) Get(buttonID string, params map[string]string) (*FeedbackResponse, error) {
	uri := fmt.Sprintf(feedbackURI, buttonURI, buttonID)

	request := &Request{
		method: "GET",
		auth:   f.auth,
		uri:    uri,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewFeedbackResponse(data)
}

// Iterate uses a FeedbackItem channel which transparently uses the HasMore field to fire
// a new api request once all items have been consumed on the channel
func (f *FeedbackItems) Iterate(buttonID string, params map[string]string) chan FeedbackItem {
	resp, err := f.Get(buttonID, params)

	if err != nil {
		panic(err)
	}

	fic := make(chan FeedbackItem)

	go items(fic, resp, f, buttonID)

	return fic
}

// items feeds a feedback item channel with items
//
// while hasMore is true and all items have been consumed in the channel
// a new request is fired using the since parameter of the response, to
// retrieve new items
//
// when HasMore is false, we close the channel
func items(fic chan FeedbackItem, resp *FeedbackResponse, f *FeedbackItems, buttonID string) {
	for {
		for _, item := range resp.Items {
			fic <- item
		}
		if !resp.HasMore {
			close(fic)
			return
		}
		params := map[string]string{
			"since": strconv.FormatInt(resp.LastTimestamp, 10),
		}

		resp, err := f.Get(buttonID, params)

		if err != nil {
			panic(err)
		}

		go items(fic, resp, f, buttonID)

		return
	}
}

// Campaigns represents the campaign resource of Usabilla API.
type Campaigns struct {
	resource
}

// Get function of Campaigns resource returns all the campaigns
// taking into account the provided query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (c *Campaigns) Get(params map[string]string) (*CampaignResponse, error) {
	request := Request{
		method: "GET",
		auth:   c.auth,
		uri:    campaignURI,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewCampaignResponse(data)
}

// Results encapsulates the campaign results resource.
func (c *Campaigns) Results() *CampaignResults {
	return &CampaignResults{
		resource: resource{
			auth: c.auth,
		},
	}
}

// CampaignResults represents the campaign result resource of Usabilla API.
type CampaignResults struct {
	resource
}

// Get function of CampaignResults resource returns all the campaign result items
// for a specific campaign, taking into account the provided query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (r *CampaignResults) Get(campaignID string, params map[string]string) (*CampaignResultResponse, error) {
	uri := fmt.Sprintf(campaignResultsURI, campaignID)

	request := Request{
		method: "GET",
		auth:   r.auth,
		uri:    uri,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewCampaignResultResponse(data)
}

// Iterate uses a CampaignResult channel which transparently uses the HasMore field to fire
// a new api request once all results have been consumed on the channel
func (r *CampaignResults) Iterate(campaignID string, params map[string]string) chan CampaignResult {
	resp, err := r.Get(campaignID, params)

	if err != nil {
		panic(err)
	}

	crc := make(chan CampaignResult)

	go campaignResults(crc, resp, r, campaignID)

	return crc
}

// campaignResults feeds a campaign results channel with items
//
// while hasMore is true and all items have been consumed in the channel
// a new request is fired using the since parameter of the response, to
// retrieve new items
//
// when HasMore is false, we close the channel
func campaignResults(crc chan CampaignResult, resp *CampaignResultResponse, r *CampaignResults, campaignID string) {
	for {
		for _, item := range resp.Items {
			crc <- item
		}
		if !resp.HasMore {
			close(crc)
			return
		}
		params := map[string]string{
			"since": strconv.FormatInt(resp.LastTimestamp, 10),
		}

		resp, err := r.Get(campaignID, params)

		if err != nil {
			panic(err)
		}

		go campaignResults(crc, resp, r, campaignID)

		return
	}
}

// Stats encapsulates the campaign statistics resource.
func (c *Campaigns) Stats() *CampaignStats {
	return &CampaignStats{
		resource: resource{
			auth: c.auth,
		},
	}
}

// CampaignStats represents the campaign statistics resource of Usabilla API.
type CampaignStats struct {
	resource
}

// Get function of CampaignStats resource returns the campaign statistics
// for a specific campaign, taking into account the provided query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (cs *CampaignStats) Get(campaignID string, params map[string]string) (*CampaignStatsResponse, error) {
	uri := fmt.Sprintf(campaignStatsURI, campaignID)

	request := Request{
		method: "GET",
		auth:   cs.auth,
		uri:    uri,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewCampaignStatsResponse(data)
}

// Iterate uses a CampaignStat channel which transparently uses the HasMore field to fire
// a new api request once all stats items have been consumed on the channel
func (cs *CampaignStats) Iterate(campaignID string, params map[string]string) chan CampaignStat {
	resp, err := cs.Get(campaignID, params)

	if err != nil {
		panic(err)
	}

	csc := make(chan CampaignStat)

	go campaignStats(csc, resp, cs, campaignID)

	return csc
}

// campagnStats feeds a campaign statistics channel with items
//
// while hasMore is true and all items have been consumed in the channel
// a new request is fired using the since parameter of the response, to
// retrieve new items
//
// when HasMore is false, we close the channel
func campaignStats(csc chan CampaignStat, resp *CampaignStatsResponse, cs *CampaignStats, campaignID string) {
	for {
		for _, item := range resp.Items {
			csc <- item
		}
		if !resp.HasMore {
			close(csc)
			return
		}
		params := map[string]string{
			"since": strconv.FormatInt(resp.LastTimestamp, 10),
		}

		resp, err := cs.Get(campaignID, params)

		if err != nil {
			panic(err)
		}

		go campaignStats(csc, resp, cs, campaignID)

		return
	}
}

// Apps represents the app resource of Usabilla API.
type Apps struct {
	resource
}

// Get function of Apps resource returns all apps
// taking into account the specified query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (a *Apps) Get(params map[string]string) (*AppResponse, error) {
	request := Request{
		method: "GET",
		auth:   a.auth,
		uri:    appsURI,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewAppResponse(data)
}

// Feedback encapsulates the app feedback item resource.
func (a *Apps) Feedback() *AppFeedbackItems {
	return &AppFeedbackItems{
		resource: resource{
			auth: a.auth,
		},
	}
}

// AppFeedbackItems represents the apps feedback item subresource of Usabilla API.
type AppFeedbackItems struct {
	resource
}

// Get function of AppFeedbackItem resource returns all the feedback items
// for a specific app, taking into account the provided query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (af *AppFeedbackItems) Get(appID string, params map[string]string) (*AppFeedbackResponse, error) {
	uri := fmt.Sprintf(feedbackURI, appsURI, appID)

	request := &Request{
		method: "GET",
		auth:   af.auth,
		uri:    uri,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewAppFeedbackResponse(data)
}

// Iterate uses an AppFeedbackItem channel which transparently uses the HasMore field to fire
// a new api request once all items have been consumed on the channel
func (af *AppFeedbackItems) Iterate(appID string, params map[string]string) chan AppFeedbackItem {
	resp, err := af.Get(appID, params)

	if err != nil {
		panic(err)
	}

	afic := make(chan AppFeedbackItem)

	go appItems(afic, resp, af, appID)

	return afic
}

// appItems feeds a feedback item channel with items.
//
// While hasMore is true and all items have been consumed in the channel
// a new request is fired using the since parameter of the response, to
// retrieve new items.
//
// When HasMore is false, we close the channel
func appItems(afic chan AppFeedbackItem, resp *AppFeedbackResponse, af *AppFeedbackItems, appID string) {
	for {
		for _, item := range resp.Items {
			afic <- item
		}
		if !resp.HasMore {
			close(afic)
			return
		}
		params := map[string]string{
			"since": strconv.FormatInt(resp.LastTimestamp, 10),
		}

		resp, err := af.Get(appID, params)

		if err != nil {
			panic(err)
		}

		go appItems(afic, resp, af, appID)

		return
	}
}

// EmailButtons represents the email button resource of Usabilla API.
type EmailButtons struct {
	resource
}

// Get function of EmailButtons resource returns all the email buttons
// taking into account the specified query params.
//
// Accepted query params are:
// - limit int
// - since string (Time stamp)
func (eb *EmailButtons) Get(params map[string]string) (*EmailButtonResponse, error) {
	request := Request{
		method: "GET",
		auth:   eb.auth,
		uri:    emailButtonURI,
		params: params,
	}

	data, err := request.Get()
	if err != nil {
		panic(err)
	}

	return NewEmailButtonResponse(data)
}

// Feedback encapsulates the email feedback item resource.
//
// We use the FeedbackItem response as it is the same with the feedback item
// response from websites, only difference is that image is contained
// in the website feedback item response, but it is omitted for the email one
func (eb *EmailButtons) Feedback() *FeedbackItems {
	return &FeedbackItems{
		resource: resource{
			auth: eb.auth,
		},
	}
}
