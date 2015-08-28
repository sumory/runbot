package main
import (
	"time"
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	t1:=time.Now().UnixNano()

	count:=64
	wg := &sync.WaitGroup{}
	wg.Add(count)

	for t:=0;t<count;t++{
		go func(){
			for i:=0;i<100000*10000;i++{
				if i==19999{

				}
			}
			wg.Done()
		}()
	}

	wg.Wait()




	t2:=time.Now().UnixNano()
	fmt.Println((t2-t1)/1000/1000)
}
