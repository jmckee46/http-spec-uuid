package main

import "strings"

type specTriplet struct {
	DesiredRequest   *request
	ExpectedResponse *response
	ActualResponse   *response
}

func (specTriplet *specTriplet) isRequestOnly() bool {
	return specTriplet.ExpectedResponse == nil
}

func (specTriplet *specTriplet) String() string {
	result := []string{}

	if specTriplet.DesiredRequest != nil {
		result = append(result, specTriplet.DesiredRequest.FirstLine.Location())
	}

	if specTriplet.ExpectedResponse != nil {
		result = append(result, specTriplet.ExpectedResponse.FirstLine.Location())
	}

	if specTriplet.ActualResponse != nil {
		result = append(result, specTriplet.ActualResponse.FirstLine.Location())
	}

	return strings.Join(result, ", ")
}