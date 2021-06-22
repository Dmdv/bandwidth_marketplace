package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
)

type (
	code uint8
)

const (
	success code = iota
	failed
)

func TestConnect(t *testing.T) {
	ctx := context.Background()

	resChannel := make(chan code)

	dClient, err := client.NewClientWithOpts()
	require.NoError(t, err)

	initDirs(t)
	defer cleanDirs(t)

	//go buildAndStart(ctx, t, dClient, magmaDockerfile, magmaContainer, magmaImage, resChannel)
	go startConsumerPostgres(ctx, dClient, t, resChannel)
	//go buildAndStart(ctx, t, dClient, consumerDockerfile, consumerContainer, consumerImage, resChannel)

	res := <-resChannel
	if res == failed {
		t.Fatal("tests are failed")
	}
}

const (
	magmaImage      = "magma_test"
	magmaDockerfile = "docker.local/magma-dockerfile"
	magmaContainer  = "magma_test_container"
	magmaPort       = "7151"
	magmaGRPCPort   = "5181"

	consumerImage             = "consumer_test"
	consumerDockerfile        = "docker.local/cons-dockerfile"
	consumerContainer         = "consumer_test_container"
	consumerPort              = "5151"
	consumerGRPCPort          = "7131"
	consumerPostgresContainer = "consumer_postgres"
	consumerPostgresScript    = "cons-postgres-entrypoint.sh"

	postgresImage = "postgres:11"
)

func buildAndStart(
	ctx context.Context, t *testing.T, cl *client.Client, dockerfile, containerName, imageName string, resChannel chan code) {

	build(t, cl, dockerfile, []string{imageName})

	_ = stopAndRemoveContainer(ctx, cl, containerName)

	var (
		cfg  *container.Config
		hCfg *container.HostConfig
	)
	switch containerName {
	case magmaContainer:
		cfg, hCfg = magmaContainerConfig(t)

	case consumerContainer:
		cfg, hCfg = consumerContainerConfig(t)

	default:
		t.Fatal("unexpected container name")
	}
	startContainer(ctx, cfg, hCfg, t, cl, containerName, resChannel)
}

const (
	magmaDirLog    = "magma/log"
	consumerDirLog = "consumer/log"
	providerDirLog = "provider/log"
)

func initDirs(t *testing.T) {
	err := os.MkdirAll(magmaDirLog, 777)
	require.NoError(t, err)

	err = os.MkdirAll(consumerDirLog, 777)
	require.NoError(t, err)

	err = os.MkdirAll(providerDirLog, 777)
	require.NoError(t, err)
}

func cleanDirs(t *testing.T) {
	err := os.RemoveAll(magmaDirLog)
	require.NoError(t, err)

	err = os.RemoveAll(magmaDirLog)
	require.NoError(t, err)

	err = os.RemoveAll(providerDirLog)
	require.NoError(t, err)
}

func build(t *testing.T, cl *client.Client, dockerFile string, tags []string) {
	rd := rootDir(t)

	tar, err := archive.TarWithOptions(rd, &archive.TarOptions{})
	require.NoError(t, err)

	opts := types.ImageBuildOptions{
		Dockerfile: dockerFile,
		Tags:       tags,
		Remove:     true,
	}
	res, err := cl.ImageBuild(context.Background(), tar, opts)
	require.NoError(t, err)
	defer res.Body.Close()

	checkBuilding(res.Body, t)
	fmt.Printf("%s is built.\n", strings.Join(tags, ""))
}

func startContainer(
	ctx context.Context, cfg *container.Config, hostCfg *container.HostConfig,
	t *testing.T, cl *client.Client, name string, resChannel chan code) {

	resp, err := cl.ContainerCreate(ctx, cfg, hostCfg, nil, nil, name)
	require.NoError(t, err)

	err = cl.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	require.NoError(t, err)

	statusCh, errCh := cl.ContainerWait(ctx, resp.ID, container.WaitConditionNextExit)
	select {
	case <-errCh:
		err := stopAndRemoveContainer(ctx, cl, name)
		require.NoError(t, err)
		resChannel <- failed
	case <-statusCh:
		err := stopAndRemoveContainer(ctx, cl, name)
		require.NoError(t, err)
		resChannel <- failed
	}
}

func stopAndRemoveContainer(ctx context.Context, client *client.Client, name string) error {
	if err := client.ContainerStop(ctx, name, nil); err != nil {
		return err
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	return client.ContainerRemove(ctx, name, removeOptions)
}

func magmaContainerConfig(t *testing.T) (*container.Config, *container.HostConfig) {
	port, err := nat.NewPort("tcp", magmaPort)
	require.NoError(t, err)

	grpcPort, err := nat.NewPort("tcp", magmaGRPCPort)
	require.NoError(t, err)

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: string(port),
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, "code", "integration-tests", "magma", "config"),
				Target: "/magma/config",
			},
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, "code", "integration-tests", "magma", "log"),
				Target: "/magma/log",
			},
		},
	}

	cfg := &container.Config{
		Image: magmaImage,
		ExposedPorts: map[nat.Port]struct{}{
			port:     {},
			grpcPort: {},
		},
		Cmd: []string{"./bin/magma"},
	}

	return cfg, hostConfig
}

func consumerContainerConfig(t *testing.T) (*container.Config, *container.HostConfig) {
	port, err := nat.NewPort("tcp", consumerPort)
	require.NoError(t, err)

	grpcPort, err := nat.NewPort("tcp", consumerGRPCPort)
	require.NoError(t, err)

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: string(port),
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, "consumer", "config"),
				Target: "/consumer/config",
			},
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, "consumer", "log"),
				Target: "/consumer/log",
			},
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, "consumer", "data"),
				Target: "/consumer/data",
			},
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, "consumer", "keys"),
				Target: "/consumer/keys_config/consumer",
			},
		},
	}

	cfg := &container.Config{
		Image: consumerImage,
		ExposedPorts: map[nat.Port]struct{}{
			port:     {},
			grpcPort: {},
		},
		Cmd: []string{
			"./bin/consumer",
			fmt.Sprintf("--port=%s", consumerPort),
			fmt.Sprintf("--grpc_port=%s", consumerGRPCPort),
			"--hostname=localhost",
			"--deployment_mode=0",
			"--keys_file=keys_config/consumer/keys.txt",
			"--log_dir=/consumer/log",
			"--db_dir=/consumer/data",
		},
	}

	return cfg, hostConfig
}

func startConsumerPostgres(ctx context.Context, cl *client.Client, t *testing.T, resChannel chan code) {
	out, err := cl.ImagePull(ctx, postgresImage, types.ImagePullOptions{})
	require.NoError(t, err)
	defer out.Close()

	checkBuilding(out, t)

	upPostgres(ctx, cl, t, consumerPostgresContainer, "consumer", resChannel)
}

func upPostgres(
	ctx context.Context, cl *client.Client, t *testing.T, containerName, nodeFolder string, resChannel chan code) {

	pCfg, pHostCfg := postgresContainerConfig(t, nodeFolder)
	pResp, err := cl.ContainerCreate(ctx, pCfg, pHostCfg, nil, nil, containerName)
	require.NoError(t, err)

	pPostContainerName := containerName + "_post"
	pPostCfg, pPostHostCfg := postgresPostContainerConfig(t, nodeFolder, consumerPostgresScript, containerName)
	pPostResp, err := cl.ContainerCreate(ctx, pPostCfg, pPostHostCfg, nil, nil, pPostContainerName)
	require.NoError(t, err)

	err = cl.ContainerStart(ctx, pResp.ID, types.ContainerStartOptions{})
	defer func() {
		err := stopAndRemoveContainer(ctx, cl, containerName)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	err = cl.ContainerStart(ctx, pPostResp.ID, types.ContainerStartOptions{})
	defer func() {
		err = stopAndRemoveContainer(ctx, cl, pPostContainerName)
		require.NoError(t, err)
	}()
	require.NoError(t, err)

	statusCh, errCh := cl.ContainerWait(ctx, pResp.ID, container.WaitConditionNextExit)
	select {
	case <-errCh:
		resChannel <- failed
	case <-statusCh:
		resChannel <- failed
	}
}

func postgresContainerConfig(t *testing.T, nodeFolder string) (*container.Config, *container.HostConfig) {
	port, err := nat.NewPort("tcp", "5432")
	require.NoError(t, err)

	hostCfg := &container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: string(port),
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: withTestRoot(t, nodeFolder, "data", "postgresql"),
				Target: "/var/lib/postgresql/data",
			},
		},
		NetworkMode: "default",
	}
	cfg := &container.Config{
		Image: postgresImage,
		Env: []string{
			"POSTGRES_PORT=5432",
			"POSTGRES_HOST=0.0.0.0",
			"POSTGRES_USER=postgres",
			"POSTGRES_HOST_AUTH_METHOD=trust",
		},
	}
	//
	//nCfg := &network.NetworkingConfig{
	//	EndpointsConfig: map[string]*network.EndpointSettings{
	//		"default": {
	//			DriverOpts: map[string]string{
	//				"driver": "bridge",
	//			},
	//		},
	//		"testnet0": {
	//			DriverOpts: map[string]string{
	//				"external": "true",
	//			},
	//		},
	//	},
	//}

	return cfg, hostCfg
}

func postgresPostContainerConfig(t *testing.T, nodeFolder, scriptName, postgresName string) (*container.Config, *container.HostConfig) {
	hostCfg := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: withRoot(t, "bin"),
				Target: strings.Join([]string{"/" + nodeFolder, "bin"}, "/"),
			},
			{
				Type:   mount.TypeBind,
				Source: withRoot(t, "sql", nodeFolder),
				Target: strings.Join([]string{"/" + nodeFolder, "sql"}, "/"),
			},
		},
		Links: []string{postgresName + ":" + postgresName},
	}
	cfg := &container.Config{
		Image: postgresImage,
		Env: []string{
			"POSTGRES_PORT=5432",
			"POSTGRES_HOST=0.0.0.0",
			"POSTGRES_USER=postgres",
		},
		Cmd: []string{
			"bash",
			strings.Join([]string{"/" + nodeFolder, "bin", scriptName}, "/"),
		},
	}

	return cfg, hostCfg
}

func rootDir(t *testing.T) string {
	wd, err := os.Getwd()
	require.NoError(t, err)

	fp := strings.Split(wd, string(os.PathSeparator))

	return strings.Join(fp[:len(fp)-2], string(os.PathSeparator))
}

func withRoot(t *testing.T, path ...string) string {
	rd := rootDir(t)
	fp := strings.Split(rd, string(os.PathSeparator))
	res := append(fp, path...)
	return strings.Join(res, string(os.PathSeparator))
}

func withTestRoot(t *testing.T, path ...string) string {
	rd := rootDir(t)
	fp := strings.Split(rd, string(os.PathSeparator))
	res := append(fp, "code", "integration-tests")
	res = append(res, path...)
	return strings.Join(res, string(os.PathSeparator))
}

type errorLine struct {
	Error       string      `json:"error"`
	ErrorDetail errorDetail `json:"errorDetail"`
}

type errorDetail struct {
	Message string `json:"message"`
}

func checkBuilding(rd io.Reader, t *testing.T) {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	errLine := &errorLine{}
	err := json.Unmarshal([]byte(lastLine), errLine)
	require.NoError(t, err)
	require.False(t, errLine.Error != "")

	require.NoError(t, scanner.Err())
}
