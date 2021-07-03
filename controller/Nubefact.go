package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulantezana/shopping/utilities"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type NubefactDocument struct {
	Operacion                      string                 `json:"operacion"`
	TipoDeComprobante              string                 `json:"tipo_de_comprobante"`
	Serie                          string                 `json:"serie"`
	Numero                         string                 `json:"numero"`
	SunatTransaction               string                 `json:"sunat_transaction"`
	ClienteTipoDeDocumento         string                 `json:"cliente_tipo_de_documento"`
	ClienteNumeroDeDocumento       string                 `json:"cliente_numero_de_documento"`
	ClienteDenominacion            string                 `json:"cliente_denominacion"`
	ClienteDireccion               string                 `json:"cliente_direccion"`
	ClienteEmail                   string                 `json:"cliente_email"`
	ClienteEmail1                  string                 `json:"cliente_email_1"`
	ClienteEmail2                  string                 `json:"cliente_email_2"`
	FechaDeEmision                 time.Time              `json:"fecha_de_emision"`
	FechaDeVencimiento             time.Time              `json:"fecha_de_vencimiento"`
	Moneda                         string                 `json:"moneda"`
	TipoDeCambio                   string                 `json:"tipo_de_cambio"`
	PorcentajeDeIgv                string                 `json:"porcentaje_de_igv"`
	DescuentoGlobal                string                 `json:"descuento_global"`
	TotalDescuento                 string                 `json:"total_descuento"`
	TotalAnticipo                  float64                `json:"total_anticipo"`
	TotalGravada                   float64                `json:"total_gravada"`
	TotalInafecta                  float64                `json:"total_inafecta"`
	TotalExonerada                 float64                `json:"total_exonerada"`
	TotalIgv                       float64                `json:"total_igv"`
	TotalGratuita                  float64                `json:"total_gratuita"`
	TotalOtrosCargos               float64                `json:"total_otros_cargos"`
	Total                          float64                `json:"total"`
	PercepcionTipo                 string                 `json:"percepcion_tipo"`
	PercepcionBaseImponible        string                 `json:"percepcion_base_imponible"`
	TotalPercepcion                float64                `json:"total_percepcion"`
	TotalIncluidoPercepcion        float64                `json:"total_incluido_percepcion"`
	Detraccion                     string                 `json:"detraccion"`
	Observaciones                  string                 `json:"observaciones"`
	DocumentoQueSeModificaTipo     string                 `json:"documento_que_se_modifica_tipo"`
	DocumentoQueSeModificaSerie    string                 `json:"documento_que_se_modifica_serie"`
	DocumentoQueSeModificaNumero   string                 `json:"documento_que_se_modifica_numero"`
	TipoDeNotaDeCredito            string                 `json:"tipo_de_nota_de_credito"`
	TipoDeNotaDeDebito             string                 `json:"tipo_de_nota_de_debito"`
	EnviarAutomaticamenteALaSunat  bool                   `json:"enviar_automaticamente_a_la_sunat"`
	EnviarAutomaticamenteAlCliente bool                   `json:"enviar_automaticamente_al_cliente"`
	CodigoUnico                    string                 `json:"codigo_unico"`
	CondicionesDePago              string                 `json:"condiciones_de_pago"`
	MedioDePago                    string                 `json:"medio_de_pago"`
	PlacaVehiculo                  string                 `json:"placa_vehiculo"`
	OrdenCompraServicio            string                 `json:"orden_compra_servicio"`
	TablaPersonalizadaCodigo       string                 `json:"tabla_personalizada_codigo"`
	FormatoDePdf                   string                 `json:"formato_de_pdf"`
	Items                          []NubefactDocumentItem `json:"items" gorm:"-"`
}

type NubefactDocumentItem struct {
	UnidadDeMedida          string  `json:"unidad_de_medida"`
	Codigo                  string  `json:"codigo"`
	Descripcion             string  `json:"descripcion"`
	Cantidad                float64 `json:"cantidad"`
	ValorUnitario           float64 `json:"valor_unitario"`
	PrecioUnitario          float64 `json:"precio_unitario"`
	Descuento               float64 `json:"descuento"`
	Subtotal                float64 `json:"subtotal"`
	TipoDeIgv               string  `json:"tipo_de_igv"`
	Igv                     float64 `json:"igv"`
	Total                   float64 `json:"total"`
	AnticipoRegularizacion  bool    `json:"anticipo_regularizacion"`
	AnticipoDocumentoSerie  string  `json:"anticipo_documento_serie"`
	AnticipoDocumentoNumero string  `json:"anticipo_documento_numero"`
}

type NubefactDocumentResponse struct {
	TipoDeComprobante  string `json:"tipo_de_comprobante"`
	Serie              string `json:"serie"`
	Numero             string `json:"numero"`
	Enlace             string `json:"enlace"`
	EnlaceDelPdf       string `json:"enlace_del_pdf"`
	EnlaceDelXml       string `json:"enlace_del_xml"`
	EnlaceDelCdr       string `json:"enlace_del_cdr"`
	AceptadaPorSunat   string `json:"aceptada_por_sunat"`
	SunatDescription   string `json:"sunat_description"`
	SunatNote          string `json:"sunat_note"`
	SunatResponsecode  string `json:"sunat_responsecode"`
	SunatSoapError     string `json:"sunat_soap_error"`
	CadenaParaCodigoQr string `json:"cadena_para_codigo_qr"`
	CodigoHash         string `json:"codigo_hash"`
}

func NubefactSendDocument(doc NubefactDocument, u string, t string, production bool) (utilities.Response, error) {
	response := utilities.Response{}

	url := "https://api.nubefact.com/api/v1/dbe6eda2-8f69-4ec1-b5a0-4f729f68eb06"
	if production {
		url = u
	}
	token := "e2df3fa8520d412b888402e556c4b75635c4815c99994ccc83c6f06650683b58"
	if production {
		token = t
	}

	b, err := json.Marshal(doc)
	if err != nil {
		return response, err
	}
	jsonData := string(b)

	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))
	if err != nil {
		return response, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=\"%s\"", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return response, errors.New(res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	stringResponse := string(body)
	nDocumentResponse := NubefactDocumentResponse{}
	err = json.Unmarshal([]byte(stringResponse), &nDocumentResponse)
	if err != nil {
		return response, err
	}

	response.Success = true
	response.Data = nDocumentResponse
	return response, nil
}

//
//func NubefactQueryDocument(doc NubefactDocument)  (utilities.Response, error) {
//
//}
//
//func NubefactCanceledDocument(doc NubefactDocument)  (utilities.Response, error)  {
//
//}
//
//func NubefactQueryCanceledDocument(doc NubefactDocument)  (utilities.Response, error) {
//
//}
