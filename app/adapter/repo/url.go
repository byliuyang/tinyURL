package repo

import (
	"database/sql"
	"fmt"
	"short/app/adapter/repo/table"
	"short/app/entity"
	"short/app/usecase/repo"
)

type UrlSql struct {
	db *sql.DB
}

func (u UrlSql) IsAliasExist(alias string) (bool, error) {
	query := fmt.Sprintf(`
SELECT "%s" 
FROM "%s" 
WHERE "%s"=$1;`,
		table.Url.ColumnAlias,
		table.Url.TableName,
		table.Url.ColumnAlias,
	)

	err := u.db.QueryRow(query, alias).Scan(&alias)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, nil
	}
	return true, err
}

func (u *UrlSql) Create(url entity.Url) error {
	statement := fmt.Sprintf(`
INSERT INTO "%s" ("%s","%s","%s","%s","%s")
VALUES ($1, $2, $3, $4, $5);`,
		table.Url.TableName,
		table.Url.ColumnAlias,
		table.Url.ColumnOriginalUrl,
		table.Url.ColumnExpireAt,
		table.Url.ColumnCreatedAt,
		table.Url.ColumnUpdatedAt,
	)
	_, err := u.db.Exec(statement, url.Alias, url.OriginalUrl, url.ExpireAt, url.CreatedAt, url.UpdatedAt)
	return err
}

func (u UrlSql) GetByAlias(alias string) (entity.Url, error) {
	statement := fmt.Sprintf(`
SELECT "%s","%s","%s","%s","%s" 
FROM "%s" 
WHERE "%s"=$1;`,
		table.Url.ColumnAlias,
		table.Url.ColumnOriginalUrl,
		table.Url.ColumnExpireAt,
		table.Url.ColumnCreatedAt,
		table.Url.ColumnUpdatedAt,
		table.Url.TableName,
		table.Url.ColumnAlias,
	)

	row := u.db.QueryRow(statement, alias)

	url := entity.Url{}
	err := row.Scan(&url.Alias, &url.OriginalUrl, &url.ExpireAt, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		return entity.Url{}, err
	}

	return url, nil
}

func NewUrlSql(db *sql.DB) repo.Url {
	return &UrlSql{
		db: db,
	}
}
