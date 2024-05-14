package main

import (
  "database/sql"
  "log"
  "net/http"
  "os"
  "time"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  _ "github.com/lib/pq"
  _ "github.com/joho/godotenv/autoload"
)

var (
  userInserts = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "user_inserts",
      Help: "New Synapse user signups",
    },
    []string{"app"},
  )
)

func init() {
  prometheus.MustRegister(userInserts)
}

func main() {
  dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
  if dbConnectionString == "" {
    log.Fatal("DB_CONNECTION_STRING not set")
  }

  db, err := sql.Open("postgres", dbConnectionString)
  if ( err != nil ) {
    log.Fatalf("Error opening database connection: %v", err)
  }
  defer db.Close()

  http.Handle("/metrics", promhttp.Handler())

  go func() {
    log.Fatal(http.ListenAndServe(":9188", nil))
  }()

  for {
    // query := "SELECT n_tup_ins FROM pg_stat_user_tables WHERE relname = 'users'"
    query := "SELECT COUNT(*) FROM users"
    var count int
    err = db.QueryRow(query).Scan(&count)
    if err != nil {
      log.Fatalf("Error executing query: %v", err)
    }
  
    userInserts.WithLabelValues("synapse").Set(float64(count))

    time.Sleep(30 * time.Second)
  }
}
