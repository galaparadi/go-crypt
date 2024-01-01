<p style="font-size:48px">Go-crypt</p>
My personal crypt tool.

## Encryption method
This program using AES specification. You can use 192 or 256 key length to encrypt the file. Make sure the generate a key file first. To generate key file, go to [generate key](#generate-key) section for more info.

## Encrypting file
To Encrypt a file : 
```
go run main.go encrypt <source file> <destination file>
```

source file : your plain file name & location

destination file : your encrypted file name & location

## Decrypting file
To Decrypt a file : 
```
go run main.go decrypt <source file> <destination file>
```

source file : your encrypted file name & location

destination file : your plain file name & location

## Check file
To check encrypted file : 
```
go run main.go check <source file> <destination file>
```

source file : your plain file name & location

destination file : your encrypted file name & location

## Generate key
To generate a key : 
```
go run main.go keygen [options]
```

options :
- -length = key length. Default : 256

key file with name `q` will generated. You should put it with same path as the program.
