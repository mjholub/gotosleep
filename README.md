~~Since running `swaylock && echo mem | sudo tee /sys/power/state` might fail because the interactive shell becomes unavailable as a result of running swaylock`, I needed to write a small program to work around this by using polkit/elogind.~~

nevermind, found a much simpler solution by tiling two terminal windows together and then acquiring doas persistence, then running `sleep 10 && echo mem | doas tee /sys/power/state` and in another one launching the locker during sleep. Still keeping this repo public as a possible reference in low-level programming in Go.
