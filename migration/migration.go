package migration

import (
	"crypto/sha256"
	"fmt"

	"github.com/paulantezana/shopping/models"
	"github.com/paulantezana/shopping/provider"
)

// Migrate function
func Migrate() {
	db := provider.GetConnection()
	defer db.Close()

	db.AutoMigrate(
		&models.UtilAdditionalLegendType{},
		&models.UtilCatAffectationIgvType{},
		&models.UtilCreditDebitType{},
		&models.UtilCurrencyType{},
		&models.UtilDocumentType{},
		&models.UtilGeographicalLocation{},
		&models.UtilIdentityDocumentType{},
		&models.UtilOperationType{},
		&models.UtilPerceptionType{},
		&models.UtilProductType{},
		&models.UtilSubjectDetractionType{},
		&models.UtilSystemIscType{},
		&models.UtilTransferReasonType{},
		&models.UtilTransportModeType{},
		&models.UtilTributeType{},
		&models.UtilUnitMeasureType{},

		&models.Invoice{},
		&models.InvoiceItem{},
		&models.InvoiceCustomer{},
		&models.InvoiceSunat{},
		&models.InvoiceCreditDebit{},

		&models.Company{},
		&models.CompanyLocal{},
		&models.CompanySerie{},
		&models.CompanySalePoint{},
		&models.CompanyWareHouse{},
		&models.CompanySalePointDocumentAuth{},

		&models.User{},
		&models.UserRole{},
		&models.UserRoleAuthorization{},

		&models.App{},
		&models.AppAuthorization{},
	)

	// Default data
	docType := models.UtilDocumentType{}
	db.First(&docType)
	if docType.ID == 0 {
		// Document type
		db.Create(&models.UtilDocumentType{Code: "01", Description: "FACTURA"})
		db.Create(&models.UtilDocumentType{Code: "03", Description: "BOLETA DE VENTA"})
		db.Create(&models.UtilDocumentType{Code: "07", Description: "NOTA DE CREDITO"})
		db.Create(&models.UtilDocumentType{Code: "08", Description: "NOTA DE DEBITO"})
		db.Create(&models.UtilDocumentType{Code: "09", Description: "GUIA DE REMISIÓN REMITENTE"})

		// Currency type
		db.Create(&models.UtilCurrencyType{Code: "PEN", Description: "SOLES", Symbol: "S/"})
		db.Create(&models.UtilCurrencyType{Code: "USD", Description: "DÓLARES AMERICANOS", Symbol: "$"})
		db.Create(&models.UtilCurrencyType{Code: "EUR", Description: "EURO", Symbol: "€"})
		db.Create(&models.UtilCurrencyType{Code: "JPY", Description: "YEN", Symbol: "¥"})

		// Unit Measure type
		db.Create(&models.UtilUnitMeasureType{Code: "4A", Description: "BOBINAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "BJ", Description: "BALDE"})
		db.Create(&models.UtilUnitMeasureType{Code: "BLL", Description: "BARRILES"})
		db.Create(&models.UtilUnitMeasureType{Code: "BG", Description: "BOLSA"})
		db.Create(&models.UtilUnitMeasureType{Code: "BO", Description: "BOTELLAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "BX", Description: "CAJA"})
		db.Create(&models.UtilUnitMeasureType{Code: "CT", Description: "CARTONES"})
		db.Create(&models.UtilUnitMeasureType{Code: "CMK", Description: "CENTIMETROCUADRADO"})
		db.Create(&models.UtilUnitMeasureType{Code: "CMQ", Description: "CENTIMETROCUBICO"})
		db.Create(&models.UtilUnitMeasureType{Code: "CMT", Description: "CENTIMETROLINEAL"})
		db.Create(&models.UtilUnitMeasureType{Code: "CEN", Description: "CIENTODEUNIDADES"})
		db.Create(&models.UtilUnitMeasureType{Code: "CY", Description: "CILINDRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "CJ", Description: "CONOS"})
		db.Create(&models.UtilUnitMeasureType{Code: "DZN", Description: "DOCENA"})
		db.Create(&models.UtilUnitMeasureType{Code: "DZP", Description: "DOCENAPOR10**6"})
		db.Create(&models.UtilUnitMeasureType{Code: "BE", Description: "FARDO"})
		db.Create(&models.UtilUnitMeasureType{Code: "GLI", Description: "GALONINGLES(4,545956L)"})
		db.Create(&models.UtilUnitMeasureType{Code: "GRM", Description: "GRAMO"})
		db.Create(&models.UtilUnitMeasureType{Code: "GRO", Description: "GRUESA"})
		db.Create(&models.UtilUnitMeasureType{Code: "HLT", Description: "HECTOLITRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "LEF", Description: "HOJA"})
		db.Create(&models.UtilUnitMeasureType{Code: "SET", Description: "JUEGO"})
		db.Create(&models.UtilUnitMeasureType{Code: "KGM", Description: "KILOGRAMO"})
		db.Create(&models.UtilUnitMeasureType{Code: "KTM", Description: "KILOMETRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "KWH", Description: "KILOVATIOHORA"})
		db.Create(&models.UtilUnitMeasureType{Code: "KT", Description: "KIT"})
		db.Create(&models.UtilUnitMeasureType{Code: "CA", Description: "LATAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "LBR", Description: "LIBRAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "LTR", Description: "LITRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MWH", Description: "MEGAWATTHORA"})
		db.Create(&models.UtilUnitMeasureType{Code: "MTR", Description: "METRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MTK", Description: "METROCUADRADO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MTQ", Description: "METROCUBICO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MGM", Description: "MILIGRAMOS"})
		db.Create(&models.UtilUnitMeasureType{Code: "MLT", Description: "MILILITRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MMT", Description: "MILIMETRO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MMK", Description: "MILIMETROCUADRADO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MMQ", Description: "MILIMETROCUBICO"})
		db.Create(&models.UtilUnitMeasureType{Code: "MLL", Description: "MILLARES"})
		db.Create(&models.UtilUnitMeasureType{Code: "UM", Description: "MILLONDEUNIDADES"})
		db.Create(&models.UtilUnitMeasureType{Code: "ONZ", Description: "ONZAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "PF", Description: "PALETAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "PK", Description: "PAQUETE"})
		db.Create(&models.UtilUnitMeasureType{Code: "PR", Description: "PAR"})
		db.Create(&models.UtilUnitMeasureType{Code: "FOT", Description: "PIES"})
		db.Create(&models.UtilUnitMeasureType{Code: "FTK", Description: "PIESCUADRADOS"})
		db.Create(&models.UtilUnitMeasureType{Code: "FTQ", Description: "PIESCUBICOS"})
		db.Create(&models.UtilUnitMeasureType{Code: "C62", Description: "PIEZAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "PG", Description: "PLACAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "ST", Description: "PLIEGO"})
		db.Create(&models.UtilUnitMeasureType{Code: "INH", Description: "PULGADAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "RM", Description: "RESMA"})
		db.Create(&models.UtilUnitMeasureType{Code: "DR", Description: "TAMBOR"})
		db.Create(&models.UtilUnitMeasureType{Code: "STN", Description: "TONELADACORTA"})
		db.Create(&models.UtilUnitMeasureType{Code: "LTN", Description: "TONELADALARGA"})
		db.Create(&models.UtilUnitMeasureType{Code: "TNE", Description: "TONELADAS"})
		db.Create(&models.UtilUnitMeasureType{Code: "TU", Description: "TUBOS"})
		db.Create(&models.UtilUnitMeasureType{Code: "NIU", Description: "UNIDAD(BIENES)"})
		db.Create(&models.UtilUnitMeasureType{Code: "ZZ", Description: "UNIDAD(SERVICIOS)"})
		db.Create(&models.UtilUnitMeasureType{Code: "GLL", Description: "USGALON(3,7843L)"})
		db.Create(&models.UtilUnitMeasureType{Code: "YRD", Description: "YARDA"})
		db.Create(&models.UtilUnitMeasureType{Code: "YDK", Description: "YARDACUADRADA"})

		// Tribute type
		db.Create(&models.UtilTributeType{Code: "1000", Description: "IGV Impuesto General a las Ventas", InternationalCode: "VAT", Name: "IGV"})
		db.Create(&models.UtilTributeType{Code: "1016", Description: "Impuesto a la Venta Arroz Pilado", InternationalCode: "VAT", Name: "IVAP"})
		db.Create(&models.UtilTributeType{Code: "2000", Description: "ISC Impuesto Selectivo al Consumo", InternationalCode: "EXC", Name: "ISC"})
		db.Create(&models.UtilTributeType{Code: "7152", Description: "Impuesto a la bolsa plastica", InternationalCode: "OTH", Name: "ICBPER"})
		db.Create(&models.UtilTributeType{Code: "9995", Description: "Exportación", InternationalCode: "FRE", Name: "EXP"})
		db.Create(&models.UtilTributeType{Code: "9996", Description: "Gratuito", InternationalCode: "FRE", Name: "GRA"})
		db.Create(&models.UtilTributeType{Code: "9997", Description: "Exonerado", InternationalCode: "VAT", Name: "EXO"})
		db.Create(&models.UtilTributeType{Code: "9998", Description: "Inafecto", InternationalCode: "FRE", Name: "INA"})
		db.Create(&models.UtilTributeType{Code: "9999", Description: "Otros tributos", InternationalCode: "OTH", Name: "OTROS"})

		// Identity type
		db.Create(&models.UtilIdentityDocumentType{Code: "0", Description: "0 NO DOMICILIADO, SIN RUC (EXPORTACIÓN)"})
		db.Create(&models.UtilIdentityDocumentType{Code: "1", Description: "1 DNI - DOC. NACIONAL DE IDENTIDAD"})
		db.Create(&models.UtilIdentityDocumentType{Code: "4", Description: "4 CARNET DE EXTRANJERIA"})
		db.Create(&models.UtilIdentityDocumentType{Code: "6", Description: "6 RUC - REG. UNICO DE CONTRIBUYENTES"})
		db.Create(&models.UtilIdentityDocumentType{Code: "7", Description: "7 PASAPORTE"})
		db.Create(&models.UtilIdentityDocumentType{Code: "A", Description: "A CED. DIPLOMATICA DE IDENTIDAD"})
		db.Create(&models.UtilIdentityDocumentType{Code: "B", Description: "B DOC.IDENT.PAIS.RESIDENCIA-NO.D"})
		db.Create(&models.UtilIdentityDocumentType{Code: "C", Description: "C Tax Identification Number - TIN – Doc Trib PP.NN"})
		db.Create(&models.UtilIdentityDocumentType{Code: "D", Description: "D Identification Number - IN – Doc Trib PP. JJ"})
		db.Create(&models.UtilIdentityDocumentType{Code: "-", Description: "- VARIOS - VENTAS MENORES A S/.700.00 Y OTROS"})

		// Isc type
		db.Create(&models.UtilSystemIscType{Code: "01", Description: "Sistema al valor (Apéndice IV, lit. A – T.U.O IGV e ISC)"})
		db.Create(&models.UtilSystemIscType{Code: "02", Description: "Aplicación del Monto Fijo ( Sistema específico, bienes en el apéndice III, Apéndice IV, lit. B – T.U.O IGV e ISC)"})
		db.Create(&models.UtilSystemIscType{Code: "03", Description: "Sistema de Precios de Venta al Público (Apéndice IV, lit. C – T.U.O IGV e ISC)"})

		// Credit And Debit Type
		db.Create(&models.UtilCreditDebitType{Code: "01", Description: "Anulación de la operación", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "02", Description: "Anulación por error en el RUC", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "03", Description: "Corrección por error en la descripción", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "04", Description: "Descuento global", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "05", Description: "Descuento por ítem", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "06", Description: "Devolución total", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "07", Description: "Devolución por ítem", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "08", Description: "Bonificación", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "09", Description: "Disminución en el valor", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "10", Description: "Otros Conceptos", UtilDocumentTypeId: 3})
		db.Create(&models.UtilCreditDebitType{Code: "01", Description: "Intereses por mora", UtilDocumentTypeId: 4})
		db.Create(&models.UtilCreditDebitType{Code: "02", Description: "Aumento en el valor", UtilDocumentTypeId: 4})
		db.Create(&models.UtilCreditDebitType{Code: "03", Description: "Penalidades/ otros conceptos", UtilDocumentTypeId: 4})

		// Additional Legend
		db.Create(&models.UtilAdditionalLegendType{Code: "1000", Description: "Monto en Letras"})
		db.Create(&models.UtilAdditionalLegendType{Code: "1002", Description: "TRANSFERENCIA GRATUITA DE UN BIEN Y/O SERVICIO PRESTADO GRATUITAMENTE"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2000", Description: "COMPROBANTE DE PERCEPCIÓN”"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2001", Description: "BIENES TRANSFERIDOS EN LA AMAZONÍA REGIÓN SELVAPARA SER CONSUMIDOS EN LA MISMA"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2002", Description: "SERVICIOS PRESTADOS EN LA AMAZONÍA  REGIÓN SELVA PARA SER CONSUMIDOS EN LA MISMA"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2003", Description: "CONTRATOS DE CONSTRUCCIÓN EJECUTADOS  EN LA AMAZONÍA REGIÓN SELVA"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2004", Description: "Agencia de Viaje - Paquete turístico"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2005", Description: "Venta realizada por emisor itinerante"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2006", Description: "Operación sujeta a detracción"})
		db.Create(&models.UtilAdditionalLegendType{Code: "2007", Description: "Operación sujeta a IVAP"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3000", Description: "CODIGO DE BB Y SS SUJETOS A DETRACCION"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3001", Description: "NUMERO DE CTA EN EL BN"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3002", Description: "Recursos Hidrobiológicos-Nombre y matrícula de la embarcación"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3003", Description: "Recursos Hidrobiológicos-Tipo y cantidad de especie vendida"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3004", Description: "Recursos Hidrobiológicos -Lugar de descarga"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3005", Description: "Recursos Hidrobiológicos -Fecha de descarga"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3006", Description: "Transporte Bienes vía terrestre – Numero Registro MTC"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3007", Description: "Transporte Bienes vía terrestre – configuración vehicular"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3008", Description: "Transporte Bienes vía terrestre – punto de origen"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3009", Description: "Transporte Bienes vía terrestre – punto destino"})
		db.Create(&models.UtilAdditionalLegendType{Code: "3010", Description: "Transporte Bienes vía terrestre – valor referencial preliminar"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4000", Description: "Código País de emisión del pasaporte"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4001", Description: "Código País de residencia del sujeto no domiciliado"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4002", Description: "Fecha de ingreso al país"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4003", Description: "Fecha de ingreso al establecimiento"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4004", Description: "Fecha de salida del establecimiento"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4005", Description: "Número de días de permanencia"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4006", Description: "Fecha de consumo"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4007", Description: "Paquete turístico - Nombres y Apellidos del Huésped"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4008", Description: "Paquete turístico – Tipo documento identidad del huésped"})
		db.Create(&models.UtilAdditionalLegendType{Code: "4009", Description: "Paquete turístico – Numero de documento identidad de huésped"})
		db.Create(&models.UtilAdditionalLegendType{Code: "5000", Description: "Número de Expediente"})
		db.Create(&models.UtilAdditionalLegendType{Code: "5001", Description: "Código de unidad ejecutora"})
		db.Create(&models.UtilAdditionalLegendType{Code: "5002", Description: "N° de proceso de selección"})
		db.Create(&models.UtilAdditionalLegendType{Code: "5003", Description: "N° de contrato"})
		db.Create(&models.UtilAdditionalLegendType{Code: "6000", Description: "Comercialización de Oro Código Unico Concesión Minera"})
		db.Create(&models.UtilAdditionalLegendType{Code: "6001", Description: "Comercialización de Oro N° declaración compromiso"})
		db.Create(&models.UtilAdditionalLegendType{Code: "6002", Description: "Comercialización de Oro N° Reg. Especial .Comerci. Oro"})
		db.Create(&models.UtilAdditionalLegendType{Code: "6003", Description: "Comercialización de Oro N° Resolución que autoriza Planta de Beneficio"})
		db.Create(&models.UtilAdditionalLegendType{Code: "6004", Description: "Comercialización de Oro Ley Mineral (% concent. oro)"})
		db.Create(&models.UtilAdditionalLegendType{Code: "7000", Description: "Primera venta de mercancia identificable entre usuarios de la zona comercial"})
		db.Create(&models.UtilAdditionalLegendType{Code: "7001", Description: "Venta exonerada del IGV-ISC-IPM. Prohibida la venta fuera de la zona comercial de Tacna"})

		// Operation Type
		db.Create(&models.UtilOperationType{Code: "0101", Description: "Venta lnterna"})
		db.Create(&models.UtilOperationType{Code: "0200", Description: "Exportación de Bienes"})
		db.Create(&models.UtilOperationType{Code: "0401", Description: "Ventas no domiciliados que no califican como exportación"})
		db.Create(&models.UtilOperationType{Code: "1001", Description: "Operación Sujeta a Detracción"})
		db.Create(&models.UtilOperationType{Code: "2001", Description: "Operación Sujeta a Percepción"})
		db.Create(&models.UtilOperationType{Code: "1004", Description: "Operación Sujeta a Detracción- Servicios de Transporte Carga"})

		// Transport Mode type
		db.Create(&models.UtilTransportModeType{Code: "01", Description: "Transporte público"})
		db.Create(&models.UtilTransportModeType{Code: "02", Description: "Transporte privado"})

		// Summary state code
		//db.Create(&models.UtilS{Code: "02", Description: "Transporte privado"})

		// Transfer Reason type
		db.Create(&models.UtilTransferReasonType{Code: "01", Description: "Venta"})
		db.Create(&models.UtilTransferReasonType{Code: "02", Description: "Compra"})
		db.Create(&models.UtilTransferReasonType{Code: "04", Description: "Traslado entre establecimientos de la misma empresa"})
		db.Create(&models.UtilTransferReasonType{Code: "08", Description: "Importación"})
		db.Create(&models.UtilTransferReasonType{Code: "09", Description: "Exportación"})
		db.Create(&models.UtilTransferReasonType{Code: "13", Description: "Otros"})
		db.Create(&models.UtilTransferReasonType{Code: "14", Description: "Venta sujeta a confirmación del comprador"})
		db.Create(&models.UtilTransferReasonType{Code: "18", Description: "Traslado emisor itinerante CP"})
		db.Create(&models.UtilTransferReasonType{Code: "19", Description: "Traslado a zona primaria"})

		// Perception type
		db.Create(&models.UtilPerceptionType{Code: "01", Description: "Percepción Venta Interna", Percentage: 2})
		db.Create(&models.UtilPerceptionType{Code: "02", Description: "Percepción a la adquisición de combustible", Percentage: 1})
		db.Create(&models.UtilPerceptionType{Code: "03", Description: "Percepción realizada al agente de percepción con tasa especial", Percentage: 0.5})

		// Invoice State
		//db.Create(&models.UtilInacacas{Code: "01", Description: "Percepción Venta Interna", Percentage: 2})

		// Subject Detraction
		db.Create(&models.UtilSubjectDetractionType{Code: "001", Description: "Azúcar y melaza de caña"})
		db.Create(&models.UtilSubjectDetractionType{Code: "002", Description: "Arroz"})
		db.Create(&models.UtilSubjectDetractionType{Code: "003", Description: "Alcohol etílico"})
		db.Create(&models.UtilSubjectDetractionType{Code: "004", Description: "Recursos hidrobiológicos"})
		db.Create(&models.UtilSubjectDetractionType{Code: "005", Description: "Maíz amarillo duro"})
		db.Create(&models.UtilSubjectDetractionType{Code: "007", Description: "Caña de azúcar"})
		db.Create(&models.UtilSubjectDetractionType{Code: "008", Description: "Madera"})
		db.Create(&models.UtilSubjectDetractionType{Code: "009", Description: "Arena y piedra."})
		db.Create(&models.UtilSubjectDetractionType{Code: "010", Description: "Residuos, subproductos, desechos, recortes y desperdicios"})
		db.Create(&models.UtilSubjectDetractionType{Code: "011", Description: "Bienes gravados con el IGV, o renuncia a la exoneración"})
		db.Create(&models.UtilSubjectDetractionType{Code: "012", Description: "Intermediación laboral y tercerización"})
		db.Create(&models.UtilSubjectDetractionType{Code: "013", Description: "Animales vivos"})
		db.Create(&models.UtilSubjectDetractionType{Code: "014", Description: "Carnes y despojos comestibles"})
		db.Create(&models.UtilSubjectDetractionType{Code: "015", Description: "Abonos, cueros y pieles de origen animal"})
		db.Create(&models.UtilSubjectDetractionType{Code: "016", Description: "Aceite de pescado"})
		db.Create(&models.UtilSubjectDetractionType{Code: "017", Description: "Harina, polvo y “pellets” de pescado, crustáceos, moluscos y demás invertebrados acuáticos"})
		db.Create(&models.UtilSubjectDetractionType{Code: "019", Description: "Arrendamiento de bienes muebles"})
		db.Create(&models.UtilSubjectDetractionType{Code: "020", Description: "Mantenimiento y reparación de bienes muebles"})
		db.Create(&models.UtilSubjectDetractionType{Code: "021", Description: "Movimiento de carga"})
		db.Create(&models.UtilSubjectDetractionType{Code: "022", Description: "Otros servicios empresariales"})
		db.Create(&models.UtilSubjectDetractionType{Code: "023", Description: "Leche"})
		db.Create(&models.UtilSubjectDetractionType{Code: "024", Description: "Comisión mercantil"})
		db.Create(&models.UtilSubjectDetractionType{Code: "025", Description: "Fabricación de bienes por encargo"})
		db.Create(&models.UtilSubjectDetractionType{Code: "026", Description: "Servicio de transporte de personas"})
		db.Create(&models.UtilSubjectDetractionType{Code: "027", Description: "Servicio de transporte de carga"})
		db.Create(&models.UtilSubjectDetractionType{Code: "028", Description: "Transporte de pasajeros"})
		db.Create(&models.UtilSubjectDetractionType{Code: "030", Description: "Contratos de construcción"})
		db.Create(&models.UtilSubjectDetractionType{Code: "031", Description: "Oro gravado con el IGV"})
		db.Create(&models.UtilSubjectDetractionType{Code: "032", Description: "Paprika y otros frutos de los generos capsicum o pimienta"})
		db.Create(&models.UtilSubjectDetractionType{Code: "034", Description: "Minerales metálicos no auríferos"})
		db.Create(&models.UtilSubjectDetractionType{Code: "035", Description: "Bienes exonerados del IGV"})
		db.Create(&models.UtilSubjectDetractionType{Code: "036", Description: "Oro y demás minerales metálicos exonerados del IGV"})
		db.Create(&models.UtilSubjectDetractionType{Code: "037", Description: "Demás servicios gravados con el IGV"})
		db.Create(&models.UtilSubjectDetractionType{Code: "039", Description: "Minerales no metálicos"})
		db.Create(&models.UtilSubjectDetractionType{Code: "040", Description: "Bien inmueble gravado con IGV"})
		db.Create(&models.UtilSubjectDetractionType{Code: "041", Description: "Plomo"})
		db.Create(&models.UtilSubjectDetractionType{Code: "099", Description: "Ley 30737"})

		// Init Value
		db.Create(&models.UtilGeographicalLocation{Code: "010101", District: "Chachapoyas", Province: "Chachapoyas", Department: "Amazonas"})
		db.Create(&models.UtilGeographicalLocation{Code: "010102", District: "Asuncion", Province: "Chachapoyas", Department: "Amazonas"})
	}

	// Company
	company := models.Company{}
	db.First(&company)
	if company.ID == 0 {
		company = models.Company{
			DocumentNumber:             "99999999999",
			UtilGeographicalLocationId: 1,
		}
		db.Create(&company)
	}

	// App Authorization
    appAuthorization := models.AppAuthorization{}
    db.First(&appAuthorization)
    if appAuthorization.ID == 0 {
        db.Create(&models.AppAuthorization{Key: "sale", Description: "Venta", Action: "List" })
        db.Create(&models.AppAuthorization{Key: "purchase", Description: "Compra", Action: "List" })
        db.Create(&models.AppAuthorization{Key: "box", Description: "Caja", Action: "List" })
        db.Create(&models.AppAuthorization{Key: "inventory", Description: "Inventario", Action: "List" })
        db.Create(&models.AppAuthorization{Key: "maintenance", Description: "Mantenimiento", Action: "List" })
        db.Create(&models.AppAuthorization{Key: "setting", Description: "Configuracion", Action: "List" })
        db.Create(&models.AppAuthorization{Key: "setting_company", Description: "Configuracion", Action: "List", ParentId: 6 })
        db.Create(&models.AppAuthorization{Key: "setting_subsidiary", Description: "Configuracion", Action: "List", ParentId: 6 })
        db.Create(&models.AppAuthorization{Key: "setting_warehouse", Description: "Configuracion", Action: "List", ParentId: 6 })
        db.Create(&models.AppAuthorization{Key: "setting_sale_point", Description: "Configuracion", Action: "List", ParentId: 6 })
        db.Create(&models.AppAuthorization{Key: "setting_user_rol", Description: "Configuracion", Action: "List", ParentId: 6 })
        db.Create(&models.AppAuthorization{Key: "setting_user", Description: "Configuracion", Action: "List", ParentId: 6 })
    }

    // UserRole
    userRole := models.UserRole{}
    db.First(&userRole)
    if userRole.ID == 0 {
        db.Create(&models.UserRole{Description: "Administrador"})
        db.Create(&models.UserRole{Description: "Usuario"})
    }

    // User
    user := models.User{}
    db.First(&user)
    if user.ID == 0 {
        cc := sha256.Sum256([]byte("admin1"))
        pwd := fmt.Sprintf("%x", cc)

        newUser := models.User{
            UserName: "admin1",
            Password: pwd,
            Freeze:   true,
            UserRoleId: 1,
        }
        db.Create(&newUser)
    }
}
