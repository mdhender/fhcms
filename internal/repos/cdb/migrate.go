/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
******************************************************************************/

package cdb

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

func (db *DB) CreateSchema(ctx context.Context, username, password, email string) error {
	if err := db.dropSchema(ctx); err != nil {
		return err
	}
	if err := db.createSchema(ctx); err != nil {
		return err
	}

	log.Printf("[cdb] creating owner: username %q password %q\n", username, password)

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	b := &pgx.Batch{}
	b.Queue("insert into fh_version(version) values(1)")
	b.Queue("insert into site(title, active) values('Far Horizons Reactor', false)")
	b.Queue("insert into accounts(id, email, username, password, admin, active) values(0, $1, $2, $3, true, true)", email, username, password)
	for batchResults := tx.SendBatch(ctx, b); ; {
		rows, err := batchResults.Query()
		rows.Close()
		if err != nil {
			if err.Error() != "no result" {
				log.Printf("[schema] batch.create.schema: %+v\n", err)
			}
			break
		}
	}

	return tx.Commit(ctx)
}

func (db *DB) createSchema(ctx context.Context) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	b := &pgx.Batch{}
	b.Queue(`CREATE TABLE fh_version(
		version INTEGER NOT NULL)`)
	b.Queue(`CREATE TABLE site(
		title       TEXT,
		description TEXT,
		domain_name TEXT,
		active      BOOLEAN NOT NULL DEFAULT TRUE)`)

	b.Queue("CREATE SEQUENCE seq_accounts")
	b.Queue(`CREATE TABLE accounts(
		id       INTEGER NOT NULL,
		email    TEXT    NOT NULL,
		username TEXT    NOT NULL,
		password TEXT    NOT NULL,
		notes    TEXT,
		admin    BOOLEAN NOT NULL DEFAULT FALSE,
		active   BOOLEAN NOT NULL DEFAULT FALSE,
		CONSTRAINT pk_accounts__id PRIMARY KEY (id))`)

	for batchResults := tx.SendBatch(ctx, b); ; {
		rows, err := batchResults.Query()
		rows.Close()
		if err != nil {
			if err.Error() != "no result" {
				log.Printf("[schema] batch.create.schema: %+v\n", err)
			}
			break
		}
	}

	return tx.Commit(ctx)
}

func (db *DB) dropSchema(ctx context.Context) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	b := &pgx.Batch{}
	for _, sequence := range []string{"seq_accounts"} {
		b.Queue("drop sequence if exists " + sequence)
	}
	for _, table := range []string{"accounts", "site", "fh_version"} {
		b.Queue("drop table if exists " + table)
	}
	for batchResults := tx.SendBatch(ctx, b); ; {
		rows, err := batchResults.Query()
		rows.Close()
		if err != nil {
			if err.Error() != "no result" {
				log.Printf("[schema] batch.drop.schema: %+v\n", err)
			}
			break
		}
	}

	return tx.Commit(ctx)
}
