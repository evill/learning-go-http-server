package main

type HttpTextBody struct {
	text string
}

func (textBody HttpTextBody) ToString() string {
	return string(textBody.text)
}

func (textBody HttpTextBody) ContentLength() int {
	return len(textBody.text)
}

func (textBody HttpTextBody) ContentType() string {
	return "plain/text"
}
