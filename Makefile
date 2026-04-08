tag = v1.1.883

prod:
	git commit -am "f" && git push
	git tag $(tag)
	git push origin $(tag)


.PHONY: prod