rsync -r -a -v -e "ssh -l ec2-user" --exclude ".git" --exclude "onchange.sh" --exclude "*.swo" --exclude "*.swp" --delete ~/code/RemindmeSlackIntegration/ reminder:~/go/src/RemindmeSlackIntegration
afplay "/System/Library/Sounds/Morse.aiff"

