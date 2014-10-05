package ec2

import "strconv"

// Response to a DescribeVolumes request.
//
// See http://goo.gl/Bm91er for more details.
type VolumesResp struct {
	RequestId string   `xml:"requestId"`
	Volumes   []Volume `xml:"volumeSet>item"`
}

type VolumeAttachment struct {
	VolumeID            string `xml:"volumeId"`
	InstanceID          string `xml:"instanceId"`
	Device              string `xml:"device"`
	Status              string `xml:"status"`
	AttachTime          string `xml:"attachTime"`
	DeleteOnTermination bool   `xml:"deleteOnTermination"`
}

type Volume struct {
	VolumeID         string             `xml:"volumeId"`
	Size             string             `xml:"size"`
	SnapshotID       string             `xml:"snapshotId"`
	AvailabilityZone string             `xml:"availabilityZone"`
	Status           string             `xml:"status"`
	CreateTime       string             `xml:"createTime"`
	Attachments      []VolumeAttachment `xml:"attachmentSet>item"`
	Tags             []Tag              `xml:"tagSet>item"`
	VolumeType       string             `xml:"volumeType"`
	IOPS             int                `xml:"iops"`
	Encrypted        bool               `xml:"encrypted"`
}

// Volumes returns details about EBS volumes.
// The ids and filter parameters, if provided, will limit the EBS volumes returned.
// For example, to get all the ssd volumes associated with this account set
// the filter "volume-type" to gp2.
//
// Note: calling this function with nil ids and filter parameters will result in
// a very large number of volumes being returned.
//
// See http://goo.gl/Bm91er for more details.
func (ec2 *EC2) Volumes(ids []string, filter *Filter) (resp *VolumesResp, err error) {
	params := makeParams("DescribeVolumes")
	for i, id := range ids {
		params["VolumeId."+strconv.Itoa(i+1)] = id
	}
	filter.addParams(params)

	resp = &VolumesResp{}
	err = ec2.query(params, resp)
	return
}

type AttachVolumeRequest struct {
	VolumeID   string
	InstanceID string
	Device     string
}

// See http://goo.gl/mgC5e1 for more details
type AttachVolumeResp struct {
	RequestId string `xml:"requestId"`
	Volume
}

// Attaches an Amazon EBS volume to a running or stopped instance and
// exposes it to the instance with the specified device name.
//
// Encrypted Amazon EBS volumes can be attached only to instances that support Amazon EBS encryption.
// For more information, see Amazon EBS encryption in the Amazon Elastic Compute Cloud User Guide for Linux.
//
// For a list of supported device names, see Attaching the Volume to an Instance.
// Any device names that aren't reserved for instance store volumes can be used for Amazon EBS volumes.
// For more information, see Amazon EC2 Instance Store in the Amazon Elastic Compute Cloud User Guide for Linux.
//
// Note:
// If a volume has an AWS Marketplace product code:
// * The volume can be attached only to the root device of a stopped instance.
// * AWS Marketplace product codes are copied from the volume to the instance.
// * You must be subscribed to the product.
// * The instance type and operating system of the instance must support the product.
//   For example, you can't detach a volume from a Windows instance and attach it to a Linux instance.
//
// See http://goo.gl/mgC5e1 for more details
func (ec2 *EC2) AttachVolume(req AttachVolumeRequest) (resp *AttachVolumeResp, err error) {
	params := makeParams("DescribeVolumes")
	params["VolumeId"] = req.VolumeID
	params["InstanceId"] = req.InstanceID
	params["Device"] = req.Device

	resp = &AttachVolumeResp{}
	err = ec2.query(params, resp)
	return
}
