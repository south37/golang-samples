// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Hello world is a sample program demonstrating use of the Bigtable client
// library to perform basic CRUD operations
package main

// [START bigtable_hw_imports]
import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	_ "sync"
	"time"

	"cloud.google.com/go/bigtable"
	_ "google.golang.org/api/option"
)

// [END bigtable_hw_imports]

// User-provided constants.
const (
	tableName        = "Hello-Bigtable"
	columnFamilyName = "cf1"
	columnName       = "greeting"
)

var greetings = []string{"Hello World!", "Hello Cloud Bigtable!", "Hello golang!"}

// sliceContains reports whether the provided string is present in the given slice of strings.
func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func main() {
	project := flag.String("project", "", "The Google Cloud Platform project ID. Required.")
	instance := flag.String("instance", "", "The Google Cloud Bigtable instance ID. Required.")
	flag.Parse()

	for _, f := range []string{"project", "instance"} {
		if flag.Lookup(f).Value.String() == "" {
			log.Fatalf("The %s flag is required.", f)
		}
	}

	ctx := context.Background()

	// Set up admin client, tables, and column families.
	// NewAdminClient uses Application Default Credentials to authenticate.
	// [START bigtable_hw_connect]
	adminClient, err := bigtable.NewAdminClient(ctx, *project, *instance)
	if err != nil {
		log.Fatalf("Could not create admin client: %v", err)
	}
	// [END bigtable_hw_connect]

	// [START bigtable_hw_create_table]
	tables, err := adminClient.Tables(ctx)
	if err != nil {
		log.Fatalf("Could not fetch table list: %v", err)
	}

	if !sliceContains(tables, tableName) {
		log.Printf("Creating table %s", tableName)
		if err := adminClient.CreateTable(ctx, tableName); err != nil {
			log.Fatalf("Could not create table %s: %v", tableName, err)
		}
	}

	tblInfo, err := adminClient.TableInfo(ctx, tableName)
	if err != nil {
		log.Fatalf("Could not read info for table %s: %v", tableName, err)
	}

	if !sliceContains(tblInfo.Families, columnFamilyName) {
		if err := adminClient.CreateColumnFamily(ctx, tableName, columnFamilyName); err != nil {
			log.Fatalf("Could not create column family %s: %v", columnFamilyName, err)
		}
	}
	// [END bigtable_hw_create_table]

	// Set up Bigtable data operations client.
	// NewClient uses Application Default Credentials to authenticate.
	// [START bigtable_hw_connect_data]
	// client, err := bigtable.NewClient(ctx, *project, *instance, option.WithGRPCConnectionPool(10))
	client, err := bigtable.NewClient(ctx, *project, *instance)
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}
	// [END bigtable_hw_connect_data]

	// NOTE: Multipli greetings
	for i, greeting := range greetings {
		greetings[i] = strings.Repeat(greeting, 10000) // Here, result will be about 100kB
	}
	// fmt.Printf("greetings: %v\n", greetings)

	// [START bigtable_hw_write_rows]
	tbl := client.Open(tableName)
	multiplier := 1000 // 100kB*300,000 = 30GB
	chunkSize := 1000
	chunkedLen := len(greetings) * multiplier / chunkSize

	log.Printf("Writing greeting rows to table")

	s := time.Now()
	for i := 0; i < chunkedLen; i++ {
		runtime.GC() // Trigger GC explicitly

		muts := make([]*bigtable.Mutation, chunkSize)
		rowKeys := make([]string, chunkSize)

		for j := 0; j < chunkSize; j++ {
			idx := i*chunkSize + j

			muts[j] = bigtable.NewMutation()
			muts[j].Set(columnFamilyName, columnName, bigtable.Now(), []byte(greetings[idx%len(greetings)]))

			// Each row has a unique row key.
			//
			// Note: This example uses sequential numeric IDs for simplicity, but
			// this can result in poor performance in a production application.
			// Since rows are stored in sorted order by key, sequential keys can
			// result in poor distribution of operations across nodes.
			//
			// For more information about how to design a Bigtable schema for the
			// best performance, see the documentation:
			//
			//     https://cloud.google.com/bigtable/docs/schema-design

			rowKeys[j] = fmt.Sprintf("%s%d", columnName, idx)
		}

		func() {
			rowErrs, err := tbl.ApplyBulk(ctx, rowKeys, muts)
			if err != nil {
				log.Fatalf("Could not apply bulk row mutation: %v", err)
			}
			if rowErrs != nil {
				for _, rowErr := range rowErrs {
					log.Printf("Error writing row: %v", rowErr)
				}
				log.Fatalf("Could not write some rows")
			}
		}()
	}

	diff := time.Now().Sub(s)
	fmt.Printf("duration of applyBulk: %v\n", diff)
	// [END bigtable_hw_write_rows]

	size := len(greetings) * multiplier

	for j := 0; j < 2; j++ { // NOTE: Repeat twice
		// [START bigtable_hw_get_by_key]
		s := time.Now()
		log.Printf("Getting a single greeting by row key:")
		ssize := size // copy
		if ssize > 100 {
			ssize = 100
		}
		for i := 0; i < ssize; i++ {
			rowKey := fmt.Sprintf("%s%d", columnName, i)
			// rowKey := rowKeys[i]
			_, err := tbl.ReadRow(ctx, rowKey, bigtable.RowFilter(bigtable.ColumnFilter(columnName)))
			if err != nil {
				log.Fatalf("Could not read row with key %s: %v", rowKey, err)
			}
			// log.Printf("\t%s = %s\n", rowKeys[i], string(row[columnFamilyName][0].Value))
		}
		diff := time.Now().Sub(s)
		fmt.Printf("duration of read all: %v\n", diff)
		fmt.Printf("duration of read all in average: %v\n", time.Millisecond*time.Duration(int(diff/time.Millisecond)/ssize))
		// [END bigtable_hw_get_by_key]
	}

	// [START bigtable_hw_scan_all]
	{
		ssize := 1000

		s2 := time.Now()
		log.Printf("Reading all greeting rows:")

		i := 0
		err = tbl.ReadRows(ctx, bigtable.PrefixRange(columnName), func(row bigtable.Row) bool {
			// item := row[columnFamilyName][0]
			// log.Printf("\t%s = %s\n", item.Row, string(item.Value))
			// return true
			i++
			return i <= ssize
		}, bigtable.RowFilter(bigtable.ColumnFilter(columnName)))

		diff2 := time.Now().Sub(s2)
		fmt.Printf("duration of scan all: %v\n", diff2)
		fmt.Printf("duration of scan all in average: %v\n", time.Millisecond*time.Duration(int(diff2/time.Millisecond)/ssize))
	}

	if err = client.Close(); err != nil {
		log.Fatalf("Could not close data operations client: %v", err)
	}

	// [END bigtable_hw_scan_all]

	// // [START bigtable_hw_delete_table]
	// log.Printf("Deleting the table")
	// if err = adminClient.DeleteTable(ctx, tableName); err != nil {
	// 	log.Fatalf("Could not delete table %s: %v", tableName, err)
	// }

	// if err = adminClient.Close(); err != nil {
	// 	log.Fatalf("Could not close admin client: %v", err)
	// }
	// // [END bigtable_hw_delete_table]
}
