# StorageProblemTask

For running the app need implement following steps

**Create database promotionsdb**

**Create Table**

`CREATE TABLE promotions (
    ID  SERIAL PRIMARY KEY,
    cvs_id         TEXT ,
    price          TEXT,
    expiration_date TEXT
);`

In the Main.go file change the database connections
`"postgres://dgalstyan:dgalstyan@localhost:5432/promotionsdb?sslmode=disable"
` 

Run Go application

http://localhost:1321/promotions/1

returns following json

{"id":"d018ef0b-dbd9-48f1-ac1a-eb4d90e57118","price":"60.683466","expiration_date":"2018-08-04 05:32:31 +0200 CEST"}

# **Additionally, consider:**

1. The .csv file could be very big (billions of entries) - how would your application
perform?

* Database Indexing
* Database Configuration
* Load Balancing and Scaling
* Memory Usage
* Database Transactions
* Performance Monitoring and Optimization

To handle billions of entries in CSV files efficiently, we need to adopt streaming techniques, 
optimize database interactions, 
configure the database appropriately, 
consider load balancing and scaling, 
and continuously monitor and optimize performance 
based on our specific requirements and infrastructure.

2. How would your application perform in peak periods (millions of requests per
   minute)?

* Hardware Resources
* Load Balancing
* Caching
* Database Optimization
* Asynchronous Processing. This can be achieved using message queues or job scheduling frameworks like RabbitMQ, Apache Kafka.
* Scaling Horizontally
* Performance Testing and Monitoring


How would you operate this app in production (e.g. deployment, scaling, monitoring)?

1. Deployment:

   * Containerization: Use containerization technologies like Docker to package the application and its dependencies into portable containers.
   * Orchestration: Utilize container orchestration tools like Kubernetes to manage and deploy containers across a cluster of servers.
   * Deployment Strategy: Implement a deployment strategy such as rolling updates or blue/green deployments to ensure seamless updates without downtime.

2. Scaling:

   * Horizontal Scaling: As the application load increases, scale horizontally by adding more instances of the application to distribute the workload. Utilize auto-scaling features provided by cloud platforms or implement scaling based on metrics like CPU usage or request rate.
   * Load Balancing: Employ a load balancer to evenly distribute incoming requests across multiple application instances, ensuring efficient utilization of resources.

3. Monitoring:

   * Application Monitoring: Utilize application performance monitoring tools like Prometheus, New Relic, or Datadog to track key metrics such as response time, error rates, and resource utilization. Set up alerts and notifications for abnormal behavior or performance degradation.
   * Logging: Implement structured logging to capture relevant application logs. Centralize log collection using tools like Elasticsearch, Logstash, and Kibana (ELK stack) or use cloud-based logging services.
   * Tracing: Implement distributed tracing to gain insights into application performance across different components. Tools like Jaeger or Zipkin can help trace and analyze requests as they traverse the system.

4. High Availability and Fault Tolerance:

    * Data Replication: Ensure database replication or employ a database clustering solution to maintain data redundancy and enable failover in case of a server or database failure.
    * Redundancy and Failover: Use techniques like server clustering or multiple availability zones in cloud environments to ensure high availability and fault tolerance.

5. Security:

    * Implement security best practices such as securing network communication with HTTPS, input validation, and user authentication/authorization.
    * Regularly apply security patches and updates to all components of the application stack.
    * Employ security monitoring tools and intrusion detection systems to detect and respond to potential security threats.
    * Backup and Disaster Recovery:
    * Regularly back up critical data and configuration files to prevent data loss.
    * Implement disaster recovery mechanisms such as off-site backups or replication to different geographic regions.

6. Continuous Integration and Delivery (CI/CD):

    * Set up a CI/CD pipeline to automate the build, test, and deployment processes.
    * Utilize version control systems (e.g., Git) and automation tools (e.g., Jenkins, CircleCI, or GitLab CI) to ensure consistent and reliable application updates.