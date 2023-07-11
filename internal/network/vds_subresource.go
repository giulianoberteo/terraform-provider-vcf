/*
 *  Copyright 2023 VMware, Inc.
 *    SPDX-License-Identifier: MPL-2.0
 */

package network

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	validation_utils "github.com/vmware/terraform-provider-vcf/internal/validation"
	"github.com/vmware/vcf-sdk-go/models"
)

// VdsSchema this helper function extracts the VDS Schema, so that
// it's made available for both Domain and Cluster creation.
// This specification contains vSphere distributed switch configurations.
func VdsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "vSphere Distributed Switch name",
				ValidateFunc: validation.NoZeroValues,
			},
			"is_used_by_nsx": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Identifies if the vSphere distributed switch is used by NSX-T",
			},
			"portgroup": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of portgroups to be associated with the vSphere Distributed Switch",
				Elem:        PortgroupSchema(),
			},
			"nioc_bandwidth_allocations": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "List of Network I/O Control Bandwidth Allocations for System Traffic based on" +
					" shares, reservation, and limit, you can configure Network I/O Control to allocate certain amount" +
					" of bandwidth for traffic generated by vSphere Fault Tolerance, iSCSI storage, vSphere vMotion, and so on." +
					" You can use Network I/O Control on a distributed switch to configure bandwidth allocation for the traffic " +
					" that is related to the main system features in vSphere",
				Elem: NiocBandwidthAllocationSchema(),
			},
		},
	}
}

func TryConvertToVdsSpec(object map[string]interface{}) (*models.VdsSpec, error) {
	result := &models.VdsSpec{}
	if object == nil {
		return nil, fmt.Errorf("cannot conver to VdsSpec, object is nil")
	}
	name := object["name"].(string)
	if len(name) == 0 {
		return nil, fmt.Errorf("cannot conver to VdsSpec, name is required")
	}
	result.Name = &name
	if isUsedByNsx, ok := object["is_used_by_nsx"]; ok && !validation_utils.IsEmpty(isUsedByNsx) {
		result.IsUsedByNSXT = isUsedByNsx.(bool)
	}
	if portgroupsRaw, ok := object["portgroup"]; ok && !validation_utils.IsEmpty(portgroupsRaw) {
		portgroupsList := portgroupsRaw.([]interface{})
		if len(portgroupsList) > 0 {
			result.PortGroupSpecs = []*models.PortgroupSpec{}
			for _, portgroupListEntry := range portgroupsList {
				portgroupSpec, err := tryConvertToPortgroupSpec(portgroupListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.PortGroupSpecs = append(result.PortGroupSpecs, portgroupSpec)
			}
		}
	}
	if niocBandwidthAllocationsRaw, ok := object["nioc_bandwidth_allocations"]; ok && !validation_utils.IsEmpty(niocBandwidthAllocationsRaw) {
		niocBandwidthAllocationsList := niocBandwidthAllocationsRaw.([]interface{})
		if len(niocBandwidthAllocationsList) > 0 {
			result.NiocBandwidthAllocationSpecs = []*models.NiocBandwidthAllocationSpec{}
			for _, niocBandwidthAllocationListEntry := range niocBandwidthAllocationsList {
				niocBandwidthAllocationSpec, err := tryConvertToNiocBandwidthAllocationSpec(
					niocBandwidthAllocationListEntry.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
				result.NiocBandwidthAllocationSpecs = append(result.NiocBandwidthAllocationSpecs,
					niocBandwidthAllocationSpec)
			}
		}
	}

	return result, nil
}
