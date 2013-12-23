cryptoad
========

Ever wanted to send a (probably large) file to someone
over a cloud service?? Ever get that icky feeling knowing
that their servers will probably retain the file forever
even after you delete it? You could encrypt it but that
means the person to which you're sending must have the 
necessary software/know-how to decrypt it! OHNOTHEYDONT?
Expletives!

So... try summoning a toad freind, then make to his mouth
a file. Kuchiyose no Jutsu!

<pre>
  _____________________
 /                     \
|  YUM! S3KR3T DATAS!  /
 \____________________/   
                      \    ___
                        (-)   (-)\
          ______       / ____    o\_
 file -> @))))))) --->  \_____\     \
                         / ---  /    |
                      __/    __/  __/
                       /\     /\   /\
</pre>

Cryptoad takes any single input file, encrypts it using a
supplied password, then packages the file into a binary
for a given target architecture (386, amd64, arm) and 
a given target OS (darwin, freebsd, linux, windows). The output
binary can later be executed, at which point it will prompt
for a password and attempt to decrypt the file using the
given password. If successful, the encrypted file will be
copied to the current working directory.

In summary (the following was run on linux):  
<pre>
$> ./cryptoad -arch amd64 -os windows -pass mountmyoboku -v 1 README.md toad-agent.exe
generate key from password
seal message
message sealed! size = 1330 bytes
summon toad...
finished.
$> stat -c '%s' toad-agent.exe
4177167
$> file toad-agent.exe
toad-agent.exe: PE32+ executable (console) x86-64 (stripped to external PDB), for MS Windows
</pre>

Now `toad-agent.exe` is shipped off to someone via drop box/google drive/etc
with the encrypted file embedded within. Its recipient can then download it,
run it, enter the password (presumably transmitted out of band by phone,
email or text), and finally receive the original un-encrypted file.

###Installation
WARNING: this is beta software, do not use it in production systems or trust
it with super secret stuff.

Due to the magic of [go](http://golang.org/) binaries, the output binary of 
cryptoad is completely self-sufficient (no dynamic linking). This means that
recipients of cryptoad binaries will be able to simply run them without 
installing anything.

On the other hand, cryptoad requires a golang development environment to
actually generate cryptoad output binaries. Specifically, the golang
environment must support cross compilation to the various operating
systems and architectures for which a user intends to target with cryptoad.
[This](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1)
tutorial details the steps required to set up a go cross compilation environment.

After a proper go environment is set up, build the cryptoad binary by running:  
`go build`  
finally, you can test that cryptoad detects the cross compilation support
in the go environment:  
```
$> ./cryptoad -h
usage: cryptoad INPUT OUTPUT
  INPUT:           the file to encrypt
  OUTPUT:          your newly summoned, self-decrypting toad friend /.0 _0}
  -arch="amd64":   target architecture. one of 386, amd64, arm
  -os="linux":     target OS. one of darwin, freebsd, linux, windows
  -pass="":        the password with which to encrypt. if not specified, you
                   will be prompted for one. make it a good one whydontcha!?
  -v=0:            verbosity level. higher level is more verbose
```  
the usage help lists `386`, `amd64` and `arm` for architectures and
lists `darwin`, `freebsd`, `linux` and `windows` for operating systems
because cryptoad has detected that they are available for cross 
compilation in the go environment.
