# db package

## Usage
Just some high-level usage highlights to get you started.  Read the godocs for a complete list of capabilities and their documentation.

```go
	conn, _ := gorm.Open("postgres", url)
    
    // GormReadWriter implements both the Reader and Writer interfaces
    rw := GormReadWriter{Tx: conn}
    
    // There are writer methods like: Create, Update and Delete
    // that will write Gorm struct to the db.  These writer methods
    // all support options for writing Oplog entries for the 
    // caller: WithOplog(true), WithWrapper(yourWrapper), WithMetadata(yourMetadata)
    err = rw.Create(context.Background(), user)
   
    // There are reader methods like: LookupByPublicId, LookupById, 
    // LookupByFriendlyName, SearchBy, LookupBy, etc
    // which will lookup resources for you and scan them into your Gorm struct
    err = rw.LookupByPublicId(context.Background(), foundUser)

    // There's reader ScanRows that facilitates scanning rows from 
    // a "raw" SQL query into your Gorm struct
    tx, err := rw.DB()
    where := "select * from test_users where friendly_name in ($1, $2)"
    rows, err := tx.Query(where, "alice", "bob")
	defer rows.Close()
	for rows.Next() {
        user := db_test.NewTestUser()
        // scan the row into your Gorm struct
		if err := rw.ScanRows(rows, &user); err != nil {
            return err
        }
        // Do something with the Gorm user struct
    }

    // DoTx is a writer function that wraps a TxHandler 
    // in a retryable transaction.  You simply implement a
    // TxHandler that does your writes and hand the handler
    // to DoTx, which will wrap the writes in a retryable 
    // transaction with the retry attempts and backoff
    // strategy you specify via options.
    _, err = rw.DoTx(
			context.Background(),
			10,             // 10 retries
			ExpBackoff{},   // exponential backoff 
			func(w Writer) error {
                // my handler will update the user's friendly name
				return w.Update(context.Background(), user, []string{"FriendlyName"},
					WithOplog(true), // write an oplog entry for this update
					WithWrapper(InitTestWrapper(t)),
					WithMetadata(oplog.Metadata{
						"deployment": []string{"amex"},
						"project":    []string{"central-info-systems", "local-info-systems"},
					}),
				)
			})


```