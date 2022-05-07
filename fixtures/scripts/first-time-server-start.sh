#!/usr/bin/expect -f

set PASSWORD [lindex $argv 0]

spawn /opt/pzserver/start-server.sh

expect "Enter new administrator password:\r"

send -- "$PASSWORD\r"

expect "Confirm the password:\r"

send -- "$PASSWORD\r"

expect "Administrator account 'admin' created.\r"