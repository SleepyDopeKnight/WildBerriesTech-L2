Что выведет программа? Объяснить вывод программы.

```
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Вывод: error. Так как функция тест возвращает нам значение nil, типа *customError,
для того чтобы вывело "ok" нужно чтобы и значение, и тип были равны nil.   
```