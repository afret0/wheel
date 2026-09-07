tag = v1.1.906

build:
	git commit -am "build" && git push || true
	git tag $(tag)
	git push origin $(tag)


.PHONY: build