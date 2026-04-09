tag = v1.1.887

build:
	git commit -am "f" && git push
	git tag $(tag)
	git push origin $(tag)


.PHONY: build