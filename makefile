# Default target to generate Hugo site into /docs directory
default: publish

# Target to generate Hugo site
publish:
	hugo -d ./docs