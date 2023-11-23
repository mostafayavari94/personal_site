package main

import (
	// "io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexAnotherMethod (t *testing.T){
rr := httptest.NewRecorder()

req, err := http.NewRequest(http.MethodGet,"",nil)
if err != nil {
	t.Error(err)
}

index(rr, req)

if rr.Result().StatusCode != http.StatusOK {
	t.Errorf("expected %v, got %v", http.StatusOK, rr.Result().Status)
}

defer rr.Result().Body.Close()

// expected := `Foo`

// body, err := io.ReadAll(rr.Result().Body)
// if err != nil {
// 	t.Error(err)
// }

// if string(body) != expected {
// 	t.Errorf("expected %v, got %v", expected, string(body))
// }

}


func TestIndex (t *testing.T){
	server := httptest.NewServer(http.HandlerFunc(index))
	res, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, res.Status)
	}

	defer res.Body.Close()

	// expected := `Foo`

	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	t.Error(err)
	// }
	
	// if string(body) != expected {
	// 	t.Errorf("expected %v, got %v", expected, string(body))
	// }

}