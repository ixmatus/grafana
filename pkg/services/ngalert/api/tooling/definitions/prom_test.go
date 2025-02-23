package definitions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortAlertsByImportance(t *testing.T) {
	tm1, tm2 := time.Now(), time.Now().Add(time.Second)
	tc := []struct {
		name     string
		input    []Alert
		expected []Alert
	}{{
		name:     "alerts are ordered in expected importance",
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}, {State: "pending"}, {State: "error"}, {State: "nodata"}, {State: "normal"}},
	}, {
		name:     "alerts with same importance are ordered active first",
		input:    []Alert{{State: "normal"}, {State: "normal", ActiveAt: &tm1}},
		expected: []Alert{{State: "normal", ActiveAt: &tm1}, {State: "normal"}},
	}, {
		name:     "active alerts with same importance are ordered newest first",
		input:    []Alert{{State: "alerting", ActiveAt: &tm2}, {State: "alerting", ActiveAt: &tm1}},
		expected: []Alert{{State: "alerting", ActiveAt: &tm1}, {State: "alerting", ActiveAt: &tm2}},
	}, {
		name: "inactive alerts with same importance are ordered by labels",
		input: []Alert{
			{State: "normal", Labels: map[string]string{"c": "d"}},
			{State: "normal", Labels: map[string]string{"a": "b"}},
		},
		expected: []Alert{
			{State: "normal", Labels: map[string]string{"a": "b"}},
			{State: "normal", Labels: map[string]string{"c": "d"}},
		},
	}, {
		name: "active alerts with same importance and active time are ordered fewest labels first",
		input: []Alert{
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"a": "b", "c": "d"}},
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"e": "f"}},
		},
		expected: []Alert{
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"e": "f"}},
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"a": "b", "c": "d"}},
		},
	}, {
		name: "active alerts with same importance and active time are ordered by labels",
		input: []Alert{
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"c": "d"}},
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"a": "b"}},
		},
		expected: []Alert{
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"a": "b"}},
			{State: "alerting", ActiveAt: &tm1, Labels: map[string]string{"c": "d"}},
		},
	}}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			AlertsBy(AlertsByImportance).Sort(tt.input)
			assert.EqualValues(t, tt.expected, tt.input)
		})
	}
}

func TestTopKAlertsByImportance(t *testing.T) {
	//	tm1, tm2 := time.Now(), time.Now().Add(time.Second)
	tc := []struct {
		name     string
		k        int
		input    []Alert
		expected []Alert
	}{{
		name:     "alerts are ordered in expected importance (k=1)",
		k:        1,
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}},
	}, {
		name:     "alerts are ordered in expected importance (k=2)",
		k:        2,
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}, {State: "pending"}},
	}, {
		name:     "alerts are ordered in expected importance (k=3)",
		k:        3,
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}, {State: "pending"}, {State: "error"}},
	}, {
		name:     "alerts are ordered in expected importance (k=4)",
		k:        4,
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}, {State: "pending"}, {State: "error"}, {State: "nodata"}},
	}, {
		name:     "alerts are ordered in expected importance (k=5)",
		k:        5,
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}, {State: "pending"}, {State: "error"}, {State: "nodata"}, {State: "normal"}},
	}, {
		name:     "alerts are ordered in expected importance (k=6)",
		k:        6,
		input:    []Alert{{State: "normal"}, {State: "nodata"}, {State: "error"}, {State: "pending"}, {State: "alerting"}},
		expected: []Alert{{State: "alerting"}, {State: "pending"}, {State: "error"}, {State: "nodata"}, {State: "normal"}},
	},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			result := AlertsBy(AlertsByImportance).TopK(tt.input, tt.k)
			assert.EqualValues(t, tt.expected, result)
		})
	}
}
