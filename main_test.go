package main

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCommonHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/common", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2, err := http.NewRequest("POST", "/common", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(commonSetHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusSeeOther {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusSeeOther)
		}

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestCustomHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/custom", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2, err := http.NewRequest("POST", "/custom", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(customSetHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusSeeOther {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusSeeOther)
		}

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestTodoHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	req2, err := http.NewRequest("POST", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusSeeOther {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusSeeOther)
		}

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestRandomPageHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req1, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req3, err := http.NewRequest("GET", "/20", nil)
	if err != nil {
		t.Fatal(err)
	}

	req4, err := http.NewRequest("GET", "/testindex.html", nil)
	if err != nil {
		t.Fatal(err)
	}

	req5, err := http.NewRequest("GET", "/Thisshouldnotbeaurl", nil)
	if err != nil {
		t.Fatal(err)
	}

	req6, err := http.NewRequest("GET", "/This should be & a ) bad / request", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req3, req4, req5, req6}

	for _, req := range reqs {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(randomPageHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the response body is what we expect.
		// expected := `{"alive": true}`
		// if rr.Body.String() != expected {
		// 	t.Errorf("handler returned unexpected body: got %v want %v",
		// 		rr.Body.String(), expected)
		// }
	}
}

func TestHomeHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestGetPort(t *testing.T) {
	// Save value so we can set it back
	orig, isSet := os.LookupEnv("PORT")

	err := os.Unsetenv("PORT")
	if err != nil {
		t.Fatal(err)
	}
	expected := ":8080"

	if getPort() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			getPort(), expected)
	}

	err = os.Setenv("PORT", "3000")
	if err != nil {
		t.Fatal(err)
	}
	expected = ":3000"

	if getPort() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			getPort(), expected)
	}

	if isSet {
		err = os.Setenv("PORT", orig)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		err = os.Unsetenv("PORT")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestRoll(t *testing.T) {
	numTests := 10000
	minVal := 1
	maxVal := 20
	expectedNumberOfHits := numTests / ((maxVal - minVal) + 1)
	criticalChi := 36028.79701896 // This is for a 5% with a DF of 9999

	die := Dice{High: maxVal, Low: minVal}

	totalValues := make(map[int]int)

	for i := 0; i < numTests; i++ {
		val := die.roll()
		totalValues[val] = totalValues[val] + 1
	}

	keys := make([]int, 0, len(totalValues))
	for k := range totalValues {
		keys = append(keys, k)
	}

	for j := minVal; j < maxVal; j++ {
		_, ok := totalValues[j]
		if !ok {
			t.Errorf("%d value was not reached", j)
		}
	}

	sort.Ints(keys)

	chiSquared := 0.0
	for _, val := range keys {
		numHits := totalValues[val]
		subChi := math.Pow((float64(numHits)-float64(expectedNumberOfHits)), 2) / float64(expectedNumberOfHits)
		chiSquared = chiSquared + subChi
		fmt.Printf("%d:\n", val)
		fmt.Printf("    numHits  = %d\n", numHits)
		fmt.Printf("    expected = %d\n", expectedNumberOfHits)
		fmt.Printf("    subChi   = %.6f\n", subChi)
	}

	fmt.Printf("\nchi-squared = %.6f\n", chiSquared)
	fmt.Printf("critical-ch = %.6f\n", criticalChi)
	if chiSquared > criticalChi {
		t.Errorf("the chi-squared value was too large")
	}

}
