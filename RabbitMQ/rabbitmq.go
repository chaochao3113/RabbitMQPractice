package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//url 格式 amqp:账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://dyc:dyc@192.168.56.22:5672/dyc"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel

	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	Key string
	// 连接信息
	Mqurl string
}

// 创建RabbitMQ结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	// 创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

//断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

//简单模式Step: 1.创建简单模式下rabbitmq实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return 	NewRabbitMQ(queueName, "", "")  // 使用默认exchange,key是没有
}

//简单模式Step: 2.简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列,如果队列不存在会自动创建,如果存在则跳过创建
	//保证队列存在,消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
		)
	if err != nil {
		fmt.Println(err)
	}
	//2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,   //使用default交换机,类型是direct
		r.QueueName,
		//如果为true,根据exchange类型和routkey规则,如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者,则会吧消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})

}

func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列,如果队列不存在会自动创建,如果存在则跳过创建
	//保证队列存在,消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true,表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//是否阻塞
		false,
		nil)

	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)

	//启用协程处理消息
	go func() {
		for d := range msgs {
			//实现我们要处理的逻辑函数
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for message, To exit press CTRL+C")
	<-forever
}

//订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//订阅模式生产
func (r *RabbitMQ) PublishPub(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",  //订阅模式下,交换机要设置成fanout广播类型
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息,仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil)

	r.failOnErr(err, "Failed to declare an exchange")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})

	r.failOnErr(err, "Failed to send message: " + message)
}

//订阅模式消费端代码
func (r *RabbitMQ) ReceiveSub() {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",  //订阅模式下,交换机要设置成fanout广播类型
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息,仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil)

	r.failOnErr(err, "Failed to declare an exchange")

	//2.试探性创建队列,这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",  //随机生成队列名称
		false,
		false,
		true,
		false,
		nil)

	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下,这里的key要为空
		"",
		r.Exchange,
		false,
		nil)

	//消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			//实现我们要处理的逻辑函数
			log.Printf("Receiverd a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}

//路由模式
//创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//路由模式要改成direct
		"direct",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息,仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil)

	r.failOnErr(err, "Failed to declare an exchange")

	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})

	r.failOnErr(err, "Failed to send message: " + message)
}

func (r *RabbitMQ) ReceiveRouting() {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//路由模式,要设置
		"direct",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息,仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil)

	r.failOnErr(err, "Failed to declare an exchange")

	//2.试探性创建队列,这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",  //随机生成队列名称
		false,
		false,
		true,
		false,
		nil)

	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定key
		r.Key,
		r.Exchange,
		false,
		nil)

	//消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			//实现我们要处理的逻辑函数
			log.Printf("Receiverd a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C")
	<-forever
}