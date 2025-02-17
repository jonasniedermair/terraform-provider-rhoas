package kafkas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rhoasAPI "redhat.com/rhoas/rhoas-terraform-provider/m/rhoas/api"
	"redhat.com/rhoas/rhoas-terraform-provider/m/rhoas/utils"
)

func DataSourceKafka() *schema.Resource {
	return &schema.Resource{
		Description: "`rhoas_kafka` provides a Kafka accessible to your organization in Red Hat OpenShift Streams for Apache Kafka.",
		ReadContext: dataSourceKafkaRead,
		Schema: map[string]*schema.Schema{
			CloudProviderField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud provider to use. A list of available cloud providers can be obtained using `data.rhoas_cloud_providers`.",
			},
			RegionField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region to use. A list of available regions can be obtained using `data.rhoas_cloud_providers_regions`.",
			},
			NameField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Kafka instance",
			},
			HrefField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path to the Kafka instance in the REST API",
			},
			StatusField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Kafka instance",
			},
			OwnerField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the Red Hat account that owns the Kafka instance",
			},
			BootstrapServerHostField: {
				Description: "The bootstrap server (host:port)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			CreatedAtField: {
				Description: "The RFC3339 date and time at which the Kafka instance was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			UpdatedAtField: {
				Description: "The RFC3339 date and time at which the Kafka instance was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
			IDField: {
				Description: "The unique identifier for the Kafka instance",
				Type:        schema.TypeString,
				Required:    true,
			},
			KindField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The kind of resource in the API",
			},
			VersionField: {
				Description: "The version of Kafka the instance is using",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceKafkaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	api, ok := m.(rhoasAPI.Clients)
	if !ok {
		return diag.Errorf("unable to cast %v to rhoasAPI.Clients)", m)
	}

	val := d.Get(IDField)
	id, ok := val.(string)
	if !ok {
		return diag.Errorf("unable to cast %v to string for use as for kafka id", val)
	}

	kafka, resp, err := api.KafkaMgmt().GetKafkaById(ctx, id).Execute()
	if err != nil {
		if apiErr := utils.GetAPIError(resp, err); apiErr != nil {
			return diag.FromErr(apiErr)
		}
	}

	err = setResourceDataFromKafkaData(d, &kafka)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
