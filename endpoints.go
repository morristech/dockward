package main

import (
	"github.com/abiosoft/dockward/balancer"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
)

func endpointsFromLabel(containerPort int, label string) (balancer.Endpoints, error) {
	return endpointsFromFilter(containerPort, "label", label)
}

func endpointsFromName(containerPort int, name string) (balancer.Endpoints, error) {
	return endpointsFromFilter(containerPort, "name", name)
}

func endpointsFromId(containerPort int, id string) (balancer.Endpoints, error) {
	return endpointsFromFilter(containerPort, "id", id)
}

func endpointsFromFilter(containerPort int, key, value string) (balancer.Endpoints, error) {
	filter := filters.NewArgs()
	filter.Add(key, value)
	containers, err := client.ContainerList(types.ContainerListOptions{Filter: filter})
	if err != nil {
		return nil, err
	}
	endpoints := make(balancer.Endpoints, len(containers))
	for i, c := range containers {
		if err := connectContainer(c.ID); err != nil {
			return nil, err
		}
		ip, err := ipFromContainer(c.ID)
		if err != nil {
			return nil, err
		}
		endpoints[i] = balancer.Endpoint{
			Id:   c.ID,
			Ip:   ip,
			Port: containerPort,
		}
	}
	return endpoints, nil
}
