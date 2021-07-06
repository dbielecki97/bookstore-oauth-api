package cassandra

import (
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
	session      *gocql.Session
	user         = os.Getenv(cassandraUser)
	password     = os.Getenv(cassandraPass)
	clusterHosts = strings.Split(os.Getenv(cassandraCluster), ",")
	keyspace     = os.Getenv(cassandraKeyspace)
)

func init() {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: user,
		Password: password,
	}

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}

}

func GetSession() *gocql.Session {
	return session
}
