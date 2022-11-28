# Part 1 : Architectural Challenge

## Question 1
How would you improve the current design to achieve better:
- High availability 
- Resilience  
- Performance
- Cost efficiency




## Question 2
The number of scan requests can increase/decrease randomly in a day and on most weekends the system receives almost no requests at all.
What strategy would you suggest to save cost while still maintaining the best possible performance and scan completion times?


## Question 3
In step 6, each job needs to mount the source code folder into every engine that needs to run. How would you store the source code and make sure that engines can run in a scalable way?




## Question 4
Propose a high-level disaster recovery plan for the current architecture.



# Part 2 : Technical Challenge

I have created a script in go , which will read batch.json , fetch parameters jobName, image, requestMem & requestCpu values for each job and using podaffinity feature of kubernetes deploy the jobs on minimum number of nodes as we required. 

Script can be found at [here](https://github.com/mkushal/k8sclient/tree/main/go-client)
