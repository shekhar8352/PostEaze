package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Call(request *http.Request, c http.Client) error {
	resp, err := c.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// You can add logic to treat non-2xx as error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-success HTTP status: %s", resp.Status)
	}

	return nil
}

func CallAndGetResponse(request *http.Request, c http.Client) ([]byte, error) {
	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-success HTTP status: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func CallAndBind(request *http.Request, response interface{}, c http.Client) error {
	bodyBytes, err := CallAndGetResponse(request, c)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
