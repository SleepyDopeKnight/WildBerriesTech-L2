Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
Вывод: <nil> false. Переменная типа *os.PathError преобразуется в интерфейс error.
Этот интерфейс содрежит значение <nil> и тип *os.PathError, по этому он не является <nil>,
так как содрежит информацию о типе.

Интерфейс - структура с двумя полями. Первое поле ссылается на значение интерфейса, а второе
ссылается на тип и методы, которые реализуют этот тип. 
Пустые интерфейс содержат и значение и тип <nil> и ничего не весят.
```