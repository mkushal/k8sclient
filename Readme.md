# Architectural Challenge

**Question 1**
How would you improve the current design to achieve better:
High availability 
Resilience - 
Performance - For the performance from Infrastructure side , we should use auto scalling 
Cost efficiency

> => For high availability we should have multi site implementation of our services.
  To make application deployments highly available and fault-tolerant, it’s a good practice to run pods on nodes deployed in separate availability zones.
  To improve the performance of certain services, avoid placing them with other services that consume a lot of resources.


**Question 2**
The number of scan requests can increase/decrease randomly in a day and on most weekends the system receives almost no requests at all.
What strategy would you suggest to save cost while still maintaining the best possible performance and scan completion times?

> We should resize the k8s cluster based on traffic & resource consumption , We should use kubernetes cluster autoscaler to extend the node capacity automatically when more parrallel jobs are getting executed to get performance in pick hours. During the time when we have lesser traffic autoscaler will shrink kubernetes cluster size to minimial number of nodes.

**Question 3**
In step 6, each job needs to mount the source code folder into every engine that needs to run. How would you store the source code and make sure that engines can run in a scalable way?

> You should always store the artifact(code) as a part of image & image should have version tag based on your scm commit/tag for easy tracking.


**Question 4**
Propose a high-level disaster recovery plan for the current architecture.

> Regular backups server including software and configuration on the server
Multisite application Deployment ( Active Datacenter ==> Standby Datacenter )
