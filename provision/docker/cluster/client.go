package cluster

import (
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"
)

const (
	DOCKER_NETWORK = "/docker/networks"
	DOCKER_LOGS    = "/docker/logs"
	HTTP = "http://"
)

type Gulp struct {
	Port string
}

type DockerClient struct {
	ContainerName string
	ContainerId   string
	Bridge        string
	IpAddr        string
	Gateway       string
}

func (d *DockerClient) LogsRequest(url string, port string) error {
	url = HTTP + url + port + DOCKER_LOGS
	err := request(d, url)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) NetworkRequest(url string, port string) error {
	url = HTTP + url + port + DOCKER_NETWORK
	err := request(d, url)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

/*
 * Request to gulp
 */
func request(d *DockerClient, url string) error {
	res, _ := json.Marshal(&d)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}
	return nil
}
