package cassandra

import (
	"github.com/dbielecki97/bookstore-utils-go/logger"
	"github.com/gocql/gocql"
	"os"
	"strings"
)

const (
	cassandraUser     = "CASSANDRA_USER"
	cassandraPass     = "CASSANDRA_PASSWORD"
	cassandraCluster  = "CASSANDRA_CLUSTER"
	cassandraKeyspace = "CASSANDRA_KEYSPACE"
)

var (
	session *gocql.Session
)

func init() {
	user := os.Getenv(cassandraUser)
	password := os.Getenv(cassandraPass)
	clusterHosts := strings.Split(os.Getenv(cassandraCluster), ",")
	keyspace := os.Getenv(cassandraKeyspace)

	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: user,
		Password: password,
	}

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		logger.Fatal("could not create cassandra cluster session", err)
	}

}

func GetSession() *gocql.Session {
	return session
}
