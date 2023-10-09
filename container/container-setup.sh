#!/bin/bash -xe

function start(){
    podman run -dt --rm --name main --network host localhost/bluechi:latest
    podman run -dt --rm --name worker1 --network host localhost/bluechi:latest
    podman run -dt --rm --name worker2 --network host localhost/bluechi:latest
    podman run -dt --rm --name worker3 --network host localhost/bluechi:latest

    podman exec main bash -c "echo 'Port 2020' >> /etc/ssh/sshd_config"
    podman exec worker1 bash -c "echo 'Port 2021' >> /etc/ssh/sshd_config"
    podman exec worker2 bash -c "echo 'Port 2022' >> /etc/ssh/sshd_config"
    podman exec worker3 bash -c "echo 'Port 2023' >> /etc/ssh/sshd_config"

    podman exec main bash -c "systemctl restart sshd"
    podman exec worker1 bash -c "systemctl restart sshd"
    podman exec worker2 bash -c "systemctl restart sshd"
    podman exec worker3 bash -c "systemctl restart sshd"
}

function stop() {
    podman stop main
    podman stop worker1
    podman stop worker2
    podman stop worker3
}

$1
