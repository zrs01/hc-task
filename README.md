# Concurrent tasks execution
The program execute the defined tasks executing concurrently.

A text file is provided to feed the parameter to the tasks.

Example:

If the content in the text (input.txt) as below:
```txt
P1
P2
P3
```

Execute below command:
```batch
C:\etc -i input.txt -c "cmd /c check.bat"
```

it equals to below commands executing concurrently
```batch
C:\cmd /c check.bat P1
C:\cmd /c check.bat P2
C:\cmd /c check.bat P3
```

## Options

```sh
NAME:
   task - Execute concurrent tasks

USAGE:
   task [global options] command [command options] [arguments...]

VERSION:
   development

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value, -i value    input file with content separated by enter
   --command value, -c value  command to be executed
   --tasks value              max tasks execute concurrently (default: 5)
   --output value, -o value   output file (default: "output.txt")
   --help, -h                 show help
   --version, -v              print the version
```