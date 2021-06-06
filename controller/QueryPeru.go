package controller

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulantezana/shopping/utilities"
	"io/ioutil"
	"net/http"
	"strings"
)

type queryPeru struct {
	DocumentNumber             string `json:"document_number"`
	SocialReason               string `json:"social_reason"`
	Address                    string `json:"address"`
	TaxpayerState              string `json:"taxpayer_state"`
	DomicileCondition          string `json:"domicile_condition"`
	UtilIdentityDocumentTypeId uint   `json:"util_identity_document_type_id"`
}

type censusPeruRuc struct {
	Ruc               string `json:"ruc"`
	SocialReason      string `json:"social_reason"`
	TaxpayerState     string `json:"taxpayer_state"`
	DomicileCondition string `json:"domicile_condition"`
	Ubigeo            string `json:"ubigeo"`
	TypeRoad          string `json:"type_road"`
	NameRoad          string `json:"name_road"`
	ZoneCode          string `json:"zone_code"`
	TypeZone          string `json:"type_zone"`
	Number            string `json:"number"`
	Inside            string `json:"inside"`
	Lot               string `json:"lot"`
	Department        string `json:"department"`
	Kilometer         string `json:"kilometer"`
	Address           string `json:"address"`
	FullAddress       string `json:"full_address"`
	LastUpdateSunat   string `json:"last_update_sunat"`
}

type censusPeruDni struct {
	Name           string `json:"name"`
	MotherLastName string `json:"motherLastName"`
	LastName       string `json:"lastName"`
	DocumentNumber string `json:"documentNumber"`
	Sex            string `json:"sex"`
	BirthDate      string `json:"birthDate"`
}

type resultRuc struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Result  censusPeruRuc `json:"result"`
}

type resultDni struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Result  censusPeruDni `json:"result"`
}

// CreatePurchase function create new purchase
func QueryDocument(c echo.Context) error {
	// Get data request
	queryPeruData := queryPeru{}
	if err := c.Bind(&queryPeruData); err != nil {
		return c.JSON(http.StatusBadRequest, utilities.Response{
			Message: "La estructura no es válida",
		})
	}
	if len(queryPeruData.DocumentNumber) == 8 {
		result, err := queryDni(queryPeruData.DocumentNumber)
		if err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		if !result.Success {
			return c.JSON(http.StatusOK, utilities.Response{Message: result.Message})
		}
		queryPeruData.SocialReason = result.Result.Name + " " + result.Result.LastName + " " + result.Result.MotherLastName
		queryPeruData.UtilIdentityDocumentTypeId = 2
	} else if len(queryPeruData.DocumentNumber) == 11 {
		result, err := queryRuc(queryPeruData.DocumentNumber)
		if err != nil {
			return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("%s", err)})
		}
		if !result.Success {
			return c.JSON(http.StatusOK, utilities.Response{Message: result.Message})
		}
		queryPeruData.SocialReason = result.Result.SocialReason
		queryPeruData.Address = result.Result.Address
		queryPeruData.TaxpayerState = result.Result.TaxpayerState
		queryPeruData.DomicileCondition = result.Result.DomicileCondition
		queryPeruData.UtilIdentityDocumentTypeId = 4
	} else {
		return c.JSON(http.StatusOK, utilities.Response{Message: fmt.Sprintf("El %s no es un número de documento válido para la SUNAT o la RENIEC", queryPeruData.DocumentNumber)})
	}

	// Return response
	return c.JSON(http.StatusCreated, utilities.Response{
		Success: true,
		Data:    queryPeruData,
		Message: fmt.Sprintf("Consulta exitosa"),
	})
}

func queryRuc(ruc string) (resultRuc, error) {
	result := resultRuc{}
	token := "eyJ1c2VySWQiOjEsInVzZXJUb2tlbklkIjoiMSJ9.rthwrahfCQXPncvfd2t8ZJBp8AIwTPfSslKNCUenxm5MV0AWtggjyOkFsTRC8NbLFPop0bAXyoyCtR8TUI1GCUlvWADIyIHGGw3Wuss-ODvC6mLLvjrvXRarXrYzhGq2rtQHhECY_10kLXU1EDj5zcBHpYg4yX1iLPQ8jz-A4iM"
	payload := strings.NewReader(`{ "ruc": "` + ruc + `", "token": "` + token + `" }`)

	req, err := http.NewRequest("POST", "https://ruc.paulantezana.com/api/v1/ruc", payload)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Type", "application/json")

	// Send Query
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	// Read
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	// string to struct
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func queryDni(dni string) (resultDni, error) {
	result := resultDni{}
	token := "eyJ1c2VySWQiOjEsInVzZXJUb2tlbklkIjoiMSJ9.rthwrahfCQXPncvfd2t8ZJBp8AIwTPfSslKNCUenxm5MV0AWtggjyOkFsTRC8NbLFPop0bAXyoyCtR8TUI1GCUlvWADIyIHGGw3Wuss-ODvC6mLLvjrvXRarXrYzhGq2rtQHhECY_10kLXU1EDj5zcBHpYg4yX1iLPQ8jz-A4iM"
	payload := strings.NewReader(`{ "dni": "` + dni + `", "token": "` + token + `" }`)

	req, err := http.NewRequest("POST", "https://ruc.paulantezana.com/api/v1/dni", payload)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Type", "application/json")

	// Send Query
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	// Read
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	// string to struct
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
