package goplex

type User struct {
	Email        string       `xml:"email,attr"`
	Id           int          `xml:"id,attr"`
	Thumb        HttpUrl      `xml:"thumb,attr"`
	Username     string       `xml:"username,attr"`
	Title        string       `xml:"title,attr"`
	Locale       string       `xml:"locale,attr"`
	AuthToken    string       `xml:"authenticationToken,attr"`
	QueueEmail   string       `xml:"queueEmail,attr"`
	Subscription Subscription `xml:"subscription"`
}

type Subscription struct {
	Active IntAsBool `xml:"active,attr"`
	Plan   string    `xml:"plan,attr"`
}
