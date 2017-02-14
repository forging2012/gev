package swagger

import (
	"testing"
)

func TestSwagger(t *testing.T) {
	swagger := NewSwagger()
	t.Logf("%v\n", swagger.WriteJson("swagger.json"))
}
