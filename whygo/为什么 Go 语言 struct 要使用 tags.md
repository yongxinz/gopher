**原文链接：**[为什么 Go 语言 struct 要使用 tags](https://mp.weixin.qq.com/s/L7-TJ-CzYfuVrIBWP7Ebaw)

在 Go 语言中，struct 是一种常见的数据类型，它可以用来表示复杂的数据结构。在 struct 中，我们可以定义多个字段，每个字段可以有不同的类型和名称。

除了这些基本信息之外，Go 还提供了 struct tags，它可以用来指定 struct 中每个字段的元信息。

在本文中，我们将探讨为什么 Go 语言中需要使用 struct tags，以及 struct tags 的使用场景和优势。

## struct tags 的使用

struct tags 使用还是很广泛的，特别是在 json 序列化，或者是数据库 ORM 映射方面。

在定义上，它以 `key:value` 的形式出现，跟在 struct 字段后面，除此之外，还有以下几点需要注意：

### 使用反引号

在声明 struct tag 时，使用反引号 `` ` `` 包围 tag 的值，可以防止转义字符的影响，使 tag 更容易读取和理解。例如：

```go
type User struct {
    ID    int    `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    Email string `json:"email" db:"email"`
}
```

### 避免使用空格

在 struct tag 中，应该避免使用空格，特别是在 tag 名称和 tag 值之间。使用空格可能会导致编码或解码错误，并使代码更难以维护。例如：

```go
// 不规范的写法
type User struct {
    ID    int    `json: "id" db: "id"`
    Name  string `json: "name" db: "name"`
    Email string `json: "email" db: "email"`
}

// 规范的写法
type User struct {
    ID    int    `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    Email string `json:"email" db:"email"`
}
```

### 避免重复

在 struct 中，应该避免重复使用同一个 tag 名称。如果重复使用同一个 tag 名称，编译器可能会无法识别 tag，从而导致编码或解码错误。例如：

```go
// 不规范的写法
type User struct {
    ID    int    `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    Email string `json:"email" db:"name"`
}

// 规范的写法
type User struct {
    ID    int    `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    Email string `json:"email" db:"email"`
}
```

### 使用标准化的 tag 名称

为了使 struct tag 更加标准化和易于维护，应该使用一些标准化的 tag 名称。

例如，对于序列化和反序列化，可以使用 `json`、`xml`、`yaml` 等；对于数据库操作，可以使用 `db`。

```go
type User struct {
    ID       int    `json:"id" db:"id"`
    Name     string `json:"name" db:"name"`
    Password string `json:"-" db:"password"` // 忽略该字段
    Email    string `json:"email" db:"email"`
}
```

其中，`Password` 字段后面的 `-` 表示忽略该字段，也就是说该字段不会被序列化或反序列化。

### 多个 tag 值

如果一个字段需要指定多个 tag 值，可以使用 `,` 将多个 tag 值分隔开。例如：

```go
type User struct {
    ID        int    `json:"id" db:"id"`
    Name      string `json:"name" db:"name"`
    Email     string `json:"email,omitempty" db:"email,omitempty"`
}
```

其中 `omitempty` 表示如果该字段值为空，则不序列化该字段。

## struct tags 的原理

Go 的反射库提供了一些方法，可以让我们在程序运行时获取和解析结构体标签。

介绍这些方法之前，先来看看 `reflect.StructField` ，它是描述结构体字段的数据类型。定义如下：

```go
type StructField struct {
    Name      string      // 字段名
    Type      Type        // 字段类型
    Tag       StructTag   // 字段标签
}
```

结构体中还有一些其他字段，被我省略了，只保留了和本文相关的。

在结构体的反射中，我们经常使用 `reflect.TypeOf` 获取类型信息，然后使用 `Type.Field` 或 `Type.FieldByName()`  获取结构体字段的 `reflect.StructField`，然后根据 `StructField` 中的信息做进一步处理。

例如，可以通过 `StructField.Tag.Get` 方法获取结构体字段的标签值。

下面看一段代码：

```go
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type Manager struct {
    Title string `json:"title"`
    User
}

func main() {
    m := Manager{Title: "Manager", User: User{Name: "Alice", Age: 25}}

    mt := reflect.TypeOf(m)

    // 获取 User 字段的 reflect.StructField
    userField, _ := mt.FieldByName("User")
    fmt.Println("Field 'User' exists:", userField.Name, userField.Type)

    // 获取 User.Name 字段的 reflect.StructField
    nameField, _ := userField.Type.FieldByName("Name")
    tag := nameField.Tag.Get("json")
    fmt.Println("User.Name tag:", tag)
}
```

运行以上代码，输出结果如下：

```go
Field 'User' exists: User {string int}
User.Name tag: "name"
```

## struct tags 的优势

使用 struct tag 的主要优势之一是可以在**运行时通过反射来访问和操作 struct 中的字段**。

比如在 Go Web 开发中，常常需要将 HTTP 请求中的参数绑定到一个 struct 中。这时，我们可以使用 struct tag 指定每个字段对应的参数名称、验证规则等信息。在接收到 HTTP 请求时，就可以使用反射机制读取这些信息，并根据信息来验证参数是否合法。

另外，在将 struct 序列化为 JSON 或者其他格式时，我们也可以使用 struct tag 来指定每个字段在序列化时的名称和规则。

此外，使用 struct tag 还可以提高代码的**可读性**和**可维护性**。在一个大型的项目中，struct 中的字段通常会包含很多不同的元信息，比如数据库中的表名、字段名、索引、验证规则等等。

如果没有 struct tag，我们可能需要将这些元信息放在注释中或者在代码中进行硬编码。这样会让代码变得难以维护和修改。而使用 struct tag 可以将这些元信息与 struct 字段紧密关联起来，使代码更加清晰和易于维护。

## 常用的 struct tags

在 Go 的官方 wiki 中，有一个常用的 struct tags 的库的列表，我复制在下面了，感兴趣的同学可以看看源码，再继续深入学习。

Tag       | Documentation
----------|---------------
xml       | https://pkg.go.dev/encoding/xml
json      | https://pkg.go.dev/encoding/json
asn1      | https://pkg.go.dev/encoding/asn1
reform    | https://pkg.go.dev/gopkg.in/reform.v1
dynamodb  | https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbattribute/#Marshal
bigquery  | https://pkg.go.dev/cloud.google.com/go/bigquery
datastore | https://pkg.go.dev/cloud.google.com/go/datastore
spanner   | https://pkg.go.dev/cloud.google.com/go/spanner
bson      | https://pkg.go.dev/labix.org/v2/mgo/bson, https://pkg.go.dev/go.mongodb.org/mongo-driver/bson/bsoncodec
gorm      | https://pkg.go.dev/github.com/jinzhu/gorm
yaml      | https://pkg.go.dev/gopkg.in/yaml.v2
toml      | https://pkg.go.dev/github.com/pelletier/go-toml
validate  | https://github.com/go-playground/validator
mapstructure | https://pkg.go.dev/github.com/mitchellh/mapstructure
parser    | https://pkg.go.dev/github.com/alecthomas/participle
protobuf  | https://github.com/golang/protobuf
db        | https://github.com/jmoiron/sqlx
url       | https://github.com/google/go-querystring
feature   | https://github.com/nikolaydubina/go-featureprocessing

以上就是本文的全部内容，如果觉得还不错的话欢迎**点赞**，**转发**和**关注**，感谢支持。

***

**参考文章：**

- https://github.com/golang/go/wiki/Well-known-struct-tags

**推荐阅读：**

- [为什么 Go 不支持 []T 转换为 []interface](https://mp.weixin.qq.com/s/cwDEgnicK4jkuNpzulU2bw)