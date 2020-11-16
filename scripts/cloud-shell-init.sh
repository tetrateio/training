#!/bin/bash
set -e

red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
blue=`tput setaf 4`
magenta=`tput setaf 5`
cyan=`tput setaf 6`
reset=`tput sgr0`

parseCmdLine(){
    export GOOGLE_PROJECT=""
    export ISTIO_VERSION=1.7.4
	while true; do
		case "$1" in
			-p|--project) export GOOGLE_PROJECT=$2; shift 2;;
			-i|--istio-version) export ISTIO_VERSION=$2; shift 2;;
			*) [ $# -gt 0 ] && { echo "Unrecognized argument $1"; usage; return 1; }; break
		esac;
	done
    if [[ "${GOOGLE_PROJECT}" == "" ]]; then
      export GOOGLE_PROJECT=$(gcloud projects list --format json | jq -r ".[0].projectId")
		  echo "${green}Using GOOGLE_PROJECT ${yellow}'${GOOGLE_PROJECT}'${reset}"
    fi
    
	return 0
}

# main
parseCmdLine "$@"

# download and install istio
echo -e "\n${green}Installing Istio${yellow} ${ISTIO_VERSION}${reset}"
curl --silent -L https://istio.io/downloadIstio | ISTIO_VERSION="${ISTIO_VERSION}" sh - > /dev/null 2>&1
echo "${green}Istio CLI available at${green}${yellow} istioctl${reset}"

# Add Istio to PATH
export ISTIO_DIR=$PWD/istio-${ISTIO_VERSION}
export PATH=${ISTIO_DIR}/bin:$PATH

# setup gcloud project and cluster, cluster name is the same as project name
echo -e "\n${green}Setting google cloud project to ${GOOGLE_PROJECT} and cluster ${GOOGLE_PROJECT}${reset}"
gcloud config set project "${GOOGLE_PROJECT}" > /dev/null 2>&1
gcloud config set container/cluster "${GOOGLE_PROJECT}" > /dev/null 2>&1

# users should only have 1 cluster in their project
ZONE=$(gcloud container clusters list --format json | jq -r ".[0].locations[0]")
gcloud container clusters get-credentials "${GOOGLE_PROJECT}" --zone ${ZONE} > /dev/null 2>&1

# useful kubectl aliases
cat <<EOF

${green}Adding kubectl aliases
  ${yellow}kat='cat <<EOF | kubectl apply -f -'
  k='kubectl'
  kg='kubectl get'
  kgp='kubectl get pods'
  kd='kubectl describe'${reset}

  ${green}To use aliases run ${yellow}source ~/.bashrc${reset}

EOF
echo "alias kat='cat <<EOF | kubectl apply -f -'" >> ~/.bashrc
echo "alias k='kubectl'" >> ~/.bashrc
echo "alias kg='kubectl get'" >> ~/.bashrc
echo "alias kgp='kubectl get pods'" >> ~/.bashrc
echo "alias kd='kubectl describe'" >> ~/.bashrc
echo "function kubectl() { echo "+ kubectl \$@">&2; command kubectl \$@; }" >> ~/.bashrc

source ~/.bashrc

# install k9s (not enabled due to the fact it tanks the cloud shell vm)
# echo -e "\n\n${green}Installing k9s..."
# GO111MODULE=on go get github.com/derailed/k9s@v0.23.10 > /dev/null 2>&1
# echo "K9s (github.com/derailed/k9s) available with command ${yellow}k9s${reset}"

cat <<EOF

${green}If you would like to install k9s 
    ${yellow}GO111MODULE=on go get github.com/derailed/k9s@v0.23.10${reset}
EOF