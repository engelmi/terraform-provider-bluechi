#!/bin/bash -xe

PUBKEYPATH=~/.ssh/id_rsa.pub
PUBKEY=$( cat $PUBKEYPATH )

SCRIPT_DIR=$( realpath "$0"  )
SCRIPT_DIR=$(dirname "$SCRIPT_DIR")

function build_image(){
    if [[ "$1" == "bluechi" ]]; then
        podman build -t localhost/bluechi -f $SCRIPT_DIR/bluechi.image
    elif [[ "$1" == "centos" ]]; then
        podman build -t localhost/centos -f $SCRIPT_DIR/centos.image
    else
        echo "Unknown image: '$1'"
    fi
}

function start(){
    if [[ "$1" != "bluechi" && "$1" != "centos" ]]; then
        echo "Unknown container image: '$1'"
        exit 1
    fi
    # start all containers
    podman run -dt --rm --name main --network host localhost/$1:latest
    podman run -dt --rm --name worker1 --network host localhost/$1:latest
    podman run -dt --rm --name worker2 --network host localhost/$1:latest
    podman run -dt --rm --name worker3 --network host localhost/$1:latest

    # inject public key
    podman exec main bash -c "echo $PUBKEY >> ~/.ssh/authorized_keys"
    podman exec worker1 bash -c "echo $PUBKEY >> ~/.ssh/authorized_keys"
    podman exec worker2 bash -c "echo $PUBKEY >> ~/.ssh/authorized_keys"
    podman exec worker3 bash -c "echo $PUBKEY >> ~/.ssh/authorized_keys"

    # change the port for the ssh config
    podman exec main bash -c "echo 'Port 2020' >> /etc/ssh/sshd_config"
    podman exec main bash -c "systemctl restart sshd"

    podman exec worker1 bash -c "echo 'Port 2021' >> /etc/ssh/sshd_config"
    podman exec worker1 bash -c "systemctl restart sshd"

    podman exec worker2 bash -c "echo 'Port 2022' >> /etc/ssh/sshd_config"
    podman exec worker2 bash -c "systemctl restart sshd"

    podman exec worker3 bash -c "echo 'Port 2023' >> /etc/ssh/sshd_config"
    podman exec worker3 bash -c "systemctl restart sshd"
}

function stop() {
    podman stop main
    podman stop worker1
    podman stop worker2
    podman stop worker3
}

$1 $2
