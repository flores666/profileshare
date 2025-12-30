module content

go 1.25

require (
	github.com/jackc/pgx/v5 v5.8.0
	github.com/joho/godotenv v1.5.1
)

require github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/flores666/profileshare/tree/master/lib/api v0.0.0
	github.com/flores666/profileshare/tree/master/lib/config v0.0.0
	github.com/flores666/profileshare/tree/master/lib/logger v0.0.0
	github.com/flores666/profileshare/tree/master/lib/utils v0.0.0
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/render v1.0.3
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jmoiron/sqlx v1.4.0
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace (
	github.com/flores666/profileshare/tree/master/lib/api v0.0.0 => ../../profileshare/lib/api
	github.com/flores666/profileshare/tree/master/lib/config v0.0.0 => ../../profileshare/lib/config
	github.com/flores666/profileshare/tree/master/lib/logger v0.0.0 => ../../profileshare/lib/logger
	github.com/flores666/profileshare/tree/master/lib/utils v0.0.0 => ../../profileshare/lib/utils
)
