package jellyapi

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

type Server struct {
	Host       string
	AuthKey    string
	AuthHeader string
	UserId     string

	client *http.Client
}

func NewServer(host, authKey, userId string) *Server {
	return &Server{
		Host:       host,
		AuthKey:    authKey,
		client:     &http.Client{},
		AuthHeader: fmt.Sprintf("MediaBrowser Client=\"JellyCli\", Device=\"Pavilion\", DeviceId=\"1\", Version=\"1.0\", Token=\"%s\"", authKey),
		UserId:     userId,
	}
}

type Collection struct {
	Id          string
	Name        string
	IsFolder    bool
	Type        string
	IndexNumber int
	VideoType   string
}

func (i Collection) Title() string       { return i.Name }
func (i Collection) Description() string { return "" }
func (i Collection) FilterValue() string { return i.Name }

func (s *Server) processRequest(url string, response any) error {
	//log.Println("processing request for %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("X-Emby-Authorization", s.AuthHeader)
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %s", err)
	}
	return nil
}

type CollectionResponse struct {
	Collections []Collection `json:"Items"`
}

func (s *Server) GetCollections() ([]Collection, error) {
	url := fmt.Sprintf(s.Host+"/Items?userId=%s", s.UserId)
	var cr CollectionResponse
	err := s.processRequest(url, &cr)
	if err != nil {
		return nil, fmt.Errorf("error getting collections: %s", err)
	}
	return cr.Collections, nil
}

func (s *Server) GetChildItems(parentId string) ([]Collection, error) {
	url := fmt.Sprintf(s.Host+"/Items?parentId=%s", parentId)
	var cr CollectionResponse
	err := s.processRequest(url, &cr)
	if err != nil {
		return nil, fmt.Errorf("error getting collections: %s", err)
	}
	if len(cr.Collections) > 0 {
		switch cr.Collections[0].Type {
		case "Episode":
			slices.SortStableFunc(cr.Collections, func(a, b Collection) int {
				return cmp.Compare(a.IndexNumber, b.IndexNumber)
			})
		}
	}
	return cr.Collections, nil
}
