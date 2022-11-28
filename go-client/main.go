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
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"	
)

// BatchJobs struct which contains
// an array of baychJob1
type BatchJobs struct {
	BatchJobs []Job `json:"batchJob1"`
}

// User struct which contains a name
// a type and a list of social links
type Job struct {
	JobName    string `json:"jobName"`
	Image      string `json:"image"`
	RequestMem string `json:"requestMem"`
	RequestCpu string `json:"requestCpu"`
}

func connectToK8s() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/c/Users/kushal"
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

func launchK8sJob(clientset *kubernetes.Clientset, jobName *string, image *string, requestCpu *string, requestMem *string) {
	jobs := clientset.BatchV1().Jobs("default")
	var backOffLimit int32 = 0
	var completions int32 = 2
	var parallelism int32 = 2
	labels := make(map[string]string)
	labels["batchjobs-group"] = "batchjob1"	
	var topologykey string = "kubernetes.io/hostname"

    
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *jobName,
			Namespace: "default",
			Labels : labels,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{				
				ObjectMeta: metav1.ObjectMeta{
					Name:      *jobName,
					Namespace: "default",
					Labels : labels,
				},
				Spec: v1.PodSpec{					
					Containers: []v1.Container{
						{
							Name:  *jobName,
							Image: *image,							
							Command: strings.Split("ls", " "),
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									"cpu":    resource.MustParse(*requestCpu),
									"memory": resource.MustParse(*requestMem),
								},
							},							
						},
					},
					Affinity: &v1.Affinity{
						PodAffinity: &v1.PodAffinity {
							PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm {
								{
									Weight: 100,
									PodAffinityTerm: v1.PodAffinityTerm {
										LabelSelector: &metav1.LabelSelector {
											MatchLabels: labels,
										},
										TopologyKey: topologykey,
									},									
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

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	_, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln("Failed to create K8s job." + *jobName)
		log.Println("panic occurred:", err)
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

	// we initialize our UBatchJobs array
	var batchJobs BatchJobs

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'batchjob1' which we defined above
	json.Unmarshal(byteValue, &batchJobs)

	// we iterate through every user within our batchjobs array and
	// print out the user Type, their name, and their facebook url
	// as just an example
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

