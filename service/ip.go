package service

import (
	"context"
	"github.com/zcong1993/ip2region-service/pb"
	"github.com/zcong1993/ip2region-service/pkg"
)

type IP2RegionService struct {
	client *ip2region.Ip2Region
}

func NewIP2RegionService(p string) *IP2RegionService {
	c, err := ip2region.New(p)
	if err != nil {
		panic(err)
	}
	return &IP2RegionService{
		client: c,
	}
}

func (irs *IP2RegionService) Search(ctx context.Context, in *pb.IP) (*pb.IpInfo, error) {
	ip := in.Ip
	info, err := irs.client.MemorySearch(ip)
	if err != nil {
		return nil, err
	}

	return &pb.IpInfo{
		CityId:   info.CityId,
		Country:  info.Country,
		Region:   info.Region,
		Province: info.Province,
		City:     info.City,
		Isp:      info.ISP,
	}, nil
}
