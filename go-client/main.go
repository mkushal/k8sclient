package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

// BatchJobs struct which contains
// an array of batchJob1
type BatchJobs struct {
	BatchJobs []Job `json:"batchJob1"`
}

// Job struct which contains a name
type Job struct {
	JobName    string `json:"jobName"`
	Image      string `json:"image"`
	RequestMem string `json:"requestMem"`
	RequestCpu string `json:"requestCpu"`
}

func connectToK8s() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/c/Users/kusha"
	}

	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Fatalln("failed to create K8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Failed to create K8s clientset")
	}

	return clientset
}

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, requestMem *string, requestCpu *string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0
	var completions int32 = 2
	var parallelism int32 = 2

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  *jobName,
							Image: *image,
							//Command: strings.Split(*cmd, " "),
							Command: strings.Split("ls", " "),
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									"cpu":    resource.MustParse(*requestCpu),
									"memory": resource.MustParse(*requestMem),
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
			Completions:  &completions,
			Parallelism:  &parallelism,
		},
	}

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job.")
	}

	//print job details
	log.Println("Created K8s job successfully" + *jobName)
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("batch.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened batch.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our BatchJobs array
	var batchJobs BatchJobs

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'batchJob1' which we defined above
	json.Unmarshal(byteValue, &batchJobs)

	// we iterate through every Job our BatchJobs array and
	// print out the batchJob1 their name
	for i := 0; i < len(batchJobs.BatchJobs); i++ {
		fmt.Println("Job Name: " + batchJobs.BatchJobs[i].JobName)
		fmt.Println("Job Image: " + batchJobs.BatchJobs[i].Image)
		fmt.Println("Requested Memory: " + batchJobs.BatchJobs[i].RequestMem)
		fmt.Println("Requested CPU: " + batchJobs.BatchJobs[i].RequestCpu)

		//jobName := flag.String("jobname", "test-job", "The name of the job")
		jobName := (batchJobs.BatchJobs[i].JobName)
		//containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
		containerImage := (batchJobs.BatchJobs[i].Image)
		//entryCommand := flag.String("command", "ls", "The command to run inside the container")
		//requestCpu := flag.String("requestcpu", "1", "cpu requested")
		requestCpu := (batchJobs.BatchJobs[i].RequestCpu)
		//requestMem := flag.String("requestmem", "500Mi", "memory requested")
		requestMem := (batchJobs.BatchJobs[i].RequestMem)

		flag.Parse()

		clientset := connectToK8s()
		launchK8sJob(clientset, &jobName, &containerImage, &requestCpu, &requestMem)
	}

}

