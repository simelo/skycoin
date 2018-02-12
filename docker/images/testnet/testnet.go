// Package testnet provides a Docker test net for SkyCoin
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
)

// SkyCoinTestNetwork encapsulates the test data and functionality
type SkyCoinTestNetwork struct {
	Compose      dockerCompose
	Services     []dockerService
	Peers        []string
	BuildContext string
}
type dockerService struct {
	SkyCoinParameters []string
	ImageName         string
	ImageTag          string
	NodesNum          int
	Ports             []string
}
type serviceBuild struct {
	Context    string
	Dockerfile string
}
type serviceNetwork struct {
	IPv4Address string `yaml:"ipv4_address"`
}
type service struct {
	Image       string
	Build       serviceBuild
	Networks    map[string]serviceNetwork
	Volumes     []string `yaml:"volumes,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
	Environment []string `yaml:"environment,omitempty"`
}
type networkIpamConfig struct {
	Subnet string
}
type networkIpam struct {
	Driver string
	Config []networkIpamConfig
}
type network struct {
	Driver string
	Ipam   networkIpam
}
type dockerCompose struct {
	Version  string
	Services map[string]service
	Networks map[string]network
}

// GetCurrentGitCommit returns the current git commit SHA
func GetCurrentGitCommit() string {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "git"
	cmdArgs := []string{"rev-parse", "--verify", "HEAD"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		log.Print("There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}
	sha := string(cmdOut)
	firstSix := sha[:6]
	return firstSix
}

// CreateDockerFile makes the Dockerfiles needed to build the images
// for the testnet
func (d *dockerService) CreateDockerFile(tempDir string) {
	dockerfileTemplate := path.Join("templates", "Dockerfile")
	_, err := os.Stat(dockerfileTemplate)
	if err != nil {
		log.Print(err)
		return
	}
	buildTemplate, err := template.ParseFiles(dockerfileTemplate)
	f, err := os.Create(path.Join(tempDir, "Dockerfile-"+d.ImageName))
	if err != nil {
		log.Print(err)
		return
	}
	err = buildTemplate.Execute(f, d)
	if err != nil {
		log.Print(err)
		return
	}
	f.Close()
}

// BuildImage builds the node docker image
func runTest(tempDir string) {
	cmdName := "docker-compose"
	composeFile := path.Join(tempDir, "docker-compose.yml")
	cmdArgs := []string{"-f", composeFile, "up"}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
	fmt.Println("Cleaning up...")
	cmdArgs = []string{"-f", composeFile, "down"}
	err = os.RemoveAll(tempDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}

// NewSkyCoinTestNetwork is the SkyCoinTestNetwork factory function
func NewSkyCoinTestNetwork(nodesNum int, buildContext string, tempDir string) SkyCoinTestNetwork {
	t := SkyCoinTestNetwork{}
	pubKey, secKey := cipher.GenerateKeyPair()
	ipHostNum := 2
	networkAddr := "172.16.200."
	commonParameters := []string{
		"--launch-browser=false",
		"--gui-dir=/usr/local/skycoin/static",
		"--master-public-key=" + pubKey.Hex(),
		"--testchain",
	}
	currentCommit := GetCurrentGitCommit()
	networkName := "skycoin-" + currentCommit
	t.BuildContext = buildContext
	t.Services = []dockerService{
		dockerService{
			ImageName: "skycoin-gui",
			SkyCoinParameters: []string{
				"--web-interface-addr=0.0.0.0",
				"--master",
				"--master-secret-key=" + secKey.Hex(),
			},
			ImageTag: currentCommit,
			NodesNum: 1,
			Ports:    []string{"6420:6420"},
		},
		dockerService{
			ImageName: "skycoin-nogui",
			SkyCoinParameters: []string{
				"--web-interface=false",
			},
			ImageTag: currentCommit,
			NodesNum: nodesNum,
		},
	}
	t.Compose = dockerCompose{
		Version:  "3",
		Services: make(map[string]service),
		Networks: map[string]network{
			string(networkName): network{
				Driver: "bridge",
				Ipam: networkIpam{
					Driver: "default",
					Config: []networkIpamConfig{
						networkIpamConfig{Subnet: networkAddr + "0/24"},
					},
				},
			},
		},
	}
	for idx, s := range t.Services {
		for i := 1; i <= s.NodesNum; i++ {
			num := strconv.Itoa(ipHostNum)
			ipAddress := networkAddr + num
			serviceName := "skycoin-" + num
			dockerfile := path.Join(tempDir, "Dockerfile-"+s.ImageName)
			dataDir := path.Join(tempDir, serviceName)
			t.Compose.Services[serviceName] = service{
				Image: s.ImageName + ":" + s.ImageTag,
				Build: serviceBuild{
					Context:    t.BuildContext,
					Dockerfile: dockerfile,
				},
				Networks: map[string]serviceNetwork{
					string(networkName): serviceNetwork{
						IPv4Address: ipAddress,
					},
				},
				Volumes: []string{dataDir + ":/root/.skycoin-test"},
				Ports:   s.Ports,
			}
			t.Peers = append(t.Peers, ipAddress+":6000")
			ipHostNum++
		}
		t.Services[idx].SkyCoinParameters = append(t.Services[idx].SkyCoinParameters, commonParameters...)
	}
	// SkyCoin Explorer
	explorerContext, err := filepath.Abs("../../../../.../../")
	if err != nil {
		log.Fatal(err)
	}
	t.Compose.Services["skycoin-explorer"] = service{
		Image: "skycoin-explorer:" + currentCommit,
		Build: serviceBuild{
			Context:    explorerContext,
			Dockerfile: path.Join(tempDir, "Dockerfile-explorer"),
		},
		Networks: map[string]serviceNetwork{
			string(networkName): serviceNetwork{
				IPv4Address: networkAddr + strconv.Itoa(ipHostNum),
			},
		},
		Ports:       []string{"8001:8001"},
		Environment: []string{"SKYCOIN_ADDR=http://172.16.200.2:6420"},
	}
	return t
}
func (t *SkyCoinTestNetwork) createComposeFile(tempDir string) {
	text, err := yaml.Marshal(t.Compose)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path.Join(tempDir, "docker-compose.yml"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(text)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

}
func (t *SkyCoinTestNetwork) prepareTestEnv(tempDir string) {
	for _, s := range t.Services {
		s.CreateDockerFile(tempDir)
	}

	peersText := []byte(strings.Join(t.Peers, "\n"))
	for k := range t.Compose.Services {
		err := os.Mkdir(path.Join(tempDir, k), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create(path.Join(tempDir, k, "connections.txt"))
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(peersText)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	// Copies the explorer dockerfile
	dockerfileExplorerSrc := path.Join("templates", "Dockerfile-explorer")
	dockerfileExplorerDst := path.Join(tempDir, "Dockerfile-explorer")
	df, err := os.Open(dockerfileExplorerSrc)
	_, err = file.CopyFile(dockerfileExplorerDst, df)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	_, callerFile, _, _ := runtime.Caller(0)
	projectPath, _ := filepath.Abs(filepath.Join(filepath.Dir(callerFile), "../../../"))
	log.Print("Source code base dir at ", projectPath)
	nodesPtr := flag.Int("-nodes", 5, "Number of nodes to launch.")
	buildContextPtr := flag.String("-buildcontext", projectPath, "Docker build context (source code root).")
	flag.Parse()
	buildContext, err := filepath.Abs(*buildContextPtr)
	tempDir, err := ioutil.TempDir("", "skycointest")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	testNet := NewSkyCoinTestNetwork(*nodesPtr, buildContext, tempDir)
	testNet.prepareTestEnv(tempDir)
	testNet.createComposeFile(tempDir)
	runTest(tempDir)
}