// Code generated by go-swagger; DO NOT EDIT.

package aks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ListAKSLocationsReader is a Reader for the ListAKSLocations structure.
type ListAKSLocationsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAKSLocationsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAKSLocationsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListAKSLocationsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAKSLocationsOK creates a ListAKSLocationsOK with default headers values
func NewListAKSLocationsOK() *ListAKSLocationsOK {
	return &ListAKSLocationsOK{}
}

/* ListAKSLocationsOK describes a response with status code 200, with default header values.

AKSLocationList
*/
type ListAKSLocationsOK struct {
	Payload models.AKSLocationList
}

func (o *ListAKSLocationsOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/locations][%d] listAKSLocationsOK  %+v", 200, o.Payload)
}
func (o *ListAKSLocationsOK) GetPayload() models.AKSLocationList {
	return o.Payload
}

func (o *ListAKSLocationsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAKSLocationsDefault creates a ListAKSLocationsDefault with default headers values
func NewListAKSLocationsDefault(code int) *ListAKSLocationsDefault {
	return &ListAKSLocationsDefault{
		_statusCode: code,
	}
}

/* ListAKSLocationsDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ListAKSLocationsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list a k s locations default response
func (o *ListAKSLocationsDefault) Code() int {
	return o._statusCode
}

func (o *ListAKSLocationsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/locations][%d] listAKSLocations default  %+v", o._statusCode, o.Payload)
}
func (o *ListAKSLocationsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAKSLocationsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
