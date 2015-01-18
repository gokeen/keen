package keen

import (
	"gopkg.in/gokeen/query.v1"
	"gopkg.in/httgo/dmx.v2"
	"gopkg.in/httgo/mock.v1"
	"gopkg.in/nowk/assert.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func h(str string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(str))
	})
}

type MyEvent struct {
	Foo string `json:"foo"`
}

func (MyEvent) CollectionName() string {
	return "awesomeness"
}

func TestWrite(t *testing.T) {
	mux := dmx.New()
	mux.Post("/3.0/projects/:project_id/events/:collection", h(`{"created":true}`))
	mo := &mock.Mock{
		Testing: t,
		Ts:      httptest.NewUnstartedServer(mux.Handler(dmx.NotFound(mux))),
	}
	mo.Start()
	defer mo.Done()
	var setMock = func(c *Client) {
		c.Client = mo
	}

	k := NewClient("12345", setMock, func(c *Client) {
		c.WriteKey = "abcdefg"
	})

	err := k.Write(MyEvent{"bar"})
	assert.Nil(t, err)

	reqs := mo.History("POST",
		"https://api.keen.io/3.0/projects/12345/events/awesomeness")
	assert.Equal(t, 1, len(reqs))

	r := reqs[0]
	hdr := r.Header
	assert.Equal(t, "abcdefg", hdr.Get("Authorization"))
	assert.Equal(t, "application/json", hdr.Get("Content-Type"))

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `{"foo":"bar"}`, string(b))
}

func TestWrite_NotCreatedReturnsError(t *testing.T) {
	mux := dmx.New()
	mux.Post("/3.0/projects/:project_id/events/:collection",
		h(`{"created":false,"message":"something went wrong","error_code":"bad"}`))
	mo := &mock.Mock{
		Testing: t,
		Ts:      httptest.NewUnstartedServer(mux.Handler(dmx.NotFound(mux))),
	}
	mo.Start()
	defer mo.Done()
	var setMock = func(c *Client) {
		c.Client = mo
	}

	k := NewClient("12345", setMock, func(c *Client) {
		c.WriteKey = "abcdefg"
	})

	err := k.Write(MyEvent{"bar"})
	assert.Equal(t, "event not created: [bad] something went wrong", err.Error())
}

func TestWrite_RequiresProperKey(t *testing.T) {
	err := NewClient("12345").Write(MyEvent{"bar"})
	assert.Equal(t, "key error: authorization key required", err.Error())
}

type MyResult struct {
	Result int `json:"result"`
}

func TestQuery(t *testing.T) {
	mux := dmx.New()
	mux.Post("/3.0/projects/:project_id/queries/:resource", h(`{"result":120}`))
	mo := &mock.Mock{
		Testing: t,
		Ts:      httptest.NewUnstartedServer(mux.Handler(dmx.NotFound(mux))),
	}
	mo.Start()
	defer mo.Done()
	var setMock = func(c *Client) {
		c.Client = mo
	}

	k := NewClient("67890", setMock, func(c *Client) {
		c.ReadKey = "hijklm"
	})

	var d MyResult
	q := query.Count("awesome-events")
	err := k.Query(q, &d)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 120, d.Result)

	reqs := mo.History("POST",
		"https://api.keen.io/3.0/projects/67890/queries/count")
	assert.Equal(t, 1, len(reqs))

	r := reqs[0]
	hdr := r.Header
	assert.Equal(t, "hijklm", hdr.Get("Authorization"))
	assert.Equal(t, "application/json", hdr.Get("Content-Type"))

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `{"event_collection":"awesome-events"}`, string(b))
}

func TestQuery_RequiresProperKey(t *testing.T) {
	err := NewClient("67890").Query(query.Count("awesome-events"), nil)
	assert.Equal(t, "key error: authorization key required", err.Error())
}
