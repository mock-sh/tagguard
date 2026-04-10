# Known Tag Keys

tagguard recognizes the following struct tag keys as valid. Tags not in this list will trigger an "unknown tag key" warning.

## Standard Library

| Key | Used By |
|-----|---------|
| `json` | `encoding/json` |
| `xml` | `encoding/xml` |

## Serialization

| Key | Used By |
|-----|---------|
| `yaml` | `gopkg.in/yaml.v3` |
| `toml` | `github.com/BurntSushi/toml` |
| `bson` | `go.mongodb.org/mongo-driver` |
| `msgpack` | `github.com/vmihailenco/msgpack` |
| `cbor` | CBOR encoding |
| `protobuf` | Protocol Buffers |
| `avro` | Apache Avro |

## Database / ORM

| Key | Used By |
|-----|---------|
| `db` | `github.com/jmoiron/sqlx` |
| `gorm` | `gorm.io/gorm` |
| `sql` | Various |
| `pg` | `github.com/go-pg/pg` |
| `bigquery` | `cloud.google.com/go/bigquery` |
| `dynamodbav` | AWS SDK Go |
| `firestore` | Google Cloud Firestore |
| `spanner` | Google Cloud Spanner |
| `datastore` | Google Cloud Datastore |

## Web Frameworks

| Key | Used By |
|-----|---------|
| `form` | gin, echo, fiber |
| `query` | URL query params |
| `uri` | gin URI params |
| `header` | HTTP header binding |
| `cookie` | cookie binding |
| `binding` | gin (alias for validate) |
| `param` | fiber path params |
| `path` | path params |

## Validation

| Key | Used By |
|-----|---------|
| `validate` | `github.com/go-playground/validator` |
| `valid` | `github.com/asaskevich/govalidator` |

## Configuration / Environment

| Key | Used By |
|-----|---------|
| `env` | `github.com/kelseyhightower/envconfig` |
| `envconfig` | envconfig |
| `mapstructure` | `github.com/mitchellh/mapstructure` (used by Viper) |
| `flag` | standard `flag` package |

## Other

| Key | Used By |
|-----|---------|
| `redis` | Redis clients |
| `csv` | CSV encoding |
| `xlsx` | Excel files |
| `default` | Default value libraries |
| `description` | Documentation generators |
| `example` | OpenAPI / Swagger |
| `swaggertype` | Swagger type override |

---

If you use a tag key not on this list that is widely used, please [open an issue](https://github.com/mokshg/tagguard/issues) to have it added.
