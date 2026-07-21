#!/bin/bash

SUCCESS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0

COLOR_RED="\033[31m"
COLOR_GREEN="\033[32m"
COLOR_YELLOW="\033[33m"
COLOR_BLUE="\033[36m"
COLOR_RESET="\033[0m"

log_info() {
    echo -e "${COLOR_BLUE}[.]${COLOR_RESET} $1"
}

log_success() {
    echo -e "${COLOR_GREEN}[+]${COLOR_RESET} $1"
}

log_warn() {
    echo -e "${COLOR_YELLOW}[!]${COLOR_RESET} $1"
}

log_question() {
    echo -e "${COLOR_YELLOW}[?]${COLOR_RESET} $1"
}

log_error() {
    echo -e "${COLOR_RED}[-]${COLOR_RESET} $1"
}

separate_stars() {
    local cols=$(tput cols 2>/dev/null || echo 10)
    printf "%*s\n" "$cols" "" | tr ' ' '*'
}

install_docker() {
    log_info "Installing Docker..."

    if docker --version > /dev/null 2>&1; then
        log_info "Docker already installed"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        return
    fi

    sudo apt update
    sudo apt install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

    sudo apt update
    sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    sudo groupadd docker
    sudo usermod -aG docker $USER

    log_success "Docker installed"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
}

print_counts() {
    log_info "Task Summary:"
    log_success "  Success: $SUCCESS_COUNT"
    log_warn "  Skip: $SKIP_COUNT"
    log_error "  Fail: $FAIL_COUNT"
}

main() {
    install_docker
    separate_stars

    print_counts
    separate_stars

    log_info "Docker is installed. Please log out and log in again to use Docker without sudo."
}

main "$@"