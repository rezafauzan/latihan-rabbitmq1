# Rabbit MQ
Rabbit MQ merupakan salah satu message broker yang populer untuk membangun sistem yang saling berkomunikasi secara asynchronous

Singkatnya dengan RabbitMQ dua aplikasi atau lebih bisa menerima, menyimpan, dan meneruskan pesan antar aplikasi

## Installasi Rabbit MQ Container
`docker run -d --name rabbitmq --hostname rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management`

menggunakan image `rabbitmq:management` port `5672` merupakan port default untuk client RabbitMQ menggunakan protokol AMQP, port `15672` merupakan port web ui Dashboard RabbitMQ

## Publisher & Consumer
`Publisher` biasa juga disebut juga `Producer` yaitu sebuah aplikasi yang mem-publish atau memproduksi pesan

`Consumer` yaitu sebuah aplikasi yang `mengambil`, `membaca`, `menerima`, `memproses` pesan yang dipublish oleh Producer

terkadang satu aplikasi yang sama bisa menjadi Publisher dan pesan yang di produce bisa diconsume diaplikasi itu juga 

## Channel
Dikarenakan membuat koneksi itu relatif `mahal` rabbitMQ didesain untuk buka 1 connection dan di connection ini bisa dibuat banyak channel
TCP Connection
```     
        |
+----+----+----+
|    |    |    |
ch1 ch2 ch3   ch4
```
nantinya di channel inilah kita akan mendeklarasikan `Queue`, `Publish message`, `Consume message`

## Queue
anggaplah `Queue` ini sebagai wadah untuk menyimpan pesan pesan yang diproduce

Bayangkan seperti tabel database.

```
Database
 ├── users
 ├── orders
 └── products
```

Di RabbitMQ:

```
RabbitMQ
 ├── email_queue
 ├── order_queue
 └── payment_queue
 ```

## Exchange
Exchange adalah komponen yang menerima pesan dari Producer, lalu menentukan pesan itu harus dikirim ke queue mana.

Jadi Producer tidak mengirim langsung ke Queue, melainkan mengirim ke Exchange. Exchange yang bertugas melakukan routing pesan.

```
Producer
   |
   v
Exchange
   |
   +--> Queue A
   |
   +--> Queue B
   |
   +--> Queue C
```

## Routing Key
`Routing Key` adalah label/alamat tujuan pesan yang digunakan oleh Exchange untuk menentukan ke queue mana pesan harus dikirim.

### Membuat Pesan
untuk memproduce pesan di golang pertama tambahkan dependency RabbitMQ dengan perintah

`go get github.com/rabbitmq/amqp091-go`

kemudian buat koneksi ke client rabbitmq dengan memanggil `amqp.Dial()` function ini memiliki 1 parameter yaitu url ke client rabbitmq function ini mengembalikan `*amqp.Connection` dan `error`

selanjutnya adalah membuat `channel` diatas koneksi yang dibuat dengan menggunakan function `connection.Channel()` mengembalikan `*amqp.Channel` dan `error` 

dengan channel ini kita bisa mendeklarasikan kita bisa membuat `Queue` dengan menggunakan function `QueueDeclare()` function ini menerima 6 parameter yaitu 
- `name string` 
nama queue
- `durable bool`
jika diisi true ketika aplikasi direstart queue yang dibuat tidak akan hilang
- `autoDelete bool`
jika diisi true queue akan dihapus ketika queue tidak ada consumer lagi
- `exclusive bool`
jika diisi true maka `queue` hanya akan bisa digunakan oleh connection yang membuatnya 
- `noWait bool`
jika diisi false akan menunggu response RabbitMQ
- `args amqp.Table`
ini diisi dengan konfigurasi tambahan untuk queue yang dibuat sebagai contoh konfigurasi untuk message expiration(TTL) 

tahap berikutnya adalah membuat pesan di queue yang dibuat tadi dengan function `PublishWithContext()` menerima 6 parameter yaitu : 
- `ctx context.Context`
Context Go standar Digunakan untuk timeout, cancellation, tracing
- `exchange string`
RabbitMQ selalu punya `exchange` bawaan yaitu `name = "" type = direct` jadi parameter ini bisa diisi dengan string kosong (`""`)
- `key string`
dikarenakan menggunakan default exchange yang memiliki aturan routing key == queue name parameter ini bisa diisi dengan nama queue yang dituju
- `mandatory bool`
jika diisi false ketika queue yang dituju tidak ada maka pesan akan dibuang jika diisi true maka akan mengembalikan pesan ke producer
- `immediate bool`
parameter ini sudah deprecated di RabbitMQ modern dan bisa diabaikan diisi dengan false
- `msg amqp.Publishing`
ini merupakan metadata pesan contohnya 
```
amqp.Publishing{
	ContentType: "text/plain",
	Body: []byte("test"),
}
```

function `PublishWithContext()` ini mengembalikan error yang bisa ditangkap