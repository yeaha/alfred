// 解析输入的时间字符串，打印出另外几种时间格式输出
//
package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ruedap/go-alfred"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: time [time]")
		os.Exit(1)
	}

	input := strings.Join(os.Args[1:], " ")

	t, err := parse(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	u := t.UTC()

	r := alfred.NewResponse()

	r.AddItem(&alfred.ResponseItem{
		Valid: true,
		Title: fmt.Sprint(t.Unix()),
	})

	if t.Nanosecond() > 0 {
		r.AddItem(&alfred.ResponseItem{
			Valid: true,
			Title: fmt.Sprint(t.UnixNano()),
		})
	}

	r.AddItem(&alfred.ResponseItem{
		Valid: true,
		Title: fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()),
	})

	r.AddItem(&alfred.ResponseItem{
		Valid: true,
		Title: fmt.Sprint(t.Format(time.RFC3339)),
	})
	r.AddItem(&alfred.ResponseItem{
		Valid: true,
		Title: fmt.Sprint(u.Format(time.RFC3339)),
	})

	if t.Nanosecond() > 0 {
		r.AddItem(&alfred.ResponseItem{
			Valid: true,
			Title: fmt.Sprint(t.Format(time.RFC3339Nano)),
		})
		r.AddItem(&alfred.ResponseItem{
			Valid: true,
			Title: fmt.Sprint(u.Format(time.RFC3339Nano)),
		})
	}

	xml, err := r.ToXML()
	if err != nil {
		title := fmt.Sprintf("Error: %v", err.Error())
		subtitle := "Time Workflow Error"
		arg := title
		errXML := alfred.ErrorXML(title, subtitle, arg)
		fmt.Println(errXML)
	}

	fmt.Println(xml)
}

func parse(input string) (time.Time, error) {
	switch {
	case input == "now":
		return time.Now(), nil
	// unix timestamp
	case regexp.MustCompile(`^\d{10,19}$`).Match([]byte(input)):
		return parseUnixTimestamp(input)
	case regexp.MustCompile(`^\d{4}\-\d{1,2}\-\d{1,2}$`).Match([]byte(input)):
		return time.Parse("2006-01-02", input)
	case regexp.MustCompile(`^\d{4}\-\d{1,2}\-\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}$`).Match([]byte(input)):
		return time.Parse("2006-01-02 15:04:05", input)
	}

	layouts := []string{time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.RFC822, time.RFC822Z, time.RFC850}
	for _, l := range layouts {
		if t, err := time.Parse(l, input); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("unknown time format")
}

func parseUnixTimestamp(input string) (time.Time, error) {
	switch len(input) {
	case 10: // second
		i, _ := strconv.ParseInt(input, 10, 64)
		return time.Unix(i, 0), nil
	case 13: // millisecond
		s, _ := strconv.ParseInt(input[:10], 10, 64)
		m, _ := strconv.ParseInt(input[10:], 10, 64)
		return time.Unix(s, m*1000*1000), nil
	case 16: // microsecond
		s, _ := strconv.ParseInt(input[:10], 10, 64)
		m, _ := strconv.ParseInt(input[10:], 10, 64)
		return time.Unix(s, m*1000), nil
	case 19: // nanosecond
		s, _ := strconv.ParseInt(input[:10], 10, 64)
		n, _ := strconv.ParseInt(input[10:], 10, 64)
		return time.Unix(s, n), nil
	}

	return time.Time{}, errors.New("invalid unix timestamp")
}
