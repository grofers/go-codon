if ! hash glide 2>/dev/null; then
	curl https://glide.sh/get | sh
fi
glide install
