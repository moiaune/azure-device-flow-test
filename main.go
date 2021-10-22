package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func main() {

	deviceConfig := auth.NewDeviceFlowConfig("<app_id>", "<tenant_id>")
	authorizer, err := deviceConfig.Authorizer()
	if err != nil {
		log.Fatalf("error: could not get authorizer: %v", err)
	}

	ctx := context.Background()

	// We initialize a new GroupsClient with our subscription ID
	// and the clients authorizer to the one we got from the
	// device flow authentication, so that the GroupsClient can
	// authenticate as us
	c := resources.NewGroupsClient("<subscription_id>")
	c.Authorizer = authorizer

	// ListComplete will give us all resource groups within the
	// subscription that the GroupsClient is initialized for
	groupList, err := c.ListComplete(ctx, "", nil)
	if err != nil {
		log.Fatalf("error: could not list resource groups, %v", err)
	}

	// Loop through the result
	for groupList.NotDone() {
		// Retrieve the current Group struct from the iterator
		group := groupList.Value()

		// Let's print the name and location of the resource group
		fmt.Printf("- %s (Location: %s)\n", *group.Name, *group.Location)

		// NextWithContext() will return an error if there are noe more results.
		// We then want to exist
		if err := groupList.NextWithContext(ctx); err != nil {
			break
		}
	}
}
