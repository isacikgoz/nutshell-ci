package github

var queryTpl = `{
	repository(name: "{{.Repository}}", owner: "{{.Owner}}") {
	  pullRequest(number: {{.ID}}) {
		headRefName
		headRefOid
	  }
	}
  }`
