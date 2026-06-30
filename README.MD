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

## Membuat Pesan
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

## Consume Message

setelah message berhasil dipublish ke dalam queue, message tersebut sudah dapat diterima oleh consumer. sama seperti saat melakukan publish, langkah pertama adalah membuat `channel` dari `connection` yang sudah ada menggunakan function `connection.Channel()`.

setelah mendapatkan `channel`, kita dapat mulai mendengarkan queue menggunakan function `Consume()` yang menerima 7 parameter yaitu :

- `queue string`  
  nama queue yang akan didengarkan.

- `consumer string`  
  nama atau identitas consumer. nilai ini digunakan RabbitMQ untuk mengidentifikasi consumer yang sedang terkoneksi. jika diisi string kosong (`""`) maka RabbitMQ akan membuat nama consumer secara otomatis.

- `autoAck bool`  
  jika diisi `true` maka RabbitMQ akan langsung menganggap message berhasil diproses setelah dikirim ke consumer. jika diisi `false` maka consumer harus memanggil `Ack()` secara manual setelah proses selesai.

- `exclusive bool`  
  jika diisi `true` maka hanya consumer tersebut yang dapat mengonsumsi message dari queue tersebut selama koneksi masih aktif.

- `noLocal bool`  
  fitur ini sudah tidak digunakan oleh RabbitMQ sehingga biasanya selalu diisi `false`.

- `noWait bool`  
  jika diisi `false` maka client akan menunggu konfirmasi dari RabbitMQ bahwa consumer berhasil dibuat.

- `args amqp.Table`  
  berisi konfigurasi tambahan untuk consumer. pada penggunaan dasar biasanya diisi `nil`.

function `Consume()` akan mengembalikan dua nilai yaitu :

- `<-chan amqp.Delivery`  
  channel Go yang berisi setiap message yang diterima dari RabbitMQ.

- `error`  
  berisi error apabila proses pembuatan consumer gagal.

karena nilai yang dikembalikan berupa **channel**, maka consumer cukup melakukan perulangan menggunakan `for range` untuk menerima setiap message yang masuk.

```go
for msg := range msgs {
    log.Printf("%s", string(msg.Body))
}
```

setiap kali ada message baru di queue, RabbitMQ akan mengirimkan object `amqp.Delivery` melalui channel tersebut. object ini berisi seluruh informasi mengenai message seperti isi pesan (`Body`), header, routing key, exchange, delivery tag, dan informasi lainnya.

karena pada contoh berikut parameter `autoAck` diisi `false`

```go
msgs, err := channel.Consume(
    queueName,
    consumerName,
    false, // autoAck
    false,
    false,
    false,
    nil,
)
```

maka setelah message selesai diproses consumer wajib mengirimkan konfirmasi ke RabbitMQ menggunakan function `Ack()`.

```go
msg.Ack(true)
```

function `Ack()` menerima 1 parameter yaitu :

- `multiple bool`  
  jika diisi `false` maka hanya message tersebut yang akan di-acknowledge. jika diisi `true` maka RabbitMQ akan menganggap seluruh message dengan `DeliveryTag` yang lebih kecil atau sama telah berhasil diproses. pada kebanyakan kasus consumer memproses satu message dalam satu waktu sehingga parameter ini biasanya diisi `false`.

setelah RabbitMQ menerima `Ack()`, message akan dihapus dari queue sehingga tidak akan dikirim kembali.

apabila consumer berhenti atau mengalami error sebelum memanggil `Ack()`, RabbitMQ akan menganggap message tersebut belum selesai diproses dan akan mengirimkannya kembali ke consumer lain yang tersedia. mekanisme ini membuat RabbitMQ mampu menjaga agar message tidak hilang ketika terjadi kegagalan pada consumer.