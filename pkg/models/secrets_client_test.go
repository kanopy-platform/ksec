package models

import (
	"strings"
	"testing"
)

var secretsClient *SecretsClient

// helpers
func testErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

// unit tests
func TestNewSecretsClient(t *testing.T) {
	secretsClient = MockNewSecretsClient()
}

func TestCreate(t *testing.T) {
	_, err := secretsClient.Create("test-secret")
	testErr(err, t)
}

func TestList(t *testing.T) {
	secrets, err := secretsClient.List()
	testErr(err, t)

	if secrets.Items[0].Name != "test-secret" {
		t.Fatal(err.Error())
	}
}

func TestCreateWithData(t *testing.T) {
	dataArgs := []string{
		"key=value",
		"secret-key=secret-value",
		"ENV_VAR=~!@#$%^&*()_+",
		"DB_URL=mongodb://host1.example.com:27017,host2.example.com:27017/prod?replicaSet=prod",
	}
	data := make(map[string][]byte)

	for _, item := range dataArgs {
		split := strings.SplitN(item, "=", 2)
		if len(split) != 2 {
			t.Errorf("Data is not formatted correctly: %s", item)
		}
		data[split[0]] = []byte(split[1])
	}

	_, err := secretsClient.CreateWithData("test-secret-with-data", data)
	testErr(err, t)
}

func TestGet(t *testing.T) {
	secret, err := secretsClient.Get("test-secret-with-data")
	testErr(err, t)

	if string(secret.Data["DB_URL"]) != "mongodb://host1.example.com:27017,host2.example.com:27017/prod?replicaSet=prod" {
		t.Fatal(err.Error())
	}
}

func TestGetKey(t *testing.T) {
	value, err := secretsClient.GetKey("test-secret-with-data", "secret-key")
	testErr(err, t)

	if value != "secret-value" {
		t.Fatal(err.Error())
	}

	value, err = secretsClient.GetKey("test-secret-with_data", "thiskeydoesnotexist")
	if err == nil {
		t.Fatal("non-existent key should have received an error")
	}
}

func TestUpdate(t *testing.T) {
	secret, err := secretsClient.Get("test-secret-with-data")
	testErr(err, t)

	data := map[string][]byte{
		"key": []byte("newvalue"),
	}

	secretsClient.Update(secret, data)

	if string(secret.Data["key"]) != "newvalue" {
		t.Fatal(err.Error())
	}
}

func TestUpsert(t *testing.T) {
	data := map[string][]byte{
		"key": []byte("upsert"),
	}
	secret, err := secretsClient.Upsert("upsert-secret", data)
	testErr(err, t)

	if string(secret.Data["key"]) != "upsert" {
		t.Fatal(err.Error())
	}
}
