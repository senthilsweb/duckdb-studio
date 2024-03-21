package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"templrjs/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func jiraApiVersion() string {
	return viper.GetString("integration.jira.api_version")
}

func jiraURL() string {
	return viper.GetString("integration.jira.url")
}

func jiraUsername() string {
	return viper.GetString("integration.jira.username")
}

func jiraToken() string {
	return viper.GetString("integration.jira.token")
}

func timezone() string {
	return viper.GetString("integration.jira.timezone")
}

//var DoneTransitionName = os.Getenv("DONE_TransitionName")

var loc, _ = time.LoadLocation(timezone())

var UserIds = map[string]string{
	"user1": "1t6881uu9d23ad4ef096g78k",
	"user2": "2rt881uu9d23ad4ef096g78k",
	"user3": "378881uu9d23ad4ef096g78k",
	"user4": "4ad881uu9d23ad4ef096g78k",
}

var AssigneeMap = map[string]string{
	"Title1": UserIds["user1"],
	"Title2": UserIds["user2"],
	"Title3": UserIds["user3"],
	"Title4": UserIds["user4"],
	"Title5": UserIds["user1"],
	"Title6": UserIds["user2"],
}

type Issue struct {
	Project     string
	Key         string
	Summary     string
	Description string
	Assignee    string
	Type        string
	Status      string
	Labels      []string
	Delta       time.Duration
	Updaded     time.Time
	Created     time.Time
}

type Desc struct {
	Content   string
	Extension bool
}

var templateAssignee = "\"assignee\": {\"id\": \"%s\"},"

func New(issue Issue) Issue {
	return issue
}

func (issue *Issue) AddIssue() (err error) {

	client := &http.Client{}
	templateJson := `{
		"fields": {
			{{.assignee}}
			{{.labels}}
			"project": {
				"key": "{{.project}}"
			},
			"issuetype": {
				"name": "{{.issuetype}}"
			},
			"summary": "{{.summary}}",
			"description": "{{.description}} \n"
		}
	}`

	data := map[string]interface{}{
		"project":     issue.Project,
		"assignee":    "",
		"labels":      "",
		"issuetype":   issue.Type,
		"summary":     issue.Summary,
		"description": issue.Description,
	}

	if issue.Assignee == "auto" {
		if utils.KeyInStringMap(AssigneeMap, issue.Summary) {
			data["assignee"] = fmt.Sprintf(templateAssignee, AssigneeMap[issue.Summary])
		}
	} else if issue.Assignee != "" {
		data["assignee"] = fmt.Sprintf(templateAssignee, issue.Assignee)
	}

	if len(issue.Labels) > 0 {
		data["labels"] = "\"labels\": [\"" + strings.Join(issue.Labels, "\",\"") + "\"],"
	}

	t := template.Must(template.New("json").Parse(templateJson))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	dataJson := buf.String()

	//log.Printf("Add Issue JSON: %s \n", dataJson)

	requestBody := strings.NewReader(dataJson)

	req, err := http.NewRequest("POST", jiraURL()+"/rest/api/"+jiraApiVersion()+"/issue/", requestBody)

	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("---- ERROR ------")
		//log.Fatal(err)
		return err
	}

	if resp.StatusCode == 201 {
		return nil
	} else {
		return errors.New("Error:" + resp.Status)
	}

	//log.Printf("Response : %+v", resp)

	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	var parsed map[string]interface{}
	json.Unmarshal([]byte(s), &parsed)

	issue.Key = parsed["key"].(string)

	return nil
}

func (issue *Issue) DelIssue() (err error) {

	client := &http.Client{}
	url := jiraURL() + "/rest/api/" + jiraApiVersion() + "/issue/" + issue.Key
	req, err := http.NewRequest("DELETE", url, nil)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("---- ERROR ------")
		//log.Fatal(err)
		return err
	}

	if resp.StatusCode == 201 {
		return nil
	} else {
		return errors.New("Error:" + resp.Status)
	}

	//log.Printf("Response : %+v", resp)

	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	var parsed map[string]interface{}
	json.Unmarshal([]byte(s), &parsed)

	issue.Key = parsed["key"].(string)

	return nil
}

func (issue Issue) GetIssues() (issues []Issue) {

	jsondata := issue.getIssuesJson()
	//fmt.Printf(jsondata)
	return issue.parseIssuesJson(jsondata)
}

func (issue Issue) GetIssue() Issue {

	jsondata := issue.getIssueJson()
	//fmt.Printf(jsondata)
	return issue.parseIssueJson(jsondata)
}

func (issue Issue) getIssueJson() string {

	client := &http.Client{}
	url := jiraURL() + "/rest/api/" + jiraApiVersion() + "/issue/" + issue.Key
	log.Printf("Response : %s", url)
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("---- ERROR ------")
		log.Fatal(err)
	}
	//log.Printf("Response : %+v", resp)
	bodyText, err := ioutil.ReadAll(resp.Body)

	return string(bodyText)
}

func (issue Issue) getIssuesJson() string {

	client := &http.Client{}
	//url := fmt.Sprintf(jiraURL() + "/rest/api/"+jiraApiVersion+"/search?jql=project=%s%%20and%%20status!=Done&maxResults=250", issue.Project)
	//log.Printf("Response : %s", url)
	url := fmt.Sprintf(jiraURL()+"/rest/api/"+jiraApiVersion()+"/search?jql=project=%s&maxResults=250", issue.Project)

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("---- ERROR ------")
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	return string(bodyText)
}

func (issue Issue) parseIssuesJson(content string) (issues []Issue) {

	var parsed map[string]interface{}
	json.Unmarshal([]byte(content), &parsed)

	//log.Printf(">>> Content : %+v",  parsed["issues"])

	arr := parsed["issues"].([]interface{})

	for _, val := range arr {
		level1 := val.(map[string]interface{})
		key := level1["key"].(string)
		level2 := level1["fields"].(map[string]interface{})
		summary := level2["summary"].(string)

		labels := level2["labels"].([]interface{})

		var labelsArray []string

		for _, item := range labels {
			labelsArray = append(labelsArray, item.(string))
		}

		status := level2["status"].(map[string]interface{})["name"].(string)
		status = strings.ToUpper(strings.ReplaceAll(status, " ", ""))

		//if status == DoneTransitionName {
		//	continue
		//}

		desc := level2["description"].(string)

		updaded, _ := time.Parse("2006-01-02T15:04:05.999-0700", level2["updated"].(string))
		created, _ := time.Parse("2006-01-02T15:04:05.999-0700", level2["created"].(string))

		now := time.Now().In(loc)

		delta := now.Sub(created)

		issue := Issue{
			Key:         key,
			Summary:     summary,
			Delta:       delta,
			Status:      status,
			Labels:      labelsArray,
			Description: desc,
			Updaded:     updaded,
			Created:     created,
		}

		issues = append(issues, issue)
	}

	return
}

func (issue Issue) parseIssueJson(content string) Issue {

	var parsed map[string]interface{}
	json.Unmarshal([]byte(content), &parsed)

	key := parsed["key"].(string)

	fields := parsed["fields"].(map[string]interface{})

	summary := fields["summary"].(string)

	labels := fields["labels"].([]interface{})

	var labelsArray []string

	for _, item := range labels {
		labelsArray = append(labelsArray, item.(string))
	}

	status := fields["status"].(map[string]interface{})["name"].(string)
	status = strings.ToUpper(strings.ReplaceAll(status, " ", ""))

	desc := fields["description"].(string)

	updaded, _ := time.Parse("2006-01-02T15:04:05.999-0700", fields["updated"].(string))
	created, _ := time.Parse("2006-01-02T15:04:05.999-0700", fields["created"].(string))

	now := time.Now().In(loc)

	delta := now.Sub(updaded)

	issue = Issue{
		Key:         key,
		Summary:     summary,
		Labels:      labelsArray,
		Delta:       delta,
		Status:      status,
		Description: desc,
		Updaded:     updaded,
		Created:     created,
	}

	return issue
}

func (issue *Issue) FindIssueBySummary(issues []Issue) bool {

	if issues == nil {
		issues = issue.GetIssues()
	}

	for _, item := range issues {
		if item.Summary == issue.Summary {
			*issue = item
			return true
		}
	}

	return false
}

func (issue Issue) UpdateLabel2Issue() (bool, error) {

	client := &http.Client{}
	templateJson := "{\"update\":{\"labels\":[{\"set\":%s}]}}"
	labelsArray := "[\"" + strings.Join(issue.Labels, "\",\"") + "\"]"
	dataJson := fmt.Sprintf(templateJson, labelsArray)
	requestBody := strings.NewReader(dataJson)
	url := jiraURL() + "/rest/api/" + jiraApiVersion() + "/issue/" + issue.Key

	req, err := http.NewRequest("PUT", url, requestBody)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("---- ERROR ------")
		return false, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else {
		log.Printf("---- Fail ------")
		return false, errors.New("Error:" + resp.Status)
	}
}

func (issue Issue) AddLabel2Issue(label string) (bool, error) {

	client := &http.Client{}
	templateJson := "{\"update\":{\"labels\":[{\"add\":\"%s\"}]}}"
	dataJson := fmt.Sprintf(templateJson, label)
	fmt.Printf("\n %s \n", dataJson)
	requestBody := strings.NewReader(dataJson)
	url := jiraURL() + "/rest/api/" + jiraApiVersion() + "/issue/" + issue.Key
	fmt.Printf("\n %v \n", issue)

	req, err := http.NewRequest("PUT", url, requestBody)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("---- ERROR ------")
		return false, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else {
		log.Printf("---- Fail ------")
		return false, errors.New("Error:" + resp.Status)
	}
}

func (issue Issue) AddComment2Issue(commment string) (bool, error) {

	client := &http.Client{}
	templateJson := `{ "body": "%s" }`
	dataJson := fmt.Sprintf(templateJson, commment)
	fmt.Printf("\n %s \n", dataJson)

	requestBody := strings.NewReader(dataJson)
	url := jiraURL() + "/rest/api/" + jiraApiVersion() + "/issue/" + issue.Key + "/comment"

	req, err := http.NewRequest("POST", url, requestBody)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("---- ERROR ------")
		//log.Fatal(err)
		return false, err
	}

	if resp.StatusCode == 201 {
		return true, nil
	} else {
		return false, errors.New("Error:" + resp.Status)
	}

}

func (issue Issue) UpdateDescription(desc Desc) (bool, error) {

	client := &http.Client{}
	templateJson := `{
				"fields": {
					"description":" %s \n"
				}
			}`
	var dataJson string

	if desc.Extension == true {
		dataJson = fmt.Sprintf(templateJson, strings.ReplaceAll(issue.Description, "\n", "\\n")+desc.Content)
	} else {
		dataJson = fmt.Sprintf(templateJson, desc.Content)
	}

	requestBody := strings.NewReader(dataJson)
	url := jiraURL() + "/rest/api/" + jiraApiVersion() + "/issue/" + issue.Key

	req, err := http.NewRequest("PUT", url, requestBody)
	req.SetBasicAuth(jiraUsername(), jiraToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("---- ERROR ------")
		return false, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else {
		return false, errors.New("Error:" + resp.Status)
	}
}
