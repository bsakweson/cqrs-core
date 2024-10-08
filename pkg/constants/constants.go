package constants

const (
	GrpcPort                   = "GRPC_PORT"
	HttpPort                   = "HTTP_PORT"
	ConfigPath                 = "CONFIG_PATH"
	KafkaBrokers               = "KAFKA_BROKERS"
	JaegerHostPort             = "JAEGER_HOST"
	RedisAddr                  = "REDIS_ADDR"
	MongoDbURI                 = "MONGO_URI"
	EventStoreConnectionString = "EVENT_STORE_CONNECTION_STRING"
	ElasticUrl                 = "ELASTIC_URL"

	ReaderServicePort = "READER_SERVICE"

	Yaml          = "yaml"
	Tcp           = "tcp"
	Redis         = "redis"
	Kafka         = "kafka"
	Postgres      = "postgres"
	MongoDB       = "mongo"
	ElasticSearch = "elasticSearch"

	GRPC     = "GRPC"
	SIZE     = "SIZE"
	URI      = "URI"
	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"

	Topic        = "topic"
	Partition    = "partition"
	Message      = "message"
	WorkerID     = "workerID"
	Offset       = "offset"
	Time         = "time"
	GroupName    = "GroupName"
	StreamID     = "StreamID"
	EventID      = "EventID"
	EventType    = "EventType"
	EventNumber  = "EventNumber"
	CreatedDate  = "CreatedDate"
	UserMetadata = "UserMetadata"

	Page   = "page"
	Size   = "size"
	Search = "search"
	ID     = "id"

	EsAll = "$all"

	Validate        = "validate"
	FieldValidation = "field validation"
	RequiredHeaders = "required header"
	Base64          = "base64"
	Unmarshal       = "unmarshal"
	Uuid            = "uuid"
	Cookie          = "cookie"
	Token           = "token"
	Bcrypt          = "bcrypt"
	SQLState        = "sqlstate"

	MongoProjection   = "(MongoDB Projection)"
	ElasticProjection = "(Elastic Projection)"

	ProductIdIndex = "product_id_index"
	ProductIndex   = "product_index"
	ProductId      = "product_id"
	VariantId      = "variantId"
	Category       = "category"
	Name           = "name"
	Description    = "description"
	Highlight      = "highlight"
	Price          = "price"
	Barcode        = "barcode"
	Sku            = "sku"
	SupplierId     = "supplier_id"
	ShopID         = "shop_id"
	LowStock       = "lowstock"
	IsComposite    = "isComposite"
	Cost           = "cost"
	Variants       = "variants"
	Product        = "product"
	Products       = "products"
	Currency       = "currency"
	Reason         = "reason"

	ErrProductIDIsRequired   = "product id is required"
	ErrCategoryIsRequired    = "category is required"
	ErrNameIsRequired        = "name is required"
	ErrDescriptionIsRequired = "description is required"
	ErrHighlightIsRequired   = "highlight is required"
	ErrPriceIsRequired       = "price is required"
	ErrBarcodeIsRequired     = "barcode is required"
	ErrSkuIsRequired         = "sku is required"
	ErrSupplierIdIsRequired  = "supplier id is required"
	ErrShopIDIsRequired      = "shop id is required"
	ErrLowStockIsRequired    = "low stock is required"
	ErrIsCompositeIsRequired = "is composite is required"
	ErrCostIsRequired        = "cost is required"
	ErrVariantsIsRequired    = "variants is required"
	ErrProductIsRequired     = "product is required"
	ErrProductsIsRequired    = "products is required"
	ErrCurrencyIsRequired    = "currency is required"
	ErrReasonIsRequired      = "reason is required"
)