package network

// import (
// 	"log"
// 	"testing"

// 	"github.com/joho/godotenv"
// 	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// 	"github.com/stretchr/testify/assert"
// )

// type mocks struct{}

// func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
// 	return args.Name + "_id", args.Inputs, nil
// }

// func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
// 	return args.Args, nil
// }

// func TestInfrastructure(t *testing.T) {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Println("⚠️  .env file not found or failed to load")
// 	}

// 	err = pulumi.RunErr(func(ctx *pulumi.Context) error {
// 		err := deploy(ctx)
// 		assert.NoError(t, err)

// 		// Here you can make controlls

// 		return nil
// 	}, pulumi.WithMocks("myproject", "devstack", mocks{}))
// 	assert.NoError(t, err)
// }
