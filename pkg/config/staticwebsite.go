package config

type StaticWebSite struct {
	Enable    bool   `default:"false"`
	WebFolder string `default:"/opt/{{app}}/www"`
	IndexPage string `default:"index.html"`
	Route     string `default:"/"`
}
