# Part 1 : Architectural Challenge

## Question 1
How would you improve the current design to achieve better:
- High availability 
- Resilience  
- Performance
- Cost efficiency

### Answer 1
```
Considering the architecture mentioned in the diagram , with few assumtions I can suggest below.
 - For more high availability on whole deployment region we should have multi region implementation of our services.
 - To make application deployments highly available and fault-tolerant, it’s a good practice to run pods having more than 1 replicas on nodes deployed in separate availability zones. 
 - To improve the performance of certain services, avoid placing them with other services that consume a lot of resources.
 - For the better performance , we should use horizontal pod autoscaler.
 - For the cost optimization we should implment cluster autoscaler / karpenter or similar solution which can automatically resize the cluster capacity in case of no usage.
 - Also we can leverage AWS autoscaling group dynamic scalling policies to save the cost based on network traffic.
 - we should also implment automatic monitoring & alerting for our deployments residing in kubernetes cluster. ( Prometheus + Grafana + All required exporters )
 - Closely monitoring the exceptions in logs of deployments could help us proactively to avoid downtime later. (EFK / ELK Logging stack)
 - If we are using aws than for optimize the cost we can use gp3 ebs volumes, spot instances for the stateless node groups etc..
 - Apart from above things , we should have standard devops practices like IaaC (terraform or similar), CI, Container Images Security Scanning etc.. 
 - Also we should have automated API testing , Performance testing for our Services involved in CI tool.
 - Messaging Queue should be also clustered for high availability.
```



## Question 2
The number of scan requests can increase/decrease randomly in a day and on most weekends the system receives almost no requests at all.
What strategy would you suggest to save cost while still maintaining the best possible performance and scan completion times?

### Answer 2
---
We should resize the k8s cluster based on traffic & resource consumption , 
We should use kubernetes cluster autoscaler to extend the node capacity automatically when more parrallel jobs are getting executed to get performance in pick hours.
During the time when we have lesser traffic autoscaler will shrink kubernetes cluster size to minimial number of nodes.

## Question 3
In step 6, each job needs to mount the source code folder into every engine that needs to run. How would you store the source code and make sure that engines can run in a scalable way?

### Answer 3
If We are using aws than We should use EFS , in case of azure it's  Microsoft Azure File Storage which can be used to pull & mount source code with Scanning job. 
Also we may use same EFS source code directory across multiple jobs. Also it's a cost effective option for the storage.



## Question 4
Propose a high-level disaster recovery plan for the current architecture.

### Answer 4
- Regular database backup with specific retention period.
- Regular Service configuration backup with specific retention period.
- Multisite (region) application Deployment ( Active Datacenter ==> Standby Datacenter )
- Continuous Monitoring for service resources utilization & performance.
- API SourceCode backup, IaaC backup
- Also make restore plan & test it on preproduction environment.
- Rotation of security keys, tokens should be done at regular intervals.



# Part 2 : Technical Challenge

I have created a script in go , which will read batch.json , fetch parameters jobName, image, requestMem & requestCpu values for each job and using podaffinity feature of kubernetes deploy the jobs on minimum number of nodes as we required. 

Script can be found at [here](https://github.com/mkushal/k8sclient/tree/main/go-client)
