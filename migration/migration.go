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
	config := provider.GetConfig()

	if config.Global.Develop == true {
		db.AutoMigrate(
			&models.UtilAdditionalLegendType{},
			&models.UtilAffectationIgvType{},
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

			&models.PaymentType{},

			&models.Sale{},
			&models.SaleItem{},
			&models.SaleSunat{},

			&models.Company{},
			&models.CompanyLocal{},
			&models.CompanySalePointSerie{},
			&models.CompanySalePoint{},
			&models.CompanyWareHouse{},
			&models.CompanySalePointDocumentAuth{},
			&models.SaleConf{},
			&models.PurchaseConf{},
			&models.IncomeType{},
			&models.Income{},
			&models.ExpenseType{},
			&models.Expense{},

			&models.User{},
			&models.UserForgot{},
			&models.UserRole{},
			&models.UserRoleAuthorization{},
			&models.UserSalePoint{},
			&models.UserLocalAuth{},
			&models.UserWareHouseAuth{},

			&models.Category{},
			&models.Product{},
			&models.ProductMedia{},

			&models.Purchase{},
			&models.PurchaseItem{},
			&models.Provider{},
			&models.Customer{},

			&models.Kardex{},

			&models.App{},
			&models.AppAuthorization{},
		)

		// Default data
		docType := models.UtilDocumentType{}
		if db.First(&docType).RowsAffected == 0 {
			// App Authorization
			db.Create(&models.AppAuthorization{Key: "operation", Title: "Operaciones", Icon: "trophy", Description: "Operaciones", Action: "menu"})
			db.Create(&models.AppAuthorization{Key: "process", Title: "Proceso", Icon: "control", Description: "Proceso", Action: "menu"})
			db.Create(&models.AppAuthorization{Key: "report", Title: "Reportes", Icon: "line-chart", Description: "Reportes", Action: "menu"})
			db.Create(&models.AppAuthorization{Key: "statistic", Title: "Estadisticas", Icon: "pie-chart", Description: "Estadisticas", Action: "menu"})
			db.Create(&models.AppAuthorization{Key: "setting", Title: "Configuración", Icon: "setting", Description: "Configuración", Action: "menu"})

			db.Create(&models.AppAuthorization{Key: "operation_new_sale", Title: "Nueva venta", Icon: "shopping", To: "/admin/operation/newSale", Description: "Venta", Action: "new", ParentId: 1})
			db.Create(&models.AppAuthorization{Key: "operation_new_purchase", Title: "Nueva compra", Icon: "shop", To: "/admin/operation/newPurchase", Description: "Compra", Action: "new", ParentId: 1})
			db.Create(&models.AppAuthorization{Key: "operation_customer", Title: "Clientes", Icon: "user", To: "/admin/operation/customer", Description: "Cliente", Action: "list", ParentId: 1})
			db.Create(&models.AppAuthorization{Key: "operation_provider", Title: "Proveedores", Icon: "user", To: "/admin/operation/provider", Description: "Proveedor", Action: "list", ParentId: 1})
			//db.Create(&models.AppAuthorization{Key: "operation_new_quotation", Title: "Nueva cotización", Icon: "plus", To: "/admin/operation/newQuotation", Description: "Cotización", Action: "new", ParentId: 1})
			//db.Create(&models.AppAuthorization{Key: "operation_new_order", Title: "Nueva Orden", Icon: "plus", To: "/admin/operation/newOrder", Description: "Orden", Action: "new", ParentId: 1})
			//db.Create(&models.AppAuthorization{Key: "operation_new_transfer", Title: "Nueva transferencia", Icon: "plus", To: "/admin/operation/newTransfer", Description: "Transfer", Action: "new", ParentId: 1})
			db.Create(&models.AppAuthorization{Key: "operation_product", Title: "Productos", Icon: "rest", To: "/admin/operation/product", Description: "Productos", Action: "list", ParentId: 1})
            db.Create(&models.AppAuthorization{Key: "operation_income", Title: "Ingresos", Icon: "vertical-align-top", To: "/admin/operation/income", Description: "Ingresos", Action: "List", ParentId: 1})
            db.Create(&models.AppAuthorization{Key: "operation_expense", Title: "Gastos", Icon: "vertical-align-bottom", To: "/admin/operation/expense", Description: "Gastos", Action: "List", ParentId: 1})

            db.Create(&models.AppAuthorization{Key: "process_import", Title: "Importar", Icon: "cloud-upload", To: "/admin/process/import", Description: "Importar", Action: "list", ParentId: 2})
            db.Create(&models.AppAuthorization{Key: "process_export", Title: "Exportar", Icon: "cloud-download", To: "/admin/process/export", Description: "Exportar", Action: "list", ParentId: 2})

            db.Create(&models.AppAuthorization{Key: "report_sale", Title: "Ventas", Icon: "bar-chart", To: "/admin/report/sale", Description: "Ventas", Action: "list", ParentId: 3})
			db.Create(&models.AppAuthorization{Key: "report_purchase", Title: "Compras", Icon: "bar-chart", To: "/admin/report/purchase", Description: "Compras", Action: "list", ParentId: 3})
			db.Create(&models.AppAuthorization{Key: "report_kardex", Title: "Kardex", Icon: "bar-chart", To: "/admin/report/kardex", Description: "Kardex", Action: "list", ParentId: 3})
			db.Create(&models.AppAuthorization{Key: "report_utility", Title: "Utilidad", Icon: "bar-chart", To: "/admin/report/utility", Description: "Utilidad", Action: "list", ParentId: 3})
			db.Create(&models.AppAuthorization{Key: "report_customer", Title: "Clientes", Icon: "bar-chart", To: "/admin/report/customer", Description: "Clientes", Action: "list", ParentId: 3})
			db.Create(&models.AppAuthorization{Key: "report_provider", Title: "Proveedores", Icon: "bar-chart", To: "/admin/report/provider", Description: "Proveedores", Action: "list", ParentId: 3})
            db.Create(&models.AppAuthorization{Key: "report_box", Title: "Caja", Icon: "bar-chart", To: "/admin/report/box", Description: "Caja", Action: "list", ParentId: 3})

			db.Create(&models.AppAuthorization{Key: "setting_company", Title: "Empresa", Icon: "bank", To: "/admin/setting/company", Description: "Empresa", Action: "list", ParentId: 5})
			db.Create(&models.AppAuthorization{Key: "setting_subsidiary", Title: "Sucursal", Icon: "home", To: "/admin/setting/local", Description: "Sucursal", Action: "list", ParentId: 5})
			db.Create(&models.AppAuthorization{Key: "setting_warehouse", Title: "Almacen", Icon: "hdd", To: "/admin/setting/wareHouse", Description: "Almacen", Action: "list", ParentId: 5})
			db.Create(&models.AppAuthorization{Key: "setting_sale_point", Title: "Punto de venta", Icon: "shop", To: "/admin/setting/salePoint", Description: "Punto de venta", Action: "list", ParentId: 5})
			db.Create(&models.AppAuthorization{Key: "setting_user_rol", Title: "Roles", Icon: "profile", To: "/admin/setting/role", Description: "Roles", Action: "list", ParentId: 5})
			db.Create(&models.AppAuthorization{Key: "setting_user", Title: "Usuarios", Icon: "user", To: "/admin/setting/user", Description: "Usuarios", Action: "list", ParentId: 5})
			db.Create(&models.AppAuthorization{Key: "setting_operative", Title: "Operatividad", Icon: "hdd", To: "/admin/setting/operative", Description: "Operatividad", Action: "list", ParentId: 5})
            db.Create(&models.AppAuthorization{Key: "setting_income_type", Title: "Tipo de ingreso", Icon: "vertical-align-top", To: "/admin/setting/incomeType", Description: "Tipo de ingreso", Action: "list", ParentId: 5})
            db.Create(&models.AppAuthorization{Key: "setting_expense_type", Title: "Tipo de gasto", Icon: "vertical-align-bottom", To: "/admin/setting/expenseType", Description: "Tipo de gasto", Action: "list", ParentId: 5})
            db.Create(&models.AppAuthorization{Key: "setting_category", Title: "Categoria", Icon: "deployment-unit", To: "/admin/setting/category", Description: "Categoria", Action: "list", ParentId: 5})

			// Document type
			db.Create(&models.UtilDocumentType{Code: "01", NuCode: "1", Description: "FACTURA"})
			db.Create(&models.UtilDocumentType{Code: "03", NuCode: "2", Description: "BOLETA DE VENTA"})
			db.Create(&models.UtilDocumentType{Code: "07", NuCode: "3", Description: "NOTA DE CREDITO"})
			db.Create(&models.UtilDocumentType{Code: "08", NuCode: "4", Description: "NOTA DE DEBITO"})
			db.Create(&models.UtilDocumentType{Code: "09", Description: "GUIA DE REMISIÓN REMITENTE"})
			db.Create(&models.UtilDocumentType{Code: "TK", Description: "TIKET", Sunat: false})
			db.Create(&models.UtilDocumentType{Code: "NP", Description: "NOTA PEDIDO", Sunat: false})
			db.Create(&models.UtilDocumentType{Code: "CT", Description: "COTIZACIÓN", Sunat: false})
			db.Create(&models.UtilDocumentType{Code: "OT", Description: "OTROS", Sunat: false})

			// Currency type
			db.Create(&models.UtilCurrencyType{Code: "PEN", NuCode: "1", Description: "SOLES", Symbol: "S/"})
			db.Create(&models.UtilCurrencyType{Code: "USD", NuCode: "2", Description: "DÓLARES", Symbol: "$"})
			db.Create(&models.UtilCurrencyType{Code: "EUR", NuCode: "3", Description: "EURO", Symbol: "€"})
			//db.Create(&models.UtilCurrencyType{Code: "JPY", Description: "YEN", Symbol: "¥"})

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
			db.Create(&models.UtilUnitMeasureType{Code: "NIU", NuCode: "NIU", Description: "UNIDAD(BIENES)"})
			db.Create(&models.UtilUnitMeasureType{Code: "ZZ", NuCode: "ZZ", Description: "UNIDAD(SERVICIOS)"})
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

			// UtilAffectation Igv Type
			db.Create(&models.UtilAffectationIgvType{Code: "10", NuCode: "1", Description: "Gravado - Operación Onerosa", Onerous: true, UtilTributeTypeId: 1})
			db.Create(&models.UtilAffectationIgvType{Code: "11", NuCode: "2", Description: "Gravado – Retiro por premio", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "12", NuCode: "3", Description: "Gravado – Retiro por donación", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "13", NuCode: "4", Description: "Gravado – Retiro", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "14", NuCode: "5", Description: "Gravado – Retiro por publicidad", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "15", NuCode: "6", Description: "Gravado – Bonificaciones", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "16", NuCode: "7", Description: "Gravado – Retiro por entrega a trabajadores", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "20", NuCode: "8", Description: "Exonerado - Operación Onerosa", Onerous: true, UtilTributeTypeId: 7})
			db.Create(&models.UtilAffectationIgvType{Code: "30", NuCode: "9", Description: "Inafecto - Operación Onerosa", Onerous: true, UtilTributeTypeId: 8})
			db.Create(&models.UtilAffectationIgvType{Code: "31", NuCode: "10", Description: "Inafecto – Retiro por Bonificación", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "32", NuCode: "11", Description: "Inafecto – Retiro", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "33", NuCode: "12", Description: "Inafecto – Retiro por Muestras Médicas", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "34", NuCode: "13", Description: "Inafecto - Retiro por Convenio Colectivo", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "35", NuCode: "14", Description: "Inafecto – Retiro por premio", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "36", NuCode: "15", Description: "Inafecto - Retiro por publicidad", Onerous: false, UtilTributeTypeId: 6})
			db.Create(&models.UtilAffectationIgvType{Code: "40", NuCode: "16", Description: "Exportación", Onerous: true, UtilTributeTypeId: 5})

			// Identity type
			db.Create(&models.UtilIdentityDocumentType{Code: "0", NuCode: "0", Description: "NO DOMICILIADO, SIN RUC (EXPORTACIÓN)"})
			db.Create(&models.UtilIdentityDocumentType{Code: "1", NuCode: "1", Description: "DNI"})
			db.Create(&models.UtilIdentityDocumentType{Code: "4", NuCode: "4", Description: "CARNET DE EXTRANJERIA"})
			db.Create(&models.UtilIdentityDocumentType{Code: "6", NuCode: "6", Description: "RUC"})
			db.Create(&models.UtilIdentityDocumentType{Code: "7", NuCode: "7", Description: "PASAPORTE"})
			db.Create(&models.UtilIdentityDocumentType{Code: "A", NuCode: "A", Description: "CED. DIPLOMATICA DE IDENTIDAD"})
			db.Create(&models.UtilIdentityDocumentType{Code: "B", NuCode: "", Description: "DOC.IDENT.PAIS.RESIDENCIA-NO.D", State: false})
			db.Create(&models.UtilIdentityDocumentType{Code: "C", NuCode: "", Description: "Tax Identification Number - TIN – Doc Trib PP.NN", State: false})
			db.Create(&models.UtilIdentityDocumentType{Code: "D", NuCode: "", Description: "Identification Number - IN – Doc Trib PP. JJ", State: false})
			db.Create(&models.UtilIdentityDocumentType{Code: "-", NuCode: "-", Description: "- VARIOS - VENTAS MENORES A S/.700.00 Y OTROS"})

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

			// Sale State
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

			// Types
			db.Create(&models.PaymentType{Code: "EF", Description: "Efectivo"})
			db.Create(&models.PaymentType{Code: "CH", Description: "Cheque"})
			db.Create(&models.PaymentType{Code: "CR", Description: "Crédito"})
			db.Create(&models.PaymentType{Code: "TR", Description: "Transferencia"})
			db.Create(&models.PaymentType{Code: "VA", Description: "Vales"})
			db.Create(&models.PaymentType{Code: "TA", Description: "Tarjeta"})
			db.Create(&models.PaymentType{Code: "AN", Description: "Anticipo"})

			// Init Value
			db.Create(&models.UtilGeographicalLocation{Code: "010101", District: "Chachapoyas", Province: "Chachapoyas", Department: "Amazonas"})
			db.Create(&models.UtilGeographicalLocation{Code: "010102", District: "Asuncion", Province: "Chachapoyas", Department: "Amazonas"})
		}

		// Company
		company := models.Company{}
		if db.First(&company).RowsAffected == 0 {
			// COMPANY
			company = models.Company{
				DocumentNumber:             "99999999999",
				SocialReason:               "ABC Company",
				CommercialReason:           "ABC",
				Email:                      "abc@gmail.com",
				UtilGeographicalLocationId: 1,
			}
			db.Create(&company)

			// COMPANY LOCAL
			companyLocal := models.CompanyLocal{
				Description:                "ss",
				SocialReason:               "LOCAL PRINCIPAL",
				CommercialReason:           "LOCAL",
				CompanyId:                  company.ID,
				UtilGeographicalLocationId: 1,
			}
			db.Create(&companyLocal)

			// WARE HOUSE
			db.Create(&models.CompanyWareHouse{CompanyLocalId: companyLocal.ID, CompanyId: company.ID, Description: "ALMACEN PRINCIPAL"})

			// SALE POINT
			companySalePoint := models.CompanySalePoint{
				CompanyLocalId: companyLocal.ID,
				CompanyId:      company.ID,
				Description:    "Punto de venta 1",
			}
			db.Create(&companySalePoint)

			// SERIE
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "FPP1", UtilDocumentTypeId: 1, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "FPP1", UtilDocumentTypeId: 3, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "FPP1", UtilDocumentTypeId: 4, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "BPP1", UtilDocumentTypeId: 2, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "BPP1", UtilDocumentTypeId: 3, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "BPP1", UtilDocumentTypeId: 4, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 1, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 3, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 4, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 2, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 3, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 4, Contingency: true})

			// Create Roles
			adminRole := models.UserRole{Description: "Administrador", CompanyId: company.ID}

			db.Create(&models.UserRole{Description: "Usuario", CompanyId: company.ID})
			db.Create(&adminRole)
			db.Create(&models.UserRole{Description: "Gerencia", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Contabilidad", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Ventas", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Operaciones", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Almacenero", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Ayudantes", CompanyId: company.ID})

			// Find AppAuthorizations
			appAuthorizations := make([]models.AppAuthorization, 0)
			db.Where("state = true").Find(&appAuthorizations)

			userRoles := make([]models.UserRole, 0)
			db.Where("state = true").Find(&userRoles)

			for _, authorization := range appAuthorizations {
				for _, role := range userRoles {
					db.Create(&models.UserRoleAuthorization{AppAuthorizationId: authorization.ID, UserRoleId: role.ID})
				}
			}

			cc := sha256.Sum256([]byte("admin1"))
			pwd := fmt.Sprintf("%x", cc)
			newUser := models.User{
				UserName:   "admin1",
				Password:   pwd,
				Freeze:     true,
				UserRoleId: adminRole.ID,
				CompanyId:  company.ID,
			}
			db.Create(&newUser)

			// Default Conf
			db.Create(&models.SaleConf{DocumentSize: "TICKET", CompanyId: company.ID})
			db.Create(&models.PurchaseConf{CompanyId: company.ID})

			/// --------------------------------------------------------------------------------------------------------------
			// COMPANY SECOND
			company = models.Company{
				DocumentNumber:             "99999999991",
				SocialReason:               "ABC Company - 2",
				CommercialReason:           "ABC - 2",
				Email:                      "abc2@gmail.com",
				UtilGeographicalLocationId: 1,
			}
			db.Create(&company)

			// COMPANY LOCAL
			companyLocal = models.CompanyLocal{
				Description:                "ss",
				SocialReason:               "LOCAL PRINCIPAL",
				CommercialReason:           "LOCAL",
				CompanyId:                  company.ID,
				UtilGeographicalLocationId: 1,
			}
			db.Create(&companyLocal)

			// WARE HOUSE
			db.Create(&models.CompanyWareHouse{CompanyLocalId: companyLocal.ID, CompanyId: company.ID, Description: "ALMACEN PRINCIPAL"})

			// SALE POINT
			companySalePoint = models.CompanySalePoint{
				CompanyLocalId: companyLocal.ID,
				CompanyId:      company.ID,
				Description:    "Punto de venta 1",
			}
			db.Create(&companySalePoint)

			// SERIE
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "FPP1", UtilDocumentTypeId: 1, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "FPP1", UtilDocumentTypeId: 3, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "FPP1", UtilDocumentTypeId: 4, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "BPP1", UtilDocumentTypeId: 2, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "BPP1", UtilDocumentTypeId: 3, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "BPP1", UtilDocumentTypeId: 4, Contingency: false})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 1, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 3, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 4, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 2, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 3, Contingency: true})
			db.Create(&models.CompanySalePointSerie{CompanySalePointId: companySalePoint.ID, Serie: "0001", UtilDocumentTypeId: 4, Contingency: true})

			// Create Roles
			adminRole = models.UserRole{Description: "Administrador", CompanyId: company.ID}

			db.Create(&models.UserRole{Description: "Usuario", CompanyId: company.ID})
			db.Create(&adminRole)
			db.Create(&models.UserRole{Description: "Gerencia", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Contabilidad", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Ventas", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Operaciones", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Almacenero", CompanyId: company.ID})
			db.Create(&models.UserRole{Description: "Ayudantes", CompanyId: company.ID})

			// Find AppAuthorizations
			appAuthorizations = make([]models.AppAuthorization, 0)
			db.Where("state = true").Find(&appAuthorizations)

			userRoles = make([]models.UserRole, 0)
			db.Where("state = true").Find(&userRoles)

			for _, authorization := range appAuthorizations {
				for _, role := range userRoles {
					db.Create(&models.UserRoleAuthorization{AppAuthorizationId: authorization.ID, UserRoleId: role.ID})
				}
			}

			cc = sha256.Sum256([]byte("admin2"))
			pwd = fmt.Sprintf("%x", cc)
			newUser = models.User{
				UserName:   "admin2",
				Password:   pwd,
				Freeze:     true,
				UserRoleId: adminRole.ID,
				CompanyId:  company.ID,
			}
			db.Create(&newUser)

			// Default Conf
			db.Create(&models.SaleConf{DocumentSize: "TICKET", CompanyId: company.ID})
			db.Create(&models.PurchaseConf{CompanyId: company.ID})
		}
	}
}
