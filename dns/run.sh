docker exec -i nginx1 /bin/sh -c "cat >> /etc/hosts << EOF
10.10.10.10 s22.test.service
EOF"