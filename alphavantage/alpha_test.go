package alpha

import (
	"testing"
)

func TestRetrieveHistory(t *testing.T) {

	_, err := RetrieveHistory("AAPL", 50)
	if err != nil {
		t.Errorf("RetrieveHistory() error = %v", err)
		return
	}

}
func TestRetrieveWeeklyHistory(t *testing.T) {

	_, err := RetrieveWeeklyHistory("AAPL", 10)
	if err != nil {
		t.Errorf("TestRetrieveWeeklyHistory() error = %v", err)
		return
	}

}
