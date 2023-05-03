
// 获取一个原始的 *amqp.Channel 类型的变量
ch, err := conn.GetChannel()
if err != nil {
    log.Fatal(err)
}
defer ch.Close()

// 创建一个延时队列，设置 TTL 为 10 秒
_, err = ch.QueueDeclare(
    "delay_queue",
    false,
    false,
    false,
    false,
    amqp.Table{
        "x-message-ttl": int32(10 * time.Second),
        // 设置死信交换器和路由键
        "x-dead-letter-exchange": "dlx",
        "x-dead-letter-routing-key": "dlx_key",
    },
)
if err != nil {
    log.Fatal(err)
}

// 创建一个死信交换器和死信队列
err = ch.ExchangeDeclare(
    "dlx",
    "direct",
    false,
    false,
    false,
    false,
    nil,
)
if err != nil {
    log.Fatal(err)
}

_, err = ch.QueueDeclare(
    "dlq",
    false,
    false,
    false,
    false,
    nil,
)
if err != nil {
    log.Fatal(err)
}

// 绑定死信交换器和路由键
err = ch.QueueBind(
    "dlq",
    "dlx_key",
    "dlx",
    false,
    nil,
)
if err != nil {
    log.Fatal(err)
}