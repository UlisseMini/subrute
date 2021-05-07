## subrute

A way to bruteforce user account passwords on linux, using `su`, because many systems aren't configured to prevent it.


Given the password hashes john cracks at around 3400 passwords per second.
subrute cracks at around 500 passwords/second, 7 times slower *but you don't need the hash*.

That means you could try all of `rockyou.txt` in around 8 hours!

