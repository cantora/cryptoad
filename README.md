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

###installation
WARNING: this is beta software, do not use it in production systems or trust
it with super secret stuff.

developer installation instructions coming soon.

