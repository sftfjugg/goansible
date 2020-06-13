# ansible-go

Use golang to implement some functions of ansible

- use goroutine instead of multi-process
- parallel remote bash execution on multiple Linux hosts
- parallel copy file to multiple Linux hosts

## usage
```shell script
$ ansible-go -h
Usage of ./ansible-go:
  -c string
        bash command to execute
  -dest string
        destination of the file to copy(must including filename)
  -i string
        csv file including hosts information (default "hosts.csv")
  -m string
        module, could be shell or copy, default shell (default "shell")
  -mask string
        mask of the file to copy, default 0744 (default "0744")
  -src string
        source of the file to copy
```

## hosts.csv demo
```
192.168.220.121,22,root,111111
192.168.220.122,22,root,111111
192.168.220.123,22,root,111111
```

## Related Works
I have been influenced by the following great works:
- go-scp: https://github.com/bramvdbogaerde/go-scp