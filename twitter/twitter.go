package twitter

import (
	"github.com/logashoff/anaconda"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	IMAGE_DIR       = "./static/images/twitter/"
	POLL_INTERVAL   = 15 * time.Minute
	TWITTER_USER_ID = 111507370
	TWEET_COUNT     = 300
)

var (
	twitterConsumerKey       = os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerSecret    = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	twitterAccessTokenKey    = os.Getenv("TWITTER_ACCESS_TOKEN_KEY")
	imageMap                 = make(map[string]string)
	api                      = anaconda.NewTwitterApi(twitterAccessTokenKey, twitterAccessTokenSecret)
	values                   = &url.Values{}
	instance                 = &Twitter{}
)

type BasicTweet struct {
	Id           string            `json:"id"`
	Username     string            `json:"username"`
	Url          string            `json:"url"`
	Timestamp    string            `json:"timestamp"`
	ScreenName   string            `json:"screen_name"`
	Text         string            `json:"text"`
	UserImage    string            `json:"user_image"`
	ProfileColor string            `json:"profile_color"`
	Entities     anaconda.Entities `json:"entities"`
	Protected    bool              `json:"protected"`
}

type Twitter struct {
	Tweets []anaconda.Tweet
}

func (t *Twitter) NewBasicTweet(tweet anaconda.Tweet) BasicTweet {
	at := tweet
	if tweet.Retweeted {
		at = *tweet.RetweetedStatus
	}
	user := at.User
	return BasicTweet{
		Id:           at.IdStr,
		Username:     user.Name,
		Url:          "https://twitter.com/" + user.ScreenName + "/statuses/" + at.IdStr,
		Timestamp:    at.CreatedAt,
		ScreenName:   user.ScreenName,
		Text:         at.Text,
		UserImage:    imageMap[user.ProfileImageUrlHttps],
		ProfileColor: "#" + user.ProfileLinkColor,
		Entities:     at.Entities,
		Protected:    user.Protected,
	}
}

func (t *Twitter) GetTweetId(tweet anaconda.Tweet) (id string) {
	if tweet.Retweeted {
		if len(tweet.RetweetedStatus.IdStr) > 0 {
			id = tweet.RetweetedStatus.IdStr
		} else {
			id = tweet.IdStr
		}
	} else {
		id = tweet.IdStr
	}
	return id
}

func (t *Twitter) ToBasicTweets(tweets []anaconda.Tweet) (basicTweets []BasicTweet) {
	for i := range tweets {
		basicTweets = append(basicTweets, t.NewBasicTweet(tweets[i]))
	}
	return basicTweets
}

func (t *Twitter) GetTweets() {
	t.Tweets, _ = api.GetUserTimeline(*values)
	for _, tweet := range t.Tweets {
		if tweet.Retweeted {
			getImages(&tweet.RetweetedStatus.User.ProfileImageUrlHttps)
		} else {
			getImages(&tweet.User.ProfileImageUrlHttps)
		}
	}
}

func (t *Twitter) PollTweets() {
	c := time.Tick(POLL_INTERVAL)
	for range c {
		t.GetTweets()
	}
}

func (t *Twitter) GetMaxId(maxId string, limit int) (tweets []anaconda.Tweet) {
	if len(maxId) > 0 {
		length := len(t.Tweets)
		for i := range t.Tweets {
			tweet := t.Tweets[i]
			id := t.GetTweetId(tweet)
			if id == maxId {
				if length-i < limit {
					limit = length - i - 1
				}
				min := i + 1
				max := min + limit
				tweets = t.Tweets[min:max]
				break
			}
		}
	}
	return tweets
}

func NewTwitter() Twitter {
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	values.Set("user_id", strconv.FormatInt(TWITTER_USER_ID, 10))
	values.Set("count", strconv.FormatInt(TWEET_COUNT, 10))
	values.Set("exclude_replies", "true")
	values.Set("include_rts", "true")
	instance.GetTweets()
	go instance.PollTweets()
	return *instance
}

func getImages(url *string) error {
	u := *url
	if imageMap[u] != "" {
		return nil
	}
	urlParts := strings.Split(u, "/")
	image := urlParts[len(urlParts)-1]
	imageMap[u] = IMAGE_DIR + image
	resp, getErr := http.Get(u)
	defer resp.Body.Close()
	if getErr != nil {
		log.Fatal(getErr)
		return getErr
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return readErr
	}
	writeErr := ioutil.WriteFile(imageMap[u], body, 0444)
	if writeErr != nil {
		log.Fatal(writeErr)
		return writeErr
	}
	return nil
}