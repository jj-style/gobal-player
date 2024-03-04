package resty

import "encoding/json"

// shim - pure function acting as a method/delegate to client.post to allow generics
func Post[T, U any](c Client, uri string, req T) (U, error) {
	var res U
	resp, err := c.Post(uri, req)
	if err != nil {
		return res, err
	}
	if err = json.Unmarshal(resp, &res); err != nil {
		return res, err
	}
	return res, nil
}

// shim - pure function acting as a method/delegate to client.http.get to allow generics
func Get[T any](c Client, uri string) (T, error) {
	var res T
	resp, err := c.Get(uri)
	if err != nil {
		return res, err
	}
	if err = json.Unmarshal(resp, &res); err != nil {
		return res, err
	}
	return res, nil
}
