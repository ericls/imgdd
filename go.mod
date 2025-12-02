module github.com/ericls/imgdd

go 1.24.0

toolchain go1.24.9

// License: AGPL-3.0-only

require (
	github.com/99designs/gqlgen v0.17.84
	github.com/aws/aws-sdk-go v1.55.8
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/ory/dockertest/v3 v3.12.0
	github.com/redis/go-redis/v9 v9.17.1
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.11.1
	github.com/urfave/cli/v2 v2.27.7
	github.com/vektah/gqlparser/v2 v2.5.31
	golang.org/x/time v0.14.0
)

require (
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.17.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.5.3 // indirect
	cloud.google.com/go/longrunning v0.7.0 // indirect
	cloud.google.com/go/monitoring v1.24.3 // indirect
	cloud.google.com/go/spanner v1.86.1 // indirect
	cloud.google.com/go/storage v1.57.2 // indirect
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.6.3 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.24 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.2 // indirect
	github.com/Azure/go-autorest/tracing v0.6.1 // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.30.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.54.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.54.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/apache/arrow-go/v18 v18.4.1 // indirect
	github.com/apache/thrift v0.22.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.40.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.20.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.92.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.2 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/golz4 v0.0.0-20240916140612-caecf3c00c06 // indirect
	github.com/cncf/xds/go v0.0.0-20251110193048-8bfbf64dc13e // indirect
	github.com/cockroachdb/cockroach-go/v2 v2.4.2 // indirect
	github.com/containerd/continuity v0.4.5 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/danieljoos/wincred v1.2.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v29.1.1+incompatible // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsoprea/go-exif/v2 v2.0.0-20230826092837-6579e82b732d // indirect
	github.com/dsoprea/go-iptc v0.0.0-20200610044640-bc9ca208b413 // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200610045659-121dd752914d // indirect
	github.com/dsoprea/go-utility v0.0.0-20221003172846-a3e1774ef349 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.8.0 // indirect
	github.com/edsrzf/mmap-go v1.2.0 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.36.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.0 // indirect
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/geo v0.0.0-20251125140653-09e2dd3603dd // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/flatbuffers v25.9.23+incompatible // indirect
	github.com/google/go-github/v39 v39.2.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.7 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgerrcode v0.0.0-20250907135507-afb5586c32a6 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/klauspost/asmfmt v1.3.2 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/ktrysmt/go-bitbucket v0.9.87 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/microsoft/go-mssqldb v1.9.4 // indirect
	github.com/minio/asm2plan9s v0.0.0-20200509001527-cdd76441f9d8 // indirect
	github.com/minio/c2goasm v0.0.0-20190812172519-36a3d3bbc4f3 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/moby/api v1.52.0 // indirect
	github.com/moby/moby/client v0.2.1 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/mutecomm/go-sqlcipher/v4 v4.4.2 // indirect
	github.com/nakagami/chacha20 v0.1.0 // indirect
	github.com/nakagami/firebirdsql v0.9.15 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/neo4j/neo4j-go-driver v1.8.3 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/opencontainers/runc v1.3.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rqlite/gorqlite v0.0.0-20250609141355-ac86a4a1c9a8 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/snowflakedb/gosnowflake v1.18.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	github.com/urfave/cli/v3 v3.6.1 // indirect
	github.com/xanzy/go-gitlab v0.115.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	gitlab.com/nyarla/go-crypt v0.0.0-20160106005555-d9a5dc2b789b // indirect
	go.mongodb.org/mongo-driver v1.17.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.38.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/exp v0.0.0-20251125195548-87e1e737ad39 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/oauth2 v0.33.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/telemetry v0.0.0-20251128220624-abf20d0e57ec // indirect
	golang.org/x/term v0.37.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	golang.org/x/tools/godoc v0.1.0-deprecated // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/api v0.256.0 // indirect
	google.golang.org/genproto v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/b v1.1.0 // indirect
	modernc.org/db v1.0.28 // indirect
	modernc.org/file v1.0.18 // indirect
	modernc.org/fileutil v1.3.40 // indirect
	modernc.org/golex v1.1.0 // indirect
	modernc.org/internal v1.1.9 // indirect
	modernc.org/libc v1.67.1 // indirect
	modernc.org/lldb v1.0.8 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/ql v1.4.29 // indirect
	modernc.org/sortutil v1.2.1 // indirect
	modernc.org/sqlite v1.40.1 // indirect
	modernc.org/strutil v1.2.1 // indirect
	modernc.org/zappy v1.1.0 // indirect
)

require (
	github.com/alicebob/miniredis/v2 v2.34.0
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20221012074422-4f3f7e934102
	github.com/dsoprea/go-png-image-structure v0.0.0-20210512210324-29b889a6093d
	github.com/felixge/httpsnoop v1.0.4
	github.com/go-jet/jet/v2 v2.14.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/sessions v1.4.0
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/rbcervilla/redisstore/v9 v9.0.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/studio-b12/gowebdav v0.11.0
	github.com/vikstrous/dataloadgen v0.0.10
	github.com/xrash/smetrics v0.0.0-20250705151800-55b8f293f342 // indirect
	golang.org/x/crypto v0.45.0
	golang.org/x/image v0.33.0
	golang.org/x/sys v0.38.0 // indirect
)

exclude (
	github.com/envoyproxy/go-control-plane v0.13.1
	github.com/envoyproxy/go-control-plane v0.13.2
)
