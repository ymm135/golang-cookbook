# rabbitmq 基础  
[官网](https://www.rabbitmq.com/)  

## docker 安装

```
# 携带管理界面的
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672  -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin rabbitmq:3.7.7-management  

# 进入管理界面
http://localhost:15672
```  

## java 发送与订阅  

[java spring demo](https://www.rabbitmq.com/tutorials/tutorial-three-spring-amqp.html)    
```
@Configuration
public class MQConfig {

    public final static String QUEUE_NAME = "queue.test";
    public final static String TOPIC_ROUTINGKEY = "exchange.topic.stream";
    public final static String EXCHANGE_STREAM = "exchange.stream";

    @Bean
    public TopicExchange topic() {
        return new TopicExchange(EXCHANGE_STREAM);
    }

    @Bean
    public Queue packetQueue() {
        return new AnonymousQueue();
    }


    @Bean
    public Binding binding1a(TopicExchange topic, Queue packetQueue) {
        return BindingBuilder.bind(packetQueue).to(topic).with(TOPIC_ROUTINGKEY);
    }


    @Profile("receiver")
    @Bean
    public PacketRecvHandler receiver() {
        return new PacketRecvHandler();
    }

}
```  

> 可以直接订阅 `exchange`，但是在绑定时，仍要增加一个`queue name`  

```
 @RabbitListener(queues = MQConfig.QUEUE_NAME)
    public void rawPackerRecv(String msg) { //(byte[] msg)
    
    }
```  

> 结束数据的类型要根据发送的数据类型变化，如果发送的是`byte[]`,接收的是`String`，接收的数据可能会有逗号分割符  


## go 发送与订阅  
[官网demo](https://www.rabbitmq.com/tutorials/tutorial-three-go.html)  


