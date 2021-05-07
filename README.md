## subrute

A way to bruteforce user account passwords on linux, using `su`, because many systems aren't configured to prevent it.


Given the password hashes john cracks at around 3400 passwords per second.
subrute cracks at around 500 passwords/second, 7 times slower *but you don't need the hash*.

That means you could try all of `rockyou.txt` in around 8 hours!



### Ideas for improvement

1. Calibrate sleep time based on testing with known user/pass
Tried this for a while, wierdly on success `su` returns in ~100-300ms but if killed in ~50ms it always returns success.
This needs more investigation


2. Rewrite in C, use faster syscalls
I'm not sure if go does this in the fastest way, I read that more memory = slower fork() so maybe
C will give an improvement there. I noticed even when increasing the of workers it doesn't improve after 500p/s,
But decreasing `sleepTime` gives a small improvement. perhaps this hints at a bottleneck somewhere?

