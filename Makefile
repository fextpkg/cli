.SILENT:

install-builder:
	pip install setuptools wheel > /dev/null

build: install-builder
	python setup.py bdist_wheel --plat-name=${FEXT_PLATFORM_TAG} --dist-dir=${PY_DIST_DIRECTORY}